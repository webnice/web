package context

import (
	"net/http"

	"github.com/webnice/web/v2/context/errors"
	"github.com/webnice/web/v2/context/handlers"
	"github.com/webnice/web/v2/context/route"
)

var (
	constContextKey    = &key{`2EA70633-A7DF-44AD-9DE4-C7593085B685`}
	globalVerifyPlugin VerifyPlugin // Registered interface of data verification external library
)

// This is the default routing context implementation
type impl struct {
	route    route.Interface    // Routing interface
	errors   errors.Interface   // Errors interface
	handlers handlers.Interface // Handlers interface
	Request  *http.Request      // net/http request
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

	// Data Extracting from a request and decoding data to structure of obj
	Data(obj interface{}) ([]byte, error)
}

// VerifyPlugin External data verification library interface
type VerifyPlugin interface {
	// Verify Check data function
	Verify(data interface{}) ([]byte, error)

	// Error400 Create response data for HTTP error 400 in the library format
	Error400(err error) []byte
}
