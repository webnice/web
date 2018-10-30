package nocache

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"
	"time"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/status"
)

type impl struct {
	http.ResponseWriter
	isHeaderWritten bool
	Writer          io.Writer
}

func cleanHeaders(wr http.ResponseWriter) {
	var key string
	var headers = []string{
		header.ETag,
		header.IfModifiedSince,
		header.IfMatch,
		header.IfNoneMatch,
		header.IfRange,
		header.IfUnmodifiedSince,
	}
	for _, key = range headers {
		if wr.Header().Get(key) != "" {
			wr.Header().Del(key)
		}
	}
}

func setHeaders(wr http.ResponseWriter) {
	var key, value string
	var headers = map[string]string{
		"Expires":         time.Unix(0, 0).Format(time.RFC1123),
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}
	for key, value = range headers {
		wr.Header().Set(key, value)
	}
}

// NoCache Middleware set headers to disable cache
func NoCache(hndl http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		var nch = &impl{
			ResponseWriter: wr,
			Writer:         wr, // by default
		}
		setHeaders(wr)
		hndl.ServeHTTP(nch, rq)
	}
	return http.HandlerFunc(fn)
}

// WriteHeader Write header code
func (nch *impl) WriteHeader(code int) {
	if nch.isHeaderWritten {
		return
	}
	cleanHeaders(nch.ResponseWriter)
	nch.isHeaderWritten = true
}

// Write Implementation of an interface io.Writer
func (nch *impl) Write(p []byte) (n int, err error) {
	if !nch.isHeaderWritten {
		nch.WriteHeader(status.Ok)
	}
	n, err = nch.Writer.Write(p)
	return
}
