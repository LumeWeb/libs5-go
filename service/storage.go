package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/metadata"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
)

type StorageService interface {
	SetProviderStore(store storage.ProviderStore)
	ProviderStore() storage.ProviderStore
	GetCachedStorageLocations(hash *encoding.Multihash, kinds []types.StorageLocationType) (map[string]storage.StorageLocation, error)
	AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location storage.StorageLocation, message []byte) error
	DownloadBytesByHash(hash *encoding.Multihash) ([]byte, error)
	DownloadBytesByCID(cid *encoding.CID) ([]byte, error)
	GetMetadataByCID(cid *encoding.CID) (metadata.Metadata, error)
	Service
}
