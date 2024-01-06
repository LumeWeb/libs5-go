package libs5_go

import (
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type NodeConfig struct {
	P2P     P2PConfig
	KeyPair ed25519.KeyPairEd25519
	DB      bolt.DB
	Logger  *zap.Logger
}
type P2PConfig struct {
	Network string
	Peers   PeersConfig
}

type PeersConfig struct {
	Initial []string
}
