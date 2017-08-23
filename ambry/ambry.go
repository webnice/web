package ambry

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// New is a constructor
func New() Interface {
	var itm = new(impl)
	itm.items = map[key][]*item{}
	return itm
}

// Set sets the params entries associated with key to
// the single element value. It replaces any existing
// values associated with key
func (itm *impl) Set(k interface{}, v interface{}) {
	itm.Lock()
	itm.items[key(k)] = []*item{{Value: v}}
	itm.Unlock()
}

// Add adds the key, value pair to the params.
// It appends to any existing values associated with key
func (itm *impl) Add(k interface{}, v interface{}) {
	var ok bool
	itm.RLock()
	_, ok = itm.items[key(k)]
	itm.RUnlock()
	if ok {
		itm.Lock()
		itm.items[key(k)] = append(itm.items[key(k)], &item{Value: v})
		itm.Unlock()
	} else {
		itm.Set(k, v)
	}
}

// Has will return a boolean, which will be true or
// false depending on whether the item exists
func (itm *impl) Has(k interface{}) (ret bool) {
	itm.RLock()
	_, ret = itm.items[key(k)]
	itm.RUnlock()
	return
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, or to use access the map directly
func (itm *impl) Get(k interface{}) (ret interface{}) {
	var value []*item
	var ok bool

	itm.RLock()
	value, ok = itm.items[key(k)]
	itm.RUnlock()
	if !ok {
		return
	}
	ret = value[0].Value
	return
}

// Del deletes the values associated with key
func (itm *impl) Del(k interface{}) (ret interface{}) {
	var value []*item
	var ok bool

	itm.RLock()
	value, ok = itm.items[key(k)]
	itm.RUnlock()
	if !ok {
		return
	}
	itm.Lock()
	delete(itm.items, key(k))
	itm.Unlock()
	ret = value[0].Value
	return
}

// Keys gets the all keys.
// If there are no values, Keys returns empty slice
func (itm *impl) Keys() (ret []interface{}) {
	var i key
	itm.RLock()
	for i = range itm.items {
		ret = append(ret, i)
	}
	itm.RUnlock()
	return
}

// Get gets the all value associated with the given key.
// If there are no values associated with the key, Get returns empty slice
func (itm *impl) GetAll(k interface{}) (ret []interface{}) {
	var value []*item
	var ok bool
	var i int

	itm.RLock()
	value, ok = itm.items[key(k)]
	itm.RUnlock()
	if !ok {
		return
	}
	for i = range value {
		ret = append(ret, value[i].Value)
	}
	return
}
