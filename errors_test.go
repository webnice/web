package web

import (
	"testing"
)

func TestErrAlreadyRunning(t *testing.T) {
	if ErrAlreadyRunning() != errAlreadyRunning {
		t.Errorf("Error func ErrAlreadyRunning()")
	}
}

func TestErrNoConfiguration(t *testing.T) {
	if ErrNoConfiguration() != errNoConfiguration {
		t.Errorf("Error func ErrNoConfiguration()")
	}
}
