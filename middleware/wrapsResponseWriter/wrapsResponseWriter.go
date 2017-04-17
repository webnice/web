package wrapsResponseWriter //import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"
)

func New(wr http.ResponseWriter, protoMajor int) WrapsResponseWriter {
	var ba = basic{ResponseWriter: wr}
	_, cn := wr.(http.CloseNotifier)
	_, fl := wr.(http.Flusher)
	if protoMajor == 2 {
		_, ps := wr.(http.Pusher)
		if cn && fl && ps {
			return &http2FancyWriter{ba}
		}
	} else {
		_, hj := wr.(http.Hijacker)
		_, rf := wr.(io.ReaderFrom)
		if cn && fl && hj && rf {
			return &httpFancyWriter{ba}
		}
	}
	if fl {
		return &flush{ba}
	}
	return &ba
}
