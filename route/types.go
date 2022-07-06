package route

import (
	"net/http"
	"sync"

	"github.com/webnice/web/v3/context"
	"github.com/webnice/web/v3/context/errors"
	"github.com/webnice/web/v3/context/handlers"
)

// Interface is an interface of router consisting of the core routing methods, using only the standard net/http
type Interface interface {
	http.Handler
	Routes

	// Use appends one of more middlewares onto the Router stack
	Use(...func(http.Handler) http.Handler)

	// With adds inline middlewares for an endpoint handler
	With(...func(http.Handler) http.Handler) Interface

	// Group adds a new inline-Router along the current routing
	// path, with a fresh middleware stack for the inline-router
	Group(func(r Interface)) Interface

	// Subroute mounts a sub-Router along a `pattern` string
	Subroute(string, func(r Interface)) Interface

	// Mount attaches another http.Handler along ./pattern/*
	Mount(string, http.Handler)

	// Handle and HandleFunc adds routes for `pattern` that matches all HTTP methods
	Handle(string, http.Handler)
	HandleFunc(string, http.HandlerFunc)

	// All HTTP-methods

	// Connect http method
	Connect(string, http.HandlerFunc)
	// Delete http method
	Delete(string, http.HandlerFunc)
	// Get http method
	Get(string, http.HandlerFunc)
	// Head http method
	Head(string, http.HandlerFunc)
	// Options http method
	Options(string, http.HandlerFunc)
	// Patch http method
	Patch(string, http.HandlerFunc)
	// Post http method
	Post(string, http.HandlerFunc)
	// Put http method
	Put(string, http.HandlerFunc)
	// Trace http method
	Trace(string, http.HandlerFunc)

	// Errors interface
	Errors() errors.Interface

	// Handlers interface
	Handlers() handlers.Interface
}

// Routes interface adds two methods for router traversal
type Routes interface {
	// Routes returns the routing tree in an easily traversable structure
	Routes() []Route

	// Middlewares returns the list of middlewares in use by the router
	Middlewares() Middlewares
}

// Middlewares type is a slice of standard middleware handlers with methods to compose middleware chains and http.Handler's
type Middlewares []func(http.Handler) http.Handler

// Is an private route implementation
type impl struct {
	// Accumulation of errors when configuring the router
	errors []error

	// Context interface
	context context.Interface

	// The radix trie router
	tree *node

	// The middleware stack
	middlewares []func(http.Handler) http.Handler

	// Controls the behaviour of middleware chain generation when a route is registered as an inline group inside another route
	inline bool
	parent *impl

	// The computed route handler made of the chained middleware stack and the tree router
	handler http.Handler

	// Routing context pool
	pool sync.Pool
}
