package libs5_go

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

type Metadata interface {
	ToJson() map[string]interface{}
}

type Services struct {
	p2p *service.P2P
}

func (s *Services) P2P() *service.P2P {
	return s.p2p
}

const cacheBucketName = "object-cache"

type Node struct {
	nodeConfig            *NodeConfig
	metadataCache         *structs.Map
	started               bool
	hashQueryRoutingTable *structs.Map
	services              Services
	cacheBucket           *bolt.Bucket
}

func (n *Node) Services() *Services {
	return &n.services
}

func NewNode(config *NodeConfig) *Node {
	return &Node{
		nodeConfig:            config,
		metadataCache:         structs.NewMap(),
		started:               false,
		hashQueryRoutingTable: structs.NewMap(),
	}
}
func (n *Node) HashQueryRoutingTable() *structs.Map {
	return n.hashQueryRoutingTable
}

func (n *Node) IsStarted() bool {
	return n.started
}

func (n *Node) Config() *NodeConfig {
	return n.nodeConfig
}

func (n *Node) Logger() *zap.Logger {
	if n.nodeConfig != nil {
		return n.nodeConfig.Logger
	}
	return nil
}

func (n *Node) Db() *bolt.DB {
	if n.nodeConfig != nil {
		return &n.nodeConfig.DB
	}
	return nil
}

func (n *Node) Start() error {
	err :=
		utils.CreateBucket(cacheBucketName, n.Db(), func(bucket *bolt.Bucket) {
			n.cacheBucket = bucket
		})

	if err != nil {
		return err
	}

	n.started = true
	return nil
}

/*
	func (n *Node) Services() *S5Services {
		if n.nodeConfig != nil {
			return n.nodeConfig.Services
		}
		return nil
	}

	func (n *Node) Start() error {
		n.started = true
		return nil
	}

	func (n *Node) Stop() error {
		n.started = false
		return nil
	}
*/
func (n *Node) GetCachedStorageLocations(hash *encoding.Multihash, types []int) (map[encoding.NodeIdCode]*StorageLocation, error) {
	locations := make(map[encoding.NodeIdCode]*StorageLocation)

	locationMap, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return nil, err
	}
	if len(locationMap) == 0 {
		return make(map[encoding.NodeIdCode]*StorageLocation), nil
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

			storageLocation := NewStorageLocation(t, addresses, expiry)
			if len(value) > 4 {
				if providerMessage, ok := value[4].([]byte); ok {
					storageLocation.ProviderMessage = providerMessage
				}
			}

			locations[key] = storageLocation
		}
	}
	return locations, nil
}
func (n *Node) readStorageLocationsFromDB(hash *encoding.Multihash) (storageLocationMap, error) {
	locationMap := newStorageLocationMap()

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
func (n *Node) AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location *StorageLocation, message []byte, config *NodeConfig) error {
	// Read existing storage locations
	locationDb, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return err
	}

	// Get or create the inner map for the specific type
	innerMap, exists := locationDb[location.Type]
	if !exists {
		innerMap = make(nodeStorage, 1)
		innerMap[nodeId.HashCode()] = make(nodeDetailsStorage, 1)
	}

	// Create location map with new data
	locationMap := make(map[int]interface{}, 3)
	locationMap[1] = location.Parts
	locationMap[3] = location.Expiry
	locationMap[4] = message

	// Update the inner map with the new location
	innerMap[nodeId.HashCode()] = locationMap
	locationDb[location.Type] = innerMap

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

func (n *Node) DownloadBytesByHash(hash Multihash) ([]byte, error) {
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

func (n *Node) GetMetadataByCID(cid CID) (Metadata, error) {
	var metadata Metadata
	var ok bool

	if metadata, ok = n.MetadataCache[cid.Hash]; !ok {
		bytes, err := n.DownloadBytesByHash(cid.Hash)
		if err != nil {
			return Metadata{}, err
		}

		switch cid.Type {
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
