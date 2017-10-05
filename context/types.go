package context

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/route"
import "gopkg.in/webnice/web.v1/context/errors"
import "gopkg.in/webnice/web.v1/context/handlers"
import (
	"net/http"
)

var (
	constContextKey = &key{`2EA70633-A7DF-44AD-9DE4-C7593085B685`}
)

// This is the default routing context implementation
type impl struct {
	route    route.Interface    // Routing interface
	errors   errors.Interface   // Errors interface
	handlers handlers.Interface // Handlers interface
}

// Interface is an interface of package
type Interface interface {
	// Route context interface
	Route() route.Interface

	// Errors interface
	Errors(is ...errors.Interface) errors.Interface

	// Handlers interface
	Handlers(is ...handlers.Interface) handlers.Interface

	// NewRequest Create new http.Request and copy context from parent request to new request
	NewRequest(*http.Request) *http.Request
}

// key is a value for use with context.WithValue
// It's used as a pointer so it fits in an interface{} without allocation
type key struct {
	Name string
}

// String convert type to string
func (k *key) String() string { return k.Name }
