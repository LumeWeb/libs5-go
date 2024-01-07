package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"net/url"
)

type P2PService interface {
	Node() Node
	Peers() structs.Map
	Start() error
	Stop() error
	Init() error
	ConnectToNode(connectionUris []*url.URL, retried bool) error
	OnNewPeer(peer net.Peer, verifyId bool) error
	OnNewPeerListen(peer net.Peer, verifyId bool)
	ReadNodeScore(nodeId *encoding.NodeId) (NodeVotes, error)
	GetNodeScore(nodeId *encoding.NodeId) (float64, error)
	SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error)
	SignMessageSimple(message []byte) ([]byte, error)
}
