package node

import (
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/metadata"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/go-resty/resty/v2"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"sync"
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
	httpClient            *resty.Client
	connections           sync.WaitGroup
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
	protocol.Init()
	signed.Init()
	err :=
		utils.CreateBucket(cacheBucketName, n.Db())

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

			addressesInterface, ok := value[1].([]interface{})
			if !ok {
				continue
			}

			// Create a slice to hold the strings
			addresses := make([]string, len(addressesInterface))

			// Convert each element to string
			for i, v := range addressesInterface {
				str, ok := v.(string)
				if !ok {
					// Handle the error, maybe skip this element or set a default value
					continue
				}
				addresses[i] = str
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
	var locationMap storage.StorageLocationMap

	err := n.Db().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cacheBucketName)) // Replace with your actual bucket name
		if b == nil {
			return fmt.Errorf("bucket %s not found", cacheBucketName)
		}

		bytes := b.Get(hash.FullBytes())
		if bytes == nil {
			// If no data found, return an empty locationMap but no error
			locationMap = storage.NewStorageLocationMap()
			return nil
		}

		return msgpack.Unmarshal(bytes, &locationMap)
	})

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
	err = n.Db().Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(cacheBucketName))

		return b.Put(hash.FullBytes(), packedBytes)
	})
	if err != nil {
		return err
	}

	return nil
}

func (n *NodeImpl) DownloadBytesByHash(hash *encoding.Multihash) ([]byte, error) {
	// Initialize the download URI provider
	dlUriProvider := storage.NewStorageLocationProvider(n, hash, types.StorageLocationTypeFull, types.StorageLocationTypeFile)
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

		n.Logger().Debug("Trying to download from", zap.String("url", dlUri.Location().BytesURL()))

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

func (n *NodeImpl) GetMetadataByCID(cid *encoding.CID) (md metadata.Metadata, err error) {
	hashStr, err := cid.Hash.ToString()
	if err != nil {
		return nil, err
	}

	if n.metadataCache.Contains(hashStr) {
		md, _ := n.metadataCache.Get(hashStr)

		return md.(metadata.Metadata), nil
	}

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

	return md, nil
}
func (n *NodeImpl) WaitOnConnectedPeers() {
	n.connections.Wait()
}

func (n *NodeImpl) ConnectionTracker() *sync.WaitGroup {
	return &n.connections
}
