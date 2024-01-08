package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"sync"
)

var (
	messageTypes sync.Map
)

var (
	_ base.IncomingMessage = (*base.IncomingMessageImpl)(nil)
)

func init() {
	messageTypes = sync.Map{}

	// Register factory functions instead of instances
	RegisterMessageType(types.ProtocolMethodHandshakeOpen, func() base.IncomingMessage {
		return NewHandshakeOpen([]byte{}, "")
	})
	RegisterMessageType(types.ProtocolMethodHashQuery, func() base.IncomingMessage {
		return NewHashQuery()
	})
	RegisterMessageType(types.ProtocolMethodSignedMessage, func() base.IncomingMessage {
		return signed.NewSignedMessage()
	})

}

func RegisterMessageType(messageType types.ProtocolMethod, factoryFunc func() base.IncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes.Store(messageType, factoryFunc)
}

func GetMessageType(kind types.ProtocolMethod) (base.IncomingMessage, bool) {
	value, ok := messageTypes.Load(kind)
	if !ok {
		return nil, false
	}

	factoryFunc, ok := value.(func() base.IncomingMessage)
	if !ok {
		return nil, false
	}

	return factoryFunc(), true
}
