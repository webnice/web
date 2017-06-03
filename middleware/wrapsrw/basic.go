package wrapsrw // import "gopkg.in/webnice/web.v1/middleware/wrapsrw"

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

func (b *basic) Write(buf []byte) (ret int, err error) {
	b.WriteHeader(status.Ok)
	ret, err = b.ResponseWriter.Write(buf)
	if b.tee != nil {
		if _, e := b.tee.Write(buf[:ret]); err == nil {
			err = e
		}
	}
	b.bytes += uint64(ret)
	return
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
