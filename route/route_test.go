package route

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context"
import "gopkg.in/webnice/web.v1/status"
import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var ok bool
	var rou = New().(*impl)

	if rou.tree == nil {
		t.Errorf("Error New(), tree is nil")
	}
	if rou.context == nil {
		t.Errorf("Error New(), context is nil")
	}
	if rou.parent != nil {
		t.Errorf("Error New(), parent is not nil")
	}
	if _, ok = rou.pool.Get().(context.Interface); !ok {
		t.Errorf("Error New(), error in sync.Pool")
	}
}

func TestErrorsHandlers(t *testing.T) {
	var r = New()
	if r.Errors() == nil {
		t.Errorf("Errors() returns nil")
	}
	if r.Handlers() == nil {
		t.Errorf("Handlers() returns nil")
	}
}

func TestSetErrors(t *testing.T) {
	const (
		testErrorString = `dm3T36w7jKnQvG74Gzm6y74yMBZaCnVeyvfzEMR97PF8wHCs9KvuBzEwHjBVXN4T`
	)
	var r = New().(*impl)

	r.setErrors()
	if r.context.Errors().RouteConfigurationError(nil) != nil {
		t.Errorf("Error in setErrors()")
	}
	r.errors = append(r.errors, errors.New(testErrorString))

	r.setErrors()
	if r.context.Errors().RouteConfigurationError(nil) == nil {
		t.Errorf("Error in setErrors()")
	}
	if !strings.Contains(r.context.Errors().RouteConfigurationError(nil).Error(), testErrorString) {
		t.Errorf("Error in setErrors()")
	}
}

func testServeHTTP(wr http.ResponseWriter, rq *http.Request) {
	if context.IsContext(rq) {
		wr.WriteHeader(status.NoContent)
	} else {
		wr.WriteHeader(status.NotFound)
	}
}

func TestServeHTTP(t *testing.T) {
	var err error
	var w1 Interface
	var srv *httptest.Server
	var rsp *http.Response
	var buf *bytes.Buffer

	w1 = New()
	srv = httptest.NewServer(w1)
	if rsp, err = http.Get(srv.URL); err != nil {
		t.Errorf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	buf = bytes.NewBufferString(``)
	_ = rsp.Write(buf)
	if rsp.StatusCode != status.InternalServerError {
		t.Errorf("Error ServeHTTP(), running server without handlers, request must returned internal server error")
	}
	if !strings.Contains(buf.String(), "route with no handlers") {
		t.Errorf("Error ServeHTTP(), returns incorrect error description: %q", buf.String())
	}
	_ = rsp.Body.Close()
	srv.Close()

	// Add incorrect handler
	w1 = New()
	w1.Get("/", testServeHTTP)
	w1.Post("", testServeHTTP)
	srv = httptest.NewServer(w1)
	if rsp, err = http.Get(srv.URL); err != nil {
		t.Errorf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	buf = bytes.NewBufferString(``)
	_ = rsp.Write(buf)
	if rsp.StatusCode != status.InternalServerError {
		t.Errorf("Error ServeHTTP(), running server without handlers, request must returned internal server error")
	}
	if !strings.Contains(buf.String(), "must begin with '/'") {
		t.Errorf("Error ServeHTTP(), returns incorrect error description: %q", buf.String())
	}
	_ = rsp.Body.Close()
	srv.Close()

	// Add handler
	w1 = New()
	w1.Get("/", testServeHTTP)
	srv = httptest.NewServer(w1)
	if rsp, err = http.Get(srv.URL); err != nil {
		t.Errorf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	buf = bytes.NewBufferString(``)
	_ = rsp.Write(buf)
	if rsp.StatusCode != status.NoContent {
		t.Errorf("Error test ServeHTTP(), returns %d expected %d", rsp.StatusCode, status.NoContent)
	}
	_ = rsp.Body.Close()
	srv.Close()
}
