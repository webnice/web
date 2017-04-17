package wrapsResponseWriter //import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"bufio"
	"io"
	"net"
	"net/http"
)

func (f *httpFancyWriter) CloseNotify() <-chan bool {
	cn := f.basic.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}

func (f *httpFancyWriter) Flush() {
	fl := f.basic.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (f *httpFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basic.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (f *httpFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	if f.basic.tee != nil {
		return io.Copy(&f.basic, r)
	}
	rf := f.basic.ResponseWriter.(io.ReaderFrom)
	f.basic.maybeWriteHeader()
	return rf.ReadFrom(r)
}

// HTTP2

func (f *http2FancyWriter) CloseNotify() <-chan bool {
	cn := f.basic.ResponseWriter.(http.CloseNotifier)
	return cn.CloseNotify()
}
func (f *http2FancyWriter) Flush() {
	fl := f.basic.ResponseWriter.(http.Flusher)
	fl.Flush()
}

// Go 1.8 http Push

func (f *http2FancyWriter) Push(target string, opts *http.PushOptions) error {
	return f.basic.ResponseWriter.(http.Pusher).Push(target, opts)
}
