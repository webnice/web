package contexterror // import "gopkg.in/webnice/web.v1/context/contexterror"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/param"
import ()

// This is an inplementation
type impl struct {
	params *param.Params
}

// Interface is an interface of package
type Interface interface {
	// Add key value to heap
	Add(key string, value string)

	// Get value by key from heap
	Get(key string) string

	// Set value by key to heap
	Set(key string, value string)

	// Del Delete value by key from heap
	Del(key string) string
}
