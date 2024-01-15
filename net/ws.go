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
	socket *websocket.Conn
	abused bool
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
		socket: options.Socket.(*websocket.Conn),
	}

	return peer, nil
}

func (p *WebSocketPeer) SendMessage(message []byte) error {
	err := p.socket.Write(context.Background(), websocket.MessageBinary, message)
	if err != nil {
		return err
	}

	return nil
}

func (p *WebSocketPeer) RenderLocationURI() string {
	return "WebSocket client"
}

func (p *WebSocketPeer) ListenForMessages(callback EventCallback, options ListenerOptions) {
	errChan := make(chan error, 10)

	for {
		_, message, err := p.socket.Read(context.Background())
		if err != nil {
			if options.OnError != nil {
				(*options.OnError)(err)
			}
			break
		}

		// Process each message in a separate goroutine
		go func(msg []byte) {
			// Call the callback and send any errors to the error channel
			if err := callback(msg); err != nil {
				errChan <- err
			}
		}(message)

		// Non-blocking error check
		select {
		case err := <-errChan:
			if options.OnError != nil {
				(*options.OnError)(err)
			}
		default:
		}
	}

	if options.OnClose != nil {
		(*options.OnClose)()
	}

	// Handle remaining errors
	close(errChan)
	for err := range errChan {
		if options.OnError != nil {
			(*options.OnError)(err)
		}
	}
}

func (p *WebSocketPeer) End() error {
	err := p.socket.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return err
	}

	return nil
}
func (p *WebSocketPeer) EndForAbuse() error {
	err := p.socket.Close(websocket.StatusPolicyViolation, "")
	if err != nil {
		return err
	}

	p.abused = true

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

func (b *WebSocketPeer) GetIP() string {
	ctx, cancel := context.WithCancel(context.Background())
	netConn := websocket.NetConn(ctx, b.socket, websocket.MessageBinary)

	ipAddr := netConn.RemoteAddr().String()

	cancel()

	return ipAddr
}
func (p *WebSocketPeer) Abused() bool {
	return p.abused
}
