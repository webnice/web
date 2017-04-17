package wrapsResponseWriter //import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

func (f *flush) Flush() { f.basic.ResponseWriter.(http.Flusher).Flush() }
