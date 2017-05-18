package param // import "gopkg.in/webnice/web.v1/param"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

type Params map[string][]string

// Add adds the key, value pair to the params.
// It appends to any existing values associated with key
func (h Params) Add(key, value string) {
	h[key] = append(h[key], value)
}

// Set sets the params entries associated with key to
// the single element value. It replaces any existing
// values associated with key
func (h Params) Set(key, value string) {
	h[key] = []string{value}
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, or to use access the map directly
func (h Params) Get(key string) string {
	if h == nil {
		return ""
	}
	v := h[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

// Del deletes the values associated with key
func (h Params) Del(key string) string {
	var val = h.Get(key)
	delete(h, key)
	return val
}
