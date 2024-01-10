package node

import "git.lumeweb.com/LumeWeb/libs5-go/interfaces"

var (
	_ interfaces.Services = (*ServicesImpl)(nil)
)

type ServicesImpl struct {
	p2p      interfaces.P2PService
	registry interfaces.RegistryService
}

func (s *ServicesImpl) All() []interfaces.Service {
	services := make([]interfaces.Service, 0)
	services = append(services, s.p2p)
	services = append(services, s.registry)

	return services
}

func (s *ServicesImpl) Registry() interfaces.RegistryService {
	return s.registry
}

func NewServices(p2p interfaces.P2PService, registry interfaces.RegistryService) interfaces.Services {
	return &ServicesImpl{
		p2p:      p2p,
		registry: registry,
	}
}

func (s *ServicesImpl) P2P() interfaces.P2PService {
	return s.p2p
}
