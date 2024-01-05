package metadata

import "git.lumeweb.com/LumeWeb/libs5-go/encoding"

type WebAppMetadataFileReference struct {
	ContentType *string
	Cid         *encoding.CID
}

func NewWebAppMetadataFileReference(cid *encoding.CID, contentType *string) *WebAppMetadataFileReference {
	return &WebAppMetadataFileReference{
		Cid:         cid,
		ContentType: contentType,
	}
}
