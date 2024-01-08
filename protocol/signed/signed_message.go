package signed

import (
	"crypto/ed25519"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

var (
	_ base.IncomingMessageTyped = (*SignedMessage)(nil)
	_ msgpack.CustomDecoder     = (*signedMessagePayoad)(nil)
	_ msgpack.CustomEncoder     = (*SignedMessage)(nil)
)

var (
	errInvalidSignature = errors.New("Invalid signature found")
)

type SignedMessage struct {
	nodeId    *encoding.NodeId
	signature []byte
	message   []byte
	base.IncomingMessageTypedImpl
}

func (s *SignedMessage) SetNodeId(nodeId *encoding.NodeId) {
	s.nodeId = nodeId
}

func (s *SignedMessage) SetSignature(signature []byte) {
	s.signature = signature
}

func (s *SignedMessage) SetMessage(message []byte) {
	s.message = message
}

func NewSignedMessageRequest(message []byte) *SignedMessage {
	return &SignedMessage{message: message}
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

	message, err := io.ReadAll(dec.Buffered())

	if err != nil {
		return err
	}

	s.message = message

	return nil
}

func NewSignedMessage() *SignedMessage {
	return &SignedMessage{}
}

func (s *SignedMessage) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	var payload signedMessagePayoad

	err := msgpack.Unmarshal(s.message, &payload)
	if err != nil {
		return err
	}

	if msgHandler, valid := GetMessageType(types.ProtocolMethod(payload.kind)); valid {
		msgHandler.SetIncomingMessage(s)
		msgHandler.SetSelf(msgHandler)
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

	if !ed25519.Verify(s.nodeId.Raw()[1:], s.message, s.signature) {
		return errInvalidSignature
	}

	return nil
}
func (s *SignedMessage) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeInt(int64(types.ProtocolMethodSignedMessage))

	if err != nil {
		return err
	}

	err = enc.EncodeBytes(s.nodeId.Raw())

	if err != nil {
		return err
	}

	err = enc.EncodeBytes(s.signature)

	if err != nil {
		return err
	}

	err = enc.EncodeBytes(s.message)

	if err != nil {
		return err
	}

	return nil
}
func (s *SignedMessage) Sign(node interfaces.Node) error {
	if s.nodeId == nil {
		panic("nodeId is nil")
	}

	if s.message == nil {
		panic("message is nil")
	}

	s.signature = ed25519.Sign(node.Config().KeyPair.ExtractBytes(), s.message)

	return nil
}
