package protocol

import (
	ed25519p "crypto/ed25519"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
)

var (
	_ interfaces.SignedRegistryEntry = (*SignedRegistryEntryImpl)(nil)
	_ interfaces.SignedRegistryEntry = (*SignedRegistryEntryImpl)(nil)
)

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

func NewSignedRegistryEntry(pk []byte, revision uint64, data []byte, signature []byte) interfaces.SignedRegistryEntry {
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

func NewRegistryEntry(kp ed25519.KeyPairEd25519, data []byte, revision uint64) interfaces.RegistryEntry {
	return &RegistryEntryImpl{
		kp:       kp,
		data:     data,
		revision: revision,
	}
}

func (r *RegistryEntryImpl) Sign() interfaces.SignedRegistryEntry {
	return SignRegistryEntry(r.kp, r.data, r.revision)
}

func SignRegistryEntry(kp ed25519.KeyPairEd25519, data []byte, revision uint64) interfaces.SignedRegistryEntry {
	buffer := MarshalRegistryEntry(data, revision)

	privateKey := kp.ExtractBytes()
	signature := ed25519p.Sign(privateKey, buffer)

	return NewSignedRegistryEntry(kp.PublicKey(), uint64(revision), data, signature)
}
func VerifyRegistryEntry(sre interfaces.SignedRegistryEntry) bool {
	buffer := MarshalRegistryEntry(sre.Data(), sre.Revision())
	publicKey := sre.PK()[1:]

	return ed25519p.Verify(publicKey, buffer, sre.Signature())
}

func MarshalSignedRegistryEntry(sre interfaces.SignedRegistryEntry) []byte {
	buffer := MarshalRegistryEntry(sre.Data(), sre.Revision())
	buffer = append(buffer, sre.Signature()...)

	return buffer
}
func MarshalRegistryEntry(data []byte, revision uint64) []byte {
	var buffer []byte
	buffer = append(buffer, byte(types.RecordTypeRegistryEntry))

	revBytes := utils.EncodeEndian(uint32(revision), 8)
	buffer = append(buffer, revBytes...)

	buffer = append(buffer, byte(len(data)))
	buffer = append(buffer, data...)

	return buffer
}

func UnmarshalSignedRegistryEntry(event []byte) (sre interfaces.SignedRegistryEntry, err error) {
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
