package libs5_go

import (
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type Metadata interface {
	ToJson() map[string]interface{}
}

type services struct {
	p2p *service.P2P
}

type Node struct {
	nodeConfig            *NodeConfig
	metadataCache         *structs.Map
	started               bool
	hashQueryRoutingTable *structs.Map
	services              services
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
func (n *Node) GetCachedStorageLocations(hash Multihash, types []int) (map[NodeId]*StorageLocation, error) {
	locations := make(map[NodeId]*StorageLocation)

	mapFromDB, err := n.readStorageLocationsFromDB(hash)
	if err != nil {
		return nil, err
	}
	if len(mapFromDB) == 0 {
		return make(map[NodeId]*StorageLocation), nil
	}

	ts := time.Now().Unix()

	for _, t := range types {
		nodeMap, ok := mapFromDB[t]
		if !ok {
			continue
		}

		for key, value := range nodeMap {
			if len(value) < 4 {
				continue // or handle error
			}

			expiry, ok := value[3].(int64)
			if !ok || expiry < ts {
				continue
			}

			addresses, ok := value[1].([]string)
			if !ok {
				continue // or handle error
			}

			storageLocation := NewStorageLocation(t, addresses, expiry)
			if len(value) > 4 {
				if providerMessage, ok := value[4].(string); ok {
					storageLocation.ProviderMessage = providerMessage
				}
			}

			locations[NodeId(key)] = storageLocation
		}
	}

	return locations, nil
}
func (n *Node) ReadStorageLocationsFromDB(hash Multihash) (map[int]map[NodeId]map[int]interface{}, error) {
	locations := make(map[int]map[NodeId]map[int]interface{})

	bytes, err := n.config.CacheDb.Get(StringifyHash(hash)) // Assume StringifyHash and CacheDb.Get are implemented
	if err != nil {
		return locations, nil
	}
	if bytes == nil {
		return locations, nil
	}

	unpacker := NewUnpacker(bytes) // Assume NewUnpacker is implemented to handle the unpacking
	mapLength, err := unpacker.UnpackMapLength()
	if err != nil {
		return nil, err
	}

	for i := 0; i < mapLength; i++ {
		t, err := unpacker.UnpackInt()
		if err != nil {
			continue // or handle error
		}

		innerMap := make(map[NodeId]map[int]interface{})
		locations[t] = innerMap

		innerMapLength, err := unpacker.UnpackMapLength()
		if err != nil {
			continue // or handle error
		}

		for j := 0; j < innerMapLength; j++ {
			nodeIdBytes, err := unpacker.UnpackBinary()
			if err != nil {
				continue // or handle error
			}
			nodeId := NodeId(nodeIdBytes)

			// Assuming unpacker.UnpackMap() returns a map[string]interface{} and is implemented
			unpackedMap, err := unpacker.UnpackMap()
			if err != nil {
				continue // or handle error
			}

			convertedMap := make(map[int]interface{})
			for key, value := range unpackedMap {
				intKey, err := strconv.Atoi(key)
				if err != nil {
					continue // or handle error
				}
				convertedMap[intKey] = value
			}
			innerMap[nodeId] = convertedMap
		}
	}
	return locations, nil
}
func (n *Node) AddStorageLocation(hash Multihash, nodeId NodeId, location StorageLocation, message []byte, config S5Config) error {
	// Read existing storage locations
	mapFromDB, err := n.ReadStorageLocationsFromDB(hash)
	if err != nil {
		return err
	}

	// Get or create the inner map for the specific type
	innerMap, exists := mapFromDB[location.Type]
	if !exists {
		innerMap = make(map[NodeId]map[int]interface{})
		mapFromDB[location.Type] = innerMap
	}

	// Create location map with new data
	locationMap := make(map[int]interface{})
	locationMap[1] = location.Parts
	// locationMap[2] = location.BinaryParts // Uncomment if BinaryParts is a field of StorageLocation
	locationMap[3] = location.Expiry
	locationMap[4] = message

	// Update the inner map with the new location
	innerMap[nodeId] = locationMap

	// Serialize the updated map and store it in the database
	packedBytes, err := NewPacker().Pack(mapFromDB) // Assuming NewPacker and Pack are implemented
	if err != nil {
		return err
	}

	err = config.CacheDb.Put(StringifyHash(hash), packedBytes) // Assume CacheDb.Put and StringifyHash are implemented
	if err != nil {
		return err
	}

	return nil
}

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
