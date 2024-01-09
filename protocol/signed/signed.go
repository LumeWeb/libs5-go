package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
)

var (
	messageTypes map[int]func() base.SignedIncomingMessage
)

func Init() {
	messageTypes = make(map[int]func() base.SignedIncomingMessage)

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
	messageTypes[messageType] = factoryFunc
}

func GetMessageType(kind int) (base.SignedIncomingMessage, bool) {
	value, ok := messageTypes[kind]
	if !ok {
		return nil, false
	}

	return value(), true
}

var (
	_ base.SignedIncomingMessage = (*IncomingMessageImpl)(nil)
)

type IncomingMessageImpl struct {
	base.IncomingMessageImpl
	message []byte
}
