package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
)

var (
	messageTypes map[int]func() base.IncomingMessage
)

var (
	_ base.IncomingMessage = (*base.IncomingMessageImpl)(nil)
)

func Init() {
	messageTypes = make(map[int]func() base.IncomingMessage)

	// Register factory functions instead of instances
	RegisterMessageType(int(types.ProtocolMethodHandshakeOpen), func() base.IncomingMessage {
		return NewHandshakeOpen([]byte{}, "")
	})
	RegisterMessageType(int(types.ProtocolMethodHashQuery), func() base.IncomingMessage {
		return NewHashQuery()
	})
	RegisterMessageType(int(types.RecordTypeStorageLocation), func() base.IncomingMessage {
		return NewStorageLocation()
	})
	RegisterMessageType(int(types.RecordTypeRegistryEntry), func() base.IncomingMessage {
		return NewEmptyRegistryEntryRequest()
	})
	RegisterMessageType(int(types.ProtocolMethodRegistryQuery), func() base.IncomingMessage {
		return NewEmptyRegistryQuery()
	})
	RegisterMessageType(int(types.ProtocolMethodSignedMessage), func() base.IncomingMessage {
		return signed.NewSignedMessage()
	})

}

func RegisterMessageType(messageType int, factoryFunc func() base.IncomingMessage) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes[messageType] = factoryFunc
}

func GetMessageType(kind int) (base.IncomingMessage, bool) {
	value, ok := messageTypes[kind]
	if !ok {
		return nil, false
	}

	return value(), true
}
