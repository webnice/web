package context // import "gopkg.in/webnice/web.v1/context"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"net/http"
)

// IsContext Check if a routing context exists in http context
func IsContext(rq *http.Request) (ret bool) {
	if ContextFromRequest(rq) != nil {
		ret = true
	}
	return
}

// ContextFromRequest Get the routing Context object from a http context
func ContextFromRequest(rq *http.Request) Interface { return Context(rq.Context()) }

// Context Get the routing Context object from a http context
func Context(ctx context.Context) (ret Interface) {
	var ok bool
	if ret, ok = ctx.Value(_ContextKey).(Interface); !ok {
		ret = nil
		return
	}
	return
}
