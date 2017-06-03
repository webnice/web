package route

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import (
	"net/http"
)

// MethodNotAllowedHandler returns the default 405 responder whenever a method cannot be resolved for a route
func (rou *impl) MethodNotAllowedHandler() http.HandlerFunc {
	if rou.methodNotAllowedHandler != nil {
		return rou.methodNotAllowedHandler
	}
	return methodNotAllowedHandler
}

// MethodNotAllowed sets a custom http.HandlerFunc for routing paths where the method is unresolved
func (rou *impl) MethodNotAllowed(handlerFn http.HandlerFunc) {
	var tmp = rou
	var hFn = handlerFn
	// Build MethodNotAllowed handler chain
	if rou.inline && rou.parent != nil {
		tmp = rou.parent
		hFn = rou.Chain(rou.middlewares...).HandlerFunc(hFn).ServeHTTP
	}
	// Update the methodNotAllowedHandler from this point forward
	tmp.methodNotAllowedHandler = hFn
	tmp.updateSubRoutes(func(subMux *impl) {
		if subMux.methodNotAllowedHandler == nil {
			subMux.MethodNotAllowed(hFn)
		}
	})
}

// is a helper function to respond with a 405, method not allowed
func methodNotAllowedHandler(wr http.ResponseWriter, rq *http.Request) {
	wr.WriteHeader(status.MethodNotAllowed)
	wr.Write(status.Bytes(status.MethodNotAllowed))
}
