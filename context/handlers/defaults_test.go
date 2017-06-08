package handlers

import "gopkg.in/webnice/web.v1/context/errors"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testDefaults(t *testing.T, handler http.HandlerFunc, hfName string, response string) {
	var err error
	var srv *httptest.Server
	var rsp *http.Response
	var buf []byte

	srv = httptest.NewServer(handler)
	defer srv.Close()

	rsp, err = http.Get(srv.URL)
	if err != nil {
		t.Errorf("Error response HandlerFunc: %s", err.Error())
	}
	defer rsp.Body.Close()

	if buf, err = ioutil.ReadAll(rsp.Body); err != nil {
		t.Errorf("Error read response: %s", err.Error())
	}
	if rsp.StatusCode == 200 {
		t.Errorf("%s error: return staus code: %d, text: %s", hfName, rsp.StatusCode, string(buf))
	}
	if response == "" {
		return
	}

	if string(buf) != response {
		t.Errorf("%s error: return staus code: %d, text: %s", hfName, rsp.StatusCode, string(buf))
	}
}

func TestDefaultsAll(t *testing.T) {
	const testError = `ce8449d43faeb80d9365a916d1e3e1931b5f684979fa772ac658c985679d81047f0cc342228808b590166fcc257079816552abc66305a2cce3e2155076a75cca`
	var obj *impl
	var err = errors.New()

	obj = New(err).(*impl)

	testDefaults(t, obj.defaultInternalServerError, "defaultInternalServerError()", "")
	obj.errors.InternalServerError(fmt.Errorf("%s", testError))
	testDefaults(t, obj.defaultInternalServerError, "defaultInternalServerError()", testError)

	testDefaults(t, obj.defaultMethodNotAllowed, "defaultMethodNotAllowed()", "")
	obj.errors.MethodNotAllowed(fmt.Errorf("%s", testError))
	testDefaults(t, obj.defaultMethodNotAllowed, "defaultMethodNotAllowed()", testError)

	testDefaults(t, obj.defaultNotFound, "defaultNotFound()", "")
	obj.errors.NotFound(fmt.Errorf("%s", testError))
	testDefaults(t, obj.defaultNotFound, "defaultNotFound()", testError)
}
