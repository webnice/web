package route

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "net/http"

// ChainHandler is a type of chain handlers
type ChainHandler struct {
	Middlewares Middlewares
	Endpoint    http.Handler
	chain       http.Handler
}

// ServeHTTP Inplementation of http.Handler interface
func (c *ChainHandler) ServeHTTP(wr http.ResponseWriter, rq *http.Request) { c.chain.ServeHTTP(wr, rq) }

// Handler builds and returns a http.Handler from the chain of middlewares,
// with `h http.Handler` as the final handler
func (mws Middlewares) Handler(hfn http.Handler) http.Handler {
	return &ChainHandler{mws, hfn, chain(mws, hfn)}
}

// HandlerFunc builds and returns a http.Handler from the chain of middlewares,
// with `h http.Handler` as the final handler
func (mws Middlewares) HandlerFunc(hfn http.HandlerFunc) http.Handler {
	return &ChainHandler{mws, hfn, chain(mws, hfn)}
}

// chain builds a http.Handler composed of an inline middleware stack and endpoint
// handler in the order they are passed
func chain(middlewares []func(http.Handler) http.Handler, endpoint http.Handler) (ret http.Handler) {
	var i int
	// Return ahead of time if there aren't any middlewares for the chain
	if len(middlewares) == 0 {
		ret = endpoint
		return
	}
	// Wrap the end handler with the middleware chain
	ret = middlewares[len(middlewares)-1](endpoint)
	for i = len(middlewares) - 2; i >= 0; i-- {
		ret = middlewares[i](ret)
	}
	return
}
