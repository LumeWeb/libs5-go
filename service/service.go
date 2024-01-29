package service

import "git.lumeweb.com/LumeWeb/libs5-go/node"

type Service interface {
	Node() *node.Node
	Start() error
	Stop() error
	Init() error
}
type Services interface {
	P2P() *P2PService
	Registry() *RegistryService
	HTTP() *HTTPService
	All() []Service
}
