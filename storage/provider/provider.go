package provider

import (
	"bytes"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"sync"
	"time"
)

var _ storage.StorageLocationProvider = (*StorageLocationProviderImpl)(nil)

type StorageLocationProviderImpl struct {
	services        storage.StorageLocationProviderServices
	hash            *encoding.Multihash
	types           []types.StorageLocationType
	timeoutDuration time.Duration
	availableNodes  []*encoding.NodeId
	uris            map[string]storage.StorageLocation
	timeout         time.Time
	isTimedOut      bool
	isWaitingForUri bool
	mutex           sync.Mutex
	logger          *zap.Logger
}

func (s *StorageLocationProviderImpl) Start() error {
	var err error

	s.uris, err = s.services.Storage().GetCachedStorageLocations(s.hash, s.types)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.availableNodes = make([]*encoding.NodeId, 0, len(s.uris))
	for k := range s.uris {
		nodeId, err := encoding.DecodeNodeId(k)
		if err != nil {
			continue
		}

		s.availableNodes = append(s.availableNodes, nodeId)
	}

	s.availableNodes, err = s.services.P2P().SortNodesByScore(s.availableNodes)
	if err != nil {
		s.mutex.Unlock()
		return err
	}

	s.timeout = time.Now().Add(s.timeoutDuration)
	s.isTimedOut = false
	s.mutex.Unlock()
	go func() {
		requestSent := false

		for {
			s.mutex.Lock()
			if time.Now().After(s.timeout) {
				s.isTimedOut = true
				s.mutex.Unlock()
				break
			}

			newUris, err := s.services.Storage().GetCachedStorageLocations(s.hash, s.types)
			if err != nil {
				s.mutex.Unlock()
				break
			}

			if len(s.availableNodes) == 0 && len(newUris) < 2 && !requestSent {
				s.logger.Debug("Sending hash request")
				err := s.services.P2P().SendHashRequest(s.hash, s.types)
				if err != nil {
					s.logger.Error("Error sending hash request", zap.Error(err))
					continue
				}
				requestSent = true
			}

			hasNewNode := false
			for k, v := range newUris {
				if _, exists := s.uris[k]; !exists || s.uris[k] != v {
					s.uris[k] = v
					nodeId, err := encoding.DecodeNodeId(k)
					if err != nil {
						s.logger.Error("Error decoding node id", zap.Error(err))
						continue
					}
					if !containsNode(s.availableNodes, nodeId) {
						s.availableNodes = append(s.availableNodes, nodeId)
						hasNewNode = true
					}
				}
			}

			if hasNewNode {
				score, err := s.services.P2P().SortNodesByScore(s.availableNodes)
				if err != nil {
					s.logger.Error("Error sorting nodes by score", zap.Error(err))
				} else {
					s.availableNodes = score
				}
			}
			s.mutex.Unlock()

			time.Sleep(10 * time.Millisecond)
		}
	}()
	return nil
}
func (s *StorageLocationProviderImpl) Next() (storage.SignedStorageLocation, error) {
	s.timeout = time.Now().Add(s.timeoutDuration)

	for {
		if len(s.availableNodes) > 0 {
			s.isWaitingForUri = false
			nodeId := s.availableNodes[0]
			s.availableNodes = s.availableNodes[1:]

			nodIdStr, err := nodeId.ToString()
			if err != nil {
				return nil, err
			}

			uri, exists := s.uris[nodIdStr]
			if !exists {
				s.logger.Error("Could not find uri for node id", zap.String("nodeId", nodIdStr))
				continue
			}

			return storage.NewSignedStorageLocation(nodeId, uri), nil
		}

		s.isWaitingForUri = true
		if s.isTimedOut {
			hashStr, err := s.hash.ToString()
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("Could not download raw file: Timed out after %s %s", s.timeoutDuration.String(), hashStr)
		}

		time.Sleep(10 * time.Millisecond) // Replace with a proper wait/notify mechanism if applicable
	}
}

func (s *StorageLocationProviderImpl) Upvote(uri storage.SignedStorageLocation) error {
	err := s.services.P2P().UpVote(uri.NodeId())
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageLocationProviderImpl) Downvote(uri storage.SignedStorageLocation) error {
	err := s.services.P2P().DownVote(uri.NodeId())
	if err != nil {
		return err
	}
	return nil
}

func NewStorageLocationProvider(params StorageLocationProviderParams) *StorageLocationProviderImpl {
	if params.LocationTypes == nil {
		params.LocationTypes = []types.StorageLocationType{
			types.StorageLocationTypeFull,
		}
	}

	return &StorageLocationProviderImpl{
		services:        params.Services,
		hash:            params.Hash,
		types:           params.LocationTypes,
		timeoutDuration: 60 * time.Second,
		uris:            make(map[string]storage.StorageLocation),
		logger:          params.Logger,
	}
}
func containsNode(slice []*encoding.NodeId, item *encoding.NodeId) bool {
	for _, v := range slice {
		if bytes.Equal(v.Bytes(), item.Bytes()) {
			return true
		}
	}
	return false
}

type StorageLocationProviderParams struct {
	Services      storage.StorageLocationProviderServices
	Hash          *encoding.Multihash
	LocationTypes []types.StorageLocationType
	Logger        *zap.Logger
	Config        *config.NodeConfig
	Db            *bolt.DB
}
