package node

import (
	"git.lumeweb.com/LumeWeb/libs5-go/service"
)

var (
	_ service.Services = (*ServicesImpl)(nil)
)

type ServicesParams struct {
	P2P      *service.P2PService
	Registry *service.RegistryService
	HTTP     *service.HTTPService
}

type ServicesImpl struct {
	p2p      *service.P2PService
	registry *service.RegistryService
	http     *service.HTTPService
	started  bool
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

func NewServices(params ServicesParams) service.Services {
	sc := &ServicesImpl{
		p2p:      params.P2P,
		registry: params.Registry,
		http:     params.HTTP,
		started:  false,
	}

	for _, svc := range sc.All() {
		svc.SetServices(sc)
	}

	return sc
}

func (s *ServicesImpl) P2P() *service.P2PService {
	return s.p2p
}

func (s *ServicesImpl) IsStarted() bool {
	return s.started
}

func (s *ServicesImpl) Start() error {
	for _, svc := range s.All() {
		err := svc.Start()
		if err != nil {
			return err
		}
	}

	s.started = true

	return nil
}
func (s *ServicesImpl) Stop() error {
	for _, svc := range s.All() {
		err := svc.Stop()
		if err != nil {
			return err
		}
	}

	s.started = false

	return nil
}
