package storage

import (
	"bytes"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	_node "git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

var (
	_ msgpack.CustomDecoder   = (*StorageLocationMap)(nil)
	_ msgpack.CustomEncoder   = (*StorageLocationMap)(nil)
	_ StorageLocation         = (*StorageLocationImpl)(nil)
	_ StorageLocationProvider = (*StorageLocationProviderImpl)(nil)
	_ SignedStorageLocation   = (*SignedStorageLocationImpl)(nil)
)

type StorageLocationMap map[int]NodeStorage
type NodeStorage map[string]NodeDetailsStorage
type NodeDetailsStorage map[int]interface{}

type StorageLocationImpl struct {
	kind            int
	parts           []string
	binaryParts     [][]byte
	expiry          int64
	providerMessage []byte
}

func (s *StorageLocationImpl) Type() int {
	return s.kind
}

func (s *StorageLocationImpl) Parts() []string {
	return s.parts
}

func (s *StorageLocationImpl) BinaryParts() [][]byte {
	return s.binaryParts
}

func (s *StorageLocationImpl) Expiry() int64 {
	return s.expiry
}

func (s *StorageLocationImpl) SetType(t int) {
	s.kind = t
}

func (s *StorageLocationImpl) SetParts(p []string) {
	s.parts = p
}

func (s *StorageLocationImpl) SetBinaryParts(bp [][]byte) {
	s.binaryParts = bp
}

func (s *StorageLocationImpl) SetExpiry(e int64) {
	s.expiry = e
}

func (s *StorageLocationImpl) SetProviderMessage(msg []byte) {
	s.providerMessage = msg
}

func (s *StorageLocationImpl) ProviderMessage() []byte {
	return s.providerMessage
}

func NewStorageLocation(Type int, Parts []string, Expiry int64) StorageLocation {
	return &StorageLocationImpl{
		kind:   Type,
		parts:  Parts,
		expiry: Expiry,
	}
}

func (s *StorageLocationImpl) BytesURL() string {
	return s.parts[0]
}

func (s *StorageLocationImpl) OutboardBytesURL() string {
	if len(s.parts) == 1 {
		return s.parts[0] + ".obao"
	}
	return s.parts[1]
}

func (s *StorageLocationImpl) String() string {
	expiryDate := time.Unix(s.expiry, 0)
	return "StorageLocationImpl(" + strconv.Itoa(s.Type()) + ", " + fmt.Sprint(s.parts) + ", expiry: " + expiryDate.Format(time.RFC3339) + ")"
}

type SignedStorageLocationImpl struct {
	nodeID   *encoding.NodeId
	location StorageLocation
}

func NewSignedStorageLocation(NodeID *encoding.NodeId, Location StorageLocation) SignedStorageLocation {
	return &SignedStorageLocationImpl{
		nodeID:   NodeID,
		location: Location,
	}
}

func (ssl *SignedStorageLocationImpl) String() string {
	nodeString, _ := ssl.nodeID.ToString()

	if nodeString == "" {
		nodeString = "failed to decode node id"
	}

	return "SignedStorageLocationImpl(" + ssl.location.String() + ", " + nodeString + ")"
}

func (ssl *SignedStorageLocationImpl) NodeId() *encoding.NodeId {
	return ssl.nodeID
}
func (ssl *SignedStorageLocationImpl) Location() StorageLocation {
	return ssl.location
}

func (s *StorageLocationMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	if *s == nil {
		*s = make(StorageLocationMap)
	}

	// Decode directly into a temp map
	temp := make(map[int]map[string]map[int]interface{})
	err := dec.Decode(&temp)
	if err != nil {
		return fmt.Errorf("error decoding msgpack: %w", err)
	}

	// Convert temp map to StorageLocationMap
	for k, v := range temp {
		nodeStorage, exists := (*s)[k]
		if !exists {
			nodeStorage = make(NodeStorage, len(v)) // preallocate if size is known
			(*s)[k] = nodeStorage
		}

		for nk, nv := range v {
			nodeDetailsStorage, exists := nodeStorage[nk]
			if !exists {
				nodeDetailsStorage = make(NodeDetailsStorage, len(nv)) // preallocate if size is known
				nodeStorage[nk] = nodeDetailsStorage
			}

			for ndk, ndv := range nv {
				nodeDetailsStorage[ndk] = ndv
			}
		}
	}

	return nil
}

func (s StorageLocationMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	// Create a temporary map to hold the encoded data
	tempMap := make(map[int]map[string]map[int]interface{})

	// Populate the temporary map with data from storageLocationMap
	for storageKey, nodeStorages := range s {
		tempNodeStorages := make(map[string]map[int]interface{})
		for nodeId, nodeDetails := range nodeStorages {
			tempNodeStorages[nodeId] = nodeDetails
		}
		tempMap[storageKey] = tempNodeStorages
	}

	// Encode the temporary map using MessagePack
	return enc.Encode(tempMap)
}

func NewStorageLocationMap() StorageLocationMap {
	return StorageLocationMap{}
}

type StorageLocationProviderImpl struct {
	node            *_node.Node
	hash            *encoding.Multihash
	types           []types.StorageLocationType
	timeoutDuration time.Duration
	availableNodes  []*encoding.NodeId
	uris            map[string]StorageLocation
	timeout         time.Time
	isTimedOut      bool
	isWaitingForUri bool
	mutex           sync.Mutex
}

func (s *StorageLocationProviderImpl) Start() error {
	var err error

	s.uris, err = s.node.GetCachedStorageLocations(s.hash, s.types)
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

	s.availableNodes, err = s.node.Services().P2P().SortNodesByScore(s.availableNodes)
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

			newUris, err := s.node.GetCachedStorageLocations(s.hash, s.types)
			if err != nil {
				s.mutex.Unlock()
				break
			}

			//s.node.Logger().Debug("New URIs", zap.Any("uris", newUris), zap.Any("availableNodes", s.availableNodes), zap.Any("timeout", s.timeout), zap.Any("isTimedOut", s.isTimedOut), zap.Any("isWaitingForUri", s.isWaitingForUri), zap.Any("requestSent", requestSent))
			if len(s.availableNodes) == 0 && len(newUris) < 2 && !requestSent {
				s.node.Logger().Debug("Sending hash request")
				err := s.node.Services().P2P().SendHashRequest(s.hash, s.types)
				if err != nil {
					s.node.Logger().Error("Error sending hash request", zap.Error(err))
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
						s.node.Logger().Error("Error decoding node id", zap.Error(err))
						continue
					}
					if !containsNode(s.availableNodes, nodeId) {
						s.availableNodes = append(s.availableNodes, nodeId)
						hasNewNode = true
					}
				}
			}

			if hasNewNode {
				score, err := s.node.Services().P2P().SortNodesByScore(s.availableNodes)
				if err != nil {
					s.node.Logger().Error("Error sorting nodes by score", zap.Error(err))
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
func (s *StorageLocationProviderImpl) Next() (SignedStorageLocation, error) {
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
				s.node.Logger().Error("Could not find uri for node id", zap.String("nodeId", nodIdStr))
				continue
			}

			return NewSignedStorageLocation(nodeId, uri), nil
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

func (s *StorageLocationProviderImpl) Upvote(uri SignedStorageLocation) error {
	err := s.node.Services().P2P().UpVote(uri.NodeId())
	if err != nil {
		return err
	}

	return nil
}

func (s *StorageLocationProviderImpl) Downvote(uri SignedStorageLocation) error {
	err := s.node.Services().P2P().DownVote(uri.NodeId())
	if err != nil {
		return err
	}
	return nil
}

func NewStorageLocationProvider(node *_node.Node, hash *encoding.Multihash, locationTypes ...types.StorageLocationType) StorageLocationProvider {
	if locationTypes == nil {
		locationTypes = []types.StorageLocationType{
			types.StorageLocationTypeFull,
		}
	}

	return &StorageLocationProviderImpl{
		node:            node,
		hash:            hash,
		types:           locationTypes,
		timeoutDuration: 60 * time.Second,
		uris:            make(map[string]StorageLocation),
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

type StorageLocationProvider interface {
	Start() error
	Next() (SignedStorageLocation, error)
	Upvote(uri SignedStorageLocation) error
	Downvote(uri SignedStorageLocation) error
}

type StorageLocation interface {
	BytesURL() string
	OutboardBytesURL() string
	String() string
	ProviderMessage() []byte
	Type() int
	Parts() []string
	BinaryParts() [][]byte
	Expiry() int64
	SetProviderMessage(msg []byte)
	SetType(t int)
	SetParts(p []string)
	SetBinaryParts(bp [][]byte)
	SetExpiry(e int64)
}
type SignedStorageLocation interface {
	String() string
	NodeId() *encoding.NodeId
	Location() StorageLocation
}

type ProviderStore interface {
	CanProvide(hash *encoding.Multihash, kind []types.StorageLocationType) bool
	Provide(hash *encoding.Multihash, kind []types.StorageLocationType) (StorageLocation, error)
}
