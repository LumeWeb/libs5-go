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
	raw       []byte
	hash      *encoding.Multihash
	kind      int
	expiry    int64
	parts     []string
	publicKey []byte
	signature []byte

	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
}

func (s *StorageLocation) DecodeMessage(dec *msgpack.Decoder) error {
	data, err := dec.DecodeRaw()

	if err != nil {
		return err
	}

	s.raw = data

	return nil
}
func (s *StorageLocation) HandleMessage(node interfaces.Node, peer *net.Peer, verifyId bool) error {
	hash := encoding.NewMultihash(s.raw[1:34]) // Replace NewMultihash with appropriate function
	fmt.Println("Hash:", hash)

	typeOfData := s.raw[34]

	expiry := utils.DecodeEndian(s.raw[35:39])

	partCount := s.raw[39]

	parts := []string{}
	cursor := 40
	for i := 0; i < int(partCount); i++ {
		length := utils.DecodeEndian(s.raw[cursor : cursor+2])
		cursor += 2
		part := string(s.raw[cursor : cursor+int(length)])
		parts = append(parts, part)
		cursor += int(length)
	}

	publicKey := s.raw[cursor : cursor+33]
	signature := s.raw[cursor+33:]

	if types.HashType(publicKey[0]) != types.HashTypeEd25519 { // Replace CID_HASH_TYPES_ED25519 with actual constant
		return fmt.Errorf("Unsupported public key type %d", publicKey[0])
	}

	if !ed25519.Verify(publicKey[1:], s.raw[:cursor], signature) {
		return fmt.Errorf("Signature verification failed")
	}

	nodeId := encoding.NewNodeId(publicKey)

	// Assuming `node` is an instance of your NodeImpl structure
	err := node.AddStorageLocation(hash, nodeId, storage.NewStorageLocation(int(typeOfData), parts, int64(expiry)), s.raw, node.Config()) // Implement AddStorageLocation

	if err != nil {
		return fmt.Errorf("Failed to add storage location: %s", err)
	}

	var list *hashset.Set
	listVal, ok := node.HashQueryRoutingTable().Get(hash.HashCode()) // Implement HashQueryRoutingTable method
	if !ok {
		list = hashset.New()
	}

	list = listVal.(*hashset.Set)

	for _, peerIdVal := range list.Values() {
		peerId := peerIdVal.(*encoding.NodeId)

		if peerId.Equals(nodeId) || peerId.Equals(peer) {
			continue
		}
		if peerVal, ok := node.Services().P2P().Peers().Get(peerId.HashCode()); ok {
			foundPeer := peerVal.(net.Peer)
			err := foundPeer.SendMessage(s.raw)
			if err != nil {
				node.Logger().Error("Failed to send message", zap.Error(err))
				continue
			}
		}

		node.HashQueryRoutingTable().Remove(hash.HashCode())
	}

	return nil
}
