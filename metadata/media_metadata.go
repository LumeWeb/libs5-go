package metadata

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/serialize"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ Metadata              = (*MediaMetadata)(nil)
	_ msgpack.CustomDecoder = (*MediaMetadata)(nil)
	_ msgpack.CustomEncoder = (*MediaMetadata)(nil)
	_ msgpack.CustomDecoder = (*mediaMap)(nil)
)

type mediaMap map[string][]MediaFormat

type MediaMetadata struct {
	Name          string
	MediaTypes    mediaMap
	Parents       []MetadataParentLink
	Details       MediaMetadataDetails
	Links         *MediaMetadataLinks
	ExtraMetadata ExtraMetadata
	BaseMetadata
}

func NewMediaMetadata(name string, details MediaMetadataDetails, parents []MetadataParentLink, mediaTypes map[string][]MediaFormat, links *MediaMetadataLinks, extraMetadata ExtraMetadata) *MediaMetadata {
	return &MediaMetadata{
		Name:          name,
		Details:       details,
		Parents:       parents,
		MediaTypes:    mediaTypes,
		Links:         links,
		ExtraMetadata: extraMetadata,
	}
}
func NewEmptyMediaMetadata() *MediaMetadata {
	return &MediaMetadata{}
}

func (m *MediaMetadata) EncodeMsgpack(enc *msgpack.Encoder) error {
	return errors.New("Not implemented")
}

func (m *MediaMetadata) DecodeMsgpack(dec *msgpack.Decoder) error {
	kind, err := serialize.InitUnmarshaller(dec, types.MetadataTypeProof, types.MetadataTypeMedia)
	if err != nil {
		return err
	}

	switch kind {
	case types.MetadataTypeProof:
		return m.decodeProof(dec)
	case types.MetadataTypeMedia:
		return m.decodeMedia(dec)
	default:
		return errors.New("Invalid metadata type")
	}
}

func (m *MediaMetadata) decodeProof(dec *msgpack.Decoder) error {
	return errors.New("Not implemented")
}

func (m *MediaMetadata) decodeMedia(dec *msgpack.Decoder) error {
	_, err := dec.DecodeArrayLen()
	if err != nil {
		return err
	}

	err = dec.Decode(&m.Name)
	if err != nil {
		return err
	}

	err = dec.Decode(&m.Details)
	if err != nil {
		return err
	}

	err = dec.Decode(&m.Parents)
	if err != nil {
		return err
	}

	err = dec.Decode(&m.MediaTypes)
	if err != nil {
		return err
	}

	err = dec.Decode(&m.Links)
	if err != nil {
		return err
	}

	err = dec.Decode(&m.ExtraMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (m *mediaMap) DecodeMsgpack(dec *msgpack.Decoder) error {
	mapLen, err := dec.DecodeMapLen()
	if err != nil {
		return err
	}

	for i := 0; i < mapLen; i++ {
		typ, err := dec.DecodeString()
		if err != nil {
			return err
		}

		var formats []MediaFormat

		err = dec.Decode(&formats)

		if err != nil {
			return err
		}

		(*m)[typ] = formats
	}

	return nil
}
