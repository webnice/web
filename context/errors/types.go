package errors

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/ambry"

const (
	keyRouteConfigurationError uint32 = iota
	keyInternalServerError
	keyMethodNotAllowed
	keyNotFound
)

// Interface is an interface of package
type Interface interface {
	// RouteConfigurationError Set description about roting configuration error
	RouteConfigurationError(err error) error

	// InternalServerError Set description about internal server error
	InternalServerError(err error) error

	// MethodNotAllowed Set description about method not allowed error
	MethodNotAllowed(err error) error

	// NotFound Set description about not found error
	NotFound(err error) error

	// Reset all stored errors
	Reset()
}

// This is an inplementation
type impl struct {
	errors ambry.Interface
}
