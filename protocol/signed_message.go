package protocol

import (
	"crypto/ed25519"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	libs5_go "git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ IncomingMessageTyped  = (*SignedMessage)(nil)
	_ msgpack.CustomDecoder = (*signedMessagePayoad)(nil)
)

var (
	errInvalidSignature = errors.New("Invalid signature found")
)

type SignedMessage struct {
	nodeId    *encoding.NodeId
	signature []byte
	message   []byte
	IncomingMessageTypedImpl
}

type signedMessagePayoad struct {
	kind    int
	message msgpack.RawMessage
}

func (s *signedMessagePayoad) DecodeMsgpack(dec *msgpack.Decoder) error {
	kind, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	s.kind = kind

	message, err := dec.DecodeRaw()
	if err != nil {
		return err
	}

	s.message = message

	return nil
}

func NewSignedMessage() *SignedMessage {
	return &SignedMessage{}
}

func (s *SignedMessage) HandleMessage(node *libs5_go.NodeImpl, peer *net.Peer, verifyId bool) error {
	var payload signedMessagePayoad

	err := msgpack.Unmarshal(s.message, &payload)
	if err != nil {
		return err
	}

	if msgHandler, valid := signed.GetMessageType(types.ProtocolMethod(payload.kind)); valid {
		msgHandler.SetIncomingMessage(s)
		err := msgpack.Unmarshal(payload.message, &msgHandler)
		if err != nil {
			return err
		}

		err = msgHandler.HandleMessage(node, peer, verifyId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SignedMessage) DecodeMessage(dec *msgpack.Decoder) error {
	nodeId, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.nodeId = encoding.NewNodeId(nodeId)

	signature, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.signature = signature

	message, err := dec.DecodeBytes()
	if err != nil {
		return err
	}

	s.message = message

	if !ed25519.Verify(s.nodeId.Raw(), s.message, s.signature) {
		return errInvalidSignature
	}

	return nil

}
