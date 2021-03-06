package nocache

import (
	"io"
	"net/http"
	"time"

	"github.com/webnice/web/v2/header"
	"github.com/webnice/web/v2/status"
)

type impl struct {
	http.ResponseWriter
	isHeaderWritten bool
	Writer          io.Writer
}

func cleanHeaders(wr http.ResponseWriter) {
	var (
		key     string
		headers = []string{
			header.ETag,
			header.IfModifiedSince,
			header.IfMatch,
			header.IfNoneMatch,
			header.IfRange,
			header.IfUnmodifiedSince,
		}
	)

	for _, key = range headers {
		if wr.Header().Get(key) != "" {
			wr.Header().Del(key)
		}
	}
}

func setHeaders(wr http.ResponseWriter) {
	const (
		xAccelExpires    = `X-Accel-Expires`
		keyCacheControl  = `no-cache, private, max-age=0`
		keyPragmaNoCache = `no-cache`
		keyExpires       = `0`
	)
	var (
		key, value string
		headers    = map[string]string{
			header.Expires:      time.Unix(0, 0).Format(time.RFC1123),
			header.CacheControl: keyCacheControl,
			header.Pragma:       keyPragmaNoCache,
			xAccelExpires:       keyExpires,
		}
	)

	for key, value = range headers {
		wr.Header().Set(key, value)
	}
}

// Handler Middleware set headers to disable cache
func Handler(hndl http.Handler) http.Handler {
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
