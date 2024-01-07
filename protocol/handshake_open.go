package protocol

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

type HandshakeOpen struct {
	challenge []byte
	networkId string
	IncomingMessageTypedImpl
	IncomingMessageHandler
}

func (m HandshakeOpen) Challenge() []byte {
	return m.challenge
}

func (m HandshakeOpen) NetworkId() string {
	return m.networkId
}

var _ EncodeableMessage = (*HandshakeOpen)(nil)
var (
	errInvalidChallenge = errors.New("Invalid challenge")
)

func NewHandshakeOpen(challenge []byte, networkId string) *HandshakeOpen {
	return &HandshakeOpen{
		challenge: challenge,
		networkId: networkId,
	}
}
func (m HandshakeOpen) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeUint(uint64(types.ProtocolMethodHandshakeOpen))
	if err != nil {
		return err
	}

	err = enc.EncodeBytes(m.challenge)
	if err != nil {
		return err
	}

	if m.networkId != "" {
		err = enc.EncodeString(m.networkId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *HandshakeOpen) HandleMessage(node interfaces.Node, peer *net.Peer, verifyId bool) error {

	return nil
}
