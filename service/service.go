package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/node"
)

type Service interface {
	Node() *node.NodeImpl
	Start() error
	Stop() error
	Init() error
}
