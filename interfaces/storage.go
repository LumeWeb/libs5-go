package interfaces

import "git.lumeweb.com/LumeWeb/libs5-go/encoding"

//go:generate mockgen -source=storage.go -destination=../mocks/interfaces/storage.go -package=interfaces

type StorageLocationProvider interface {
	Start() error
	Next() (SignedStorageLocation, error)
	Upvote(uri SignedStorageLocation) error
	Downvote(uri SignedStorageLocation) error
}

type StorageLocation interface {
	BytesURL() string
	OutboardBytesURL() string
	String() string
	ProviderMessage() []byte
	Type() int
	Parts() []string
	BinaryParts() [][]byte
	Expiry() int64
	SetProviderMessage(msg []byte)
	SetType(t int)
	SetParts(p []string)
	SetBinaryParts(bp [][]byte)
	SetExpiry(e int64)
}
type SignedStorageLocation interface {
	String() string
	NodeId() *encoding.NodeId
	Location() StorageLocation
}
