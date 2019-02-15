package contentTypeDefault

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/webnice/web.v1/header"
	"gopkg.in/webnice/web.v1/mime"
	"gopkg.in/webnice/web.v1/route"
	"gopkg.in/webnice/web.v1/status"
)

func testContentTypeDefaultHandlerFunc1(wr http.ResponseWriter, rq *http.Request) {
	wr.WriteHeader(status.Ok)
	_, _ = fmt.Fprint(wr, string(status.Bytes(status.Ok)))
	wr.WriteHeader(status.Forbidden)
}

func testContentTypeDefaultHandlerFunc2(wr http.ResponseWriter, rq *http.Request) {
	_, _ = fmt.Fprint(wr, string("ok"))
}

func testContentTypeDefault(t *testing.T, rou route.Interface, hf http.HandlerFunc) (contentType string, err error) {
	var srv *httptest.Server
	var rsp *http.Response
	var buf []byte

	srv = httptest.NewServer(hf)
	if rou != nil {
		srv.Config.Handler = rou
	}
	defer srv.Close()

	rsp, err = http.Get(srv.URL)
	if err != nil {
		err = fmt.Errorf("response HandlerFunc error: %s", err)
		return
	}
	defer rsp.Body.Close()

	if buf, err = ioutil.ReadAll(rsp.Body); err != nil {
		err = fmt.Errorf("read response error: %s", err)
		return
	}
	if rsp.StatusCode != 200 {
		err = fmt.Errorf("response staus code: %d, text: %s", rsp.StatusCode, string(buf))
		return
	}
	contentType = rsp.Header.Get(header.ContentType)

	return
}

func TestDefaultContentType(t *testing.T) {
	var rou route.Interface
	var err error
	var ctd Interface
	var ctv string

	rou = route.New()
	ctd = New(mime.TextRfc822Headers)
	rou.Use(ctd.DefaultContentType)
	rou.Get("/", testContentTypeDefaultHandlerFunc1)
	ctv, err = testContentTypeDefault(t, rou, testContentTypeDefaultHandlerFunc1)
	if err != nil || ctv != mime.TextRfc822Headers {
		t.Errorf("error: %v", err)
	}

	rou = route.New()
	ctd = New(mime.ApplicationMsgpack)
	rou.Use(ctd.DefaultContentType)
	rou.Get("/", testContentTypeDefaultHandlerFunc1)
	ctv, err = testContentTypeDefault(t, rou, testContentTypeDefaultHandlerFunc1)
	if err != nil || ctv != mime.ApplicationMsgpack {
		t.Errorf("error: %v", err)
	}

	rou = route.New()
	ctd = New(mime.ImageXICON)
	rou.Use(ctd.DefaultContentType)
	rou.Get("/", testContentTypeDefaultHandlerFunc2)
	ctv, err = testContentTypeDefault(t, rou, testContentTypeDefaultHandlerFunc2)
	if err != nil || ctv != mime.ImageXICON {
		t.Errorf("error: %v", err)
	}
}
