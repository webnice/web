// +build go1.8

package wrapsResponseWriter // import "gopkg.in/webnice/web.v1/middleware/wrapsResponseWriter"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "net/http"

var (
	_ http.Pusher = &http2FancyWriter{}
)
