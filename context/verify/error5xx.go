package verify

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

// E5xx HTTP error code 500-599
func E5xx() Interface {
	var err = new(Response)
	err.Error.Code = -1
	err.Error.Message = `Internal server error`
	return err
}
