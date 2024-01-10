package interfaces

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"net/url"
)

//go:generate mockgen -source=p2p.go -destination=../mocks/interfaces/p2p.go -package=interfaces

type P2PService interface {
	Node() Node
	Peers() structs.Map
	ConnectToNode(connectionUris []*url.URL, retried bool) error
	OnNewPeer(peer net.Peer, verifyId bool) error
	OnNewPeerListen(peer net.Peer, verifyId bool)
	GetNodeScore(nodeId *encoding.NodeId) (float64, error)
	SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error)
	SignMessageSimple(message []byte) ([]byte, error)
	AddPeer(peer net.Peer) error
	SendPublicPeersToPeer(peer net.Peer, peersToSend []net.Peer) error
	SendHashRequest(hash *encoding.Multihash, kinds []types.StorageLocationType) error
	UpVote(nodeId *encoding.NodeId) error
	DownVote(nodeId *encoding.NodeId) error
	NodeId() *encoding.NodeId
	Service
}
