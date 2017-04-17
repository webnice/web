package wrapsResponseWriter //import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import (
	"io"
	"net/http"
)

func (b *basic) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basic) Write(buf []byte) (int, error) {
	b.WriteHeader(status.Ok)
	n, err := b.ResponseWriter.Write(buf)
	if b.tee != nil {
		_, err2 := b.tee.Write(buf[:n])
		// Prefer errors generated by the proxied writer.
		if err == nil {
			err = err2
		}
	}
	b.bytes += uint64(n)
	return n, err
}

func (b *basic) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(status.Ok)
	}
}

func (b *basic) Status() int { return b.code }

func (b *basic) Len() uint64 { return b.bytes }

func (b *basic) Tee(w io.Writer) { b.tee = w }

func (b *basic) Unwrap() http.ResponseWriter { return b.ResponseWriter }
