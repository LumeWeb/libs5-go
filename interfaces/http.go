package interfaces

import (
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
)

//go:generate mockgen -source=http.go -destination=../mocks/interfaces/http.go -package=interfaces

type HTTPService interface {
	Service
	GetHttpRouter() *httprouter.Router
	SetHttpHandler(handler HTTPHandler)
}

type HTTPHandler interface {
	SmallFileUpload(context *jape.Context)
}
