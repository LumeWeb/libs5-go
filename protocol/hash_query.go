package protocol

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/vmihailenco/msgpack/v5"
	"go.uber.org/zap"
	"log"
)

var _ base.IncomingMessageTyped = (*HashQuery)(nil)

type HashQuery struct {
	hash  *encoding.Multihash
	kinds []int

	base.IncomingMessageTypedImpl
	base.IncomingMessageHandler
}

func NewHashQuery() *HashQuery {
	return &HashQuery{}
}

func (h HashQuery) Hash() *encoding.Multihash {
	return h.hash
}

func (h HashQuery) Kinds() []int {
	return h.kinds
}

func (h *HashQuery) DecodeMessage(dec *msgpack.Decoder) error {
	hash, err := dec.DecodeBytes()

	if err != nil {
		return err
	}

	h.hash = encoding.NewMultihash(hash)

	var kinds []int
	err = dec.Decode(&kinds)
	if err != nil {
		return err
	}

	h.kinds = kinds

	return nil
}
func (h *HashQuery) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	mapLocations, err := node.GetCachedStorageLocations(h.hash, h.kinds)
	if err != nil {
		log.Printf("Error getting cached storage locations: %v", err)
		return err
	}

	if len(mapLocations) > 0 {
		availableNodes := make([]*encoding.NodeId, 0, len(mapLocations))
		for key := range mapLocations {
			nodeId, err := encoding.DecodeNodeId(key)
			if err != nil {
				node.Logger().Error("Error decoding node id", zap.Error(err))
				continue
			}

			availableNodes = append(availableNodes, nodeId)
		}

		score, err := node.Services().P2P().SortNodesByScore(availableNodes)
		if err != nil {
			return err
		}

		sortedNodeId, err := (*score[0]).ToString()
		if err != nil {
			return err
		}

		entry, exists := mapLocations[sortedNodeId]
		if exists {
			err := peer.SendMessage(entry.ProviderMessage())
			if err != nil {
				return err
			}
		}
	}

	var peers *hashset.Set
	hashString, err := h.hash.ToString()
	if err != nil {
		return err
	}
	peersVal, ok := node.HashQueryRoutingTable().Get(hashString) // Implement HashQueryRoutingTable method
	if ok {
		peers = peersVal.(*hashset.Set)
		if !peers.Contains(peer.Id()) {
			peers.Add(peer.Id())
		}

		return nil
	}

	peerList := hashset.New()
	peerList.Add(peer.Id())

	node.HashQueryRoutingTable().Put(hashString, peerList)

	for _, val := range node.Services().P2P().Peers().Values() {
		peerVal := val.(net.Peer)
		if !peerVal.Id().Equals(peer.Id()) {
			err := peerVal.SendMessage(h.IncomingMessageImpl.Original())
			if err != nil {
				node.Logger().Error("Failed to send message", zap.Error(err))
			}
		}
	}

	return nil
}
