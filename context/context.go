package context

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/route"
import "gopkg.in/webnice/web.v1/context/errors"
import "gopkg.in/webnice/web.v1/context/handlers"
import (
	stdContext "context"
	"net/http"
)

// New returns a new routing context object
// You can pass the following types of objects as arguments:
// from "net/http" type *http.Request;
// from "context" interface context.Context;
// from Interface of this package.
// If an invalid argument type is passed, the function will return nil
func New(obj ...interface{}) Interface {
	var ctx *impl
	var i int

	for i = range obj {
		switch val := obj[i].(type) {
		case *http.Request:
			ctx = request(val)
		case stdContext.Context:
			ctx = context(val)
		case Interface, *impl:
			ctx = val.(*impl)
			if ctx.errors == nil {
				ctx.errors = errors.New()
			}
			if ctx.handlers == nil {
				ctx.handlers = handlers.New(ctx.errors)
			}
			// Allways new route object
			ctx.route = route.New()
		default:
			// invalid argument type is passed
			return nil
		}
		if ctx != nil {
			return ctx
		}
	}

	ctx = new(impl)
	ctx.route = route.New()
	ctx.errors = errors.New()
	ctx.handlers = handlers.New(ctx.errors)

	return ctx
}

// Get the routing Context object from a http context
func context(cx stdContext.Context) (ret *impl) {
	var ok bool
	if ret, ok = cx.Value(constContextKey.String()).(*impl); !ok {
		return nil
	}
	return
}

// Get the routing context object from a http context
func request(rq *http.Request) *impl { return context(rq.Context()) }

// IsContext Check if a context not empty in net/http context
func IsContext(rq *http.Request) bool { return request(rq) != nil }

// Route interface
func (ctx *impl) Route() route.Interface { return ctx.route }

// Error interface
func (ctx *impl) Errors(is ...errors.Interface) (ret errors.Interface) {
	ret = ctx.errors
	for i := range is {
		if is[i] != nil {
			ctx.errors = is[i]
			break
		}
	}
	return
}

// Handlers interface
func (ctx *impl) Handlers(is ...handlers.Interface) (ret handlers.Interface) {
	ret = ctx.handlers
	for i := range is {
		if is[i] != nil {
			ctx.handlers = is[i]
			break
		}
	}
	return
}

// NewRequest Creates new http request and copy context from parent request to new request
func (ctx *impl) NewRequest(rq *http.Request) *http.Request {
	return rq.WithContext(stdContext.WithValue(rq.Context(), constContextKey.String(), ctx))
}
