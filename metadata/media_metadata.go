package metadata

type MediaMetadata struct {
	Metadata
	Name          string
	MediaTypes    map[string][]MediaFormat
	Parents       []MetadataParentLink
	Details       MediaMetadataDetails
	Links         *MediaMetadataLinks
	ExtraMetadata ExtraMetadata
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
