package multihash

import (
	"bytes"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/multiformats/go-multibase"
	"unicode/utf8"
)

var (
	errorNotBase64Url = errors.New("not a base64url string")
)

type Multihash struct {
	FullBytes []byte
}

func New(fullBytes []byte) *Multihash {
	return &Multihash{FullBytes: fullBytes}
}

func (m *Multihash) FunctionType() types.HashType {
	return types.HashType(m.FullBytes[0])
}

func (m *Multihash) HashBytes() []byte {
	return m.FullBytes[1:]
}

func FromBase64Url(hash string) (*Multihash, error) {
	encoder, _ := multibase.EncoderByName("base64url")
	encoding, err := getEncoding(hash)

	if encoding != encoder.Encoding() {
		return nil, errorNotBase64Url
	}

	_, ret, err := multibase.Decode(hash)
	if err != nil {
		return nil, err
	}
	return New(ret), nil
}

func (m *Multihash) ToBase64Url() (string, error) {
	return bases.ToBase64Url(m.FullBytes)
}

func (m *Multihash) ToBase32() (string, error) {
	return bases.ToBase32(m.FullBytes)
}

func (m *Multihash) ToString() (string, error) {
	if m.FunctionType() == types.HashType(types.CIDTypeBridge) {
		return string(m.FullBytes), nil // Assumes the bytes are valid UTF-8
	}
	return m.ToBase64Url()
}

func (m *Multihash) Equals(other *Multihash) bool {
	return bytes.Equal(m.FullBytes, other.FullBytes)
}

func (m *Multihash) HashCode() int {
	return int(m.FullBytes[0]) +
		int(m.FullBytes[1])<<8 +
		int(m.FullBytes[2])<<16 +
		int(m.FullBytes[3])<<24
}

func getEncoding(hash string) (multibase.Encoding, error) {
	r, _ := utf8.DecodeRuneInString(hash)
	enc := multibase.Encoding(r)

	_, ok := multibase.EncodingToStr[enc]
	if !ok {
		return -1, fmt.Errorf("unsupported multibase encoding: %d", enc)

	}
	return enc, nil
}
