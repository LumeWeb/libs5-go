package libs5_go

import (
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"github.com/vmihailenco/msgpack/v5"
	"strconv"
	"time"
)

var (
	_ msgpack.CustomDecoder = (*storageLocationMap)(nil)
)

type StorageLocation struct {
	Type            int
	Parts           []string
	BinaryParts     [][]byte
	Expiry          int64
	ProviderMessage []byte
}

func NewStorageLocation(Type int, Parts []string, Expiry int64) *StorageLocation {
	return &StorageLocation{
		Type:   Type,
		Parts:  Parts,
		Expiry: Expiry,
	}
}

func (s *StorageLocation) BytesURL() string {
	return s.Parts[0]
}

func (s *StorageLocation) OutboardBytesURL() string {
	if len(s.Parts) == 1 {
		return s.Parts[0] + ".obao"
	}
	return s.Parts[1]
}

func (s *StorageLocation) String() string {
	expiryDate := time.Unix(s.Expiry, 0)
	return "StorageLocation(" + strconv.Itoa(s.Type) + ", " + fmt.Sprint(s.Parts) + ", expiry: " + expiryDate.Format(time.RFC3339) + ")"
}

type SignedStorageLocation struct {
	NodeID   encoding.NodeId
	Location StorageLocation
}

func NewSignedStorageLocation(NodeID encoding.NodeId, Location StorageLocation) *SignedStorageLocation {
	return &SignedStorageLocation{
		NodeID:   NodeID,
		Location: Location,
	}
}

func (ssl *SignedStorageLocation) String() string {
	nodeString, _ := ssl.NodeID.ToString()

	if nodeString == "" {
		nodeString = "failed to decode node id"
	}

	return "SignedStorageLocation(" + ssl.Location.String() + ", " + nodeString + ")"
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
