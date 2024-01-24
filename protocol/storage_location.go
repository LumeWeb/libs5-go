package protocol

import (
	"crypto/ed25519"
	"fmt"
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/storage"
	"git.lumeweb.com/LumeWeb/libs5-go/types"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
)

var _ base.IncomingMessageTyped = (*StorageLocation)(nil)

type StorageLocation struct {
	hash      *encoding.Multihash
	kind      int
	expiry    int64
	parts     []string
	publicKey []byte
	signature []byte

	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
}

func NewStorageLocation() *StorageLocation {
	sl := &StorageLocation{}

	sl.SetRequiresHandshake(true)

	return sl
}

func (s *StorageLocation) DecodeMessage(dec *msgpack.Decoder) error {
	// nop, we use the incoming message -> original already stored
	return nil
}
func (s *StorageLocation) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	msg := s.IncomingMessage().Original()

	hash := encoding.NewMultihash(msg[1:34]) // Replace NewMultihash with appropriate function

	typeOfData := msg[34]

	expiry := utils.DecodeEndian(msg[35:39])

	partCount := msg[39]

	parts := []string{}
	cursor := 40
	for i := 0; i < int(partCount); i++ {
		length := utils.DecodeEndian(msg[cursor : cursor+2])
		cursor += 2
		if len(msg) < cursor+int(length) {
			return fmt.Errorf("Invalid message")
		}
		part := string(msg[cursor : cursor+int(length)])
		parts = append(parts, part)
		cursor += int(length)
	}

	cursor++

	publicKey := msg[cursor : cursor+33]
	signature := msg[cursor+33:]

	if types.HashType(publicKey[0]) != types.HashTypeEd25519 { // Replace CID_HASH_TYPES_ED25519 with actual constant
		return fmt.Errorf("Unsupported public key type %d", publicKey[0])
	}

	if !ed25519.Verify(publicKey[1:], msg[:cursor], signature) {
		return fmt.Errorf("Signature verification failed")
	}

	nodeId := encoding.NewNodeId(publicKey)

	// Assuming `node` is an instance of your NodeImpl structure
	err := node.AddStorageLocation(hash, nodeId, storage.NewStorageLocation(int(typeOfData), parts, int64(expiry)), msg) // Implement AddStorageLocation

	if err != nil {
		return fmt.Errorf("Failed to add storage location: %s", err)
	}

	hashStr, err := hash.ToString()
	if err != nil {
		return err
	}

	var list *hashset.Set
	listVal, ok := node.HashQueryRoutingTable().Get(hashStr) // Implement HashQueryRoutingTable method
	if !ok {
		list = hashset.New()
	} else {
		list = listVal.(*hashset.Set)
	}

	for _, peerIdVal := range list.Values() {
		peerId := peerIdVal.(*encoding.NodeId)

		if peerId.Equals(nodeId) || peerId.Equals(peer) {
			continue
		}
		peerIdStr, err := peerId.ToString()
		if err != nil {
			return err
		}
		if peerVal, ok := node.Services().P2P().Peers().Get(peerIdStr); ok {
			foundPeer := peerVal.(net.Peer)
			err := foundPeer.SendMessage(msg)
			if err != nil {
				node.Logger().Error("Failed to send message", zap.Error(err))
				continue
			}
		}

		node.HashQueryRoutingTable().Remove(hash.HashCode())
	}

	return nil
}
