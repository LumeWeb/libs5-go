package interfaces

import (
	"github.com/julienschmidt/httprouter"
	"go.sia.tech/jape"
)

//go:generate mockgen -source=http.go -destination=../mocks/interfaces/http.go -package=interfaces

type HTTPService interface {
	Service
	GetHttpRouter(inject map[string]jape.Handler) *httprouter.Router
}
