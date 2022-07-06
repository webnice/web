package errors

import "github.com/webnice/web/v3/ambry"

// New returns new context object
func New() Interface {
	var ce = new(impl)

	ce.errors = ambry.New()

	return ce
}

// Set and get values by key
func (ce *impl) do(key uint32, value error) (ret error) {
	var tmp interface{}

	if value != nil {
		ce.errors.Set(key, value)
	}
	if tmp = ce.errors.Get(key); tmp != nil {
		ret = tmp.(error)
	}
	return
}

// Reset all stored errors
func (ce *impl) Reset() { ce.errors = ambry.New() }

// RouteConfigurationError Set description about roting configuration error
func (ce *impl) RouteConfigurationError(value error) (ret error) {
	return ce.do(keyRouteConfigurationError, value)
}

// InternalServerError Set description about internal server error
func (ce *impl) InternalServerError(value error) (ret error) {
	return ce.do(keyInternalServerError, value)
}

// MethodNotAllowed Set description about method not allowed error
func (ce *impl) MethodNotAllowed(value error) (ret error) {
	return ce.do(keyMethodNotAllowed, value)
}

// NotFound Set description about not allowed error
func (ce *impl) NotFound(value error) (ret error) {
	return ce.do(keyNotFound, value)
}
