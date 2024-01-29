package fx

import (
	"git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
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
	node.ServicesParams
}

func newP2P(params ServiceParams) *service.P2PService {
	return service.NewP2P(params.ServiceParams)
}

func newRegistry(params ServiceParams) *service.RegistryService {
	return service.NewRegistry(params.ServiceParams)
}
func newHTTP(params ServiceParams) *service.HTTPService {
	return service.NewHTTP(params.ServiceParams)
}

func newStorage(params ServiceParams) *service.StorageService {
	return service.NewStorage(params.ServiceParams)
}

func newServices(params ServicesParams) service.Services {
	return node.NewServices(params.ServicesParams)
}
