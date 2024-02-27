package config

import (
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type NodeConfig struct {
	P2P     P2PConfig `mapstructure:"p2p"`
	KeyPair *ed25519.KeyPairEd25519
	DB      *bolt.DB
	Logger  *zap.Logger
	HTTP    HTTPConfig `mapstructure:"http"`
}
type P2PConfig struct {
	Network                 string      `mapstructure:"network"`
	Peers                   PeersConfig `mapstructure:"peers"`
	MaxOutgoingPeerFailures uint        `mapstructure:"max_outgoing_peer_failures"`
}

type PeersConfig struct {
	Initial []string `mapstructure:"initial"`
}

type HTTPAPIConfig struct {
	Domain string `mapstructure:"domain"`
	Port   uint   `mapstructure:"port"`
}

type HTTPConfig struct {
	API HTTPAPIConfig `mapstructure:"api"`
}
