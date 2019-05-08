package recovery

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context"
import (
	"fmt"
	"net/http"
	runtimeDebug "runtime/debug"
)

// Handler is a middleware that recovers from panics
func Handler(next http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				var ctx = context.New(rq)
				_ = ctx.Errors().InternalServerError(
					fmt.Errorf(
						"Catch panic: %v\nGoroutine stack is:\n%s",
						e,
						string(runtimeDebug.Stack()),
					),
				)
				ctx.Handlers().InternalServerError(nil).ServeHTTP(wr, rq)
			}
		}()
		next.ServeHTTP(wr, rq)
	}

	return http.HandlerFunc(fn)
}
