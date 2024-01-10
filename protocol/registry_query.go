package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ base.IncomingMessageTyped = (*RegistryQuery)(nil)
var _ base.EncodeableMessage = (*RegistryQuery)(nil)

type RegistryQuery struct {
	pk []byte
	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
}

func NewEmptyRegistryQuery() *RegistryQuery {
	return &RegistryQuery{}
}
func NewRegistryQuery(pk []byte) *RegistryQuery {
	return &RegistryQuery{pk: pk}
}

func (s *RegistryQuery) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeInt(int64(types.ProtocolMethodRegistryQuery))
	if err != nil {
		return err
	}

	err = enc.EncodeBytes(s.pk)
	if err != nil {
		return err
	}

	return nil
}

func (s *RegistryQuery) DecodeMessage(dec *msgpack.Decoder) error {
	pk, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.pk = pk

	return nil
}

func (s *RegistryQuery) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	sre, err := node.Services().Registry().Get(s.pk)
	if err != nil {
		return err
	}

	if sre != nil {
		err := peer.SendMessage(MarshalSignedRegistryEntry(sre))
		if err != nil {
			return err
		}
	}

	return nil
}
