package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ base.IncomingMessage = (*RegistryQuery)(nil)
var _ base.EncodeableMessage = (*RegistryQuery)(nil)

type RegistryQuery struct {
	pk []byte
	base.HandshakeRequirement
}

func NewEmptyRegistryQuery() *RegistryQuery {
	rq := &RegistryQuery{}

	rq.SetRequiresHandshake(true)

	return rq
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

func (s *RegistryQuery) DecodeMessage(dec *msgpack.Decoder, message base.IncomingMessageData) error {
	pk, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.pk = pk

	return nil
}

func (s *RegistryQuery) HandleMessage(message base.IncomingMessageData) error {
	node := message.Node
	peer := message.Peer
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
