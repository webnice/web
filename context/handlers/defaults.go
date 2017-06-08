package handlers

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import "gopkg.in/webnice/web.v1/mime"
import "gopkg.in/webnice/web.v1/header"
import (
	"net/http"
)

// Default handler for internal server errors
func (hndl *impl) defaultInternalServerError(wr http.ResponseWriter, rq *http.Request) {
	var err error

	wr.Header().Set(header.ContentType, mime.TextPlainCharsetUTF8)
	wr.WriteHeader(status.InternalServerError)
	if err = hndl.errors.InternalServerError(nil); err != nil {
		_, _ = wr.Write([]byte(err.Error()))
		return
	}
	_, _ = wr.Write(status.Bytes(status.InternalServerError))
}

// Default handler for method not allowed
func (hndl *impl) defaultMethodNotAllowed(wr http.ResponseWriter, rq *http.Request) {
	var err error

	wr.Header().Set(header.ContentType, mime.TextPlainCharsetUTF8)
	wr.WriteHeader(status.MethodNotAllowed)
	if err = hndl.errors.MethodNotAllowed(nil); err != nil {
		_, _ = wr.Write([]byte(err.Error()))
		return
	}
	_, _ = wr.Write(status.Bytes(status.MethodNotAllowed))
}

// Default handler for not found
func (hndl *impl) defaultNotFound(wr http.ResponseWriter, rq *http.Request) {
	var err error

	wr.Header().Set(header.ContentType, mime.TextPlainCharsetUTF8)
	wr.WriteHeader(status.NotFound)
	if err = hndl.errors.NotFound(nil); err != nil {
		_, _ = wr.Write([]byte(err.Error()))
		return
	}
	_, _ = wr.Write(status.Bytes(status.NotFound))
}
