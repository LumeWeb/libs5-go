package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
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

	entry, err := encoding.CIDFromRegistryPublicKey(s.pk)
	if err != nil {
		return err
	}
	pid, err := peer.Id().ToString()
	if err != nil {
		return err
	}
	b64, err := bases.ToBase64Url(s.pk)
	if err != nil {
		return err
	}
	message.Logger.Debug("Handling registry entry query request", zap.Any("entryCID", entry), zap.Any("entryBase64", b64), zap.Any("peer", pid))

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
