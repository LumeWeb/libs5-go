package protocol

import (
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var _ base.EncodeableMessage = (*HandshakeOpen)(nil)
var _ base.IncomingMessage = (*HandshakeOpen)(nil)

type HandshakeOpen struct {
	challenge []byte
	networkId string
	handshake []byte
	base.HandshakeRequirement
}

func (h *HandshakeOpen) SetHandshake(handshake []byte) {
	h.handshake = handshake
}

func (h HandshakeOpen) Challenge() []byte {
	return h.challenge
}

func (h HandshakeOpen) NetworkId() string {
	return h.networkId
}

var _ base.EncodeableMessage = (*HandshakeOpen)(nil)

func NewHandshakeOpen(challenge []byte, networkId string) *HandshakeOpen {
	ho := &HandshakeOpen{
		challenge: challenge,
		networkId: networkId,
	}

	ho.SetRequiresHandshake(false)

	return ho
}
func (h HandshakeOpen) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeUint(uint64(types.ProtocolMethodHandshakeOpen))
	if err != nil {
		return err
	}

	err = enc.EncodeBytes(h.challenge)
	if err != nil {
		return err
	}

	if h.networkId != "" {
		err = enc.EncodeString(h.networkId)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *HandshakeOpen) DecodeMessage(dec *msgpack.Decoder, message base.IncomingMessageData) error {
	handshake, err := dec.DecodeBytes()

	if err != nil {
		return err
	}

	h.handshake = handshake

	_, err = dec.PeekCode()

	networkId := ""

	if err != nil {
		if err.Error() != "EOF" {
			return err
		}
		h.networkId = networkId
		return nil
	}

	networkId, err = dec.DecodeString()
	if err != nil {
		return err
	}

	h.networkId = networkId

	return nil
}

func (h *HandshakeOpen) HandleMessage(message base.IncomingMessageData) error {
	peer := message.Peer
	services := message.Services

	if h.networkId != services.P2P().NetworkId() {
		return fmt.Errorf("Peer is in different network: %s", h.networkId)
	}

	handshake := signed.NewHandshakeDoneRequest(h.handshake, types.SupportedFeatures, services.P2P().SelfConnectionUris())
	hsMessage, err := msgpack.Marshal(handshake)

	if err != nil {
		return err
	}

	secureMessage, err := services.P2P().SignMessageSimple(hsMessage)

	if err != nil {
		return err
	}

	err = peer.SendMessage(secureMessage)
	if err != nil {
		return err
	}

	return nil
}
