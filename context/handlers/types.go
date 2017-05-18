package handlers // import "gopkg.in/webnice/web.v1/context/handlers"

// import "gopkg.in/webnice/debug.v1"
// import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/errors"
import (
	"net/http"
)

const (
	_InternalServerError uint32 = iota
)

type handlers map[uint32][]http.HandlerFunc

// This is an inplementation
type impl struct {
	handlers handlers
	errors   errors.Interface
}

// Interface is an interface of package
type Interface interface {
	// Reset all stored handlers
	Reset()

	// Set and get InternalServerError handler function
	InternalServerError(fn http.HandlerFunc) http.HandlerFunc
}
