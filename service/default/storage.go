package _default

import (
	"context"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/metadata"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/storage/provider"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/go-resty/resty/v2"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"time"
)

const cacheBucketName = "object-cache"

var (
	_ service.Service        = (*StorageService)(nil)
	_ service.StorageService = (*StorageService)(nil)
)

type StorageService struct {
	httpClient    *resty.Client
	metadataCache structs.Map
	providerStore storage.ProviderStore
	service.ServiceBase
}

func NewStorage(params service.ServiceParams) *StorageService {
	return &StorageService{
		httpClient:    resty.New(),
		metadataCache: structs.NewMap(),
		ServiceBase:   service.NewServiceBase(params.Logger, params.Config, params.Db),
	}
}

func (s *StorageService) Start(ctx context.Context) error {
	err :=
		utils.CreateBucket(cacheBucketName, s.Db())

	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) Stop(ctx context.Context) error {
	return nil
}

func (s *StorageService) Init(ctx context.Context) error {
	return nil
}

func (n *StorageService) SetProviderStore(store storage.ProviderStore) {
	n.providerStore = store
}

func (n *StorageService) ProviderStore() storage.ProviderStore {
	return n.providerStore
}

func (s *StorageService) GetCachedStorageLocations(hash *encoding.Multihash, kinds []types.StorageLocationType) (map[string]storage.StorageLocation, error) {
	locations := make(map[string]storage.StorageLocation)

	locationMap, err := s.readStorageLocationsFromDB(hash)
	if err != nil {
		return nil, err
	}
	if len(locationMap) == 0 {
		return make(map[string]storage.StorageLocation), nil
	}

	ts := time.Now().Unix()

	for _, t := range kinds {
		nodeMap, ok := (locationMap)[int(t)]
		if !ok {
			continue
		}

		for key, value := range nodeMap {
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

			if providerMessage, ok := value[4].([]byte); ok {
				(storageLocation).SetProviderMessage(providerMessage)
			}

			locations[key] = storageLocation
		}
	}
	return locations, nil
}
func (s *StorageService) readStorageLocationsFromDB(hash *encoding.Multihash) (storage.StorageLocationMap, error) {
	var locationMap storage.StorageLocationMap

	err := s.Db().View(func(tx *bbolt.Tx) error {
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

func (s *StorageService) AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location storage.StorageLocation, message []byte) error {
	// Read existing storage locations
	locationDb, err := s.readStorageLocationsFromDB(hash)
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
	err = s.Db().Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(cacheBucketName))

		return b.Put(hash.FullBytes(), packedBytes)
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageService) DownloadBytesByHash(hash *encoding.Multihash) ([]byte, error) {
	// Initialize the download URI provider
	dlUriProvider := provider.NewStorageLocationProvider(provider.StorageLocationProviderParams{
		Services: s.Services(),
		Hash:     hash,
		LocationTypes: []types.StorageLocationType{
			types.StorageLocationTypeFull,
			types.StorageLocationTypeFile,
		},
		ServiceParams: service.ServiceParams{
			Logger: s.Logger(),
			Config: s.Config(),
			Db:     s.Db(),
		},
	})
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

		s.Logger().Debug("Trying to download from", zap.String("url", dlUri.Location().BytesURL()))

		res, err := s.httpClient.R().Get(dlUri.Location().BytesURL())
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

func (s *StorageService) DownloadBytesByCID(cid *encoding.CID) (bytes []byte, err error) {
	bytes, err = s.DownloadBytesByHash(&cid.Hash)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (s *StorageService) GetMetadataByCID(cid *encoding.CID) (md metadata.Metadata, err error) {
	hashStr, err := cid.Hash.ToString()
	if err != nil {
		return nil, err
	}

	if s.metadataCache.Contains(hashStr) {
		md, _ := s.metadataCache.Get(hashStr)

		return md.(metadata.Metadata), nil
	}

	bytes, err := s.DownloadBytesByHash(&cid.Hash)
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
	case types.CIDTypeDirectory:
		md = metadata.NewEmptyDirectoryMetadata()

		err = msgpack.Unmarshal(bytes, md)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported metadata format")
	}

	s.metadataCache.Put(hashStr, md)

	return md, nil
}
