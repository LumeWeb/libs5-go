package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"github.com/vmihailenco/msgpack/v5"
)

type NodeVotesImpl struct {
	good int
	bad  int
}

func NewNodeVotes() interfaces.NodeVotes {
	return &NodeVotesImpl{
		good: 0,
		bad:  0,
	}
}

func (n *NodeVotesImpl) Good() int {
	return n.good
}

func (n *NodeVotesImpl) Bad() int {
	return n.bad
}

func (n NodeVotesImpl) EncodeMsgpack(enc *msgpack.Encoder) error {
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

func (n *NodeVotesImpl) DecodeMsgpack(dec *msgpack.Decoder) error {
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

func (n *NodeVotesImpl) Upvote() {
	n.good++
}

func (n *NodeVotesImpl) Downvote() {
	n.bad++
}
