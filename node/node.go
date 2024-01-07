package node

import (
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

var _ interfaces.Node = (*NodeImpl)(nil)

const cacheBucketName = "object-cache"

type NodeImpl struct {
	nodeConfig            *config.NodeConfig
	metadataCache         structs.Map
	started               bool
	hashQueryRoutingTable structs.Map
	services              interfaces.Services
	cacheBucket           *bolt.Bucket
}

func (n *NodeImpl) NetworkId() string {
	return n.nodeConfig.P2P.Network
}

func (n *NodeImpl) Services() interfaces.Services {
	return n.services
}

func NewNode(config *config.NodeConfig) interfaces.Node {
	n := &NodeImpl{
		nodeConfig:            config,
		metadataCache:         structs.NewMap(),
		started:               false,
		hashQueryRoutingTable: structs.NewMap(),
	}
	n.services = NewServices(service.NewP2P(n))

	return n
}
func (n *NodeImpl) HashQueryRoutingTable() structs.Map {
	return n.hashQueryRoutingTable
}

func (n *NodeImpl) IsStarted() bool {
	return n.started
}

func (n *NodeImpl) Config() *config.NodeConfig {
	return n.nodeConfig
}

func (n *NodeImpl) Logger() *zap.Logger {
	if n.nodeConfig != nil {
		return n.nodeConfig.Logger
	}
	return nil
}

func (n *NodeImpl) Db() *bolt.DB {
	if n.nodeConfig != nil {
		return n.nodeConfig.DB
	}
	return nil
}

func (n *NodeImpl) Start() error {
	err :=
		utils.CreateBucket(cacheBucketName, n.Db(), func(bucket *bolt.Bucket) {
			n.cacheBucket = bucket
		})

	if err != nil {
		return err
	}

	n.started = true

	err = n.Services().P2P().Init()
	if err != nil {
		return err
	}

	err = n.Services().P2P().Start()
	if err != nil {
		return err
	}

	n.started = true
	return nil
}

/*
	func (n *NodeImpl) Services() *S5Services {
		if n.nodeConfig != nil {
			return n.nodeConfig.Services
		}
		return nil
	}

	func (n *NodeImpl) Start() error {
		n.started = true
		return nil
	}

	func (n *NodeImpl) Stop() error {
		n.started = false
		return nil
	}
*/
func (n *NodeImpl) GetCachedStorageLocations(hash *encoding.Multihash, types []int) (map[string]interfaces.StorageLocation, error) {
	locations := make(map[string]interfaces.StorageLocation)

	locationMap, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return nil, err
	}
	if len(locationMap) == 0 {
		return make(map[string]interfaces.StorageLocation), nil
	}

	ts := time.Now().Unix()

	for _, t := range types {

		nodeMap, ok := (locationMap)[t]
		if !ok {
			continue
		}

		for key, value := range nodeMap {
			if len(value) < 4 {
				continue
			}

			expiry, ok := value[3].(int64)
			if !ok || expiry < ts {
				continue
			}

			addresses, ok := value[1].([]string)
			if !ok {
				continue
			}

			storageLocation := storage.NewStorageLocation(t, addresses, expiry)
			if len(value) > 4 {
				if providerMessage, ok := value[4].([]byte); ok {
					(storageLocation).SetProviderMessage(providerMessage)
				}
			}

			locations[key] = storageLocation
		}
	}
	return locations, nil
}
func (n *NodeImpl) readStorageLocationsFromDB(hash *encoding.Multihash) (storage.StorageLocationMap, error) {
	locationMap := storage.NewStorageLocationMap()

	bytes := n.cacheBucket.Get(hash.FullBytes())
	if bytes == nil {
		return locationMap, nil
	}

	err := msgpack.Unmarshal(bytes, locationMap)
	if err != nil {
		return nil, err
	}

	return locationMap, nil
}
func (n *NodeImpl) AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location interfaces.StorageLocation, message []byte, config *config.NodeConfig) error {
	// Read existing storage locations
	locationDb, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return err
	}

	nodeIdStr, err := nodeId.ToString()
	if err != nil {
		return err
	}

	// Get or create the inner map for the specific type
	innerMap, exists := locationDb[location.Type()]
	if !exists {
		innerMap = make(storage.NodeStorage, 1)
		innerMap[nodeIdStr] = make(storage.NodeDetailsStorage, 1)
	}

	// Create location map with new data
	locationMap := make(map[int]interface{}, 3)
	locationMap[1] = location.Parts()
	locationMap[3] = location.Expiry()
	locationMap[4] = message

	// Update the inner map with the new location
	innerMap[nodeIdStr] = locationMap
	locationDb[location.Type()] = innerMap

	// Serialize the updated map and store it in the database
	packedBytes, err := msgpack.Marshal(locationDb)
	if err != nil {
		return err
	}

	err = n.cacheBucket.Put(hash.FullBytes(), packedBytes)
	if err != nil {
		return err
	}

	return nil
} /*

func (n *NodeImpl) DownloadBytesByHash(hash Multihash) ([]byte, error) {
	dlUriProvider := NewStorageLocationProvider(n, hash, []int{storageLocationTypeFull, storageLocationTypeFile})
	dlUriProvider.Start()

	retryCount := 0
	for {
		dlUri, err := dlUriProvider.Next()
		if err != nil {
			return nil, err
		}

		n.Logger.Verbose(fmt.Sprintf("[try] %s", dlUri.Location.BytesUrl))

		client := &http.Client{
			Timeout: 30 * time.Second,
		}
		res, err := client.Get(dlUri.Location.BytesUrl)
		if err != nil {
			n.Logger.Catched(err)

			dlUriProvider.Downvote(dlUri)

			retryCount++
			if retryCount > 32 {
				return nil, errors.New("too many retries")
			}
			continue
		}
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// Assuming blake3 and equalBytes functions are available
		resHash := blake3(data)

		if !equalBytes(hash.HashBytes, resHash) {
			dlUriProvider.Downvote(dlUri)
			continue
		}

		dlUriProvider.Upvote(dlUri)
		return data, nil
	}
}

func (n *NodeImpl) GetMetadataByCID(cid CID) (Metadata, error) {
	var metadata Metadata
	var ok bool

	if metadata, ok = n.MetadataCache[cid.Hash]; !ok {
		bytes, err := n.DownloadBytesByHash(cid.Hash)
		if err != nil {
			return Metadata{}, err
		}

		switch cid.kind {
		case METADATA_MEDIA, BRIDGE: // Both cases use the same deserialization method
			metadata, err = deserializeMediaMetadata(bytes)
		case METADATA_WEBAPP:
			metadata, err = deserializeWebAppMetadata(bytes)
		default:
			return Metadata{}, errors.New("unsupported metadata format")
		}

		if err != nil {
			return Metadata{}, err
		}

		n.MetadataCache[cid.Hash] = metadata
	}

	return metadata, nil
}
*/
