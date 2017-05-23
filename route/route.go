package route // import "gopkg.in/webnice/web.v1/route"

//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/web.v1/context/handlers"
import "gopkg.in/webnice/web.v1/context/errors"
import "gopkg.in/webnice/web.v1/method"
import (
	"fmt"
	"net/http"
)

// New returns a newly initialized router object that implements the Router interface
func New() Interface {
	var rou = &impl{tree: &node{}}
	rou.context = context.New()
	rou.pool.New = func() interface{} {
		var ctx context.Interface
		ctx = context.New(rou.context)
		return ctx
	}
	debug.Nop()
	return rou
}

// Errors interface
func (rou *impl) Errors() errors.Interface { return rou.context.Errors() }

// Handlers interface
func (rou *impl) Handlers() handlers.Interface { return rou.context.Handlers() }

// ServeHTTP is the single method of the http.Handler interface that makes
// route interoperable with the standard library.
// It uses a sync.Pool to get and reuse routing contexts for each request
func (rou *impl) ServeHTTP(wr http.ResponseWriter, rq *http.Request) {
	var ctx context.Interface

	// Fetch a RouteContext object from the sync pool, and call the computed
	// handler that is comprised of middlewares + routeHTTP
	// Once the request is finished, reset the routing context and put it back
	// into the pool for reuse from another request
	if !context.IsContext(rq) {
		ctx = rou.pool.Get().(context.Interface)
		ctx.Route().Reset()
		ctx.Errors().Reset()
		ctx.Handlers().InternalServerError(rou.context.Handlers().InternalServerError(nil))
		rq = ctx.NewRequest(rq)
		defer rou.pool.Put(ctx)
	}

	// Ensure the route has some routes defined on the route
	if rou.handler == nil {
		ctx = context.New(rq)
		ctx.Errors().InternalServerError(fmt.Errorf("Attempting to route with no handlers"))
		ctx.Handlers().InternalServerError(nil)(wr, rq)
		rou.pool.Put(ctx)
		return
	}

	rou.handler.ServeHTTP(wr, rq)
}

// Use appends a middleware handler to the middleware stack
//
// The middleware stack for any route will execute before searching for a matching
// route to a specific handler, which provides opportunity to respond early,
// change the course of the request execution, or set request-scoped values for
// the next http.Handler
func (rou *impl) Use(middlewares ...func(http.Handler) http.Handler) {
	const errorMiddlewares = "all middlewares must be defined before use"
	if rou.handler != nil {
		rou.context.Errors().InternalServerError(fmt.Errorf(errorMiddlewares))
		panic(errorMiddlewares)
	}
	rou.middlewares = append(rou.middlewares, middlewares...)
}

// Handle adds the route `pattern` that matches any http method to execute the `handler` http.Handler
func (rou *impl) Handle(pattern string, handler http.Handler) {
	rou.handle(method.Any, pattern, handler)
}

// HandleFunc adds the route `pattern` that matches any http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) HandleFunc(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Any, pattern, handlerFn)
}

// Connect adds the route `pattern` that matches a CONNECT http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Connect(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Connect, pattern, handlerFn)
}

// Delete adds the route `pattern` that matches a DELETE http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Delete(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Delete, pattern, handlerFn)
}

// Get adds the route `pattern` that matches a GET http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Get(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Get, pattern, handlerFn)
}

// Head adds the route `pattern` that matches a HEAD http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Head(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Head, pattern, handlerFn)
}

// Options adds the route `pattern` that matches a OPTIONS http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Options(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Options, pattern, handlerFn)
}

// Patch adds the route `pattern` that matches a PATCH http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Patch(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Patch, pattern, handlerFn)
}

// Post adds the route `pattern` that matches a POST http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Post(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Post, pattern, handlerFn)
}

// Put adds the route `pattern` that matches a PUT http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Put(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Put, pattern, handlerFn)
}

// Trace adds the route `pattern` that matches a TRACE http method to execute the `handlerFn` http.HandlerFunc
func (rou *impl) Trace(pattern string, handlerFn http.HandlerFunc) {
	rou.handle(method.Trace, pattern, handlerFn)
}

// Chain returns a Middlewares type from a slice of middleware handlers
func (rou *impl) Chain(middlewares ...func(http.Handler) http.Handler) Middlewares {
	return Middlewares(middlewares)
}

// With adds inline middlewares for an endpoint handler
func (rou *impl) With(middlewares ...func(http.Handler) http.Handler) Interface {
	var im *impl

	// Similarly as in handle(), we must build the route handler once further
	// middleware registration isn't allowed for this stack, like now
	if !rou.inline && rou.handler == nil {
		rou.buildRouteHandler()
	}

	// Copy middlewares from parent inline route
	var mws Middlewares
	if rou.inline {
		mws = make(Middlewares, len(rou.middlewares))
		copy(mws, rou.middlewares)
	}
	mws = append(mws, middlewares...)

	im = &impl{inline: true, parent: rou, tree: rou.tree, middlewares: mws}
	return im
}

// Group creates a new inline-route with a fresh middleware stack. It's useful
// for a group of handlers along the same routing path that use an additional set of middlewares
func (rou *impl) Group(fn func(r Interface)) Interface {
	var im = rou.With().(*impl)
	if fn != nil {
		fn(im)
	}
	return im
}

// Subroute creates a new route with a fresh middleware stack and mounts it
// along the `pattern` as a subrouter. Effectively, this is a short-hand call to Mount
func (rou *impl) Subroute(pattern string, fn func(r Interface)) Interface {
	var subRouter = New()
	if fn != nil {
		fn(subRouter)
	}
	rou.Mount(pattern, subRouter)
	return subRouter
}

// Mount attaches another http.Handler or Router as a subrouter along a routing
// path. It's very useful to split up a large API as many independent routers and
// compose them as a single service using Mount.
//
// Note that Mount() simply sets a wildcard along the `pattern` that will continue
// routing at the `handler`, which in most cases is another router.
// As a result, if you define two Mount() routes on the exact same pattern the mount will panic
func (rou *impl) Mount(pattern string, handler http.Handler) {
	var subr *impl
	var ok bool
	var mtd method.Method
	var n *node
	var subHandler http.HandlerFunc
	var subroutes Routes

	// Provide runtime safety for ensuring a pattern isn't mounted on an existing routing pattern
	if rou.tree.findPattern(pattern+"*") != nil || rou.tree.findPattern(pattern+"/*") != nil {
		var mountError = fmt.Sprintf("attempting to Mount() a handler on an existing path, '%s'", pattern)
		panic(mountError)
	}

	// Assign sub-Router's with the parent not found & method not allowed handler if not specified.
	subr, ok = handler.(*impl)
	if ok && subr.notFoundHandler == nil && rou.notFoundHandler != nil {
		subr.NotFound(rou.notFoundHandler)
	}
	if ok && subr.methodNotAllowedHandler == nil && rou.methodNotAllowedHandler != nil {
		subr.MethodNotAllowed(rou.methodNotAllowedHandler)
	}

	// Wrap the sub-router in a handlerFunc to scope the request path for routing.
	subHandler = http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		var ctx context.Interface
		ctx = context.New(rq)
		ctx.Route().Path("/" + ctx.Route().UrnParams().Del("*"))
		handler.ServeHTTP(wr, rq)
	})

	if pattern == "" || pattern[len(pattern)-1] != '/' {
		rou.handle(method.Any|method.Stub, pattern, subHandler)
		rou.handle(method.Any|method.Stub, pattern+"/", rou.NotFoundHandler())
		pattern += "/"
	}

	mtd = method.Any
	subroutes, _ = handler.(Routes)
	if subroutes != nil {
		mtd |= method.Stub
	}
	n = rou.handle(mtd, pattern+"*", subHandler)
	if subroutes != nil {
		n.subroutes = subroutes
	}
}

// Middlewares return all middlewares
func (rou *impl) Middlewares() Middlewares { return rou.middlewares }

// Routes Return all routes
func (rou *impl) Routes() []Route { return rou.tree.routes() }

// builds the single handler that is a chain of the middleware stack, as defined by
// calls to Use(), and the tree router itself. After this point, no other middlewares
// can be registered on stack. But you can still compose additional middlewares via
// Group()'s or using a chained middleware handler
func (rou *impl) buildRouteHandler() {
	rou.handler = chain(rou.middlewares, http.HandlerFunc(rou.routeHTTP))
}

// registers a http.Handler in the routing tree for a particular http method and routing pattern
func (rou *impl) handle(mtd method.Method, pattern string, handler http.Handler) *node {
	if len(pattern) == 0 || pattern[0] != '/' {
		var handleError = fmt.Sprintf("routing pattern must begin with '/' in '%s'", pattern)
		panic(handleError)
	}

	// Build the final routing handler for this route
	if !rou.inline && rou.handler == nil {
		rou.buildRouteHandler()
	}

	// Build endpoint handler with inline middlewares for the route
	var h http.Handler
	if rou.inline {
		rou.handler = http.HandlerFunc(rou.routeHTTP)
		h = rou.Chain(rou.middlewares...).Handler(handler)
	} else {
		h = handler
	}

	// Add the endpoint to the tree and return the node
	return rou.tree.InsertRoute(mtd, pattern, h)
}

// routes a http.Request through the routing tree to serve the matching handler for a particular http method
func (rou *impl) routeHTTP(wr http.ResponseWriter, rq *http.Request) {
	var err error
	var ctx context.Interface
	var routePath string
	var mtd method.Method
	var hs methodHandlers
	var h http.Handler
	var ok bool
	// Grab the route context object
	ctx = context.New(rq)
	// The request routing path
	routePath = ctx.Route().Path()
	if routePath == "" {
		if rq.URL.RawPath != "" {
			routePath = rq.URL.RawPath
		} else {
			routePath = rq.URL.Path
		}
	}
	// Check if method is supported
	if mtd, err = method.Parse(rq.Method); err != nil {
		rou.MethodNotAllowedHandler().ServeHTTP(wr, rq)
		return
	}
	// Find the route
	hs = rou.tree.FindRoute(ctx, routePath)
	if hs == nil {
		rou.NotFoundHandler().ServeHTTP(wr, rq)
		return
	}
	if h, ok = hs[mtd]; !ok {
		rou.MethodNotAllowedHandler().ServeHTTP(wr, rq)
		return
	}
	// Serve it up
	h.ServeHTTP(wr, rq)
}

// Recursively update data on child routers
func (rou *impl) updateSubRoutes(fn func(subMux *impl)) {
	var r Route
	var ok bool
	var subMux *impl
	for _, r = range rou.tree.routes() {
		subMux, ok = r.SubRoutes.(*impl)
		if !ok {
			continue
		}
		fn(subMux)
	}
}
