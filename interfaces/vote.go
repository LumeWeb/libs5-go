package interfaces

import "github.com/vmihailenco/msgpack/v5"

//go:generate mockgen -source=vote.go -destination=../mocks/interfaces/vote.go -package=interfaces

type NodeVotes interface {
	msgpack.CustomEncoder
	msgpack.CustomDecoder
	Good() int
	Bad() int
}
