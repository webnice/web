package route

import "gopkg.in/webnice/web.v1/context"
import (
	"errors"
	"strings"
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
	var r = New()
	if r.Errors() == nil {
		t.Errorf("Errors() returns nil")
	}
	if r.Handlers() == nil {
		t.Errorf("Handlers() returns nil")
	}
}

func TestSetErrors(t *testing.T) {
	const (
		testErrorString = `dm3T36w7jKnQvG74Gzm6y74yMBZaCnVeyvfzEMR97PF8wHCs9KvuBzEwHjBVXN4T`
	)
	var r = New().(*impl)

	r.setErrors()
	if r.context.Errors().RouteConfigurationError(nil) != nil {
		t.Errorf("Error in setErrors()")
	}
	r.errors = append(r.errors, errors.New(testErrorString))

	r.setErrors()
	if r.context.Errors().RouteConfigurationError(nil) == nil {
		t.Errorf("Error in setErrors()")
	}
	if !strings.Contains(r.context.Errors().RouteConfigurationError(nil).Error(), testErrorString) {
		t.Errorf("Error in setErrors()")
	}
}
