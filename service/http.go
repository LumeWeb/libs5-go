package service

import (
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
)

type HTTPService interface {
	GetHttpRouter(inject map[string]jape.Handler) *httprouter.Router
	Service
}
