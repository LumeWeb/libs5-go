package service

import (
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"net/url"
	"reflect"
	"sort"
	"sync"
	"time"
)

var _ interfaces.P2PService = (*P2PImpl)(nil)
var _ interfaces.NodeVotes = (*NodeVotesImpl)(nil)

var (
	errUnsupportedProtocol       = errors.New("unsupported protocol")
	errConnectionIdMissingNodeID = errors.New("connection id missing node id")
)

const nodeBucketName = "nodes"

type P2PImpl struct {
	logger         *zap.Logger
	nodeKeyPair    *ed25519.KeyPairEd25519
	localNodeID    *encoding.NodeId
	networkID      string
	nodesBucket    *bolt.Bucket
	node           interfaces.Node
	inited         bool
	reconnectDelay structs.Map
	peers          structs.Map
}

func NewP2P(node interfaces.Node) *P2PImpl {
	service := &P2PImpl{
		logger:         node.Logger(),
		nodeKeyPair:    node.Config().KeyPair,
		networkID:      node.Config().P2P.Network,
		node:           node,
		inited:         false,
		reconnectDelay: structs.NewMap(),
		peers:          structs.NewMap(),
	}

	return service
}

func (p *P2PImpl) Node() interfaces.Node {
	return p.node
}

func (p *P2PImpl) Peers() structs.Map {
	return p.peers
}

func (p *P2PImpl) Start() error {
	config := p.Node().Config()
	if len(config.P2P.Peers.Initial) > 0 {
		initialPeers := config.P2P.Peers.Initial

		for _, peer := range initialPeers {
			u, err := url.Parse(peer)
			if err != nil {
				return err
			}
			err = p.ConnectToNode([]*url.URL{u}, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *P2PImpl) Stop() error {
	panic("implement me")
}

func (p *P2PImpl) Init() error {
	if p.inited {
		return nil
	}
	p.localNodeID = encoding.NewNodeId(p.nodeKeyPair.PublicKey())

	err := utils.CreateBucket(nodeBucketName, p.Node().Db(), func(bucket *bolt.Bucket) {
		p.nodesBucket = bucket
	})

	if err != nil {
		return err
	}

	p.inited = true

	return nil
}
func (p *P2PImpl) ConnectToNode(connectionUris []*url.URL, retried bool) error {
	if !p.Node().IsStarted() {
		return nil
	}

	unsupported, _ := url.Parse("http://0.0.0.0")
	unsupported.Scheme = "unsupported"

	var connectionUri *url.URL

	for _, uri := range connectionUris {
		if uri.Scheme == "ws" || uri.Scheme == "wss" {
			connectionUri = uri
			break
		}
	}

	if connectionUri == nil {
		for _, uri := range connectionUris {
			if uri.Scheme == "tcp" {
				connectionUri = uri
				break
			}
		}
	}

	if connectionUri == nil {
		connectionUri = unsupported
	}

	if connectionUri.Scheme == "unsupported" {
		return errUnsupportedProtocol
	}

	scheme := connectionUri.Scheme

	if connectionUri.User == nil {
		return errConnectionIdMissingNodeID
	}

	username := connectionUri.User.Username()
	id, err := encoding.DecodeNodeId(username)
	if err != nil {
		return err
	}

	idString, err := id.ToString()
	if err != nil {
		return err
	}

	reconnectDelay := p.reconnectDelay.GetInt(idString)
	if reconnectDelay == nil {
		delay := 1
		reconnectDelay = &delay
	}

	if id.Equals(p.localNodeID) {
		return nil
	}

	p.logger.Debug("connect", zap.String("node", connectionUri.String()))

	socket, err := net.CreateTransportSocket(scheme, connectionUri)
	if err != nil {
		if retried {
			p.logger.Error("failed to connect, too many retries", zap.String("node", connectionUri.String()), zap.Error(err))
			return nil
		}
		retried = true

		p.logger.Error("failed to connect", zap.String("node", connectionUri.String()), zap.Error(err))

		delay := p.reconnectDelay.GetInt(idString)
		if delay == nil {
			tmp := 1
			delay = &tmp
		}
		delayDeref := *delay
		p.reconnectDelay.PutInt(idString, delayDeref*2)

		time.Sleep(time.Duration(delayDeref) * time.Second)

		return p.ConnectToNode(connectionUris, retried)
	}

	peer, err := net.CreateTransportPeer(scheme, &net.TransportPeerConfig{
		Socket: socket,
		Uris:   []*url.URL{connectionUri},
	})

	if err != nil {
		return err
	}

	peer.SetId(id)

	go func() {
		err := p.OnNewPeer(peer, true)
		if err != nil {
			p.logger.Error("failed to add peer", zap.Error(err))
		}
	}()
	return nil

}

func (p *P2PImpl) OnNewPeer(peer net.Peer, verifyId bool) error {
	var wg sync.WaitGroup

	challenge := protocol.GenerateChallenge()
	peer.SetChallenge(challenge)

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.OnNewPeerListen(peer, verifyId)
	}()

	handshakeOpenMsg, err := msgpack.Marshal(protocol.NewHandshakeOpen(challenge, p.networkID))
	if err != nil {
		return err
	}

	err = peer.SendMessage(handshakeOpenMsg)
	if err != nil {
		return err
	}

	wg.Wait() // Wait for OnNewPeerListen goroutine to finish
	return nil
}
func (p *P2PImpl) OnNewPeerListen(peer net.Peer, verifyId bool) {
	onDone := net.CloseCallback(func() {
		peerId, err := peer.Id().ToString()
		if err != nil {
			p.logger.Error("failed to get peer id", zap.Error(err))
			return
		}

		// Handle closure of the connection
		if p.peers.Contains(peerId) {
			p.peers.Remove(peerId)
		}
	})

	onError := net.ErrorCallback(func(args ...interface{}) {
		p.logger.Error("peer error", zap.Any("args", args))
	})

	peer.ListenForMessages(func(message []byte) error {
		imsg := base.NewIncomingMessageUnknown()

		err := msgpack.Unmarshal(message, imsg)
		if err != nil {
			return err
		}

		handler, ok := protocol.GetMessageType(imsg.GetKind())

		if ok {

			imsg.SetOriginal(message)
			handler.SetIncomingMessage(imsg)
			handler.SetSelf(handler)
			fmt.Println(reflect.TypeOf(handler))
			err := msgpack.Unmarshal(imsg.Data(), handler)
			if err != nil {
				return err
			}
			err = handler.HandleMessage(p.node, peer, verifyId)
			if err != nil {
				return err
			}
		}

		return nil
	}, net.ListenerOptions{
		OnClose: &onDone,
		OnError: &onError,
		Logger:  p.logger,
	})

}

func (p *P2PImpl) readNodeVotes(nodeId *encoding.NodeId) (interfaces.NodeVotes, error) {
	node := p.nodesBucket.Get(nodeId.Raw())
	if node == nil {
		return NewNodeVotes(), nil
	}

	var score interfaces.NodeVotes
	err := msgpack.Unmarshal(node, &score)
	if err != nil {
		return nil, err
	}

	return score, nil
}

func (p *P2PImpl) saveNodeVotes(nodeId *encoding.NodeId) (interfaces.NodeVotes, error) {
	votes, err := p.readNodeVotes(nodeId)
	if err != nil {
		return nil, err
	}

	data, err := msgpack.Marshal(votes)
	if err != nil {
		return nil, err
	}

	err = p.nodesBucket.Put(nodeId.Raw(), data)
	if err != nil {
		return nil, err
	}

	return votes, nil
}

func (p *P2PImpl) GetNodeScore(nodeId *encoding.NodeId) (float64, error) {
	if nodeId.Equals(p.localNodeID) {
		return 1, nil
	}

	score, err := p.readNodeVotes(nodeId)
	if err != nil {
		return 0.5, err
	}

	return protocol.CalculateNodeScore(score.Good(), score.Bad()), nil

}
func (p *P2PImpl) SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error) {
	scores := make(map[encoding.NodeIdCode]float64)
	var errOccurred error

	for _, nodeId := range nodes {
		score, err := p.GetNodeScore(nodeId)
		if err != nil {
			errOccurred = err
			scores[nodeId.HashCode()] = 0 // You may choose a different default value for error cases
		} else {
			scores[nodeId.HashCode()] = score
		}
	}

	sort.Slice(nodes, func(i, j int) bool {
		return scores[nodes[i].HashCode()] > scores[nodes[j].HashCode()]
	})

	return nodes, errOccurred
}
func (p *P2PImpl) SignMessageSimple(message []byte) ([]byte, error) {
	signedMessage := signed.NewSignedMessageRequest(message)
	signedMessage.SetNodeId(p.localNodeID)

	err := signedMessage.Sign(p.Node())

	if err != nil {
		return nil, err
	}

	result, err := msgpack.Marshal(signedMessage)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *P2PImpl) AddPeer(peer net.Peer) error {
	peerId, err := peer.Id().ToString()
	if err != nil {
		return err
	}
	p.peers.Put(peerId, peer)
	p.reconnectDelay.Put(peerId, 1)

	return nil
}
func (p *P2PImpl) SendPublicPeersToPeer(peer net.Peer, peersToSend []net.Peer) error {
	announceRequest := signed.NewAnnounceRequest(peer, peersToSend)

	message, err := msgpack.Marshal(announceRequest)

	if err != nil {
		return err
	}

	signedMessage, err := p.SignMessageSimple(message)

	if err != nil {
		return err
	}

	err = peer.SendMessage(signedMessage)

	return nil
}
func (p *P2PImpl) SendHashRequest(hash *encoding.Multihash, kinds []types.StorageLocationType) error {
	hashRequest := protocol.NewHashRequest(hash, kinds)
	message, err := msgpack.Marshal(hashRequest)
	if err != nil {
		return err
	}

	for _, peer := range p.peers.Values() {
		peerValue, ok := peer.(net.Peer)
		if !ok {
			p.node.Logger().Error("failed to cast peer to net.Peer")
			continue
		}
		err = peerValue.SendMessage(message)
	}

	return nil
}

func (p *P2PImpl) UpVote(nodeId *encoding.NodeId) error {
	err := p.vote(nodeId, true)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2PImpl) DownVote(nodeId *encoding.NodeId) error {
	err := p.vote(nodeId, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2PImpl) vote(nodeId *encoding.NodeId, upvote bool) error {
	votes, err := p.readNodeVotes(nodeId)
	if err != nil {
		return err
	}

	if upvote {
		votes.Upvote()
	} else {
		votes.Downvote()
	}

	return nil
}
