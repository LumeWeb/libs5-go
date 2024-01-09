package metadata

var (
	_ Metadata = (*MediaMetadata)(nil)
)

type MediaMetadata struct {
	Name          string
	MediaTypes    map[string][]MediaFormat
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
