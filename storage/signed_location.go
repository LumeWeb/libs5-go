package storage

import "git.lumeweb.com/LumeWeb/libs5-go/encoding"

type SignedStorageLocation interface {
	String() string
	NodeId() *encoding.NodeId
	Location() StorageLocation
}
