package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

type IncomingMessageDataSigned struct {
	base.IncomingMessageData
	NodeId *encoding.NodeId
}

type IncomingMessageSigned interface {
	HandleMessage(message IncomingMessageDataSigned) error
	DecodeMessage(dec *msgpack.Decoder, message IncomingMessageDataSigned) error
	base.HandshakeRequirer
}

var (
	messageTypes map[int]func() IncomingMessageSigned
)

func Init() {
	messageTypes = make(map[int]func() IncomingMessageSigned)

	RegisterMessageType(int(types.ProtocolMethodHandshakeDone), func() IncomingMessageSigned {
		return NewHandshakeDone()
	})
	RegisterMessageType(int(types.ProtocolMethodAnnouncePeers), func() IncomingMessageSigned {
		return NewAnnouncePeers()
	})
}

func RegisterMessageType(messageType int, factoryFunc func() IncomingMessageSigned) {
	if factoryFunc == nil {
		panic("factoryFunc cannot be nil")
	}
	messageTypes[messageType] = factoryFunc
}

func GetMessageType(kind int) (IncomingMessageSigned, bool) {
	value, ok := messageTypes[kind]
	if !ok {
		return nil, false
	}

	return value(), true
}
