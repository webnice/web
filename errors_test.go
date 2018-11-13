package web

import (
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	if Errors() != errSingleton {
		t.Fatalf("Errors() function is damaged")
	}
}

func TestErrAlreadyRunning(t *testing.T) {
	var v interface{}

	if ErrAlreadyRunning() != &errAlreadyRunning {
		t.Errorf("Error func ErrAlreadyRunning()")
	}
	switch v = ErrAlreadyRunning().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cAlreadyRunning) {
			t.Fatalf("ErrAlreadyRunning() function is damaged")
		}
	default:
		t.Fatalf("Package errors is damaged")
	}
}

func TestErrNoConfiguration(t *testing.T) {
	var v interface{}

	if ErrNoConfiguration() != &errNoConfiguration {
		t.Errorf("Error func ErrNoConfiguration()")
	}
	switch v = ErrNoConfiguration().Error(); s := v.(type) {
	case string:
		if !strings.EqualFold(s, cNoConfiguration) {
			t.Fatalf("ErrNoConfiguration() function is damaged")
		}
	default:
		t.Fatalf("Package errors is damaged")
	}
}
