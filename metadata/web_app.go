package metadata

var (
	_ Metadata = (*WebAppMetadata)(nil)
)

type WebAppMetadata struct {
	BaseMetadata
	Name          string
	TryFiles      []string
	ErrorPages    map[int]string
	ExtraMetadata ExtraMetadata
	Paths         map[string]WebAppMetadataFileReference
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
