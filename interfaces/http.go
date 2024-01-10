package interfaces

import (
	"github.com/julienschmidt/httprouter"
)

//go:generate mockgen -source=http.go -destination=../mocks/interfaces/http.go -package=interfaces

type HTTPService interface {
	Service
	GetHandler() *httprouter.Router
}
