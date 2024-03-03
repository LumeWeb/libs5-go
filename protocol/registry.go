package protocol

import (
	ed25519p "crypto/ed25519"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
)

var (
	_ SignedRegistryEntry = (*SignedRegistryEntryImpl)(nil)
	_ SignedRegistryEntry = (*SignedRegistryEntryImpl)(nil)
)

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

type SignedRegistryEntryImpl struct {
	pk        []byte
	revision  uint64
	data      []byte
	signature []byte
}

func (s *SignedRegistryEntryImpl) Verify() bool {
	return VerifyRegistryEntry(s)
}

func (s *SignedRegistryEntryImpl) PK() []byte {
	return s.pk
}

func (s *SignedRegistryEntryImpl) SetPK(pk []byte) {
	s.pk = pk
}

func (s *SignedRegistryEntryImpl) Revision() uint64 {
	return s.revision
}

func (s *SignedRegistryEntryImpl) SetRevision(revision uint64) {
	s.revision = revision
}

func (s *SignedRegistryEntryImpl) Data() []byte {
	return s.data
}

func (s *SignedRegistryEntryImpl) SetData(data []byte) {
	s.data = data
}

func (s *SignedRegistryEntryImpl) Signature() []byte {
	return s.signature
}

func (s *SignedRegistryEntryImpl) SetSignature(signature []byte) {
	s.signature = signature
}

func NewSignedRegistryEntry(pk []byte, revision uint64, data []byte, signature []byte) SignedRegistryEntry {
	return &SignedRegistryEntryImpl{
		pk:        pk,
		revision:  revision,
		data:      data,
		signature: signature,
	}
}

type RegistryEntryImpl struct {
	kp       ed25519.KeyPairEd25519
	data     []byte
	revision uint64
}

func NewRegistryEntry(kp ed25519.KeyPairEd25519, data []byte, revision uint64) RegistryEntry {
	return &RegistryEntryImpl{
		kp:       kp,
		data:     data,
		revision: revision,
	}
}

func (r *RegistryEntryImpl) Sign() SignedRegistryEntry {
	return SignRegistryEntry(r.kp, r.data, r.revision)
}

func SignRegistryEntry(kp ed25519.KeyPairEd25519, data []byte, revision uint64) SignedRegistryEntry {
	buffer := MarshalRegistryEntry(kp.PublicKey(), data, revision)

	privateKey := kp.ExtractBytes()
	signature := ed25519p.Sign(privateKey, buffer)

	return NewSignedRegistryEntry(kp.PublicKey(), uint64(revision), data, signature)
}
func VerifyRegistryEntry(sre SignedRegistryEntry) bool {
	buffer := MarshalRegistryEntry(sre.PK(), sre.Data(), sre.Revision())
	publicKey := sre.PK()[1:]

	return ed25519p.Verify(publicKey, buffer, sre.Signature())
}

func MarshalSignedRegistryEntry(sre SignedRegistryEntry) []byte {
	buffer := MarshalRegistryEntry(sre.PK(), sre.Data(), sre.Revision())
	buffer = append(buffer, sre.Signature()...)

	return buffer
}
func MarshalRegistryEntry(pk []byte, data []byte, revision uint64) []byte {
	var buffer []byte
	buffer = append(buffer, byte(types.RecordTypeRegistryEntry))

	if pk != nil {
		buffer = append(buffer, pk...)
	}

	revBytes := utils.EncodeEndian(revision, 8)
	buffer = append(buffer, revBytes...)

	buffer = append(buffer, byte(len(data)))
	buffer = append(buffer, data...)

	return buffer
}

func UnmarshalSignedRegistryEntry(event []byte) (sre SignedRegistryEntry, err error) {
	if len(event) < 43 {
		return nil, errors.New("Invalid registry entry")
	}

	dataLength := int(event[42])
	if len(event) < 43+dataLength {
		return nil, errors.New("Invalid registry entry")
	}

	pk := event[1:34]
	revisionBytes := event[34:42]
	revision := utils.DecodeEndian(revisionBytes)
	signatureStart := 43 + dataLength
	var signature []byte

	if signatureStart < len(event) {
		signature = event[signatureStart:]
	} else {
		return nil, errors.New("Invalid signature")
	}

	return NewSignedRegistryEntry(pk, uint64(revision), event[43:signatureStart], signature), nil
}
