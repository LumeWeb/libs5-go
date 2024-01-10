package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ base.IncomingMessageTyped = (*RegistryEntryRequest)(nil)
var _ base.EncodeableMessage = (*RegistryEntryRequest)(nil)

type RegistryEntryRequest struct {
	sre interfaces.SignedRegistryEntry
	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
}

func NewEmptyRegistryEntryRequest() *RegistryEntryRequest {
	return &RegistryEntryRequest{}
}
func NewRegistryEntryRequest(sre interfaces.SignedRegistryEntry) *RegistryEntryRequest {
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

func (s *RegistryEntryRequest) DecodeMessage(dec *msgpack.Decoder) error {
	pk, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	sre, err := UnmarshalSignedRegistryEntry(pk)
	if err != nil {
		return err
	}

	s.sre = sre

	return nil
}

func (s *RegistryEntryRequest) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	return node.Services().Registry().Set(s.sre, false, peer)
}
