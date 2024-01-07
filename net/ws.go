package net

import (
	"context"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"net/url"
	"nhooyr.io/websocket"
)

var (
	_ PeerFactory = (*WebSocketPeer)(nil)
	_ PeerStatic  = (*WebSocketPeer)(nil)
	_ Peer        = (*WebSocketPeer)(nil)
)

type WebSocketPeer struct {
	BasePeer
	Socket *websocket.Conn
}

func (p *WebSocketPeer) Connect(uri *url.URL) (interface{}, error) {
	dial, _, err := websocket.Dial(context.Background(), uri.String(), nil)
	if err != nil {
		return nil, err
	}

	return dial, nil
}

func (p *WebSocketPeer) NewPeer(options *TransportPeerConfig) (Peer, error) {
	peer := &WebSocketPeer{
		BasePeer: BasePeer{
			connectionURIs: options.Uris,
			socket:         options.Socket,
		},
		Socket: options.Socket.(*websocket.Conn),
	}

	return peer, nil
}

func (p *WebSocketPeer) SendMessage(message []byte) error {
	err := p.Socket.Write(context.Background(), websocket.MessageBinary, message)
	if err != nil {
		return err
	}

	return nil
}

func (p *WebSocketPeer) RenderLocationURI() string {
	return "WebSocket client"
}

func (p *WebSocketPeer) ListenForMessages(callback EventCallback, options ListenerOptions) {
	for {
		_, message, err := p.Socket.Read(context.Background())
		if err != nil {
			if options.OnError != nil {
				(*options.OnError)(err)
			}
			break
		}

		err = callback(message)
		if err != nil {
			if options.OnError != nil {
				(*options.OnError)(err)
			}
		}
	}

	if options.OnClose != nil {
		(*options.OnClose)()
	}
}

func (p *WebSocketPeer) End() error {
	err := p.Socket.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return err
	}

	return nil
}

func (p *WebSocketPeer) SetId(id *encoding.NodeId) {
	p.id = id
}

func (p *WebSocketPeer) SetChallenge(challenge []byte) {
	p.challenge = challenge
}

func (p *WebSocketPeer) GetChallenge() []byte {
	return p.challenge
}
