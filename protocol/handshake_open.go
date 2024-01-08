package protocol

import (
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"net/url"
)

var _ base.IncomingMessageTyped = (*HandshakeOpen)(nil)

type HandshakeOpen struct {
	challenge []byte
	networkId string
	handshake []byte
	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
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
var (
	errInvalidChallenge = errors.New("Invalid challenge")
)

func NewHandshakeOpen(challenge []byte, networkId string) *HandshakeOpen {
	return &HandshakeOpen{
		challenge: challenge,
		networkId: networkId,
	}
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

func (h *HandshakeOpen) DecodeMessage(dec *msgpack.Decoder) error {
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

func (h *HandshakeOpen) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	if h.networkId != node.NetworkId() {
		return fmt.Errorf("Peer is in different network: %s", h.networkId)
	}

	handshake := signed.NewHandshakeDoneRequest(h.handshake, types.SupportedFeatures, []*url.URL{})
	message, err := msgpack.Marshal(handshake)

	if err != nil {
		return err
	}

	secureMessage, err := node.Services().P2P().SignMessageSimple(message)

	if err != nil {
		return err
	}

	err = peer.SendMessage(secureMessage)
	if err != nil {
		return err
	}

	return nil
}
