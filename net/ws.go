package net

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"github.com/gorilla/websocket"
)

type WebSocketPeer struct {
	BasePeer
	Socket *websocket.Conn
}

func (p *WebSocketPeer) SendMessage(message []byte) {
	err := p.Socket.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return
	}
}

func (p *WebSocketPeer) RenderLocationURI() string {
	return p.Socket.RemoteAddr().String()
}

func (p *WebSocketPeer) ListenForMessages(callback EventCallback, onClose func(), onError func(error)) {
	for {
		_, message, err := p.Socket.ReadMessage()
		if err != nil {
			if onError != nil {
				onError(err)
			}
			break
		}

		err = callback(message)
		if err != nil {
			if onError != nil {
				onError(err)
			}
		}
	}

	if onClose != nil {
		onClose()
	}
}

func (p *WebSocketPeer) End() {
	err := p.Socket.Close()
	if err != nil {
		return
	}
}

func (p *WebSocketPeer) SetId(id *encoding.NodeId) {
	p.Id = id
}

func (p *WebSocketPeer) SetChallenge(challenge []byte) {
	p.challenge = challenge
}

func (p *WebSocketPeer) GetChallenge() []byte {
	return p.challenge
}
