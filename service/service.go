package service

import (
	"context"
	"git.lumeweb.com/LumeWeb/libs5-go/config"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

type ServicesSetter interface {
	SetServices(services Services)
}

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Init(ctx context.Context) error
	Logger() *zap.Logger
	Config() *config.NodeConfig
	Db() *bbolt.DB
	ServicesSetter
}
type Services interface {
	P2P() P2PService
	Registry() RegistryService
	HTTP() HTTPService
	Storage() StorageService
	All() []Service
	Init(ctx context.Context) error
	IsStarted() bool
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type ServiceParams struct {
	Logger *zap.Logger
	Config *config.NodeConfig
	Db     *bbolt.DB
}

type ServiceBase struct {
	logger   *zap.Logger
	config   *config.NodeConfig
	db       *bbolt.DB
	services Services
}

func NewServiceBase(logger *zap.Logger, config *config.NodeConfig, db *bbolt.DB) ServiceBase {
	return ServiceBase{logger: logger, config: config, db: db}
}

func (s *ServiceBase) SetServices(services Services) {
	s.services = services
}
func (s *ServiceBase) Services() Services {
	return s.services
}
func (s *ServiceBase) Logger() *zap.Logger {
	return s.logger
}
func (s *ServiceBase) Config() *config.NodeConfig {
	return s.config
}
func (s *ServiceBase) Db() *bbolt.DB {
	return s.db
}
