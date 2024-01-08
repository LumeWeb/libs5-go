package signed

import (
	"git.lumeweb.com/LumeWeb/libs5-go/encoding"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/protocol/base"
	"git.lumeweb.com/LumeWeb/libs5-go/utils"
	"github.com/vmihailenco/msgpack/v5"
	"net/url"
)

var (
	_ base.IncomingMessageTyped = (*AnnouncePeers)(nil)
)

type AnnouncePeers struct {
	connected      bool
	peer           *encoding.NodeId
	connectionUris []*url.URL
	base.IncomingMessageTypedImpl
}

func NewAnnouncePeers() *AnnouncePeers {
	return &AnnouncePeers{connected: false, peer: nil, connectionUris: nil}
}

func (a *AnnouncePeers) DecodeMessage(dec *msgpack.Decoder) error {
	peerId, err := dec.DecodeBytes()

	if err != nil {
		return err
	}

	a.peer = encoding.NewNodeId(peerId)

	connected, err := dec.DecodeBool()

	if err != nil {
		return err
	}

	a.connected = connected
	connectionUriVal, err := utils.DecodeMsgpackArray(dec)

	if err != nil {
		return err
	}

	a.connectionUris = make([]*url.URL, 0, len(connectionUriVal))
	connectionUris := interface{}(connectionUriVal).([]string)

	for _, connectionUri := range connectionUris {
		uri, err := url.Parse(connectionUri)
		if err != nil {
			return err
		}
		a.connectionUris = append(a.connectionUris, uri)
	}

	return nil
}

func (a AnnouncePeers) HandleMessage(node interfaces.Node, peer net.Peer, verifyId bool) error {
	if len(a.connectionUris) > 0 {
		firstUrl := a.connectionUris[0]
		uri := new(url.URL)
		*uri = *firstUrl

		if firstUrl.User != nil {
			passwd, empty := firstUrl.User.Password()
			if empty {
				passwd = ""
			}

			nodeId, err := a.peer.ToString()
			if err != nil {
				return err
			}

			uri.User = url.UserPassword(nodeId, passwd)
		}
	}

	return nil
}
