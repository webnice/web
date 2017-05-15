package wrapsResponseWriter // import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"io"
	"net"
	"net/http"
)

func (f *httpFancyWriter) CloseNotify() <-chan bool {
	return f.basic.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (f *httpFancyWriter) Flush() { f.basic.ResponseWriter.(http.Flusher).Flush() }

func (f *httpFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return f.basic.ResponseWriter.(http.Hijacker).Hijack()
}

func (f *httpFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	var rf io.ReaderFrom
	if f.basic.tee != nil {
		return io.Copy(&f.basic, r)
	}
	rf = f.basic.ResponseWriter.(io.ReaderFrom)
	f.basic.maybeWriteHeader()
	return rf.ReadFrom(r)
}

// HTTP2

func (f *http2FancyWriter) CloseNotify() <-chan bool {
	return f.basic.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (f *http2FancyWriter) Flush() { f.basic.ResponseWriter.(http.Flusher).Flush() }
