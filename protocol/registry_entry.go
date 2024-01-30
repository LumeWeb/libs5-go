package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ IncomingMessage = (*RegistryEntryRequest)(nil)
var _ EncodeableMessage = (*RegistryEntryRequest)(nil)

type RegistryEntryRequest struct {
	sre SignedRegistryEntry
	HandshakeRequirement
}

func NewEmptyRegistryEntryRequest() *RegistryEntryRequest {
	rer := &RegistryEntryRequest{}

	rer.SetRequiresHandshake(true)

	return rer
}
func NewRegistryEntryRequest(sre SignedRegistryEntry) *RegistryEntryRequest {
	return &RegistryEntryRequest{sre: sre}
}

func (s *RegistryEntryRequest) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeInt(int64(types.RecordTypeRegistryEntry))
	if err != nil {
		return err
	}

	err = enc.EncodeBytes(MarshalSignedRegistryEntry(s.sre))
	if err != nil {
		return err
	}

	return nil
}

func (s *RegistryEntryRequest) DecodeMessage(dec *msgpack.Decoder, message IncomingMessageData) error {
	sre, err := UnmarshalSignedRegistryEntry(message.Data)
	if err != nil {
		return err
	}

	s.sre = sre

	return nil
}

func (s *RegistryEntryRequest) HandleMessage(message IncomingMessageData) error {
	return message.Mediator.RegistrySet(s.sre, false, message.Peer)
}
