package encoding

import (
	"encoding/json"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/internal/bases"
	"github.com/multiformats/go-multibase"
)

var (
	errMultibaseEncodingNotSupported = errors.New("multibase encoding not supported")
	errMultibaseDecodeZeroLength     = errors.New("cannot decode multibase for zero length string")
)

type Encoder interface {
	ToBytes() []byte
}

type multibaseImpl struct {
	Multibase
	encoder Encoder
}

var _ json.Marshaler = (*multibaseImpl)(nil)

type Multibase interface {
	ToHex() (string, error)
	ToBase32() (string, error)
	ToBase64Url() (string, error)
	ToBase58() (string, error)
	ToString() (string, error)
}

var _ Multibase = (*multibaseImpl)(nil)

func NewMultibase(encoder Encoder) Multibase {
	return &multibaseImpl{encoder: encoder}
}

func MultibaseDecodeString(data string) (bytes []byte, err error) {
	if len(data) == 0 {
		return nil, errMultibaseDecodeZeroLength
	}

	switch data[0] {
	case 'z', 'f', 'u', 'b':
		_, bytes, err = multibase.Decode(data)
	case ':':
		bytes = []byte(data)
	default:
		err = errMultibaseEncodingNotSupported
	}

	return bytes, err
}

func (m *multibaseImpl) ToHex() (string, error) {
	return bases.ToHex(m.encoder.ToBytes())
}

func (m *multibaseImpl) ToBase32() (string, error) {
	return bases.ToBase32(m.encoder.ToBytes())
}

func (m *multibaseImpl) ToBase64Url() (string, error) {
	return bases.ToBase64Url(m.encoder.ToBytes())
}

func (m *multibaseImpl) ToBase58() (string, error) {
	return bases.ToBase58BTC(m.encoder.ToBytes())
}

func (m *multibaseImpl) ToString() (string, error) {
	return m.ToBase58()
}
func (b multibaseImpl) MarshalJSON() ([]byte, error) {
	url, err := b.ToBase64Url()
	if err != nil {
		return nil, err
	}

	return []byte(url), nil

}
