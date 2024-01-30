package fx

import (
	"git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	_default "git.lumeweb.com/LumeWeb/libs5-go/service/default"
	"go.uber.org/fx"
)

var Module = fx.Module("libs5",
	fx.Provide(newP2P),
	fx.Provide(newRegistry),
	fx.Provide(newHTTP),
	fx.Provide(newStorage),
	fx.Provide(newServices),
	fx.Provide(node.NewNode),
)

type ServiceParams struct {
	fx.In
	service.ServiceParams
}

type ServicesParams struct {
	fx.In
	P2P      service.P2PService
	Registry service.RegistryService
	HTTP     service.HTTPService
	Storage  service.StorageService
}

func newP2P(params ServiceParams) service.P2PService {
	return _default.NewP2P(params.ServiceParams)
}

func newRegistry(params ServiceParams) service.RegistryService {
	return _default.NewRegistry(params.ServiceParams)
}
func newHTTP(params ServiceParams) service.HTTPService {
	return _default.NewHTTP(params.ServiceParams)
}

func newStorage(params ServiceParams) service.StorageService {
	return _default.NewStorage(params.ServiceParams)
}

func newServices(params ServicesParams) service.Services {
	return node.NewServices(node.ServicesParams{
		P2P:      params.P2P,
		Registry: params.Registry,
		HTTP:     params.HTTP,
		Storage:  params.Storage,
	})
}
