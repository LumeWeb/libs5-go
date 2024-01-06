package protocol

import "github.com/vmihailenco/msgpack/v5"

type EncodeableMessage interface {
	ToMessage() (message []byte, err error)
	msgpack.CustomEncoder
}

type EncodeableMessageImpl struct {
}

func (e EncodeableMessageImpl) ToMessage() (message []byte, err error) {
	return msgpack.Marshal(e)
}

func (e EncodeableMessageImpl) EncodeMsgpack(encoder *msgpack.Encoder) error {
	panic("this method should be implemented by the child class")
}
