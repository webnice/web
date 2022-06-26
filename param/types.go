package param

import "sync"

type (
	key  string
	item struct {
		Value string
	}
)

// This is the default routing context implementation
type impl struct {
	sync.RWMutex
	items map[key][]*item
}

// Interface is an interface of package
type Interface interface {
	// Set sets the params entries associated with key to
	// the single element value. It replaces any existing
	// values associated with key
	Set(key string, value string)

	// Add adds the key, value pair to the params.
	// It appends to any existing values associated with key
	Add(key string, value string)

	// Has will return a boolean, which will be true or
	// false depending on whether the item exists
	Has(key string) bool

	// Get gets the first value associated with the given key.
	// If there are no values associated with the key, Get returns "".
	// To access multiple values of a key, or to use access the map directly
	Get(key string) string

	// Del deletes the values associated with key
	Del(key string) string

	// Keys gets the all keys
	// If there are no values, Keys returns empty slice
	Keys() []string

	// Get gets the all value associated with the given key.
	// If there are no values associated with the key, Get returns empty slice
	GetAll(key string) []string
}
