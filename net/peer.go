package net

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"go.uber.org/zap"
	"net/url"
)

// EventCallback type for the callback function
type EventCallback func(event []byte) error

// CloseCallback type for the OnClose callback
type CloseCallback func()

// ErrorCallback type for the onError callback
type ErrorCallback func(args ...interface{})

// ListenerOptions struct for options
type ListenerOptions struct {
	OnClose *CloseCallback
	OnError *ErrorCallback
	Logger  *zap.Logger
}

type Peer interface {
	SendMessage(message []byte) error
	RenderLocationURI() string
	ListenForMessages(callback EventCallback, options ListenerOptions)
	End() error
	SetId(id *encoding.NodeId)
	Id() *encoding.NodeId
	SetChallenge(challenge []byte)
	Challenge() []byte
	SetSocket(socket interface{})
	Socket() interface{}
}

type BasePeer struct {
	connectionURIs []*url.URL
	isConnected    bool
	challenge      []byte
	socket         interface{}
	id             *encoding.NodeId
}

func (b *BasePeer) Challenge() []byte {
	return b.challenge
}

func (b *BasePeer) SetChallenge(challenge []byte) {
	b.challenge = challenge
}

func (b *BasePeer) Socket() interface{} {
	return b.socket
}

func (b *BasePeer) SetSocket(socket interface{}) {
	b.socket = socket
}

func (b *BasePeer) Id() *encoding.NodeId {
	return b.id
}

func (b *BasePeer) SetId(id *encoding.NodeId) {
	b.id = id
}
