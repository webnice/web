package errors // import "gopkg.in/webnice/web.v1/context/errors"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/param"

// New returns new context object
func New() Interface {
	var ce = new(impl)
	ce.params = new(param.Params)
	return ce
}

// Add key value to heap
func (ce *impl) Add(key string, value string) { ce.params.Add(key, value) }

// Get value by key from heap
func (ce *impl) Get(key string) string { return ce.params.Get(key) }

// Set value by key to heap
func (ce *impl) Set(key string, value string) { ce.params.Set(key, value) }

// Del Delete value by key from heap
func (ce *impl) Del(key string) string { return ce.params.Del(key) }
