package errors // import "gopkg.in/webnice/web.v1/context/errors"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// New returns new context object
func New() Interface {
	var ce = new(impl)
	ce.errors = make(errors)
	return ce
}

// Add adds the key, value pair to the errors.
// It appends to any existing values associated with key
func (h errors) Add(key uint16, value error) {
	h[key] = append(h[key], value)
}

// Set sets the params entries associated with key to
// the single element value. It replaces any existing
// values associated with key
func (h errors) Set(key uint16, value error) {
	h[key] = []error{value}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns nil.
// To access multiple values of a key, or to use access the map directly
func (h errors) Get(key uint16) error {
	if h == nil {
		return nil
	}
	v := h[key]
	if len(v) == 0 {
		return nil
	}
	return v[0]
}

// Del deletes the values associated with key
func (h errors) Del(key uint16) error {
	var val = h.Get(key)
	delete(h, key)
	return val
}

// InternalServerError Set description about internal server error
func (ce *impl) InternalServerError(value error) (ret error) {
	if value != nil {
		ce.errors.Set(_InternalServerError, value)
	}
	ret = ce.errors.Get(_InternalServerError)
	return
}
