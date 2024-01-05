package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/serialize"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/vmihailenco/msgpack/v5"
)

type directoryReferenceMap struct {
	linkedhashmap.Map
}
type fileReferenceMap struct {
	linkedhashmap.Map
}

type DirectoryMetadata struct {
	Details       DirectoryMetadataDetails `json:"details"`
	Directories   directoryReferenceMap    `json:"directories"`
	Files         fileReferenceMap         `json:"files"`
	ExtraMetadata ExtraMetadata            `json:"extraMetadata"`
	BaseMetadata
}

var _ SerializableMetadata = (*DirectoryMetadata)(nil)
var _ SerializableMetadata = (*directoryReferenceMap)(nil)

func NewDirectoryMetadata(details DirectoryMetadataDetails, directories directoryReferenceMap, files fileReferenceMap, extraMetadata ExtraMetadata) *DirectoryMetadata {
	dirMetadata := &DirectoryMetadata{
		Details:       details,
		Directories:   directories,
		Files:         files,
		ExtraMetadata: extraMetadata,
	}

	dirMetadata.Type = "directory"
	return dirMetadata
}
func (dm *DirectoryMetadata) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := serialize.InitMarshaller(types.MetadataTypeDirectory, enc)
	if err != nil {
		return err
	}

	items := make([]interface{}, 4)

	items[0] = dm.Details
	items[1] = dm.Directories
	items[2] = dm.Files
	items[3] = dm.ExtraMetadata.Data

	return enc.Encode(items)
}

func (dm *DirectoryMetadata) DecodeMsgpack(dec *msgpack.Decoder) error {
	err := serialize.InitUnmarshaller(types.MetadataTypeDirectory, dec)
	if err != nil {
		return err
	}
	val, err := dec.DecodeArrayLen()

	if err != nil {
		return err
	}

	if val != 4 {
		return errors.New(" Corrupted metadata")
	}

	for i := 0; i < val; i++ {
		switch i {
		case 0:
			err = dec.Decode(&dm.Details)
			if err != nil {
				return err
			}
		case 1:
			err = dec.Decode(&dm.Directories)
			if err != nil {
				return err
			}

		case 2:
			err = dec.Decode(&dm.Files)
			if err != nil {
				return err
			}
		case 3:
			intMap, err := decodeIntMap(dec)
			if err != nil {
				return err
			}
			dm.ExtraMetadata.Data = intMap
		}
	}

	dm.Type = "directory"

	return nil
}
func (drm directoryReferenceMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	return marshallMapMsgpack(enc, &drm.Map)
}

func (drm *directoryReferenceMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	return unmarshalMapMsgpack(dec, &drm.Map, &DirectoryReference{})
}

func (frm fileReferenceMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	return marshallMapMsgpack(enc, &frm.Map)
}

func (frm *fileReferenceMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	return unmarshalMapMsgpack(dec, &frm.Map, &FileReference{})
}

func (frm *fileReferenceMap) UnmarshalJSON(bytes []byte) error {
	createFileInstance := func() interface{} { return &FileReference{} }
	return unmarshalMapJson(bytes, &frm.Map, createFileInstance)
}

type unmarshalNewInstanceFunc func() interface{}

func (drm *directoryReferenceMap) UnmarshalJSON(bytes []byte) error {
	createDirInstance := func() interface{} { return &DirectoryReference{} }
	return unmarshalMapJson(bytes, &drm.Map, createDirInstance)
}

func unmarshalMapMsgpack(dec *msgpack.Decoder, m *linkedhashmap.Map, placeholder interface{}) error {
	*m = *linkedhashmap.New()

	l, err := dec.DecodeMapLen()
	if err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		key, err := dec.DecodeString()
		if err != nil {
			return err
		}

		fmt.Println("dir: ", key)

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

		err = json.Unmarshal(data, &instance)
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
