package nocache

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import (
	"net/http"
	"time"
)

func cleanHeaders(rq *http.Request) {
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
		if rq.Header.Get(key) != "" {
			rq.Header.Del(key)
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

func NoCache(hndl http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		cleanHeaders(rq)
		setHeaders(wr)
		hndl.ServeHTTP(wr, rq)
	}
	return http.HandlerFunc(fn)
}
