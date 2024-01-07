package signed

import (
	"bytes"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"github.com/vmihailenco/msgpack/v5"
)

var _ base.IncomingMessageTyped = (*HandshakeDone)(nil)

type HandshakeDone struct {
	challenge []byte
	networkId string
	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
	supportedFeatures int
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

func (h HandshakeDone) HandleMessage(node interfaces.Node, peer *net.Peer, verifyId bool) error {
	if !node.IsStarted() {
		err := (*peer).End()
		if err != nil {
			return nil
		}
		return nil
	}

	if !bytes.Equal((*peer).GetChallenge(), h.challenge) {
		return errors.New("Invalid challenge")
	}
	/*
		if !verifyId {
			(*peer).SetId(h)
		} else {
			if !peer.ID.Equals(pId) {
				return errInvalidChallenge
			}
		}

		peer.IsConnected = true

		supportedFeatures := data.UnpackInt()

		if supportedFeatures != 3 {
			return errors.New("Remote node does not support required features")
		}

		node.Services.P2P.Peers[peer.ID.String()] = peer
		node.Services.P2P.ReconnectDelay[peer.ID.String()] = 1

		connectionUrisCount := data.UnpackInt()

		peer.ConnectionUris = make([]*url.URL, 0)
		for i := 0; i < connectionUrisCount; i++ {
			uriStr := data.UnpackString()
			uri, err := url.Parse(uriStr)
			if err != nil {
				return err
			}
			peer.ConnectionUris = append(peer.ConnectionUris, uri)
		}

		// Log information - Assuming a logging method exists
		node.Logger.Info(fmt.Sprintf("[+] %s (%s)", peer.ID.String(), peer.RenderLocationUri().String()))

		// Send peer lists and emit 'peerConnected' event
		// Assuming appropriate methods exist in node.Services.P2P
		node.Services.P2P.SendPublicPeersToPeer(peer)
	*/
	return nil
}

func (h HandshakeDone) DecodeMessage(dec *msgpack.Decoder) error {

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
