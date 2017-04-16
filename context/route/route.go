package route // import "gopkg.in/webnice/web.v1/context/route"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/param"
import (
	"strings"
)

// New returns a new routing context object
func New() Interface {
	var rt = new(impl)
	rt.Params = new(param.Params)
	return rt
}

// Reset a routing context to its initial state
func (rt *impl) Reset() {
	rt.Params = new(param.Params)
	rt.path = ""
	rt.pattern = ""
	rt.patterns = rt.patterns[:0]
}

// UrnParams Return routing URN parameters key and values
func (rt *impl) UrnParams() *param.Params { return rt.Params }

// Path Routing path override used by subrouters
func (rt *impl) Path(str ...string) string {
	if len(str) > 0 {
		rt.path = strings.Join(str, ``)
	}
	return rt.path
}

// Pattern Routing pattern matching the path
func (rt *impl) Pattern(str ...string) string {
	if len(str) > 0 {
		rt.pattern = strings.Join(str, ``)
	}
	return rt.pattern
}

// Patterns Routing patterns throughout the lifecycle of the request, across all connected routers
func (rt *impl) Patterns(items ...[]string) []string {
	var i int
	if len(items) > 0 {
		for i = range items {
			if i == 0 {
				rt.patterns = items[i]
			} else {
				rt.patterns = append(rt.patterns, items[i]...)
			}
		}
	}
	return rt.patterns
}
