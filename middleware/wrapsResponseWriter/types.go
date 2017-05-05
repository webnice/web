package wrapsResponseWriter // import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"io"
	"net/http"
)

var (
	_ http.Pusher        = &http2FancyWriter{}
	_ http.Flusher       = &flush{}
	_ http.CloseNotifier = &httpFancyWriter{}
	_ http.Flusher       = &httpFancyWriter{}
	_ http.Hijacker      = &httpFancyWriter{}
	_ io.ReaderFrom      = &httpFancyWriter{}
	_ http.CloseNotifier = &http2FancyWriter{}
	_ http.Flusher       = &http2FancyWriter{}
)

type WrapsResponseWriter interface {
	http.ResponseWriter

	// Status The HTTP status of the request
	Status() int

	// Len The total number of bytes sent to the client
	Len() uint64

	// Tee causes the response body to be written to the given io.Writer in addition to proxying the writes through
	Tee(io.Writer)

	// Unwrap Returns the original proxied target
	Unwrap() http.ResponseWriter
}

// Wraps a http.ResponseWriter that implements the minimal http.ResponseWriter interface
type basic struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       uint64
	tee         io.Writer
}

type flush struct {
	basic
}

// Is a HTTP writer that additionally satisfies http.CloseNotifier, http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to make the proxied object support the full method set of the proxied object
type httpFancyWriter struct {
	basic
}

// Is a HTTP2 writer that additionally satisfies http.CloseNotifier, http.Flusher, and io.ReaderFrom. It exists for the common case of wrapping
// the http.ResponseWriter that package http gives you, in order to make the proxied object support the full method set of the proxied object
type http2FancyWriter struct {
	basic
}
