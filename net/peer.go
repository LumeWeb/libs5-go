package net

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"go.uber.org/zap"
	"net/url"
)

//go:generate mockgen -source=peer.go -destination=../mocks/net/peer.go -package=net

var (
	_ Peer = (*BasePeer)(nil)
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
	EndForAbuse() error
	SetId(id *encoding.NodeId)
	Id() *encoding.NodeId
	SetChallenge(challenge []byte)
	Challenge() []byte
	SetSocket(socket interface{})
	Socket() interface{}
	SetConnected(isConnected bool)
	IsConnected() bool
	SetConnectionURIs(uris []*url.URL)
	ConnectionURIs() []*url.URL
	IsHandshakeDone() bool
	SetHandshakeDone(status bool)
	GetIP() string
	Abused() bool
}

type BasePeer struct {
	connectionURIs []*url.URL
	isConnected    bool
	challenge      []byte
	socket         interface{}
	id             *encoding.NodeId
	handshaked     bool
}

func (b *BasePeer) IsConnected() bool {
	return b.isConnected
}

func (b *BasePeer) SetConnected(isConnected bool) {
	b.isConnected = isConnected
}

func (b *BasePeer) SendMessage(message []byte) error {
	panic("must implement in child class")
}

func (b *BasePeer) RenderLocationURI() string {
	panic("must implement in child class")
}

func (b *BasePeer) ListenForMessages(callback EventCallback, options ListenerOptions) {
	panic("must implement in child class")
}

func (b *BasePeer) End() error {
	panic("must implement in child class")
}
func (b *BasePeer) EndForAbuse() error {
	panic("must implement in child class")
}
func (b *BasePeer) GetIP() string {
	panic("must implement in child class")
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
func (b *BasePeer) SetConnectionURIs(uris []*url.URL) {
	b.connectionURIs = uris
}
func (b *BasePeer) ConnectionURIs() []*url.URL {
	return b.connectionURIs
}

func (b *BasePeer) IsHandshakeDone() bool {
	return b.handshaked
}

func (b *BasePeer) SetHandshakeDone(status bool) {
	b.handshaked = status
}

func (b *BasePeer) Abused() bool {
	panic("must implement in child class")
}
