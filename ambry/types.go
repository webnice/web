package ambry

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "sync"

// Interface is an interface of package
type Interface interface {
	// Set sets the params entries associated with key to
	// the single element value. It replaces any existing
	// values associated with key
	Set(key interface{}, value interface{})

	// Add adds the key, value pair to the params.
	// It appends to any existing values associated with key
	Add(key interface{}, value interface{})

	// Has will return a boolean, which will be true or
	// false depending on whether the item exists
	Has(key interface{}) bool

	// Get gets the first value associated with the given key.
	// If there are no values associated with the key, Get returns "".
	// To access multiple values of a key, or to use access the map directly
	Get(key interface{}) interface{}

	// Del deletes the values associated with key
	Del(key interface{}) interface{}

	// Keys gets the all keys
	// If there are no values, Keys returns empty slice
	Keys() []interface{}

	// Get gets the all value associated with the given key.
	// If there are no values associated with the key, Get returns empty slice
	GetAll(key interface{}) []interface{}
}

// impl is an implementation of package
type impl struct {
	sync.RWMutex
	items map[key][]*item
}

type (
	key  interface{}
	item struct {
		Value interface{}
	}
)
