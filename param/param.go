package param

// New is a constructor
func New() Interface {
	var itm = new(impl)

	itm.items = map[key][]*item{}

	return itm
}

// Set sets the params entries associated with key to
// the single element value. It replaces any existing
// values associated with key
func (itm *impl) Set(k string, v string) {
	itm.Lock()
	itm.items[key(k)] = []*item{{Value: v}}
	itm.Unlock()
}

// Add adds the key, value pair to the params.
// It appends to any existing values associated with key
func (itm *impl) Add(k string, v string) {
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
func (itm *impl) Has(k string) (ret bool) {
	itm.RLock()
	_, ret = itm.items[key(k)]
	itm.RUnlock()
	return
}

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns "".
// To access multiple values of a key, or to use access the map directly
func (itm *impl) Get(k string) (ret string) {
	var (
		value []*item
		ok    bool
	)

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
func (itm *impl) Del(k string) (ret string) {
	var (
		value []*item
		ok    bool
	)

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

// Keys gets the all keys
// If there are no values, Keys returns empty slice
func (itm *impl) Keys() (ret []string) {
	var i key

	itm.RLock()
	for i = range itm.items {
		ret = append(ret, string(i))
	}
	itm.RUnlock()

	return
}

// Get gets the all value associated with the given key.
// If there are no values associated with the key, Get returns empty slice
func (itm *impl) GetAll(k string) (ret []string) {
	var (
		value []*item
		ok    bool
		n     int
	)

	itm.RLock()
	value, ok = itm.items[key(k)]
	itm.RUnlock()
	if !ok {
		return
	}
	for n = range value {
		ret = append(ret, value[n].Value)
	}

	return
}
