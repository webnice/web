package route

import "github.com/webnice/web/v3/param"

// This is the default routing context set on the root node of a request context to track URL parameters and an optional routing path
type impl struct {
	params   param.Interface
	path     string
	pattern  string
	patterns []string
}

// Interface is an interface of package
type Interface interface {
	// Reset a routing context to its initial state
	Reset()

	// Params Return routing URN parameters key and values
	Params() param.Interface

	// Path Routing path override used by subrouters
	Path(...string) string

	// Pattern Routing pattern matching the path
	Pattern(...string) string

	// Patterns Routing patterns throughout the lifecycle of the request, across all connected routers
	Patterns(...[]string) []string
}
