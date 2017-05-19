package method

import (
	"testing"
)

func TestAll(t *testing.T) {
	var v Method

	v = Get
	if v.String() != "GET" {
		t.Errorf("Error in String method")
	}
}
