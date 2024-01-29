package service

import (
	"git.lumeweb.com/LumeWeb/libs5-go/build"
	"git.lumeweb.com/LumeWeb/libs5-go/net"
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
	"go.uber.org/zap"
	"net/url"
	"nhooyr.io/websocket"
)

var _ Service = (*HTTPService)(nil)

type P2PNodesResponse struct {
	Nodes []P2PNodeResponse `json:"nodes"`
}

type P2PNodeResponse struct {
	Id   string   `json:"id"`
	Uris []string `json:"uris"`
}

type HTTPService struct {
	ServiceBase
}

func NewHTTP(params ServiceParams) *HTTPService {
	return &HTTPService{
		ServiceBase: ServiceBase{
			logger: params.Logger,
			config: params.Config,
			db:     params.Db,
		},
	}
}

func (h *HTTPService) GetHttpRouter(inject map[string]jape.Handler) *httprouter.Router {
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

func (h *HTTPService) Start() error {
	return nil
}

func (h *HTTPService) Stop() error {
	return nil
}

func (h *HTTPService) Init() error {
	return nil
}

func (h *HTTPService) versionHandler(ctx jape.Context) {
	_, _ = ctx.ResponseWriter.Write([]byte(build.Version))
}
func (h *HTTPService) p2pHandler(ctx jape.Context) {
	c, err := websocket.Accept(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		h.logger.Error("error accepting websocket connection", zap.Error(err))
		return
	}

	peer, err := net.CreateTransportPeer("wss", &net.TransportPeerConfig{
		Socket: c,
		Uris:   []*url.URL{},
	})

	if err != nil {
		h.logger.Error("error creating transport peer", zap.Error(err))
		err := c.Close(websocket.StatusInternalError, "the sky is falling")
		if err != nil {
			h.logger.Error("error closing websocket connection", zap.Error(err))
		}
		return
	}

	h.services.P2P().ConnectionTracker().Add(1)

	go func() {
		err := h.services.P2P().OnNewPeer(peer, false)
		if err != nil {
			h.logger.Error("error handling new peer", zap.Error(err))
		}
		h.services.P2P().ConnectionTracker().Done()
	}()
}

func (h *HTTPService) p2pNodesHandler(ctx jape.Context) {
	localId, err := h.services.P2P().NodeId().ToString()

	if ctx.Check("error getting local node id", err) != nil {
		return
	}

	uris := h.services.P2P().SelfConnectionUris()

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
