package param // import "gopkg.in/webnice/web.v1/param"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// Param structure
type Param struct {
	Key   string
	Value string
}

type Params []Param

// Add key value param to heap
func (ps *Params) Add(key string, value string) {
	*ps = append(*ps, Param{key, value})
}

// Get value by key from heap
func (ps Params) Get(key string) (ret string) {
	var p Param
	for _, p = range ps {
		if p.Key == key {
			ret = p.Value
		}
	}
	return
}

// Set value by key to heap
func (ps *Params) Set(key string, value string) {
	var p Param
	var i, idx int
	idx = -1
	for i, p = range *ps {
		if p.Key == key {
			idx = i
			break
		}
	}
	if idx < 0 {
		(*ps).Add(key, value)
	} else {
		(*ps)[idx] = Param{key, value}
	}
}

// Del Delete value by key from heap
func (ps *Params) Del(key string) (ret string) {
	var p Param
	var i int
	for i, p = range *ps {
		if p.Key == key {
			*ps = append((*ps)[:i], (*ps)[i+1:]...)
			ret = p.Value
		}
	}
	return
}
