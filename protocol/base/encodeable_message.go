package base

import "github.com/vmihailenco/msgpack/v5"

var (
	_ EncodeableMessage = (*EncodeableMessageImpl)(nil)
)

//go:generate mockgen -source=encodeable_message.go -destination=../mocks/base/encodeable_message.go -package=base

type EncodeableMessage interface {
	msgpack.CustomEncoder
}

type EncodeableMessageImpl struct {
}

func (e EncodeableMessageImpl) EncodeMsgpack(encoder *msgpack.Encoder) error {
	panic("this method should be implemented by the child class")
}
