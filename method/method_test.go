package method

import (
	"strings"
	"testing"
)

func TestMethods(t *testing.T) {
	if Get.String() != "GET" {
		t.Errorf("Error method String, return '%s' expected '%s'", Get.String(), "GET")
	}
	if Post.String() != "POST" {
		t.Errorf("Error method String, return '%s' expected '%s'", Post.String(), "POST")
	}
	if Put.String() != "PUT" {
		t.Errorf("Error method String, return '%s' expected '%s'", Put.String(), "PUT")
	}
	if Delete.String() != "DELETE" {
		t.Errorf("Error method String, return '%s' expected '%s'", Delete.String(), "DELETE")
	}
	if Connect.String() != "CONNECT" {
		t.Errorf("Error method String, return '%s' expected '%s'", Connect.String(), "CONNECT")
	}
	if Head.String() != "HEAD" {
		t.Errorf("Error method String, return '%s' expected '%s'", Head.String(), "HEAD")
	}
	if Patch.String() != "PATCH" {
		t.Errorf("Error method String, return '%s' expected '%s'", Patch.String(), "PATCH")
	}
	if Options.String() != "OPTIONS" {
		t.Errorf("Error method String, return '%s' expected '%s'", Options.String(), "OPTIONS")
	}
	if Trace.String() != "TRACE" {
		t.Errorf("Error method String, return '%s' expected '%s'", Trace.String(), "TRACE")
	}
	if Stub.String() != "" {
		t.Errorf("Error method String, return '%s' expected '%s'", Stub.String(), "")
	}
}

func TestInt64(t *testing.T) {
	if Get.Int64() != 1 {
		t.Errorf("Error method Int64, return '%d' expected '%d'", Get.Int64(), 1)
	}
	if Trace.Int64() != 1<<8 {
		t.Errorf("Error method Int64, return '%d' expected '%d'", Get.Int64(), 1<<8)
	}
}

func TestAll(t *testing.T) {
	var sum int64
	for _, m := range All() {
		sum += m.Int64()
	}
	if sum != Any.Int64() {
		t.Errorf("Error method All()", sum, Any.Int64())
	}
}

func TestParse(t *testing.T) {
	var err error
	var mp Method
	for _, m := range All() {
		mp, err = Parse(strings.ToTitle(m.String()))
		if err != nil {
			t.Errorf("Error method Parse('%s'), return error: %s", strings.ToTitle(m.String()), err)
		}
		if mp != m {
			t.Errorf("Error method Parse('%s'): %s", strings.ToTitle(m.String()), mp.String())
		}
	}
	_, err = Parse(`KhGJjvhgv`)
	if err != nil {
		t.Errorf("Error method Parse()")
	}
}
