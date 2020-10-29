package handlers

import (
	"net/http"

	"github.com/webnice/web/v2/ambry"
	"github.com/webnice/web/v2/context/errors"
)

// New returns new context object
func New(errors errors.Interface) Interface {
	var hndl = new(impl)

	hndl.handlers = ambry.New()
	hndl.errors = errors

	return hndl
}

func (hndl *impl) do(key uint32, fn http.HandlerFunc, defaultFn http.HandlerFunc) (ret http.HandlerFunc) {
	var tmp interface{}

	if fn != nil {
		hndl.handlers.Set(key, fn)
	}
	if tmp = hndl.handlers.Get(key); tmp != nil {
		ret = tmp.(http.HandlerFunc)
		return
	}
	ret = defaultFn

	return
}

// Reset all stored handlers
func (hndl *impl) Reset() { hndl.handlers = ambry.New() }

// Set and get InternalServerError handler function
func (hndl *impl) InternalServerError(fn http.HandlerFunc) (ret http.HandlerFunc) {
	return hndl.do(keyInternalServerError, fn, hndl.defaultInternalServerError)
}

// MethodNotAllowed	Set and get MethodNotAllowed handler function
func (hndl *impl) MethodNotAllowed(fn http.HandlerFunc) (ret http.HandlerFunc) {
	return hndl.do(keyMethodNotAllowed, fn, hndl.defaultMethodNotAllowed)
}

// NotFound	Set and get MethodNotFound handler function
func (hndl *impl) NotFound(fn http.HandlerFunc) (ret http.HandlerFunc) {
	return hndl.do(keyNotFound, fn, hndl.defaultNotFound)
}
