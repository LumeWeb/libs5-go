package interfaces

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
}
