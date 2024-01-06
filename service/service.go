package service

import libs5_go "git.lumeweb.com/LumeWeb/libs5-go"

type Service interface {
	Node() *libs5_go.Node
	Start() error
	Stop() error
	Init() error
}
