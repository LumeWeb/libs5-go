package encoding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/multiformats/go-multibase"
	"unicode/utf8"
)

var (
	errorNotBase64Url = errors.New("not a base64url string")
)

type MultihashCode = int

type Multihash struct {
	fullBytes []byte
}

func (m *Multihash) FullBytes() []byte {
	return m.fullBytes
}

var _ json.Marshaler = (*Multihash)(nil)
var _ json.Unmarshaler = (*Multihash)(nil)

func NewMultihash(fullBytes []byte) *Multihash {
	return &Multihash{fullBytes: fullBytes}
}

func (m *Multihash) FunctionType() types.HashType {
	return types.HashType(m.fullBytes[0])
}

func (m *Multihash) HashBytes() []byte {
	return m.fullBytes[1:]
}

func MultihashFromBytes(bytes []byte, kind types.HashType) *Multihash {
	return NewMultihash(append([]byte{byte(kind)}, bytes...))
}

func MultihashFromBase64Url(hash string) (*Multihash, error) {
	encoder, _ := multibase.EncoderByName("base64url")
	encoding, err := getEncoding(hash)

	if encoding != encoder.Encoding() {
		return nil, errorNotBase64Url
	}

	_, ret, err := multibase.Decode(hash)
	if err != nil {
		return nil, err
	}
	return NewMultihash(ret), nil
}

func (m *Multihash) ToBase64Url() (string, error) {
	return bases.ToBase64Url(m.fullBytes)
}

func (m *Multihash) ToBase32() (string, error) {
	return bases.ToBase32(m.fullBytes)
}

func (m *Multihash) ToString() (string, error) {
	if m.FunctionType() == types.HashType(types.CIDTypeBridge) {
		return string(m.fullBytes), nil // Assumes the bytes are valid UTF-8
	}
	return m.ToBase64Url()
}

func (m *Multihash) Equals(other *Multihash) bool {
	return bytes.Equal(m.fullBytes, other.fullBytes)
}

func (m *Multihash) HashCode() MultihashCode {
	return utils.HashCode(m.fullBytes[:4])
}

func (b *Multihash) UnmarshalJSON(data []byte) error {
	decodedData, err := MultibaseDecodeString(string(data))
	if err != nil {
		return err
	}

	b.fullBytes = decodedData
	return nil
}
func (b Multihash) MarshalJSON() ([]byte, error) {
	url, err := b.ToBase64Url()
	if err != nil {
		return nil, err
	}

	return []byte(url), nil

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
