package handlers

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/errors"
import (
	"net/http"
)

// New returns new context object
func New(errors errors.Interface) Interface {
	var hndl = new(impl)
	hndl.handlers = new(handlers)
	hndl.errors = errors
	return hndl
}

// Add key value param to heap
func (itm *handlers) Add(key uint32, value http.HandlerFunc) {
	*itm = append(*itm, item{key, value})
}

// Get value by key from heap
func (itm handlers) Get(key uint32) (ret http.HandlerFunc) {
	var p item
	for _, p = range itm {
		if p.Key == key {
			ret = p.Value
			break
		}
	}
	return
}

// Set value by key to heap
func (itm *handlers) Set(key uint32, value http.HandlerFunc) {
	var p item
	var i, idx int
	idx = -1
	for i, p = range *itm {
		if p.Key == key {
			idx = i
			break
		}
	}
	if idx < 0 {
		(*itm).Add(key, value)
	} else {
		(*itm)[idx] = item{key, value}
	}
}

// Del Delete value by key from heap
func (itm *handlers) Del(key uint32) (ret http.HandlerFunc) {
	var p item
	var i int
	for i, p = range *itm {
		if p.Key == key {
			*itm = append((*itm)[:i], (*itm)[i+1:]...)
			ret = p.Value
		}
	}
	return
}

// Reset all stored handlers
func (hndl *impl) Reset() { hndl.handlers = new(handlers) }

// Set and get InternalServerError handler function
func (hndl *impl) InternalServerError(fn http.HandlerFunc) (ret http.HandlerFunc) {
	if fn != nil {
		hndl.handlers.Set(keyInternalServerError, fn)
	}
	if ret = hndl.handlers.Get(keyInternalServerError); ret != nil {
		return
	}
	ret = hndl.defaultInternalServerError
	return
}

// MethodNotAllowed	Set and get MethodNotAllowed handler function
func (hndl *impl) MethodNotAllowed(fn http.HandlerFunc) (ret http.HandlerFunc) {
	if fn != nil {
		hndl.handlers.Set(keyMethodNotAllowed, fn)
	}
	if ret = hndl.handlers.Get(keyMethodNotAllowed); ret != nil {
		return
	}
	ret = hndl.defaultMethodNotAllowed
	return
}

// NotFound	Set and get MethodNotFound handler function
func (hndl *impl) NotFound(fn http.HandlerFunc) (ret http.HandlerFunc) {
	if fn != nil {
		hndl.handlers.Set(keyNotFound, fn)
	}
	if ret = hndl.handlers.Get(keyNotFound); ret != nil {
		return
	}
	ret = hndl.defaultNotFound
	return
}
