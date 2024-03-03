package metadata

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/serialize"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/vmihailenco/msgpack/v5"
	"sort"
)

var (
	_ Metadata             = (*WebAppMetadata)(nil)
	_ SerializableMetadata = (*WebAppMetadata)(nil)
	_ SerializableMetadata = (*WebAppFileMap)(nil)
)

type WebAppMetadata struct {
	BaseMetadata
	Name          string         `json:"name"`
	TryFiles      []string       `json:"tryFiles"`
	ErrorPages    map[int]string `json:"errorPages"`
	ExtraMetadata ExtraMetadata  `json:"extraMetadata"`
	Paths         WebAppFileMap  `json:"paths"`
}

func NewWebAppMetadata(name string, tryFiles []string, extraMetadata ExtraMetadata, errorPages map[int]string, paths WebAppFileMap) *WebAppMetadata {
	return &WebAppMetadata{
		Name:          name,
		TryFiles:      tryFiles,
		ExtraMetadata: extraMetadata,
		ErrorPages:    errorPages,
		Paths:         paths,
	}
}
func NewEmptyWebAppMetadata() *WebAppMetadata {
	return &WebAppMetadata{}
}

func (wm *WebAppMetadata) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := serialize.InitMarshaller(enc, types.MetadataTypeWebApp)
	if err != nil {
		return err
	}

	items := make([]interface{}, 5)

	items[0] = wm.Name
	items[1] = wm.TryFiles
	items[2] = wm.ErrorPages
	items[3] = wm.Paths
	items[4] = wm.ExtraMetadata

	return enc.Encode(items)
}

func (wm *WebAppMetadata) DecodeMsgpack(dec *msgpack.Decoder) error {
	_, err := serialize.InitUnmarshaller(dec, types.MetadataTypeWebApp)
	if err != nil {
		return err
	}

	val, err := dec.DecodeArrayLen()

	if err != nil {
		return err
	}

	if val != 5 {
		return errors.New(" Corrupted metadata")
	}

	for i := 0; i < val; i++ {
		switch i {
		case 0:
			wm.Name, err = dec.DecodeString()
			if err != nil {
				return err
			}
		case 1:
			err = dec.Decode(&wm.TryFiles)
			if err != nil {
				return err
			}
		case 2:
			err = dec.Decode(&wm.ErrorPages)
			if err != nil {
				return err
			}
		case 3:
			err = dec.Decode(&wm.Paths)
			if err != nil {
				return err
			}

		case 4:
			err = dec.Decode(&wm.ExtraMetadata)
			if err != nil {
				return err
			}
		default:
			return errors.New(" Corrupted metadata")
		}
	}

	wm.Type = "web_app"

	return nil
}

type WebAppFileMap struct {
	linkedhashmap.Map
}

func NewWebAppFileMap() *WebAppFileMap {
	return &WebAppFileMap{*linkedhashmap.New()}
}

func (wafm *WebAppFileMap) Put(key string, value WebAppMetadataFileReference) {
	wafm.Map.Put(key, value)
}

func (wafm *WebAppFileMap) Get(key string) (WebAppMetadataFileReference, bool) {
	value, found := wafm.Map.Get(key)
	if !found {
		return WebAppMetadataFileReference{}, false
	}
	return value.(WebAppMetadataFileReference), true
}

func (wafm *WebAppFileMap) Remove(key string) {
	wafm.Map.Remove(key)
}

func (wafm *WebAppFileMap) Keys() []string {
	keys := wafm.Map.Keys()
	ret := make([]string, len(keys))
	for i, key := range keys {
		ret[i] = key.(string)
	}
	return ret
}

func (wafm *WebAppFileMap) Values() []WebAppMetadataFileReference {
	values := wafm.Map.Values()
	ret := make([]WebAppMetadataFileReference, len(values))
	for i, value := range values {
		ret[i] = value.(WebAppMetadataFileReference)
	}
	return ret
}

func (wafm *WebAppFileMap) Sort() {
	keys := wafm.Keys()
	newMap := NewWebAppFileMap()

	sort.Strings(keys)

	for _, key := range keys {
		value, _ := wafm.Get(key)
		newMap.Put(key, value)
	}
	wafm.Map = newMap.Map
}

func (wafm *WebAppFileMap) EncodeMsgpack(encoder *msgpack.Encoder) error {
	wafm.Sort()

	for _, key := range wafm.Keys() {
		value, _ := wafm.Get(key)
		err := encoder.EncodeString(key)
		if err != nil {
			return err
		}
		err = encoder.Encode(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (wafm *WebAppFileMap) DecodeMsgpack(decoder *msgpack.Decoder) error {
	arrLen, err := decoder.DecodeArrayLen()
	if err != nil {
		return err
	}

	for i := 0; i < arrLen; i++ {
		key, err := decoder.DecodeString()
		if err != nil {
			return err
		}
		var value WebAppMetadataFileReference
		err = decoder.Decode(&value)
		if err != nil {
			return err
		}
		wafm.Put(key, value)
	}

	wafm.Sort()

	return nil
}
