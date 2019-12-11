package route

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"strings"

	"gopkg.in/webnice/web.v1/param"
)

// New returns a new routing context object
func New() Interface {
	var rt = new(impl)
	rt.params = param.New()
	return rt
}

// Reset a routing context to its initial state
func (rt *impl) Reset() {
	rt.params = param.New()
	rt.path = ""
	rt.pattern = ""
	rt.patterns = make([]string, 0)
}

// Params Return routing URN parameters key and values
func (rt *impl) Params() param.Interface { return rt.params }

// Path Routing path override used by subrouters
func (rt *impl) Path(str ...string) string {
	// if len(str) == 0 do not override!
	if len(str) == 0 {
		return rt.path
	}
	rt.path = strings.Join(str, ``)
	return rt.path
}

// Pattern Routing pattern matching the path
func (rt *impl) Pattern(str ...string) string {
	// if len(str) == 0 do not override!
	if len(str) == 0 {
		return rt.pattern
	}
	rt.pattern = strings.Join(str, ``)
	return rt.pattern
}

// Patterns Routing patterns throughout the lifecycle of the request, across all connected routers
func (rt *impl) Patterns(items ...[]string) []string {
	var i, j int

	// if len(items) == 0 do not override!
	if len(items) == 0 {
		return rt.patterns
	}

	for i = range items {
		j += len(items[i])
	}
	rt.patterns, j = make([]string, j), 0
	for i = range items {
		copy(rt.patterns[j:j+len(items[i])], items[i])
		j += len(items[i])
	}
	return rt.patterns
}
