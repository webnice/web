package web

import (
	"testing"
)

func TestErrAlreadyRunning(t *testing.T) {
	if ErrAlreadyRunning() != _ErrAlreadyRunning {
		t.Errorf("Error func ErrAlreadyRunning()")
	}
}

func TestErrNoConfiguration(t *testing.T) {
	if ErrNoConfiguration() != _ErrNoConfiguration {
		t.Errorf("Error func ErrNoConfiguration()")
	}
}
