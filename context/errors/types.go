package errors

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

const (
	_InternalServerError uint32 = iota
)

type item struct {
	Key   uint32
	Value error
}

type errors []item

// This is an inplementation
type impl struct {
	errors *errors
}

// Interface is an interface of package
type Interface interface {
	// Reset all stored errors
	Reset()

	// InternalServerError Set description about internal server error
	InternalServerError(err error) error
}
