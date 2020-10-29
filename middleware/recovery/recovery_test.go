package recovery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/webnice/web/v2/route"
)

const (
	testPanicString = `8b56addac65f178e35a0fb560ac910e281a23c39fd6302034e90b03816a0924ab54ff3229d739ae098bdacf8fbfe55f61a08efb5ee5e74c5c5ccd519c9e15318`
)

func testPanic(wr http.ResponseWriter, rq *http.Request) {
	panic(testPanicString)
}

func testRecover(t *testing.T, rou route.Interface) (err error) {
	var (
		srv *httptest.Server
		rsp *http.Response
		buf []byte
	)

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("Panic fail. Catch: %s", e.(error))
		}
	}()
	srv = httptest.NewServer(http.HandlerFunc(testPanic))
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
	if rsp.StatusCode != 500 {
		err = fmt.Errorf("Error staus code: %d, text: %s", rsp.StatusCode, string(buf))
		return
	}
	if !strings.Contains(string(buf), testPanicString) {
		err = fmt.Errorf("error, response not contains test panic string")
		return
	}

	return
}

func TestRecover(t *testing.T) {
	var (
		rou route.Interface
		err error
	)

	rou = route.New()
	rou.Use(Handler)
	rou.HandleFunc("/", testPanic)
	if err = testRecover(t, rou); err != nil {
		t.Errorf("Error Recover(): %v", err)
	}
}
