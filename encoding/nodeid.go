package encoding

import (
	"bytes"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/multiformats/go-multibase"
)

var (
	errorNotBase58BTC = errors.New("not a base58btc string")
)

type NodeIdCode = int

type NodeId struct {
	Bytes []byte
}

func NewNodeId(bytes []byte) *NodeId {
	return &NodeId{Bytes: bytes}
}

func NodeIdDecode(nodeId string) (*NodeId, error) {
	encoding, ret, err := multibase.Decode(nodeId)
	if err != nil {
		return nil, err
	}

	if encoding != multibase.Base58BTC {
		return nil, errorNotBase58BTC
	}

	return NewNodeId(ret), nil
}

func (nodeId *NodeId) Equals(other interface{}) bool {
	if otherNodeId, ok := other.(*NodeId); ok {
		return bytes.Equal(nodeId.Bytes, otherNodeId.Bytes)
	}
	return false
}

func (nodeId *NodeId) HashCode() int {
	return utils.HashCode(nodeId.Bytes[:4])
}

func (nodeId *NodeId) ToBase58() (string, error) {
	return bases.ToBase58BTC(nodeId.Bytes)
}

func (nodeId *NodeId) ToString() (string, error) {
	return nodeId.ToBase58()
}
