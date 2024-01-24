package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
)

type ProviderStore interface {
	CanProvide(hash *encoding.Multihash, kind []types.StorageLocationType) bool
	Provide(hash *encoding.Multihash, kind []types.StorageLocationType) (StorageLocation, error)
}
