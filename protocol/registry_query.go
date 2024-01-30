package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ IncomingMessage = (*RegistryQuery)(nil)
var _ EncodeableMessage = (*RegistryQuery)(nil)

type RegistryQuery struct {
	pk []byte
	HandshakeRequirement
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

func (s *RegistryQuery) DecodeMessage(dec *msgpack.Decoder, message IncomingMessageData) error {
	pk, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.pk = pk

	return nil
}

func (s *RegistryQuery) HandleMessage(message IncomingMessageData) error {
	mediator := message.Mediator
	peer := message.Peer
	sre, err := mediator.RegistryGet(s.pk)
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
