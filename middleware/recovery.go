package middleware // import "gopkg.in/webnice/web.v1/middleware"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context"
import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// Recover is a middleware that recovers from panics
func Recover(next http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				var ctx = context.ContextFromRequest(rq)
				ctx.Error().Set("InternalServerError", fmt.Sprintf("Panic: %v\nStack:\n%s", e, string(debug.Stack())))
				ctx.InternalServerError().ServeHTTP(wr, rq)
			}
		}()
		next.ServeHTTP(wr, rq)
	}
	return http.HandlerFunc(fn)
}
