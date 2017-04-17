package middleware //import "gopkg.in/webnice/web.v1/middleware"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"
import (
	"net/http"
)

func NewWrapsResponseWriter(wr http.ResponseWriter, protoMajor int) wrapsResponseWriter.WrapsResponseWriter {
	return wrapsResponseWriter.New(wr, protoMajor)
}
