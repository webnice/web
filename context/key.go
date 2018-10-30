package context

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

// key is a value for use with context.WithValue
// It's used as a pointer so it fits in an interface{} without allocation
type key struct {
	Name string
}

// String convert type to string
func (k *key) String() string { return k.Name }
