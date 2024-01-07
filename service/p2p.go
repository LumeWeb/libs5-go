package service

import (
	"errors"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"net/url"
	"sort"
	"time"
)

var _ Service = (*P2P)(nil)
var _ msgpack.CustomEncoder = (*nodeVotes)(nil)
var _ msgpack.CustomDecoder = (*nodeVotes)(nil)

var (
	errUnsupportedProtocol       = errors.New("unsupported protocol")
	errConnectionIdMissingNodeID = errors.New("connection id missing node id")
)

const nodeBucketName = "nodes"

type P2P struct {
	logger         *zap.Logger
	nodeKeyPair    *ed25519.KeyPairEd25519
	localNodeID    *encoding.NodeId
	networkID      string
	nodesBucket    *bolt.Bucket
	node           *node.NodeImpl
	inited         bool
	reconnectDelay *structs.Map
	peers          *structs.Map
}

type nodeVotes struct {
	good int
	bad  int
}

func (n nodeVotes) EncodeMsgpack(enc *msgpack.Encoder) error {
	err := enc.EncodeInt(int64(n.good))
	if err != nil {
		return err
	}

	err = enc.EncodeInt(int64(n.bad))
	if err != nil {
		return err
	}

	return nil
}

func (n *nodeVotes) DecodeMsgpack(dec *msgpack.Decoder) error {
	good, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	bad, err := dec.DecodeInt()
	if err != nil {
		return err
	}

	n.good = good
	n.bad = bad

	return nil
}

func NewP2P(node *node.NodeImpl) *P2P {
	service := &P2P{
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

func (p *P2P) Node() *node.NodeImpl {
	return p.node
}

func (p *P2P) Peers() *structs.Map {
	return p.peers
}

func (p *P2P) Start() error {
	err := p.Init()
	if err != nil {
		return err
	}

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

func (p *P2P) Stop() error {
	panic("implement me")
}

func (p *P2P) Init() error {
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

	return nil
}
func (p *P2P) ConnectToNode(connectionUris []*url.URL, retried bool) error {
	if !p.Node().IsStarted() {
		return nil
	}

	unsupported, _ := url.Parse("http://0.0.0.0")
	unsupported.Scheme = "unsupported"

	var connectionUri *url.URL

	for _, uri := range connectionUris {
		if uri.Scheme == "ws:" || uri.Scheme == "wss:" {
			connectionUri = uri
			break
		}
	}

	if connectionUri == nil {
		for _, uri := range connectionUris {
			if uri.Scheme == "tcp:" {
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
		*reconnectDelay = 1
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

		delay := *p.reconnectDelay.GetInt(idString)
		p.reconnectDelay.PutInt(idString, delay*2)

		time.Sleep(time.Duration(delay) * time.Second)

		return p.ConnectToNode(connectionUris, retried)
	}

	peer, err := net.CreateTransportPeer(scheme, &net.TransportPeerConfig{
		Socket: socket,
		Uris:   []*url.URL{connectionUri},
	})

	if err != nil {
		return err
	}

	(*peer).SetId(id)
	return p.onNewPeer(peer, true)
}

func (p *P2P) onNewPeer(peer *net.Peer, verifyId bool) error {
	challenge := protocol.GenerateChallenge()

	pd := *peer
	pd.SetChallenge(challenge)

	p.onNewPeerListen(peer, verifyId)

	handshakeOpenMsg, err := protocol.NewHandshakeOpen(challenge, p.networkID).ToMessage()

	if err != nil {
		return err
	}

	err = pd.SendMessage(handshakeOpenMsg)
	if err != nil {
		return err
	}
	return nil
}
func (p *P2P) onNewPeerListen(peer *net.Peer, verifyId bool) {
	onDone := net.DoneCallback(func() {
		peerId, err := (*peer).GetId().ToString()
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

	(*peer).ListenForMessages(func(message []byte) error {
		imsg := protocol.NewIncomingMessageUnknown()

		err := msgpack.Unmarshal(message, imsg)
		if err != nil {
			return err
		}

		handler, ok := protocol.GetMessageType(imsg.GetKind())

		if ok {
			imsg.SetOriginal(message)
			handler.SetIncomingMessage(imsg)

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
		OnDone:  &onDone,
		OnError: &onError,
		Logger:  p.logger,
	})

}

func (p *P2P) readNodeScore(nodeId *encoding.NodeId) (nodeVotes, error) {
	node := p.nodesBucket.Get(nodeId.Raw())
	if node == nil {
		return nodeVotes{}, nil
	}

	var score nodeVotes
	err := msgpack.Unmarshal(node, &score)
	if err != nil {

	}

	return score, nil
}

func (p *P2P) getNodeScore(nodeId *encoding.NodeId) (float64, error) {
	if nodeId.Equals(p.localNodeID) {
		return 1, nil
	}

	score, err := p.readNodeScore(nodeId)
	if err != nil {
		return 0.5, err
	}

	return protocol.CalculateNodeScore(score.good, score.bad), nil

}
func (p *P2P) SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error) {
	scores := make(map[encoding.NodeIdCode]float64)
	var errOccurred error

	for _, nodeId := range nodes {
		score, err := p.getNodeScore(nodeId)
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
