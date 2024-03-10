package _default

import (
	"context"
	"git.lumeweb.com/LumeWeb/libs5-go/build"
	s5net "git.lumeweb.com/LumeWeb/libs5-go/net"
	"git.lumeweb.com/LumeWeb/libs5-go/service"
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
	"go.uber.org/zap"
	"net"
	"net/url"
	"nhooyr.io/websocket"
	"strings"
)

var _ service.Service = (*HTTPServiceDefault)(nil)

type P2PNodesResponse struct {
	Nodes []P2PNodeResponse `json:"nodes"`
}

type P2PNodeResponse struct {
	Id   string   `json:"id"`
	Uris []string `json:"uris"`
}

type HTTPServiceDefault struct {
	service.ServiceBase
}

func NewHTTP(params service.ServiceParams) *HTTPServiceDefault {
	return &HTTPServiceDefault{
		ServiceBase: service.NewServiceBase(params.Logger, params.Config, params.Db),
	}
}

func (h *HTTPServiceDefault) GetHttpRouter(inject map[string]jape.Handler) *httprouter.Router {
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

func (h *HTTPServiceDefault) Start(ctx context.Context) error {
	return nil
}

func (h *HTTPServiceDefault) Stop(ctx context.Context) error {
	return nil
}

func (h *HTTPServiceDefault) Init(ctx context.Context) error {
	return nil
}

func (h *HTTPServiceDefault) versionHandler(ctx jape.Context) {
	_, _ = ctx.ResponseWriter.Write([]byte(build.Version))
}
func (h *HTTPServiceDefault) p2pHandler(ctx jape.Context) {
	c, err := websocket.Accept(ctx.ResponseWriter, ctx.Request, nil)
	if err != nil {
		h.Logger().Error("error accepting websocket connection", zap.Error(err))
		return
	}

	peer, err := s5net.CreateTransportPeer("wss", &s5net.TransportPeerConfig{
		Socket: c,
		Uris:   []*url.URL{},
	})

	if err != nil {
		h.Logger().Error("error creating transport peer", zap.Error(err))
		err := c.Close(websocket.StatusInternalError, "the sky is falling")
		if err != nil {
			h.Logger().Error("error closing websocket connection", zap.Error(err))
		}
		return
	}

	ip := peer.GetIP()

	// Check for reverse proxy headers
	realIP := ctx.Request.Header.Get("X-Real-IP")
	forwardedFor := ctx.Request.Header.Get("X-Forwarded-For")

	var clientIP net.IP
	if realIP != "" {
		clientIP = net.ParseIP(realIP)
	} else if forwardedFor != "" {
		// X-Forwarded-For can contain multiple IP addresses separated by commas
		// We take the first IP in the list as the client's IP
		parts := strings.Split(forwardedFor, ",")
		clientIP = net.ParseIP(parts[0])
	}

	blockConnection := func(ip net.Addr) bool {
		// If we have a valid client IP from headers, use that for the loopback check
		if clientIP != nil {
			return clientIP.IsLoopback()
		}
		// Otherwise, fall back to the peer's IP
		switch v := ip.(type) {
		case *net.IPNet:
			return v.IP.IsLoopback()
		case *net.TCPAddr:
			return v.IP.IsLoopback()
		default:
			return false
		}
	}

	if blockConnection(ip) {
		err := peer.End()
		if err != nil {
			h.Logger().Error("error ending peer", zap.Error(err))
		}
		return
	}

	h.Services().P2P().ConnectionTracker().Add(1)

	go func() {
		err := h.Services().P2P().OnNewPeer(peer, false)
		if err != nil {
			h.Logger().Error("error handling new peer", zap.Error(err))
		}
		h.Services().P2P().ConnectionTracker().Done()
	}()
}

func (h *HTTPServiceDefault) p2pNodesHandler(ctx jape.Context) {
	localId, err := h.Services().P2P().NodeId().ToString()

	if ctx.Check("error getting local node id", err) != nil {
		return
	}

	uris := h.Services().P2P().SelfConnectionUris()

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
