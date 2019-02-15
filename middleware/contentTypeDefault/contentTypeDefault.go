package contentTypeDefault

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/status"
)

// Interface is an interface of package
type Interface interface {
	// Handler Middleware set default content-type header
	Handler(hndl http.Handler) http.Handler
}

type impl struct {
	Default string
}

type rewr struct {
	http.ResponseWriter
	isHeaderWritten bool
	Writer          io.Writer
	Default         string
}

// New Create object of package and return interface
// First argument is a default content-type
func New(def string) Interface {
	var dct = &impl{
		Default: def,
	}
	return dct
}

// Handler Middleware set default content-type header
func (dct *impl) Handler(hndl http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		var rwr = &rewr{
			ResponseWriter: wr,
			Writer:         wr,
			Default:        dct.Default,
		}
		hndl.ServeHTTP(rwr, rq)
	}
	return http.HandlerFunc(fn)
}

// WriteHeader Write header code
func (rwr *rewr) WriteHeader(code int) {
	if rwr.isHeaderWritten {
		return
	}
	if rwr.Header().Get(header.ContentType) == "" {
		rwr.Header().Add(header.ContentType, rwr.Default)
	}
	rwr.isHeaderWritten = true
}

// Write Implementation of an interface io.Writer
func (rwr *rewr) Write(p []byte) (n int, err error) {
	if !rwr.isHeaderWritten {
		rwr.WriteHeader(status.Ok)
	}
	n, err = rwr.Writer.Write(p)
	return
}
