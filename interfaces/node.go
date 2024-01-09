package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/metadata"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

//go:generate mockgen -source=node.go -destination=../mocks/interfaces/node.go -package=interfaces

type Node interface {
	Services() Services
	HashQueryRoutingTable() structs.Map
	IsStarted() bool
	Config() *config.NodeConfig
	Logger() *zap.Logger
	Db() *bolt.DB
	Start() error
	GetCachedStorageLocations(hash *encoding.Multihash, kinds []types.StorageLocationType) (map[string]StorageLocation, error)
	AddStorageLocation(hash *encoding.Multihash, nodeId *encoding.NodeId, location StorageLocation, message []byte, config *config.NodeConfig) error
	NetworkId() string
	DownloadBytesByHash(hash *encoding.Multihash) ([]byte, error)
	GetMetadataByCID(cid *encoding.CID) (metadata.Metadata, error)
}
