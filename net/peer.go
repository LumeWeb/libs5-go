package net

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"go.uber.org/zap"
	"net/url"
)

// EventCallback type for the callback function
type EventCallback func(event []byte) error

// DoneCallback type for the onDone callback
type DoneCallback func()

// ErrorCallback type for the onError callback
type ErrorCallback func(args ...interface{})

// ListenerOptions struct for options
type ListenerOptions struct {
	OnDone  *DoneCallback
	OnError *ErrorCallback
	Logger  *zap.Logger
}

type Peer interface {
	SendMessage(message []byte) error
	RenderLocationURI() string
	ListenForMessages(callback EventCallback, options ListenerOptions)
	End() error
	SetId(id *encoding.NodeId)
	GetId() *encoding.NodeId
	SetChallenge(challenge []byte)
	GetChallenge() []byte
}

type BasePeer struct {
	ConnectionURIs []url.URL
	IsConnected    bool
	challenge      []byte
	Socket         interface{}
	Id             *encoding.NodeId
}
