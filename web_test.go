package web

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNew(t *testing.T) {
	var wsv = New().
		Handler(echo.New()).(*web)
	if wsv == nil {
		t.Errorf("Error New(), return nil")
	}
	if wsv.handler == nil {
		t.Errorf("Error New(), Handler is nil")
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
