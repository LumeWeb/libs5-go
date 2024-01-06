package signed

import (
	"bytes"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	libs5_go "git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"github.com/vmihailenco/msgpack/v5"
)

var _ protocol.IncomingMessageTyped = (*HandshakeDone)(nil)

type HandshakeDone struct {
	protocol.HandshakeOpen
	supportedFeatures int
}

func NewHandshakeDone() *HandshakeDone {
	return &HandshakeDone{HandshakeOpen: *protocol.NewHandshakeOpen(nil, ""), supportedFeatures: -1}
}

func (h HandshakeDone) HandleMessage(node *libs5_go.Node, peer *net.Peer, verifyId bool) error {
	if !(*node).IsStarted() {
		err := (*peer).End()
		if err != nil {
			return nil
		}
		return nil
	}

	if !bytes.Equal((*peer).GetChallenge(), h.HandshakeOpen.Challenge()) {
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
	supportedFeatures, err := dec.DecodeInt()

	if err != nil {
		return err
	}

	h.supportedFeatures = supportedFeatures
	return nil
}
