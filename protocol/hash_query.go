package protocol

import (
	libs5_go "git.lumeweb.com/LumeWeb/libs5-go"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"github.com/vmihailenco/msgpack/v5"
)

var _ IncomingMessageTyped = (*HashQuery)(nil)

type HashQuery struct {
	hash  *encoding.Multihash
	kinds []int

	IncomingMessageTypedImpl
	IncomingMessageHandler
}

func (h HashQuery) Hash() *encoding.Multihash {
	return h.hash
}

func (h HashQuery) Kinds() []int {
	return h.kinds
}

func (h *HashQuery) DecodeMessage(dec *msgpack.Decoder) error {
	hash, err := dec.DecodeBytes()

	if err != nil {
		return err
	}

	h.hash = encoding.NewMultihash(hash)

	var kinds []int
	err = dec.Decode(&kinds)
	if err != nil {
		return err
	}

	h.kinds = kinds

	return nil
}
func (h *HashQuery) HandleMessage(node *libs5_go.Node, peer *net.Peer, verifyId bool) error {

	return nil
}
