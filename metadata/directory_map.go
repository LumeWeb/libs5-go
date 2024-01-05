package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/vmihailenco/msgpack/v5"
)

type directoryReferenceMap struct {
	linkedhashmap.Map
}
type fileReferenceMap struct {
	linkedhashmap.Map
}

type unmarshalNewInstanceFunc func() interface{}

var _ SerializableMetadata = (*directoryReferenceMap)(nil)

func unmarshalMapMsgpack(dec *msgpack.Decoder, m *linkedhashmap.Map, placeholder interface{}, intMap bool) error {
	*m = *linkedhashmap.New()

	l, err := dec.DecodeMapLen()
	if err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		var key interface{}
		if intMap {
			intKey, err := dec.DecodeInt()
			if err != nil {
				return err
			}
			key = intKey
		} else {
			strKey, err := dec.DecodeString()
			if err != nil {
				return err
			}
			key = strKey
		}

		switch placeholder.(type) {
		case *DirectoryReference:
			var value DirectoryReference
			if err := dec.Decode(&value); err != nil {
				return err
			}
			m.Put(key, value)

		case *FileReference:
			var file FileReference
			if err := dec.Decode(&file); err != nil {
				return err
			}
			m.Put(key, file)

		default:
			return fmt.Errorf("unsupported type for decoding")
		}
	}

	return nil
}

func marshallMapMsgpack(enc *msgpack.Encoder, m *linkedhashmap.Map) error {
	// First, encode the length of the map
	if err := enc.EncodeMapLen(m.Size()); err != nil {
		return err
	}

	iter := m.Iterator()
	for iter.Next() {
		key := iter.Key().(string)
		if err := enc.EncodeString(key); err != nil {
			return err
		}

		value := iter.Value()
		switch v := value.(type) {
		case FileReference:
			if err := enc.Encode(&v); err != nil {
				return err
			}
		case DirectoryReference:
			if err := enc.Encode(&v); err != nil {
				return err
			}
		case string:
			if err := enc.EncodeString(v); err != nil {
				return err
			}
		case int:
			if err := enc.EncodeInt(int64(v)); err != nil {
				return err
			}
		case FileVersion:
			if err := enc.Encode(&v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type for encoding")
		}
	}

	return nil
}

func unmarshalMapJson(bytes []byte, m *linkedhashmap.Map, newInstance unmarshalNewInstanceFunc) error {
	*m = *linkedhashmap.New()
	err := m.FromJSON(bytes)
	if err != nil {
		return err
	}

	iter := m.Iterator()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()

		instance := newInstance()

		data, err := json.Marshal(val)
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, instance)
		if err != nil {
			return err
		}

		// Type switch to handle different types
		switch v := instance.(type) {
		case *DirectoryReference:
			m.Put(key, *v)
		case *FileReference:
			m.Put(key, *v)
		default:
			return fmt.Errorf("unhandled type: %T", v)
		}
	}

	return nil
}
func (drm directoryReferenceMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	return marshallMapMsgpack(enc, &drm.Map)
}

func (drm *directoryReferenceMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	return unmarshalMapMsgpack(dec, &drm.Map, &DirectoryReference{}, false)
}

func (frm fileReferenceMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	return marshallMapMsgpack(enc, &frm.Map)
}

func (frm *fileReferenceMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	return unmarshalMapMsgpack(dec, &frm.Map, &FileReference{}, false)
}

func (frm *fileReferenceMap) UnmarshalJSON(bytes []byte) error {
	createFileInstance := func() interface{} { return &FileReference{} }
	return unmarshalMapJson(bytes, &frm.Map, createFileInstance)
}

func (drm *directoryReferenceMap) UnmarshalJSON(bytes []byte) error {
	createDirInstance := func() interface{} { return &DirectoryReference{} }
	return unmarshalMapJson(bytes, &drm.Map, createDirInstance)
}
