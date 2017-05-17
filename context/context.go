package context // import "gopkg.in/webnice/web.v1/context"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/route"
import "gopkg.in/webnice/web.v1/context/errors"
import (
	stdContext "context"
	"net/http"
)

// New returns a new routing context object
// You can pass the following types of objects as arguments:
// - "net/http" *http.Request
// - "context" context.Context
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
		default:
			return nil
		}
	}
	if ctx != nil {
		return ctx
	}

	ctx = new(impl)
	ctx.route = route.New()
	ctx.errors = errors.New()

	return ctx
}

// Get the routing Context object from a http context
func context(cx stdContext.Context) (ret *impl) {
	var ok bool
	if ret, ok = cx.Value(_ContextKey).(*impl); ok {
		return
	}
	return
}

// Get the routing context object from a http context
func request(rq *http.Request) *impl { return context(rq.Context()) }

// IsContext Check if a context not empty in net/http context
func IsContext(rq *http.Request) (ret bool) { ret = request(rq) != nil; return }

// Route context interface
func (ctx *impl) Route() route.Interface { return ctx.route }

// Error context interface
func (ctx *impl) Errors() errors.Interface { return ctx.errors }

// InternalServerError Set and get InternalServerError handler function
func (ctx *impl) InternalServerError(fncs ...http.HandlerFunc) http.HandlerFunc {
	var i int
	for i = range fncs {
		if fncs[i] != nil {
			ctx.internalServerError = fncs[i]
		}
	}
	return ctx.internalServerError
}

// NewRequest Creates new http request and copy context from parent request to new request
func (ctx *impl) NewRequest(rq *http.Request) *http.Request {
	return rq.WithContext(stdContext.WithValue(rq.Context(), _ContextKey, ctx))
}
