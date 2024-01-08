package base

//go:generate mockgen -source=signed.go -destination=../../mocks/base/signed.go -package=base -aux_files=git.lumeweb.com/LumeWeb/libs5-go/protocol/base=base.go

type SignedIncomingMessage interface {
	IncomingMessage
}
