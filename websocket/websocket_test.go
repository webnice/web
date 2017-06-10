package websocket

import (
	"testing"
)

func TestNew(t *testing.T) {
	var wst = New()
	if wst == nil {
		t.Fatalf("Error New(), return nil")
	}
}
