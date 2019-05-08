package pprof

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/webnice/web.v1/route"
)

const testTitleString = `<title>/debug/pprof/</title>`

func testPprofSubPath(path string, t *testing.T, rou route.Interface) (err error) {
	var srv *httptest.Server
	var rsp *http.Response
	var buf []byte

	srv = httptest.NewServer(http.HandlerFunc(Handler().ServeHTTP))
	if rou != nil {
		srv.Config.Handler = rou
	}
	defer srv.Close()

	rsp, err = http.Get(srv.URL + path)
	if err != nil {
		err = fmt.Errorf("Error response HandlerFunc: %s", err)
		return
	}
	defer func() { _ = rsp.Body.Close() }()

	if buf, err = ioutil.ReadAll(rsp.Body); err != nil {
		err = fmt.Errorf("Error read response: %s", err)
		return
	}
	if rsp.StatusCode != 200 {
		err = fmt.Errorf("Error staus code: %d, text: %s", rsp.StatusCode, string(buf))
		return
	}

	return
}

func testPprof(t *testing.T, rou route.Interface) (err error) {
	var srv *httptest.Server
	var rsp *http.Response
	var buf []byte

	srv = httptest.NewServer(http.HandlerFunc(Handler().ServeHTTP))
	if rou != nil {
		srv.Config.Handler = rou
	}
	defer srv.Close()

	rsp, err = http.Get(srv.URL + `/debug/pprof/`)
	if err != nil {
		err = fmt.Errorf("Error response HandlerFunc: %s", err)
		return
	}
	defer func() { _ = rsp.Body.Close() }()

	if buf, err = ioutil.ReadAll(rsp.Body); err != nil {
		err = fmt.Errorf("Error read response: %s", err)
		return
	}
	if rsp.StatusCode != 200 {
		err = fmt.Errorf("Error staus code: %d, text: %s", rsp.StatusCode, string(buf))
		return
	}
	if !strings.Contains(string(buf), testTitleString) {
		err = fmt.Errorf("Error, response not contains test string")
		return
	}

	return
}

func TestPprof(t *testing.T) {
	var rou route.Interface
	var err error

	rou = route.New()
	rou.Mount("/debug", Handler())
	if err = testPprof(t, rou); err != nil {
		t.Errorf("Error pprof: %v", err)
	}
	if err = testPprofSubPath(`/debug`, t, rou); err != nil {
		t.Errorf("Error pprof subvars: %v", err)
	}
	if err = testPprofSubPath(`/debug/pprof`, t, rou); err != nil {
		t.Errorf("Error pprof subvars: %v", err)
	}
	if err = testPprofSubPath(`/debug/vars`, t, rou); err != nil {
		t.Errorf("Error pprof vars: %v", err)
	}
}
