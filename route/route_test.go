package route

import "gopkg.in/webnice/web.v1/context"
import (
	"testing"
)

func TestNew(t *testing.T) {
	var ok bool
	var rou = New().(*impl)

	if rou.tree == nil {
		t.Errorf("Error New(), tree is nil")
	}
	if rou.context == nil {
		t.Errorf("Error New(), context is nil")
	}
	if rou.parent != nil {
		t.Errorf("Error New(), parent is not nil")
	}
	if _, ok = rou.pool.Get().(context.Interface); !ok {
		t.Errorf("Error New(), error in sync.Pool")
	}
}

func TestErrorsHandlers(t *testing.T) {
	var r Interface

	r = New()
	if r.Errors() == nil {
		t.Errorf("Errors() returns nil")
	}
	if r.Handlers() == nil {
		t.Errorf("Handlers() returns nil")
	}
}
