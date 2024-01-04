package serialize

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

/*
	func NewSerializer(kind types.MetadataType) *Packer {
		p := NewPacker()
		_ = p.PackUint8(types.MetadataMagicByte)
		_ = p.PackUint8(uint8(types.MetadataTypeDirectory))
		return p
	}
*/
func InitMarshaller(kind types.MetadataType, enc *msgpack.Encoder) error {
	err := enc.EncodeInt(types.MetadataMagicByte)
	if err != nil {
		return err
	}
	err = enc.EncodeInt(int64(kind))
	if err != nil {
		return err
	}

	return nil
}

func InitUnmarshaller(kind types.MetadataType, enc *msgpack.Decoder) error {
	val, err := enc.DecodeUint8()
	if err != nil {
		return err
	}

	if val != types.MetadataMagicByte {
		return errors.New("Invalid magic byte")
	}

	val, err = enc.DecodeUint8()
	if err != nil {
		return err
	}

	if val != uint8(kind) {
		return errors.New("Invalid metadata type")
	}

	return nil
}
