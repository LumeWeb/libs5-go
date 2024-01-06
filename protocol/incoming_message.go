package protocol

import (
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"net/url"
)

var (
	_ EncodeableMessage = (*EncodeableMessageImpl)(nil)
)

type IncomingMessage interface {
	HandleMessage(node *libs5_go.Node, peer *net.Peer, verifyId bool) error
	SetIncomingMessage(msg IncomingMessage)
	msgpack.CustomDecoder
}

type IncomingMessageTyped interface {
	DecodeMessage(dec *msgpack.Decoder) error
	IncomingMessage
}

type IncomingMessageImpl struct {
	kind  types.ProtocolMethod
	data  msgpack.RawMessage
	known bool
}

func (i *IncomingMessageImpl) SetIncomingMessage(msg IncomingMessage) {
	*i = interface{}(msg).(IncomingMessageImpl)
}

var _ msgpack.CustomDecoder = (*IncomingMessageImpl)(nil)
var _ IncomingMessage = (*IncomingMessageImpl)(nil)

func (i *IncomingMessageImpl) GetKind() types.ProtocolMethod {
	return i.kind
}

func (i *IncomingMessageImpl) ToMessage() (message []byte, err error) {
	return msgpack.Marshal(i)
}

func (i *IncomingMessageImpl) HandleMessage(node *libs5_go.Node, peer *net.Peer, verifyId bool) error {
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

type IncomingMessageHandler func(node *libs5_go.Node, peer *net.Peer, u *url.URL, verifyId bool) error

func (i *IncomingMessageImpl) DecodeMsgpack(dec *msgpack.Decoder) error {
	if i.known {
		if msgTyped, ok := interface{}(i).(IncomingMessageTyped); ok {
			return msgTyped.DecodeMessage(dec)
		}
		return fmt.Errorf("type assertion to IncomingMessageTyped failed")
	}

	kind, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	i.kind = types.ProtocolMethod(kind)

	raw, err := dec.DecodeRaw()
	if err != nil {
		return err
	}

	i.data = raw
	return nil
}
