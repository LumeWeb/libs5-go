package metadata

import (
	"encoding/json"
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var names = map[types.MetadataExtension]string{
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

var namesReverse = map[string]types.MetadataExtension{
	"licenses":           types.MetadataExtensionLicenses,
	"donationKeys":       types.MetadataExtensionDonationKeys,
	"wikidataClaims":     types.MetadataExtensionWikidataClaims,
	"languages":          types.MetadataExtensionLanguages,
	"sourceUris":         types.MetadataExtensionSourceUris,
	"previousVersions":   types.MetadataExtensionPreviousVersions,
	"timestamp":          types.MetadataExtensionTimestamp,
	"originalTimestamp":  types.MetadataExtensionOriginalTimestamp,
	"tags":               types.MetadataExtensionTags,
	"categories":         types.MetadataExtensionCategories,
	"basicMediaMetadata": types.MetadataExtensionBasicMediaMetadata,
	"viewTypes":          types.MetadataExtensionViewTypes,
	"bridge":             types.MetadataExtensionBridge,
	"routingHints":       types.MetadataExtensionRoutingHints,
}

type ExtraMetadata struct {
	Data map[int]interface{}
}

var _ SerializableMetadata = (*ExtraMetadata)(nil)

func NewExtraMetadata(data map[int]interface{}) *ExtraMetadata {
	return &ExtraMetadata{
		Data: data,
	}
}

func (em ExtraMetadata) MarshalJSON() ([]byte, error) {
	jsonObject := make(map[string]interface{})
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

	return json.Marshal(jsonObject)
}

func (em *ExtraMetadata) UnmarshalJSON(data []byte) error {
	// Intermediate representation of the expected JSON structure
	jsonObject := make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonObject); err != nil {
		return err
	}

	em.Data = make(map[int]interface{})
	for name, value := range jsonObject {
		if key, ok := namesReverse[name]; ok {
			if key == types.MetadataExtensionUpdateCID {
				// Convert string back to CID bytes
				cid, err := encoding.Decode(value.(string))
				if err != nil {
					return err
				}
				em.Data[int(key)] = cid
			} else {
				em.Data[int(key)] = value
			}
		} else {
			return errors.New("unknown key in JSON: " + name)
		}
	}

	return nil
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
