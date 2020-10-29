package handlers

import (
	"net/http"

	"github.com/webnice/web/v2/ambry"
	"github.com/webnice/web/v2/context/errors"
)

const (
	keyInternalServerError uint32 = iota
	keyMethodNotAllowed
	keyNotFound
)

// Interface is an interface of package
type Interface interface {
	// InternalServerError Set and get InternalServerError handler function
	InternalServerError(fn http.HandlerFunc) http.HandlerFunc

	// MethodNotAllowed	Set and get MethodNotAllowed handler function
	MethodNotAllowed(fn http.HandlerFunc) http.HandlerFunc

	// NotFound	Set and get NotFound handler function
	NotFound(fn http.HandlerFunc) http.HandlerFunc

	// Reset all stored handlers
	Reset()
}

type impl struct {
	handlers ambry.Interface
	errors   errors.Interface
}
