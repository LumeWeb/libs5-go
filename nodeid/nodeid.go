package nodeid

import (
	"bytes"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"github.com/multiformats/go-multibase"
)

var (
	errorNotBase58BTC = errors.New("not a base58btc string")
)

type NodeId struct {
	Bytes []byte
}

func New(bytes []byte) *NodeId {
	return &NodeId{Bytes: bytes}
}

func Decode(nodeId string) (*NodeId, error) {
	encoding, ret, err := multibase.Decode(nodeId)
	if err != nil {
		return nil, err
	}

	if encoding != multibase.Base58BTC {
		return nil, errorNotBase58BTC
	}

	return New(ret), nil
}

func (nodeId *NodeId) Equals(other interface{}) bool {
	if otherNodeId, ok := other.(*NodeId); ok {
		return bytes.Equal(nodeId.Bytes, otherNodeId.Bytes)
	}
	return false
}

func (nodeId *NodeId) HashCode() int {
	if len(nodeId.Bytes) < 4 {
		return 0
	}
	return int(nodeId.Bytes[0]) +
		int(nodeId.Bytes[1])<<8 +
		int(nodeId.Bytes[2])<<16 +
		int(nodeId.Bytes[3])<<24
}

func (nodeId *NodeId) ToBase58() (string, error) {
	return bases.ToBase58BTC(nodeId.Bytes)
}

func (nodeId *NodeId) ToString() (string, error) {
	return nodeId.ToBase58()
}
