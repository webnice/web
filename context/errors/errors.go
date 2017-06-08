package errors

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
//import "fmt"

// New returns new context object
func New() Interface {
	var ce = new(impl)
	ce.errors = new(errors)
	return ce
}

// Add key value param to heap
func (itm *errors) Add(key uint32, value error) {
	*itm = append(*itm, item{key, value})
}

// Get value by key from heap
func (itm errors) Get(key uint32) (ret error) {
	var p item
	for _, p = range itm {
		if p.Key == key {
			ret = p.Value
			break
		}
	}
	return
}

// Set value by key to heap
func (itm *errors) Set(key uint32, value error) {
	var p item
	var i, idx int
	idx = -1
	for i, p = range *itm {
		if p.Key == key {
			idx = i
			break
		}
	}
	if idx < 0 {
		(*itm).Add(key, value)
	} else {
		(*itm)[idx] = item{key, value}
	}
}

// Del Delete value by key from heap
func (itm *errors) Del(key uint32) (ret error) {
	var p item
	var i int
	for i, p = range *itm {
		if p.Key == key {
			*itm = append((*itm)[:i], (*itm)[i+1:]...)
			ret = p.Value
		}
	}
	return
}

// Reset all stored errors
func (ce *impl) Reset() { ce.errors = new(errors) }

// InternalServerError Set description about internal server error
func (ce *impl) InternalServerError(value error) (ret error) {
	if value != nil {
		ce.errors.Set(keyInternalServerError, value)
	}
	ret = ce.errors.Get(keyInternalServerError)
	return
}

// MethodNotAllowed Set description about method not allowed error
func (ce *impl) MethodNotAllowed(value error) (ret error) {
	if value != nil {
		ce.errors.Set(keyMethodNotAllowed, value)
	}
	ret = ce.errors.Get(keyMethodNotAllowed)
	return
}

// NotFound Set description about not allowed error
func (ce *impl) NotFound(value error) (ret error) {
	if value != nil {
		ce.errors.Set(keyNotFound, value)
	}
	ret = ce.errors.Get(keyNotFound)
	return
}
