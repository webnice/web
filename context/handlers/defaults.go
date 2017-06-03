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
	_, _ = wr.Write(status.Bytes(status.InternalServerError))
	if err = hndl.errors.InternalServerError(nil); err != nil {
		_, _ = wr.Write([]byte("\n"))
		_, _ = wr.Write([]byte(err.Error()))
	}
}
