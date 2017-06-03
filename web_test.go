package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	var wsv = New().(*web)
	if wsv == nil {
		t.Errorf("Error New(), return nil")
	}
	if wsv.route == nil {
		t.Errorf("Error New(), Route is nil")
	}
	if wsv.inCloseUp == nil {
		t.Errorf("Error New(), inCloseUp is nil")
	}
	if wsv.isRun.Load().(bool) {
		t.Errorf("Error New(), isRun is %v", wsv.isRun.Load().(bool))
	}
}

func TestError(t *testing.T) {
	const _TestString = `m7SqTD9K2FEstVjD2QR9`
	var err error
	var wsv = New().(*web)
	err = fmt.Errorf("%s", _TestString)
	if wsv.err = err; wsv.Error() != err {
		t.Errorf("Error function Error()")
	}
}

func TestRouteErrorsHandlers(t *testing.T) {
	var wsv = New().(*web)
	if wsv.Route() != wsv.route {
		t.Errorf("Error function Route()")
	}
	if wsv.route.Errors() != wsv.Errors() {
		t.Errorf("Error function Errors()")
	}
	if wsv.route.Handlers() != wsv.Handlers() {
		t.Errorf("Error function Handlers()")
	}
}
