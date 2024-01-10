package interfaces

import "git.lumeweb.com/LumeWeb/libs5-go/net"

type SignedRegistryEntry interface {
	PK() []byte
	Revision() uint64
	Data() []byte
	Signature() []byte
	SetPK(pk []byte)
	SetRevision(revision uint64)
	SetData(data []byte)
	SetSignature(signature []byte)
	Verify() bool
}

type RegistryEntry interface {
	Sign() SignedRegistryEntry
}

type RegistryService interface {
	Set(sre SignedRegistryEntry, trusted bool, receivedFrom net.Peer) error
	Get(pk []byte) (SignedRegistryEntry, error)
	BroadcastEntry(sre SignedRegistryEntry, receivedFrom net.Peer) error
	SendRegistryRequest(pk []byte) error
	Service
}
