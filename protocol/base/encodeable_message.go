package base

import "github.com/vmihailenco/msgpack/v5"

type EncodeableMessage interface {
	msgpack.CustomEncoder
}
