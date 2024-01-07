package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"net/url"
)

type P2PService interface {
	Node() *Node
	Peers() *structs.Map
	Start() error
	Stop() error
	Init() error
	ConnectToNode(connectionUris []*url.URL, retried bool) error
	onNewPeer(peer *net.Peer, verifyId bool) error
	onNewPeerListen(peer *net.Peer, verifyId bool)
	readNodeScore(nodeId *encoding.NodeId) (NodeVotes, error)
	getNodeScore(nodeId *encoding.NodeId) (float64, error)
	SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error)
}
