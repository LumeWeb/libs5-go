package interfaces

//go:generate mockgen -source=meta.go -destination=../mocks/interfaces/meta.go -package=interfaces

type Metadata interface {
	ToJson() map[string]interface{}
}
