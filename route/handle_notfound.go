package route // import "gopkg.in/webnice/web.v1/route"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import (
	"net/http"
)

// NotFoundHandler returns the default 404 responder whenever a route cannot be found
func (rou *impl) NotFoundHandler() http.HandlerFunc {
	if rou.notFoundHandler != nil {
		return rou.notFoundHandler
	}
	return notFoundHandler
}

// NotFound sets a custom http.HandlerFunc for routing paths that could not be found
func (rou *impl) NotFound(handlerFn http.HandlerFunc) {
	// Build NotFound handler chain
	tmp := rou
	hFn := handlerFn
	if rou.inline && rou.parent != nil {
		tmp = rou.parent
		hFn = rou.Chain(rou.middlewares...).HandlerFunc(hFn).ServeHTTP
	}
	// Update the notFoundHandler from this point forward
	tmp.notFoundHandler = hFn
	tmp.updateSubRoutes(func(subMux *impl) {
		if subMux.notFoundHandler == nil {
			subMux.NotFound(hFn)
		}
	})
}

// is a helper function to respond with a 404, not found
func notFoundHandler(wr http.ResponseWriter, rq *http.Request) {
	wr.WriteHeader(status.NotFound)
	wr.Write(status.Bytes(status.NotFound))
}
