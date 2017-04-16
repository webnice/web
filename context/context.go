package context //import "gopkg.in/webnice/web.v1/context"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/route"
import "gopkg.in/webnice/web.v1/context/contexterror"
import (
	"context"
	"net/http"
)

// New returns a new routing Context object
func New() Interface {
	var ctx = new(impl)
	ctx.route = route.New()
	ctx.ctxerror = contexterror.New()
	return ctx
}

// Route context interface
func (ctx *impl) Route() route.Interface { return ctx.route }

// Error context interface
func (ctx *impl) Error() contexterror.Interface { return ctx.ctxerror }

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

// NewRequest Create new http request and copy context from parent request to new request
func (ctx *impl) NewRequest(rq *http.Request) *http.Request {
	return rq.WithContext(context.WithValue(rq.Context(), _ContextKey, ctx))
}
