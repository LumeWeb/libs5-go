package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type Service interface {
	Start() error
	Stop() error
	Init() error
	SetServices(services Services)
}
type Services interface {
	P2P() P2PServiceInterface
	Registry() RegistryServiceInterface
	HTTP() HTTPServiceInterface
	Storage() StorageServiceInterface
	All() []Service
	IsStarted() bool
	Start() error
	Stop() error
}

type ServiceParams struct {
	Logger *zap.Logger
	Config *config.NodeConfig
	Db     *bolt.DB
}

type ServiceBase struct {
	logger   *zap.Logger
	config   *config.NodeConfig
	db       *bolt.DB
	services Services
}

func (s *ServiceBase) SetServices(services Services) {
	s.services = services
}
