package route // import "gopkg.in/webnice/web.v1/route"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

// InternalServerError defines a handler to respond whenever a internal server error
func (rou *impl) InternalServerError(fn http.HandlerFunc) {
	var tmp = rou
	var hFn = fn
	// Build InternalServerError handler chain
	if rou.inline && rou.parent != nil {
		tmp = rou.parent
		hFn = rou.Chain(rou.middlewares...).HandlerFunc(hFn).ServeHTTP
	}
	// Update the methodNotAllowedHandler from this point forward
	tmp.internalServerErrorHandler = hFn
	tmp.updateSubRoutes(func(subRoutes *impl) {
		if subRoutes.internalServerErrorHandler == nil {
			subRoutes.InternalServerError(hFn)
		}
	})
}
