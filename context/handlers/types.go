package handlers

// import "gopkg.in/webnice/debug.v1"
// import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/errors"
import (
	"net/http"
)

const (
	keyInternalServerError uint32 = iota
	keyMethodNotAllowed
	keyNotFound
)

// internal storage structure
type item struct {
	Key   uint32
	Value http.HandlerFunc
}

type handlers []item

// This is an inplementation
type impl struct {
	handlers *handlers
	errors   errors.Interface
}

// Interface is an interface of package
type Interface interface {
	// Reset all stored handlers
	Reset()

	// InternalServerError Set and get InternalServerError handler function
	InternalServerError(fn http.HandlerFunc) http.HandlerFunc

	// MethodNotAllowed	Set and get MethodNotAllowed handler function
	MethodNotAllowed(fn http.HandlerFunc) http.HandlerFunc

	// NotFound	Set and get NotFound handler function
	NotFound(fn http.HandlerFunc) http.HandlerFunc
}
