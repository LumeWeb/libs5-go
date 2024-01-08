package utils

import (
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"net/url"
)

func EncodeMsgpackArray(enc *msgpack.Encoder, array interface{}) error {
	switch v := array.(type) {
	case []*url.URL:
		// Handle []*url.URL slice
		err := enc.EncodeInt(int64(len(v)))
		if err != nil {
			return err
		}
		for _, item := range v {
			err = enc.Encode(item)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		// Handle generic case
		arr, ok := array.([]interface{})
		if !ok {
			return errors.New("unsupported type for EncodeMsgpackArray")
		}
		err := enc.EncodeInt(int64(len(arr)))
		if err != nil {
			return err
		}
		for _, item := range arr {
			err = enc.Encode(item)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func DecodeMsgpackArray(dec *msgpack.Decoder) ([]interface{}, error) {
	arrayLen, err := dec.DecodeInt()
	if err != nil {
		return nil, err
	}

	array := make([]interface{}, arrayLen)

	for i := 0; i < int(arrayLen); i++ {
		item, err := dec.DecodeInterface()
		if err != nil {
			return nil, err
		}

		array[i] = item
	}

	return array, nil
}
