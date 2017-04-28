package route // import "gopkg.in/webnice/web.v1/route"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/web.v1/status"
import (
	"net/http"
)

// InternalServerErrorHandler returns the default 500 responder whenever a internal server error
func (rou *impl) InternalServerErrorHandler() http.HandlerFunc {
	if rou.internalServerErrorHandler != nil {
		return rou.internalServerErrorHandler
	}
	return internalServerErrorHandler
}

// InternalServerError defines a handler to respond whenever a internal server error
func (rou *impl) InternalServerError(handlerFn http.HandlerFunc) {
	var tmp = rou
	var hFn = handlerFn
	// Build InternalServerError handler chain
	if rou.inline && rou.parent != nil {
		tmp = rou.parent
		hFn = rou.Chain(rou.middlewares...).HandlerFunc(hFn).ServeHTTP
	}
	// Update the methodNotAllowedHandler from this point forward
	tmp.internalServerErrorHandler = hFn
	tmp.updateSubRoutes(func(subMux *impl) {
		if subMux.internalServerErrorHandler == nil {
			subMux.InternalServerError(hFn)
		}
	})
}

// is a helper function to respond with a 500, internal server error
func internalServerErrorHandler(wr http.ResponseWriter, rq *http.Request) {
	wr.WriteHeader(status.InternalServerError)
	wr.Write(status.Bytes(status.InternalServerError))
	wr.Write([]byte(", " + context.ContextFromRequest(rq).Error().Get(_KeyInternalServerError)))
}
