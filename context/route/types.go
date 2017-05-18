package route // import "gopkg.in/webnice/web.v1/context/route"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/param"

// This is the default routing context set on the root node of a request context to track URL parameters and an optional routing path
type impl struct {
	Params   param.Params
	path     string
	pattern  string
	patterns []string
}

// Interface is an interface of package
type Interface interface {
	// Reset a routing context to its initial state
	Reset()

	// UrnParams Return routing URN parameters key and values
	UrnParams() param.Params

	// Path Routing path override used by subrouters
	Path(...string) string

	// Pattern Routing pattern matching the path
	Pattern(...string) string

	// Patterns Routing patterns throughout the lifecycle of the request, across all connected routers
	Patterns(...[]string) []string
}
