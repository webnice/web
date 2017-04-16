package context // import "gopkg.in/webnice/web.v1/context"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/route"
import "gopkg.in/webnice/web.v1/context/contexterror"
import (
	"net/http"
)

var (
	_ContextKey = &key{`RwAQ7mHz4gd8agG2cx7jJ9pQQGk7n3J2dsJFuAd5TjP7NSAVfCRZ9ruxCPpH6mJg`}
)

// This is the default routing context implementation
type impl struct {
	route               route.Interface        // Routing context space
	ctxerror            contexterror.Interface // Error context space
	internalServerError http.HandlerFunc       // InternalServerError handler function
}

// Interface is an interface of package
type Interface interface {
	// Route context interface
	Route() route.Interface

	// Error context interface
	Error() contexterror.Interface

	// Set and get InternalServerError handler function
	InternalServerError(...http.HandlerFunc) http.HandlerFunc

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
