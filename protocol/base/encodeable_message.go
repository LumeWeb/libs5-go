package base

import "github.com/vmihailenco/msgpack/v5"

//go:generate mockgen -source=encodeable_message.go -destination=../mocks/base/encodeable_message.go -package=base

type EncodeableMessage interface {
	msgpack.CustomEncoder
}
