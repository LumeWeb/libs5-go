package service

import (
	"bytes"
	"context"
	ed25519p "crypto/ed25519"
	"errors"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/ed25519"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	_node "git.lumeweb.com/LumeWeb/libs5-go/node"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/signed"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/structs"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"net/url"
	"sort"
	"sync"
	"time"
)

var _ Service = (*P2PService)(nil)

var (
	errUnsupportedProtocol       = errors.New("unsupported protocol")
	errConnectionIdMissingNodeID = errors.New("connection id missing node id")
)

const nodeBucketName = "nodes"

type P2PService struct {
	logger                  *zap.Logger
	nodeKeyPair             *ed25519.KeyPairEd25519
	localNodeID             *encoding.NodeId
	networkID               string
	nodesBucket             *bolt.Bucket
	node                    *_node.Node
	inited                  bool
	reconnectDelay          structs.Map
	peers                   structs.Map
	peersPending            structs.Map
	selfConnectionUris      []*url.URL
	outgoingPeerBlocklist   structs.Map
	incomingPeerBlockList   structs.Map
	incomingIPBlocklist     structs.Map
	outgoingPeerFailures    structs.Map
	maxOutgoingPeerFailures uint
}

func NewP2P(node *_node.Node) *P2PService {
	uri, err := url.Parse(fmt.Sprintf("wss://%s:%d/s5/p2p", node.Config().HTTP.API.Domain, node.Config().HTTP.API.Port))
	if err != nil {
		node.Logger().Fatal("failed to HTTP API URL Config", zap.Error(err))
	}

	service := &P2PService{
		logger:                  node.Logger(),
		nodeKeyPair:             node.Config().KeyPair,
		networkID:               node.Config().P2P.Network,
		node:                    node,
		inited:                  false,
		reconnectDelay:          structs.NewMap(),
		peers:                   structs.NewMap(),
		peersPending:            structs.NewMap(),
		selfConnectionUris:      []*url.URL{uri},
		outgoingPeerBlocklist:   structs.NewMap(),
		incomingPeerBlockList:   structs.NewMap(),
		incomingIPBlocklist:     structs.NewMap(),
		outgoingPeerFailures:    structs.NewMap(),
		maxOutgoingPeerFailures: node.Config().P2P.MaxOutgoingPeerFailures,
	}

	return service
}

func (p *P2PService) SelfConnectionUris() []*url.URL {
	return p.selfConnectionUris
}

func (p *P2PService) Node() *_node.Node {
	return p.node
}

func (p *P2PService) Peers() structs.Map {
	return p.peers
}

func (p *P2PService) Start() error {
	config := p.Node().Config()
	if len(config.P2P.Peers.Initial) > 0 {
		initialPeers := config.P2P.Peers.Initial

		for _, peer := range initialPeers {
			u, err := url.Parse(peer)
			if err != nil {
				return err
			}

			peer := peer
			go func() {
				err := p.ConnectToNode([]*url.URL{u}, false, nil)
				if err != nil {
					p.logger.Error("failed to connect to initial peer", zap.Error(err), zap.String("peer", peer))
				}
			}()
		}
	}

	return nil
}

func (p *P2PService) Stop() error {
	panic("implement me")
}

func (p *P2PService) Init() error {
	if p.inited {
		return nil
	}
	p.localNodeID = encoding.NewNodeId(p.nodeKeyPair.PublicKey())

	err := utils.CreateBucket(nodeBucketName, p.Node().Db())

	if err != nil {
		return err
	}

	p.inited = true

	return nil
}
func (p *P2PService) ConnectToNode(connectionUris []*url.URL, retried bool, fromPeer net.Peer) error {
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

	if p.peersPending.Contains(idString) || p.peers.Contains(idString) {
		p.logger.Debug("already connected", zap.String("node", connectionUri.String()))
		return nil
	}

	if p.outgoingPeerBlocklist.Contains(idString) {
		p.logger.Debug("outgoing peer is on blocklist", zap.String("node", connectionUri.String()))

		var fromPeerId string

		if fromPeer != nil {
			blocked := false
			if fromPeer.Id() != nil {
				fromPeerId, err = fromPeer.Id().ToString()
				if err != nil {
					return err
				}
				if !p.incomingPeerBlockList.Contains(fromPeerId) {
					p.incomingPeerBlockList.Put(fromPeerId, true)
					blocked = true
				}
			}

			fromPeerIP := fromPeer.GetIP()

			if !p.incomingIPBlocklist.Contains(fromPeerIP) {
				p.incomingIPBlocklist.Put(fromPeerIP, true)
				blocked = true
			}
			err = fromPeer.EndForAbuse()
			if err != nil {
				return err
			}

			if blocked {
				p.logger.Debug("blocking peer for sending peer on blocklist", zap.String("node", connectionUri.String()), zap.String("peer", fromPeerId), zap.String("ip", fromPeerIP))
			}
		}
		return nil
	}

	reconnectDelay := p.reconnectDelay.GetUInt(idString)
	if reconnectDelay == nil {
		delay := uint(1)
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
			counter := uint(0)
			if p.outgoingPeerFailures.Contains(idString) {
				tmp := *p.outgoingPeerFailures.GetUInt(idString)
				counter = tmp
			}

			counter++

			p.outgoingPeerFailures.PutUInt(idString, counter)

			if counter >= p.maxOutgoingPeerFailures {

				if fromPeer != nil {
					blocked := false
					var fromPeerId string
					if fromPeer.Id() != nil {
						fromPeerId, err = fromPeer.Id().ToString()
						if err != nil {
							return err
						}
						if !p.incomingPeerBlockList.Contains(fromPeerId) {
							p.incomingPeerBlockList.Put(fromPeerId, true)
							blocked = true
						}
					}

					fromPeerIP := fromPeer.GetIP()
					if !p.incomingIPBlocklist.Contains(fromPeerIP) {
						p.incomingIPBlocklist.Put(fromPeerIP, true)
						blocked = true
					}
					err = fromPeer.EndForAbuse()
					if err != nil {
						return err
					}

					if blocked {
						p.logger.Debug("blocking peer for sending peer on blocklist", zap.String("node", connectionUri.String()), zap.String("peer", fromPeerId), zap.String("ip", fromPeerIP))
					}
				}
				p.outgoingPeerBlocklist.Put(idString, true)
				p.logger.Debug("blocking peer for too many failures", zap.String("node", connectionUri.String()))
			}

			return nil
		}
		retried = true

		p.logger.Error("failed to connect", zap.String("node", connectionUri.String()), zap.Error(err))

		delay := p.reconnectDelay.GetUInt(idString)
		if delay == nil {
			tmp := uint(1)
			delay = &tmp
		}
		delayDeref := *delay
		p.reconnectDelay.PutUInt(idString, delayDeref*2)

		time.Sleep(time.Duration(delayDeref) * time.Second)

		return p.ConnectToNode(connectionUris, retried, fromPeer)
	}

	if p.outgoingPeerFailures.Contains(idString) {
		p.outgoingPeerFailures.Remove(idString)
	}

	peer, err := net.CreateTransportPeer(scheme, &net.TransportPeerConfig{
		Socket: socket,
		Uris:   []*url.URL{connectionUri},
	})

	if err != nil {
		return err
	}

	peer.SetId(id)

	p.Node().ConnectionTracker().Add(1)

	peerId, err := peer.Id().ToString()
	if err != nil {
		return err
	}
	p.peersPending.Put(peerId, peer)

	go func() {
		err := p.OnNewPeer(peer, true)
		if err != nil && !peer.Abuser() {
			p.logger.Error("peer error", zap.Error(err))
		}
		p.Node().ConnectionTracker().Done()
	}()

	return nil

}

func (p *P2PService) OnNewPeer(peer net.Peer, verifyId bool) error {
	var wg sync.WaitGroup

	var pid string

	if peer.Id() != nil {
		pid, _ = peer.Id().ToString()
	} else {
		pid = "unknown"
	}

	pip := peer.GetIP()

	if p.incomingIPBlocklist.Contains(pid) {
		p.logger.Error("peer is on identity blocklist", zap.String("peer", pid))
		err := peer.EndForAbuse()
		if err != nil {
			return err
		}
		return nil
	}
	if p.incomingPeerBlockList.Contains(pip) {
		p.logger.Debug("peer is on ip blocklist", zap.String("peer", pid), zap.String("ip", pip))
		err := peer.EndForAbuse()
		if err != nil {
			return err
		}
		return nil
	}

	p.logger.Debug("OnNewPeer started", zap.String("peer", pid))

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
	p.logger.Debug("OnNewPeer sent handshake", zap.String("peer", pid))

	p.logger.Debug("OnNewPeer before Wait", zap.String("peer", pid))
	wg.Wait() // Wait for OnNewPeerListen goroutine to finish
	p.logger.Debug("OnNewPeer ended", zap.String("peer", pid))
	return nil
}
func (p *P2PService) OnNewPeerListen(peer net.Peer, verifyId bool) {
	onDone := net.CloseCallback(func() {
		if peer.Id() != nil {
			pid, err := peer.Id().ToString()
			if err != nil {
				p.logger.Error("failed to get peer id", zap.Error(err))
				return
			}
			// Handle closure of the connection
			if p.peers.Contains(pid) {
				p.peers.Remove(pid)
			}
			if p.peersPending.Contains(pid) {
				p.peersPending.Remove(pid)
			}
		}
	})

	onError := net.ErrorCallback(func(args ...interface{}) {
		if !peer.Abuser() {
			p.logger.Error("peer error", zap.Any("args", args))
		}
	})

	peer.ListenForMessages(func(message []byte) error {
		var reader base.IncomingMessageReader

		err := msgpack.Unmarshal(message, &reader)
		if err != nil {
			p.logger.Error("Error decoding basic message info", zap.Error(err))
			return err
		}

		// Now, get the specific message handler based on the message kind
		handler, ok := protocol.GetMessageType(reader.Kind)
		if !ok {
			p.logger.Error("Unknown message type", zap.Int("type", reader.Kind))
			return fmt.Errorf("unknown message type: %d", reader.Kind)
		}

		data := base.IncomingMessageData{
			Original: message,
			Data:     reader.Data,
			Ctx:      context.Background(),
			Node:     p.node,
			Peer:     peer,
			VerifyId: verifyId,
		}

		dec := msgpack.NewDecoder(bytes.NewReader(reader.Data))

		err = handler.DecodeMessage(dec, data)
		if err != nil {
			p.logger.Error("Error decoding message", zap.Error(err))
			return err
		}

		// Directly decode and handle the specific message type
		if err := handler.HandleMessage(data); err != nil {
			p.logger.Error("Error handling message", zap.Error(err))
			return err
		}

		return nil
	}, net.ListenerOptions{
		OnClose: &onDone,
		OnError: &onError,
		Logger:  p.logger,
	})
}

func (p *P2PService) readNodeVotes(nodeId *encoding.NodeId) (NodeVotes, error) {
	var value []byte
	var found bool
	err := p.node.Db().View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(nodeBucketName))
		if b == nil {
			return fmt.Errorf("Bucket %s not found", nodeBucketName)
		}
		value = b.Get(nodeId.Raw())
		if value == nil {
			return nil
		}
		found = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return NewNodeVotes(), nil
	}

	var score NodeVotes
	err = msgpack.Unmarshal(value, &score)
	if err != nil {
		return nil, err
	}

	return score, nil
}

func (p *P2PService) saveNodeVotes(nodeId *encoding.NodeId, votes NodeVotes) error {
	// Marshal the votes into data
	data, err := msgpack.Marshal(votes)
	if err != nil {
		return err
	}

	// Use a transaction to save the data
	err = p.node.Db().Update(func(tx *bolt.Tx) error {
		// Get or create the bucket
		b := tx.Bucket([]byte(nodeBucketName))

		// Put the data into the bucket
		return b.Put(nodeId.Raw(), data)
	})

	return err
}

func (p *P2PService) GetNodeScore(nodeId *encoding.NodeId) (float64, error) {
	if nodeId.Equals(p.localNodeID) {
		return 1, nil
	}

	score, err := p.readNodeVotes(nodeId)
	if err != nil {
		return 0.5, err
	}

	return protocol.CalculateNodeScore(score.Good(), score.Bad()), nil

}
func (p *P2PService) SortNodesByScore(nodes []*encoding.NodeId) ([]*encoding.NodeId, error) {
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
func (p *P2PService) SignMessageSimple(message []byte) ([]byte, error) {
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

func (p *P2PService) AddPeer(peer net.Peer) error {
	peerId, err := peer.Id().ToString()
	if err != nil {
		return err
	}
	p.peers.Put(peerId, peer)
	p.reconnectDelay.PutUInt(peerId, 1)

	if p.peersPending.Contains(peerId) {
		p.peersPending.Remove(peerId)
	}

	return nil
}
func (p *P2PService) SendPublicPeersToPeer(peer net.Peer, peersToSend []net.Peer) error {
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
func (p *P2PService) SendHashRequest(hash *encoding.Multihash, kinds []types.StorageLocationType) error {
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

func (p *P2PService) UpVote(nodeId *encoding.NodeId) error {
	err := p.vote(nodeId, true)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2PService) DownVote(nodeId *encoding.NodeId) error {
	err := p.vote(nodeId, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *P2PService) vote(nodeId *encoding.NodeId, upvote bool) error {
	votes, err := p.readNodeVotes(nodeId)
	if err != nil {
		return err
	}

	if upvote {
		votes.Upvote()
	} else {
		votes.Downvote()
	}

	err = p.saveNodeVotes(nodeId, votes)
	if err != nil {
		return err
	}

	return nil
}
func (p *P2PService) NodeId() *encoding.NodeId {
	return p.localNodeID
}

func (p *P2PService) PrepareProvideMessage(hash *encoding.Multihash, location storage.StorageLocation) []byte {
	// Initialize the list with the record type.
	list := []byte{byte(types.RecordTypeStorageLocation)}

	// Append the full bytes of the hash.
	list = append(list, hash.FullBytes()...)

	// Append the location type.
	list = append(list, byte(location.Type()))

	// Append the expiry time of the location, encoded as 4 bytes.
	list = append(list, utils.EncodeEndian(uint64(location.Expiry()), 4)...)

	// Append the number of parts in the location.
	list = append(list, byte(len(location.Parts())))

	// Iterate over each part in the location.
	for _, part := range location.Parts() {
		// Convert part to bytes.
		bytes := []byte(part)

		// Encode the length of the part as 4 bytes and append.
		list = append(list, utils.EncodeEndian(uint64(len(bytes)), 2)...)

		// Append the actual part bytes.
		list = append(list, bytes...)
	}

	// Append a null byte at the end of the list.
	list = append(list, 0)

	// Sign the list using the node's private key.
	signature := ed25519p.Sign(p.nodeKeyPair.ExtractBytes(), list)

	// Append the public key and signature to the list.
	finalList := append(list, p.nodeKeyPair.PublicKey()...)
	finalList = append(finalList, signature...)

	// Return the final byte slice.
	return finalList
}
