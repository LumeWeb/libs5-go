package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"sync"
)

var (
	messageTypes sync.Map
)

func Init() {
	messageTypes = sync.Map{}

	RegisterMessageType(int(types.ProtocolMethodHandshakeDone), func() base.SignedIncomingMessage {
		return NewHandshakeDone()
	})
	RegisterMessageType(int(types.ProtocolMethodAnnouncePeers), func() base.SignedIncomingMessage {
		return NewAnnouncePeers()
	})
}

func RegisterMessageType(messageType int, factoryFunc func() base.SignedIncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes.Store(int(messageType), factoryFunc)
}

func GetMessageType(kind int) (base.SignedIncomingMessage, bool) {
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
