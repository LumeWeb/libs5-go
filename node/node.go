package node

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/metadata"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/go-resty/resty/v2"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"log"
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
	httpClient            *resty.Client
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
		httpClient:            resty.New(),
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
func (n *NodeImpl) GetCachedStorageLocations(hash *encoding.Multihash, kinds []types.StorageLocationType) (map[string]interfaces.StorageLocation, error) {
	locations := make(map[string]interfaces.StorageLocation)

	locationMap, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return nil, err
	}
	if len(locationMap) == 0 {
		return make(map[string]interfaces.StorageLocation), nil
	}

	ts := time.Now().Unix()

	for _, t := range kinds {

		nodeMap, ok := (locationMap)[int(t)]
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

			storageLocation := storage.NewStorageLocation(int(t), addresses, expiry)
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
}

func (n *NodeImpl) DownloadBytesByHash(hash *encoding.Multihash) ([]byte, error) {
	// Initialize the download URI provider
	dlUriProvider := storage.NewStorageLocationProvider(n, hash)
	err := dlUriProvider.Start()
	if err != nil {
		return nil, err
	}

	retryCount := 0
	for {
		dlUri, err := dlUriProvider.Next()
		if err != nil {
			return nil, err
		}

		// Log the attempt
		log.Printf("[try] %s", dlUri.Location().BytesURL())

		res, err := n.httpClient.R().Get(dlUri.Location().BytesURL())
		if err != nil {
			err := dlUriProvider.Downvote(dlUri)
			if err != nil {
				return nil, err
			}
			retryCount++
			if retryCount > 32 {
				return nil, errors.New("too many retries")
			}
			continue
		}

		bodyBytes := res.Body()

		return bodyBytes, nil
	}
}

func (n *NodeImpl) GetMetadataByCID(cid *encoding.CID) (metadata.Metadata, error) {
	var md metadata.Metadata

	hashStr, err := cid.Hash.ToString()
	if err != nil {
		return nil, err
	}

	if n.metadataCache.Contains(hashStr) {
		bytes, err := n.DownloadBytesByHash(&cid.Hash)
		if err != nil {
			return nil, err
		}

		switch cid.Type {
		case types.CIDTypeMetadataMedia, types.CIDTypeBridge: // Both cases use the same deserialization method
			md = metadata.NewEmptyMediaMetadata()

			err = msgpack.Unmarshal(bytes, md)
			if err != nil {
				return nil, err
			}
		case types.CIDTypeMetadataWebapp:
			md = metadata.NewEmptyWebAppMetadata()

			err = msgpack.Unmarshal(bytes, md)
			if err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("unsupported metadata format")
		}

		n.metadataCache.Put(hashStr, md)
	}

	return md, nil
}
