package encoding

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
)

type EncryptedCID struct {
	Multibase
	encryptedBlobHash   Multihash
	originalCID         CID
	encryptionAlgorithm byte
	padding             uint32
	chunkSizeAsPowerOf2 int
	encryptionKey       []byte
}

func NewEncryptedCID(encryptedBlobHash Multihash, originalCID CID, encryptionKey []byte, padding uint32, chunkSizeAsPowerOf2 int, encryptionAlgorithm byte) *EncryptedCID {
	e := &EncryptedCID{
		encryptedBlobHash:   encryptedBlobHash,
		originalCID:         originalCID,
		encryptionKey:       encryptionKey,
		padding:             padding,
		chunkSizeAsPowerOf2: chunkSizeAsPowerOf2,
		encryptionAlgorithm: encryptionAlgorithm,
	}

	m := NewMultibase(e)
	e.Multibase = m

	return e
}

func DecodeEncryptedCID(cid string) (*EncryptedCID, error) {
	data, err := MultibaseDecodeString(cid)
	if err != nil {
		return nil, err
	}
	return FromBytes(data)
}

func FromBytes(data []byte) (*EncryptedCID, error) {
	if types.CIDType(data[0]) != types.CIDTypeEncryptedStatic {
		return nil, errors.New("Invalid CID type")
	}

	cid, err := CIDFromBytes(data[72:])
	if err != nil {
		return nil, err
	}

	encryptedBlobHash := NewMultihash(data[3:36])
	encryptionKey := data[36:68]
	padding := utils.DecodeEndian(data[68:72])
	chunkSizeAsPowerOf2 := int(data[2])
	encryptionAlgorithm := data[1]

	return NewEncryptedCID(*encryptedBlobHash, *cid, encryptionKey, padding, chunkSizeAsPowerOf2, encryptionAlgorithm), nil
}

func (c *EncryptedCID) ChunkSize() int {
	return 1 << uint(c.chunkSizeAsPowerOf2)
}

func (c *EncryptedCID) ToBytes() []byte {
	data := []byte{
		byte(types.CIDTypeEncryptedStatic),
		c.encryptionAlgorithm,
		byte(c.chunkSizeAsPowerOf2),
	}
	data = append(data, c.encryptedBlobHash.FullBytes...)
	data = append(data, c.encryptionKey...)
	data = append(data, utils.EncodeEndian(c.padding, 4)...)
	data = append(data, c.originalCID.ToBytes()...)
	return data
}
