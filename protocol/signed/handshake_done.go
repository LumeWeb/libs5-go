package signed

import (
	"bytes"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	"net/url"
)

var _ base.IncomingMessageTyped = (*HandshakeDone)(nil)
var _ base.EncodeableMessage = (*HandshakeDone)(nil)

type HandshakeDone struct {
	challenge []byte
	networkId string
	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
	supportedFeatures int
	connectionUris    []*url.URL
	handshake         []byte
}

func NewHandshakeDoneRequest(handshake []byte, supportedFeatures int, connectionUris []*url.URL) *HandshakeDone {
	return &HandshakeDone{
		handshake:         handshake,
		supportedFeatures: supportedFeatures,
		connectionUris:    connectionUris,
	}
}

func (m HandshakeDone) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeUint(uint64(types.ProtocolMethodHandshakeDone))
	if err != nil {
		return err
	}

	err = enc.EncodeBytes(m.handshake)
	if err != nil {
		return err
	}

	err = enc.EncodeInt(int64(m.supportedFeatures))
	if err != nil {
		return err
	}

	err = utils.EncodeMsgpackArray(enc, m.connectionUris)
	if err != nil {
		return err
	}

	return nil
}

func (m *HandshakeDone) SetChallenge(challenge []byte) {
	m.challenge = challenge
}

func (m *HandshakeDone) SetNetworkId(networkId string) {
	m.networkId = networkId
}

func NewHandshakeDone() *HandshakeDone {
	return &HandshakeDone{challenge: nil, networkId: "", supportedFeatures: -1}
}

func (h HandshakeDone) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	if !node.IsStarted() {
		err := peer.End()
		if err != nil {
			return nil
		}
		return nil
	}

	if !bytes.Equal(peer.Challenge(), h.challenge) {
		return errors.New("Invalid challenge")
	}

	nodeId := h.IncomingMessage().(*SignedMessage).NodeId()

	if !verifyId {
		peer.SetId(nodeId)
	} else {
		if !peer.Id().Equals(nodeId) {
			return fmt.Errorf("Invalid transports id on initial list")
		}
	}

	peer.SetConnected(true)

	if h.supportedFeatures != types.SupportedFeatures {
		return fmt.Errorf("Remote node does not support required features")
	}
	err := node.Services().P2P().AddPeer(peer)
	if err != nil {
		return err
	}

	peer.SetConnectionURIs(h.connectionUris)

	peerId, err := peer.Id().ToString()

	if err != nil {
		return err
	}

	node.Logger().Info(fmt.Sprintf("[+] %s (%s)", peerId, peer.RenderLocationURI()))

	err = node.Services().P2P().SendPublicPeersToPeer(peer, []net.Peer{peer})
	if err != nil {
		return err
	}

	return nil
}

func (h *HandshakeDone) DecodeMessage(dec *msgpack.Decoder) error {
	challenge, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	h.challenge = challenge

	supportedFeatures, err := dec.DecodeInt()

	if err != nil {
		return err
	}

	h.supportedFeatures = supportedFeatures
	return nil
}
