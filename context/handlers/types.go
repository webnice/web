package handlers // import "gopkg.in/webnice/web.v1/context/handlers"

// import "gopkg.in/webnice/debug.v1"
// import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

const (
	_InternalServerError uint16 = iota
)

type handlers map[uint16][]http.HandlerFunc

// This is an inplementation
type impl struct {
	handlers handlers
}

// Interface is an interface of package
type Interface interface {
	// Set and get InternalServerError handler function
	InternalServerError(fn http.HandlerFunc) http.HandlerFunc
}
