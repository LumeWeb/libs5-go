package base

import (
	"context"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/node"
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

var (
	_ msgpack.CustomDecoder = (*IncomingMessageReader)(nil)
)

type IncomingMessage interface {
	HandleMessage(message IncomingMessageData) error
	DecodeMessage(dec *msgpack.Decoder, message IncomingMessageData) error
	HandshakeRequirer
}

type IncomingMessageData struct {
	Original []byte
	Data     []byte
	Ctx      context.Context
	Node     *node.Node
	Peer     net.Peer
	VerifyId bool
}

type IncomingMessageReader struct {
	Kind int
	Data []byte
}

func (i *IncomingMessageReader) DecodeMsgpack(dec *msgpack.Decoder) error {
	kind, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	i.Kind = kind

	raw, err := io.ReadAll(dec.Buffered())

	if err != nil {
		return err
	}

	i.Data = raw

	return nil
}

type HandshakeRequirer interface {
	RequiresHandshake() bool
	SetRequiresHandshake(value bool)
}

type HandshakeRequirement struct {
	requiresHandshake bool
}

func (hr *HandshakeRequirement) RequiresHandshake() bool {
	return hr.requiresHandshake
}

func (hr *HandshakeRequirement) SetRequiresHandshake(value bool) {
	hr.requiresHandshake = value
}
