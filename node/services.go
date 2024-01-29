package node

import (
	"git.lumeweb.com/LumeWeb/libs5-go/service"
)

var (
	_ service.Services = (*ServicesImpl)(nil)
)

type ServicesImpl struct {
	p2p      *service.P2PService
	registry *service.RegistryService
	http     *service.HTTPService
}

func (s *ServicesImpl) HTTP() *service.HTTPService {
	return s.http
}

func (s *ServicesImpl) All() []service.Service {
	services := make([]service.Service, 0)
	services = append(services, s.p2p)
	services = append(services, s.registry)
	services = append(services, s.http)

	return services
}

func (s *ServicesImpl) Registry() *service.RegistryService {
	return s.registry
}

func NewServices(p2p *service.P2PService, registry *service.RegistryService, http *service.HTTPService) service.Services {
	return &ServicesImpl{
		p2p:      p2p,
		registry: registry,
		http:     http,
	}
}

func (s *ServicesImpl) P2P() *service.P2PService {
	return s.p2p
}
