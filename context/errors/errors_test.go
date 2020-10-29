package errors

import (
	"fmt"
	"testing"
)

func TestReset(t *testing.T) {
	var (
		p1, p2 string
		obj    = New().(*impl)
	)

	p1 = fmt.Sprintf("%p", obj.errors)
	obj.Reset()
	p2 = fmt.Sprintf("%p", obj.errors)
	if p1 == p2 {
		t.Errorf("Error Reset()")
	}
}

func TestDo(t *testing.T) {
	const testKey uint32 = (1 << 32) - 1
	var (
		obj *impl
		err = fmt.Errorf("test error value")
	)

	obj = New().(*impl)
	if obj.do(testKey, err) != err {
		t.Errorf("Error do(), returns object is incorrect")
	}
	if obj.do(testKey, nil) != err {
		t.Errorf("Error do(), returns object is incorrect")
	}
	obj.Reset()
	if obj.do(testKey, nil) != nil {
		t.Errorf("Error do(), return object is incorrect")
	}
}

func testFn(t *testing.T, fn func(error) error, funcName string) {
	var err = fmt.Errorf("Test error value")

	if fn(nil) != nil {
		t.Errorf("Error in '%s', returns not nil", funcName)
	}
	if fn(err) != err {
		t.Errorf("Error in '%s', returns object not equal set object", funcName)
	}
	if fn(nil) != err {
		t.Errorf("Error in '%s', returns object not equal set object", funcName)
	}
}

func TestAll(t *testing.T) {
	var obj = New()

	obj.Reset()
	testFn(t, obj.RouteConfigurationError, "RouteConfigurationError()")
	obj.Reset()
	testFn(t, obj.InternalServerError, "InternalServerError()")
	obj.Reset()
	testFn(t, obj.MethodNotAllowed, "MethodNotAllowed()")
	obj.Reset()
	testFn(t, obj.NotFound, "NotFound()")
}
