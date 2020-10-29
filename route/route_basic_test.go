package route

import (
	"bufio"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/webnice/web/v2/context"
	"github.com/webnice/web/v2/status"
)

func testRequest(t *testing.T, method string, path string, body *bytes.Buffer) (rsp *http.Response, ret *bytes.Buffer, err error) {
	var (
		req *http.Request
		buf *bytes.Buffer
		tmp [][]byte
	)

	ret = &bytes.Buffer{}
	if req, err = http.NewRequest(method, path, body); err != nil {
		return
	}
	if rsp, err = http.DefaultClient.Do(req); err != nil {
		return
	}
	defer func() { _ = rsp.Body.Close() }()
	buf = &bytes.Buffer{}
	if err = rsp.Write(buf); err != nil {
		return
	}
	if rsp, err = http.ReadResponse(bufio.NewReader(buf), req); err != nil {
		return
	}
	buf.Reset()
	if err = rsp.Write(buf); err != nil {
		return
	}
	if tmp = bytes.SplitN(buf.Bytes(), []byte{'\r', '\n', '\r', '\n'}, 2); len(tmp) > 1 {
		ret = bytes.NewBuffer(tmp[1])
	} else if tmp = bytes.SplitN(buf.Bytes(), []byte{'\n', '\n'}, 2); len(tmp) > 1 {
		ret = bytes.NewBuffer(tmp[1])
	} else {
		ret = bytes.NewBuffer(tmp[0])
	}

	return
}

func TestNew(t *testing.T) {
	var (
		ok  bool
		rou = New().(*impl)
	)

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
	var (
		err error
		w1  Interface
		srv *httptest.Server
		rsp *http.Response
		buf *bytes.Buffer
	)

	// Empty handler
	w1 = New()
	srv = httptest.NewServer(w1)
	rsp, buf, err = testRequest(t, "GET", srv.URL, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	if rsp.StatusCode != status.InternalServerError {
		t.Errorf("Error ServeHTTP(), running server without handlers, request must returned internal server error")
	}
	if !strings.Contains(buf.String(), "route with no handlers") {
		t.Errorf("Error ServeHTTP(), returns incorrect error description: %q", buf.String())
	}
	srv.Close()
	// Incorrect routing
	w1 = New()
	w1.Get("/", testServeHTTP)
	w1.Post("", testServeHTTP) // Error URI (path must begin with '/')
	srv = httptest.NewServer(w1)
	rsp, buf, err = testRequest(t, "GET", srv.URL, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	if rsp.StatusCode != status.InternalServerError {
		t.Errorf("Error ServeHTTP(), running server without handlers, request must returned internal server error")
	}
	if !strings.Contains(buf.String(), "must begin with '/'") {
		t.Errorf("Error ServeHTTP(), returns incorrect error description: %q", buf.String())
	}
	srv.Close()
	// Correct handler
	w1 = New()
	w1.Get("/", testServeHTTP)
	srv = httptest.NewServer(w1)
	rsp, _, err = testRequest(t, "GET", srv.URL, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Error httptest get %s: %s", srv.URL, err.Error())
	}
	if rsp.StatusCode != status.NoContent {
		t.Errorf("Error test ServeHTTP(), returns %d expected %d", rsp.StatusCode, status.NoContent)
	}
	srv.Close()
}
