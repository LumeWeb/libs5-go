package node

import (
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"time"
)

var (
	_ msgpack.CustomDecoder = (*storageLocationMap)(nil)

	_ msgpack.CustomEncoder      = (*storageLocationMap)(nil)
	_ interfaces.StorageLocation = (*StorageLocationImpl)(nil)
)

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
	//TODO implement me
	panic("implement me")
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

func NewStorageLocation(Type int, Parts []string, Expiry int64) *interfaces.StorageLocation {
	sl := &StorageLocationImpl{
		kind:   Type,
		parts:  Parts,
		expiry: Expiry,
	}
	var location interfaces.StorageLocation = sl
	return &location
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
	NodeID   encoding.NodeId
	Location StorageLocationImpl
}

func NewSignedStorageLocation(NodeID encoding.NodeId, Location StorageLocationImpl) *SignedStorageLocationImpl {
	return &SignedStorageLocationImpl{
		NodeID:   NodeID,
		Location: Location,
	}
}

func (ssl *SignedStorageLocationImpl) String() string {
	nodeString, _ := ssl.NodeID.ToString()

	if nodeString == "" {
		nodeString = "failed to decode node id"
	}

	return "SignedStorageLocationImpl(" + ssl.Location.String() + ", " + nodeString + ")"
}

type storageLocationMap map[int]nodeStorage
type nodeStorage map[string]nodeDetailsStorage
type nodeDetailsStorage map[int]interface{}

func (s *storageLocationMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	temp, err := dec.DecodeUntypedMap()
	if err != nil {
		return err
	}

	if *s == nil {
		*s = make(map[int]nodeStorage)
	}

	tempMap, ok := interface{}(temp).(storageLocationMap)
	if !ok {
		return fmt.Errorf("unexpected data format from msgpack decoding")
	}

	*s = tempMap

	return nil
}

func (s storageLocationMap) EncodeMsgpack(enc *msgpack.Encoder) error {
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

func newStorageLocationMap() storageLocationMap {
	return storageLocationMap{}
}
