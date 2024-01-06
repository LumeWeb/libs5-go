package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"sync"
)

var (
	messageTypes sync.Map
)

var (
	_ IncomingMessage = (*IncomingMessageImpl)(nil)
)

type IncomingMessage interface {
	protocol.IncomingMessage
}

type IncomingMessageImpl struct {
	protocol.IncomingMessageImpl
	message []byte
}

func init() {
	messageTypes = sync.Map{}

	RegisterMessageType(types.ProtocolMethodHandshakeDone, func() IncomingMessage {
		return NewHandshakeDone()
	})
	RegisterMessageType(types.ProtocolMethodAnnouncePeers, func() IncomingMessage {
		return NewAnnouncePeers()
	})
}

func RegisterMessageType(messageType types.ProtocolMethod, factoryFunc func() IncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes.Store(messageType, factoryFunc)
}

func GetMessageType(kind types.ProtocolMethod) (protocol.IncomingMessage, bool) {
	value, ok := messageTypes.Load(kind)
	if !ok {
		return nil, false
	}

	factoryFunc, ok := value.(func() IncomingMessage)
	if !ok {
		return nil, false
	}

	return factoryFunc(), true
}
