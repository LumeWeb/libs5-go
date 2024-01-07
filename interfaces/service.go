package interfaces

type Service interface {
	Node() *Node
	Start() error
	Stop() error
	Init() error
}
type Services interface {
	P2P() *P2PService
}
