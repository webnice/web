package middleware

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/middleware/wrapsrw"
import (
	"net/http"
)

// NewWrapsResponseWriter Create new object and return interface
func NewWrapsResponseWriter(wr http.ResponseWriter, protoMajor int) wrapsrw.WrapsResponseWriter {
	return wrapsrw.New(wr, protoMajor)
}
