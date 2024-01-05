package metadata

type MediaMetadataDetails struct {
	Data map[int]interface{}
}

func NewMediaMetadataDetails(data map[int]interface{}) *MediaMetadataDetails {
	return &MediaMetadataDetails{Data: data}
}
