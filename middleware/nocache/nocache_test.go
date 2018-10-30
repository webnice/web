package nocache

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

const (
	ifModifiedSinceTimeFormat = `Mon, 02 Jan 2006 15:04:05 GMT`
)

func testNoCacheHandlerFunc(wr http.ResponseWriter, rq *http.Request) {
	wr.Header().Add(header.ETag, "686897696a7c876b7e")
	wr.Header().Set(header.LastModified, time.Now().UTC().Format(ifModifiedSinceTimeFormat))
	wr.Header().Set(header.IfUnmodifiedSince, time.Now().UTC().Format(ifModifiedSinceTimeFormat))
	fmt.Fprint(wr, string(status.Bytes(status.Ok)))
	wr.WriteHeader(status.Ok)
}

func testNoCache(t *testing.T, rou route.Interface) (err error) {
	var srv *httptest.Server
	var rsp *http.Response
	var buf []byte

	srv = httptest.NewServer(http.HandlerFunc(testNoCacheHandlerFunc))
	if rou != nil {
		srv.Config.Handler = rou
	}
	defer srv.Close()

	rsp, err = http.Get(srv.URL)
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
	if rsp.Header.Get("Pragma") != "no-cache" {
		err = fmt.Errorf("Error server header")
		return
	}
	if rsp.Header.Get("Cache-Control") != "no-cache, private, max-age=0" {
		err = fmt.Errorf("Error server header")
		return
	}
	if rsp.Header.Get("X-Accel-Expires") != "0" {
		err = fmt.Errorf("Error server header")
		return
	}
	if rsp.Header.Get(header.ETag) != "" {
		err = fmt.Errorf("Error server header")
		return
	}
	if rsp.Header.Get(header.IfUnmodifiedSince) != "" {
		err = fmt.Errorf("Error server header")
		return
	}
	if string(buf) != string(status.Bytes(status.Ok)) {
		err = fmt.Errorf("Error response body")
		return
	}

	return
}

func TestNoCache(t *testing.T) {
	var rou route.Interface
	var err error

	rou = route.New()
	rou.Use(NoCache)
	rou.Get("/", testNoCacheHandlerFunc)
	if err = testNoCache(t, rou); err != nil {
		t.Errorf("Error nocache: %v", err)
	}
}
