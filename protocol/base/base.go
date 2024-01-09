package base

import (
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"github.com/vmihailenco/msgpack/v5"
)

//go:generate mockgen -source=base.go -destination=../mocks/base/base.go -package=base

type IncomingMessage interface {
	HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error
	SetIncomingMessage(msg IncomingMessage)
	IncomingMessage() IncomingMessage
	Self() IncomingMessage
	SetSelf(self IncomingMessage)
	Kind() int
	msgpack.CustomDecoder
}

type IncomingMessageTyped interface {
	DecodeMessage(dec *msgpack.Decoder) error
	IncomingMessage
}
