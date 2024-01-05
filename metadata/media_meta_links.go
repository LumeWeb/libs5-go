package metadata

import "git.lumeweb.com/LumeWeb/libs5-go/encoding"

type MediaMetadataLinks struct {
	Count     int
	Head      []*encoding.CID
	Collapsed []*encoding.CID
	Tail      []*encoding.CID
}

func NewMediaMetadataLinks(head []*encoding.CID) *MediaMetadataLinks {
	return &MediaMetadataLinks{
		Count: len(head),
		Head:  head,
	}
}
