package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type Node interface {
	Services() *Services
	HashQueryRoutingTable() *structs.Map
	IsStarted() bool
	Config() *config.NodeConfig
	Logger() *zap.Logger
	Db() *bolt.DB
	Start() error
	GetCachedStorageLocations(hash *encoding.Multihash, types []int) (map[string]*StorageLocation, error)
	AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location *StorageLocation, message []byte, config *config.NodeConfig) error
}
