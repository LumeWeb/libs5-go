package interfaces

import "github.com/vmihailenco/msgpack/v5"

type NodeVotes interface {
	msgpack.CustomEncoder
	msgpack.CustomDecoder
	Good() int
	Bad() int
}
