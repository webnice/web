package route

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/webnice/web/v2/status"
)

var testMiddlewareCount int64

func testMiddlewareCounter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		testMiddlewareCount++
		next.ServeHTTP(wr, rq)
	})
}

func TestUse(t *testing.T) {
	var (
		err error
		srv *httptest.Server
		rsp *http.Response
		buf *bytes.Buffer
		r   = New()
	)
	var hf = func(wr http.ResponseWriter, rq *http.Request) {
		wr.WriteHeader(status.Ok)
		_, _ = wr.Write(status.Bytes(status.Ok))
	}

	// Correct middlewares
	r.Use(testMiddlewareCounter)
	r.Get("/", hf)
	srv = httptest.NewServer(r)
	rsp, buf, err = testRequest(t, "GET", srv.URL, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Request error: %s", err.Error())
	}
	if testMiddlewareCount != 1 {
		t.Errorf("Error middlewares")
	}
	if rsp.StatusCode != status.Ok || buf.String() != string(status.Bytes(status.Ok)) {
		t.Errorf("Error in handlefunc")
	}
	srv.Close()
	// Incorrect call use
	testMiddlewareCount = 0
	r = New()
	r.Get("/", hf)
	r.Use(testMiddlewareCounter)
	srv = httptest.NewServer(r)
	rsp, buf, err = testRequest(t, "GET", srv.URL, &bytes.Buffer{})
	if err != nil {
		t.Fatalf("Request error: %s", err.Error())
	}
	if testMiddlewareCount != 0 {
		t.Errorf("Error middlewares")
	}
	if rsp.StatusCode != status.InternalServerError || !strings.Contains(buf.String(), "middlewares must be defined before use") {
		t.Errorf("Error in handlefunc")
	}
	srv.Close()
}
