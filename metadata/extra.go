package metadata

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

type ExtraMetadata struct {
	Data map[int]interface{}
}

func NewExtraMetadata(data map[int]interface{}) *ExtraMetadata {
	return &ExtraMetadata{
		Data: data,
	}
}

func (em *ExtraMetadata) ToJSON() map[string]interface{} {
	jsonObject := make(map[string]interface{})
	names := map[types.MetadataExtension]string{
		types.MetadataExtensionLicenses:           "licenses",
		types.MetadataExtensionDonationKeys:       "donationKeys",
		types.MetadataExtensionWikidataClaims:     "wikidataClaims",
		types.MetadataExtensionLanguages:          "languages",
		types.MetadataExtensionSourceUris:         "sourceUris",
		types.MetadataExtensionPreviousVersions:   "previousVersions",
		types.MetadataExtensionTimestamp:          "timestamp",
		types.MetadataExtensionOriginalTimestamp:  "originalTimestamp",
		types.MetadataExtensionTags:               "tags",
		types.MetadataExtensionCategories:         "categories",
		types.MetadataExtensionBasicMediaMetadata: "basicMediaMetadata",
		types.MetadataExtensionViewTypes:          "viewTypes",
		types.MetadataExtensionBridge:             "bridge",
		types.MetadataExtensionRoutingHints:       "routingHints",
	}

	for key, value := range em.Data {
		name, ok := names[types.MetadataExtension(key)]
		if ok {
			if types.MetadataExtension(key) == types.MetadataExtensionUpdateCID {
				cid, err := encoding.CIDFromBytes(value.([]byte))
				var cidString string
				if err == nil {
					cidString, err = cid.ToString()
				}

				if err == nil {
					jsonObject["updateCID"] = cidString
				} else {
					jsonObject["updateCID"] = ""
				}

			} else {
				jsonObject[name] = value
			}
		}
	}

	return jsonObject
}
func (em *ExtraMetadata) DecodeMsgpack(dec *msgpack.Decoder) error {
	mapLen, err := dec.DecodeMapLen()

	if err != nil {
		return err
	}

	for i := 0; i < mapLen; i++ {
		key, err := dec.DecodeInt8()
		if err != nil {
			return err
		}
		value, err := dec.DecodeInterface()
		if err != nil {
			return err
		}
		em.Data[int(key)] = value
	}

	return nil
}

func (em ExtraMetadata) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(em.Data)
}
