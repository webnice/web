package handlers // import "gopkg.in/webnice/web.v1/context/handlers"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/status"
import (
	"net/http"
)

// Default handler for internal server errors
func (hndl *impl) defaultInternalServerError(wr http.ResponseWriter, rq *http.Request) {
	var err error
	wr.WriteHeader(status.InternalServerError)
	wr.Write(status.Bytes(status.InternalServerError))
	if err = hndl.errors.InternalServerError(nil); err != nil {
		wr.Write([]byte(err.Error()))
	}
}
