package handlers // import "gopkg.in/webnice/web.v1/context/handlers"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/errors"
import (
	"net/http"
)

// New returns new context object
func New(errors errors.Interface) Interface {
	var hndl = new(impl)
	hndl.handlers = make(handlers)
	hndl.errors = errors
	return hndl
}

// Add adds the key, value pair to the errors.
// It appends to any existing values associated with key
func (h handlers) Add(key uint32, fn http.HandlerFunc) {
	h[key] = append(h[key], fn)
}

// Set sets the params entries associated with key to
// the single element value. It replaces any existing
// values associated with key
func (h handlers) Set(key uint32, fn http.HandlerFunc) {
	h[key] = []http.HandlerFunc{fn}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns nil.
// To access multiple values of a key, or to use access the map directly
func (h handlers) Get(key uint32) http.HandlerFunc {
	if h == nil {
		return nil
	}
	v := h[key]
	if len(v) == 0 {
		return nil
	}
	return v[0]
}

// Del deletes the values associated with key
func (h handlers) Del(key uint32) http.HandlerFunc {
	var val = h.Get(key)
	delete(h, key)
	return val
}

// Reset all stored handlers
func (hndl *impl) Reset() { hndl.handlers = make(handlers) }

// Set and get InternalServerError handler function
func (hndl *impl) InternalServerError(fn http.HandlerFunc) (ret http.HandlerFunc) {
	if fn != nil {
		hndl.handlers.Set(_InternalServerError, fn)
	}
	if ret = hndl.handlers.Get(_InternalServerError); ret != nil {
		return
	}
	ret = hndl.defaultInternalServerError
	return
}
