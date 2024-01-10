package node

import "git.lumeweb.com/LumeWeb/libs5-go/interfaces"

var (
	_ interfaces.Services = (*ServicesImpl)(nil)
)

type ServicesImpl struct {
	p2p      interfaces.P2PService
	registry interfaces.RegistryService
}

func (s *ServicesImpl) Registry() interfaces.RegistryService {
	return s.registry
}

func NewServices(p2p interfaces.P2PService) *ServicesImpl {
	return &ServicesImpl{p2p: p2p}
}

func (s *ServicesImpl) P2P() interfaces.P2PService {
	return s.p2p
}
