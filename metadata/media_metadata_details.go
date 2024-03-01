package metadata

import (
	"errors"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ msgpack.CustomDecoder = (*MediaMetadataDetails)(nil)
	_ msgpack.CustomEncoder = (*MediaMetadataDetails)(nil)
)

type MediaMetadataDetails struct {
	Data map[int]interface{}
}

func NewMediaMetadataDetails(data map[int]interface{}) *MediaMetadataDetails {
	return &MediaMetadataDetails{Data: data}
}

func (mmd *MediaMetadataDetails) EncodeMsgpack(enc *msgpack.Encoder) error {
	return errors.New("Not implemented")
}

func (mmd *MediaMetadataDetails) DecodeMsgpack(dec *msgpack.Decoder) error {
	mapLen, err := dec.DecodeMapLen()

	if err != nil {
		return err
	}

	mmd.Data = make(map[int]interface{}, mapLen)

	for i := 0; i < mapLen; i++ {
		key, err := dec.DecodeInt8()
		if err != nil {
			return err
		}
		value, err := dec.DecodeInterface()
		if err != nil {
			return err
		}
		mmd.Data[int(key)] = value
	}

	return nil
}
