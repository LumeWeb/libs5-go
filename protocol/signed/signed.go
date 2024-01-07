package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"sync"
)

var (
	messageTypes sync.Map
)

func init() {
	messageTypes = sync.Map{}

	RegisterMessageType(types.ProtocolMethodHandshakeDone, func() base.SignedIncomingMessage {
		return NewHandshakeDone()
	})
	RegisterMessageType(types.ProtocolMethodAnnouncePeers, func() base.SignedIncomingMessage {
		return NewAnnouncePeers()
	})
}

func RegisterMessageType(messageType types.ProtocolMethod, factoryFunc func() base.SignedIncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes.Store(messageType, factoryFunc)
}

func GetMessageType(kind types.ProtocolMethod) (base.SignedIncomingMessage, bool) {
	value, ok := messageTypes.Load(kind)
	if !ok {
		return nil, false
	}

	factoryFunc, ok := value.(func() base.SignedIncomingMessage)
	if !ok {
		return nil, false
	}

	return factoryFunc(), true
}

var (
	_ base.SignedIncomingMessage = (*IncomingMessageImpl)(nil)
)

type IncomingMessageImpl struct {
	base.IncomingMessageImpl
	message []byte
}
