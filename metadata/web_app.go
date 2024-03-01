package metadata

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/serialize"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/samber/lo"
	"github.com/vmihailenco/msgpack/v5"
	"sort"
)

var (
	_ Metadata             = (*WebAppMetadata)(nil)
	_ SerializableMetadata = (*WebAppMetadata)(nil)
)

type WebAppMetadata struct {
	BaseMetadata
	Name          string                                 `json:"name"`
	TryFiles      []string                               `json:"tryFiles"`
	ErrorPages    map[int]string                         `json:"errorPages"`
	ExtraMetadata ExtraMetadata                          `json:"extraMetadata"`
	Paths         map[string]WebAppMetadataFileReference `json:"paths"`
}

func NewWebAppMetadata(name string, tryFiles []string, extraMetadata ExtraMetadata, errorPages map[int]string, paths map[string]WebAppMetadataFileReference) *WebAppMetadata {
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

	keys := lo.Keys[string, WebAppMetadataFileReference](wm.Paths)
	sort.Strings(keys)

	paths := make([]WebAppMetadataFileReference, len(wm.Paths))

	for i, v := range keys {
		paths[i] = wm.Paths[v]
	}

	items := make([]interface{}, 5)

	items[0] = wm.Name
	items[1] = wm.TryFiles
	items[2] = wm.ErrorPages
	items[3] = paths
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
			paths, err := dec.DecodeSlice()
			if err != nil {
				return err
			}

			wm.Paths = make(map[string]WebAppMetadataFileReference, len(paths))

			for _, v := range paths {
				path := v.([]interface{})
				parsedCid, err := encoding.CIDFromBytes(path[1].([]byte))
				if err != nil {
					return err
				}
				contentType := ""

				if path[2] != nil {
					contentType = path[2].(string)
				}

				wm.Paths[path[0].(string)] = WebAppMetadataFileReference{
					Cid:         parsedCid,
					ContentType: contentType,
				}
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
