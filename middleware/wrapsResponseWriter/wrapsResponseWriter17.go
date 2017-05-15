// +build go1.7,!go1.8

package wrapsResponseWriter // import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"
)

// New Proxy around an http.ResponseWriter that allows you to hook into various parts of the response process
func New(wr http.ResponseWriter, protoMajor int) WrapsResponseWriter {
	var cn, fl, hj, rf bool
	var ba = basic{
		ResponseWriter: wr,
	}

	_, cn = wr.(http.CloseNotifier)
	_, fl = wr.(http.Flusher)
	switch protoMajor {
	case 2:
		if cn && fl {
			return &http2FancyWriter{ba}
		}
	default:
		_, hj = wr.(http.Hijacker)
		_, rf = wr.(io.ReaderFrom)
		if cn && fl && hj && rf {
			return &httpFancyWriter{ba}
		}
	}
	if fl {
		return &flush{ba}
	}

	return &ba
}
