package verify

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

// E4xx HTTP error code 400-499
func E4xx() Interface {
	var err = new(Response)
	err.Error.Code = 4
	err.Error.Message = `Data is incorrect`
	return err
}
