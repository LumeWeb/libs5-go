package base

import (
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"net/url"
)

var _ msgpack.CustomDecoder = (*IncomingMessageImpl)(nil)
var _ IncomingMessage = (*IncomingMessageImpl)(nil)
var _ IncomingMessageTyped = (*IncomingMessageImpl)(nil)

type IncomingMessageHandler func(node interfaces.Node, peer *net.Peer, u *url.URL, verifyId bool) error

type IncomingMessageImpl struct {
	kind     types.ProtocolMethod
	data     msgpack.RawMessage
	original []byte
	known    bool
	self     IncomingMessage
}

func (i *IncomingMessageImpl) Self() IncomingMessage {
	return i.self
}

func (i *IncomingMessageImpl) SetSelf(self IncomingMessage) {
	i.self = self
}

func (i *IncomingMessageImpl) DecodeMessage(dec *msgpack.Decoder) error {
	panic("child class should implement this method")
}

func (i *IncomingMessageImpl) Known() bool {
	return i.known
}

func (i *IncomingMessageImpl) SetKnown(known bool) {
	i.known = known
}

func (i *IncomingMessageImpl) SetOriginal(original []byte) {
	i.original = original
}

func (i *IncomingMessageImpl) Original() []byte {
	return i.original
}

func (i *IncomingMessageImpl) SetIncomingMessage(msg IncomingMessage) {
	if msgImpl, ok := msg.(*IncomingMessageImpl); ok {
		*i = *msgImpl
		i.known = true
	} else {
		// Handle the error or panic
		panic("msg is not of type *IncomingMessageImpl")
	}
}

func (i *IncomingMessageImpl) GetKind() types.ProtocolMethod {
	return i.kind
}

func (i *IncomingMessageImpl) ToMessage() (message []byte, err error) {
	return msgpack.Marshal(i)
}

func (i *IncomingMessageImpl) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	panic("child class should implement this method")
}

func (i *IncomingMessageImpl) Kind() types.ProtocolMethod {
	return i.kind
}

func (i *IncomingMessageImpl) Data() msgpack.RawMessage {
	return i.data
}

type IncomingMessageTypedImpl struct {
	IncomingMessageImpl
}

func NewIncomingMessageUnknown() *IncomingMessageImpl {
	return &IncomingMessageImpl{
		known: false,
	}
}

func NewIncomingMessageKnown(kind types.ProtocolMethod, data msgpack.RawMessage) *IncomingMessageImpl {
	return &IncomingMessageImpl{
		kind:  kind,
		data:  data,
		known: true,
	}
}

func NewIncomingMessageTyped(kind types.ProtocolMethod, data msgpack.RawMessage) *IncomingMessageTypedImpl {
	known := NewIncomingMessageKnown(kind, data)
	return &IncomingMessageTypedImpl{*known}
}

func (i *IncomingMessageImpl) DecodeMsgpack(dec *msgpack.Decoder) error {
	if i.known {
		if msgTyped, ok := interface{}(i.Self()).(IncomingMessageTyped); ok {
			return msgTyped.DecodeMessage(dec)
		}
		return fmt.Errorf("type assertion to IncomingMessageTyped failed")
	}

	kind, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	i.kind = types.ProtocolMethod(kind)

	raw, err := io.ReadAll(dec.Buffered())

	if err != nil {
		return err
	}

	i.data = raw
	return nil
}
