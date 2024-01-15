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

type HTTPImpl struct {
	node    interfaces.Node
	handler interfaces.HTTPHandler
}

func NewHTTP(node interfaces.Node, handler interfaces.HTTPHandler) interfaces.HTTPService {
	return &HTTPImpl{
		node:    node,
		handler: handler,
	}
}

func (h *HTTPImpl) GetHandler() *httprouter.Router {
	mux := jape.Mux(map[string]jape.Handler{
		"GET /s5/version": h.versionHandler,
		"GET /s5/p2p":     h.p2pHandler,
		"POST /s5/upload": h.uploadHandler,
	})

	return mux
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

func (h *HTTPImpl) uploadHandler(context jape.Context) {
	h.handler.SmallFileUpload(&context)
}
