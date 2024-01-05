package metadata

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
)

// MetadataParentLink represents the structure for Metadata Parent Link.
type MetadataParentLink struct {
	CID    *encoding.CID
	Type   types.ParentLinkType
	Role   string
	Signed bool
}

// NewMetadataParentLink creates a new MetadataParentLink with the provided values.
func NewMetadataParentLink(cid *encoding.CID, role string, signed bool) *MetadataParentLink {
	return &MetadataParentLink{
		CID:    cid,
		Type:   types.ParentLinkTypeUserIdentity,
		Role:   role,
		Signed: signed,
	}
}
