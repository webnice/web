package errors

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/ambry"

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
