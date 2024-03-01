package serialize

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
	"slices"
)

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

func InitUnmarshaller(enc *msgpack.Decoder, kinds ...types.MetadataType) (types.MetadataType, error) {
	val, err := enc.DecodeUint8()
	if err != nil {
		return 0, err
	}

	if val != types.MetadataMagicByte {
		return 0, errors.New("Invalid magic byte")
	}

	val, err = enc.DecodeUint8()
	if err != nil {
		return 0, err
	}

	convertedKinds := make([]uint8, len(kinds))
	for i, v := range kinds {
		convertedKinds[i] = uint8(v)
	}

	if !slices.Contains(convertedKinds, val) {
		return 0, errors.New("Invalid metadata type")
	}

	return types.MetadataType(val), nil
}
