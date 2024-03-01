package metadata

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	_ msgpack.CustomDecoder = (*MediaFormat)(nil)
	_ msgpack.CustomEncoder = (*MediaFormat)(nil)
)

type MediaFormat struct {
	Subtype       string
	Role          string
	Ext           string
	Cid           *encoding.CID
	Height        int
	Width         int
	Languages     []string
	Asr           int
	Fps           int
	Bitrate       int
	AudioChannels int
	Vcodec        string
	Acodec        string
	Container     string
	DynamicRange  string
	Charset       string
	Value         []byte
	Duration      int
	Rows          int
	Columns       int
	Index         int
	InitRange     string
	IndexRange    string
	Caption       string
}

func (mmd *MediaFormat) EncodeMsgpack(encoder *msgpack.Encoder) error {
	return errors.New("Not implemented")
}

func NewMediaFormat(subtype string, role, ext, vcodec, acodec, container, dynamicRange, charset, initRange, indexRange, caption string, cid *encoding.CID, height, width, asr, fps, bitrate, audioChannels, duration, rows, columns, index int, languages []string, value []byte) *MediaFormat {
	return &MediaFormat{
		Subtype:       subtype,
		Role:          role,
		Ext:           ext,
		Cid:           cid,
		Height:        height,
		Width:         width,
		Languages:     languages,
		Asr:           asr,
		Fps:           fps,
		Bitrate:       bitrate,
		AudioChannels: audioChannels,
		Vcodec:        vcodec,
		Acodec:        acodec,
		Container:     container,
		DynamicRange:  dynamicRange,
		Charset:       charset,
		Value:         value,
		Duration:      duration,
		Rows:          rows,
		Columns:       columns,
		Index:         index,
		InitRange:     initRange,
		IndexRange:    indexRange,
		Caption:       caption,
	}
}
func (mmd *MediaFormat) DecodeMsgpack(dec *msgpack.Decoder) error {
	intMap, err := decodeIntMap(dec)
	if err != nil {
		return err
	}

	for key, value := range intMap {
		switch key {
		case 1:
			mmd.Cid, err = encoding.CIDFromBytes(value.([]byte))
			if err != nil {
				return err
			}

		case 2:
			mmd.Subtype = value.(string)
		case 3:
			mmd.Role = value.(string)
		case 4:
			mmd.Ext = value.(string)
		case 10:
			mmd.Height = value.(int)
		case 11:
			mmd.Width = value.(int)
		case 12:
			mmd.Languages = value.([]string)
		case 13:
			mmd.Asr = int(value.(uint16))
		case 14:
			mmd.Fps = value.(int)
		case 15:
			mmd.Bitrate = int(value.(uint16))
		case 18:
			mmd.AudioChannels = value.(int)
		case 19:
			mmd.Vcodec = value.(string)
		case 20:
			mmd.Acodec = value.(string)
		case 21:
			mmd.Container = value.(string)
		case 22:
			mmd.DynamicRange = value.(string)
		case 23:
			mmd.Charset = value.(string)
		case 24:
			mmd.Value = value.([]byte)
		case 25:
			mmd.Duration = value.(int)
		case 26:
			mmd.Rows = value.(int)
		case 27:
			mmd.Columns = value.(int)
		case 28:
			mmd.Index = value.(int)
		case 29:
			mmd.InitRange = value.(string)
		case 30:
			mmd.IndexRange = value.(string)
		case 31:
			mmd.Caption = value.(string)
		}
	}

	return nil
}
