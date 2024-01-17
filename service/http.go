package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/build"
	"git.lumeweb.com/LumeWeb/libs5-go/interfaces"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
	"go.uber.org/zap"
	"net/url"
	"nhooyr.io/websocket"
)

var _ interfaces.HTTPService = (*HTTPImpl)(nil)

type P2PNodesResponse struct {
	Nodes []P2PNodeResponse `json:"nodes"`
}

type P2PNodeResponse struct {
	Id   string   `json:"id"`
	Uris []string `json:"uris"`
}

type HTTPImpl struct {
	node interfaces.Node
}

func NewHTTP(node interfaces.Node) interfaces.HTTPService {
	return &HTTPImpl{
		node: node,
	}
}

func (h *HTTPImpl) GetHttpRouter(inject map[string]jape.Handler) *httprouter.Router {
	routes := map[string]jape.Handler{
		"GET /s5/version":   h.versionHandler,
		"GET /s5/p2p":       h.p2pHandler,
		"GET /s5/p2p/nodes": h.p2pNodesHandler,
	}

	for k, v := range inject {
		routes[k] = v
	}

	return jape.Mux(routes)
}

func (h *HTTPImpl) Node() interfaces.Node {
	return h.node
}

func (h *HTTPImpl) Start() error {
	return nil
}

func (h *HTTPImpl) Stop() error {
	return nil
}

func (h *HTTPImpl) Init() error {
	return nil
}

func (h *HTTPImpl) versionHandler(ctx jape.Context) {
	_, _ = ctx.ResponseWriter.Write([]byte(build.Version))
}
func (h *HTTPImpl) p2pHandler(ctx jape.Context) {
	c, err := websocket.Accept(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		h.node.Logger().Error("error accepting websocket connection", zap.Error(err))
		return
	}

	peer, err := net.CreateTransportPeer("wss", &net.TransportPeerConfig{
		Socket: c,
		Uris:   []*url.URL{},
	})

	if err != nil {
		h.node.Logger().Error("error creating transport peer", zap.Error(err))
		err := c.Close(websocket.StatusInternalError, "the sky is falling")
		if err != nil {
			h.node.Logger().Error("error closing websocket connection", zap.Error(err))
		}
		return
	}

	h.Node().ConnectionTracker().Add(1)

	go func() {
		err := h.node.Services().P2P().OnNewPeer(peer, false)
		if err != nil {
			h.node.Logger().Error("error handling new peer", zap.Error(err))
		}
		h.node.ConnectionTracker().Done()
	}()
}

func (h *HTTPImpl) p2pNodesHandler(ctx jape.Context) {
	localId, err := h.node.Services().P2P().NodeId().ToString()

	if ctx.Check("error getting local node id", err) != nil {
		return
	}

	uris := h.node.Services().P2P().SelfConnectionUris()

	nodeList := make([]P2PNodeResponse, len(uris))

	for i, uri := range uris {
		nodeList[i] = P2PNodeResponse{
			Id:   localId,
			Uris: []string{uri.String()},
		}
	}

	ctx.Encode(P2PNodesResponse{
		Nodes: nodeList,
	})
}
