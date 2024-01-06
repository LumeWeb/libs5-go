package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"sync"
)

var (
	messageTypes sync.Map
)

var (
	_ IncomingMessage = (*IncomingMessageImpl)(nil)
)

func init() {
	messageTypes = sync.Map{}

	// Register factory functions instead of instances
	RegisterMessageType(types.ProtocolMethodHandshakeOpen, func() IncomingMessage {
		return NewHandshakeOpen([]byte{}, "")
	})
	RegisterMessageType(types.ProtocolMethodSignedMessage, func() IncomingMessage {
		return NewSignedMessage()
	})
}

func RegisterMessageType(messageType types.ProtocolMethod, factoryFunc func() IncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes.Store(messageType, factoryFunc)
}

func GetMessageType(kind types.ProtocolMethod) (IncomingMessage, bool) {
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
