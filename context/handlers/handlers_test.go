package handlers

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/webnice/web.v1/context/errors"
)

type testHandler struct {
	C int64
}

func (h *testHandler) ServeHTTP(wr http.ResponseWriter, rq *http.Request) { h.C++ }

func TestReset(t *testing.T) {
	var p1, p2 string
	var obj = New(errors.New()).(*impl)

	p1 = fmt.Sprintf("%p", obj.handlers)
	obj.Reset()
	p2 = fmt.Sprintf("%p", obj.handlers)
	if p1 == p2 {
		t.Errorf("Error Reset()")
	}
}

func TestDo(t *testing.T) {
	const testKey uint32 = (1 << 32) - 1
	var obj *impl
	var hf1, hf2 http.HandlerFunc
	var p1, p2, p3 string

	hf1 = new(testHandler).ServeHTTP
	hf2 = new(testHandler).ServeHTTP
	p1 = fmt.Sprintf("%p", hf1)
	p2 = fmt.Sprintf("%p", hf2)

	obj = New(errors.New()).(*impl)
	p3 = fmt.Sprintf("%p", obj.do(testKey, hf1, hf2))
	if p3 != p1 {
		t.Errorf("Error do(), returns object is incorrect")
	}

	p3 = fmt.Sprintf("%p", obj.do(testKey, nil, hf2))
	if p3 != p1 {
		t.Errorf("Error do(), returns object is incorrect")
	}

	obj.Reset()
	p3 = fmt.Sprintf("%p", obj.do(testKey, nil, hf2))
	if p3 != p2 {
		t.Errorf("Error do(), return object is incorrect")
	}
}

func testFn(t *testing.T, fn func(http.HandlerFunc) http.HandlerFunc, funcName string) {
	var hf1 http.HandlerFunc
	var p1, p3 string

	hf1 = new(testHandler).ServeHTTP
	p1 = fmt.Sprintf("%p", hf1)

	if fn(nil) == nil {
		t.Errorf("Error in '%s', returns nil", funcName)
	}

	p3 = fmt.Sprintf("%p", fn(hf1))
	if p3 != p1 {
		t.Errorf("Error in '%s', returns object not equal set object", funcName)
	}

	p3 = fmt.Sprintf("%p", fn(nil))
	if p3 != p1 {
		t.Errorf("Error in '%s', returns object not equal set object", funcName)
	}
}

func TestAll(t *testing.T) {
	var obj = New(errors.New())

	obj.Reset()
	testFn(t, obj.InternalServerError, "InternalServerError()")
	obj.Reset()
	testFn(t, obj.MethodNotAllowed, "MethodNotAllowed()")
	obj.Reset()
	testFn(t, obj.NotFound, "NotFound()")
}
