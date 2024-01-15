package interfaces

import (
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
)

//go:generate mockgen -source=http.go -destination=../mocks/interfaces/http.go -package=interfaces

type HTTPService interface {
	Service
	GetHandler() *httprouter.Router
}

type HTTPHandler interface {
	SmallFileUpload(context *jape.Context)
}
