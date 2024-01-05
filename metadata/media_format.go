package metadata

import "git.lumeweb.com/LumeWeb/libs5-go/encoding"

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
