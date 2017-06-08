package route

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	var rt = New()

	if rt == nil {
		t.Errorf("Erorr New(), returns nil")
	}
	if rt.Params() == nil {
		t.Errorf("Erorr New(), object is not contains Params()")
	}
	testEmpty(t, rt)
}

func testEmpty(t *testing.T, rt Interface) {
	if rt.Path() != "" {
		t.Errorf("Erorr New(), Path() is not ''")
	}
	if rt.Pattern() != "" {
		t.Errorf("Erorr New(), Pattern() is not ''")
	}
	if len(rt.Patterns()) > 0 {
		t.Errorf("Erorr New(), Patterns not empty")
	}
}

func TestAll(t *testing.T) {
	var (
		testPaths    = []string{"/a", "/b", "/c", "/d"}
		testPatterns = []string{"*a", ":b"}
	)
	var rsp1 string
	var rsp2 []string
	var rt = New()

	rsp1 = rt.Path(testPaths...)
	if rsp1 != rt.Path() {
		t.Errorf("Erorr Path()")
	}
	if rsp1 != strings.Join(testPaths, ``) {
		t.Errorf("Erorr Path(), incorrest result")
	}

	rsp1 = rt.Pattern(testPatterns...)
	if rsp1 != rt.Pattern() {
		t.Errorf("Erorr Pattern()")
	}
	if rsp1 != strings.Join(testPatterns, ``) {
		t.Errorf("Erorr Pattern(), incorrest result")
	}

	rsp2 = rt.Patterns(testPaths, testPatterns)
	if len(rsp2) != len(rt.Patterns()) {
		t.Errorf("Erorr Patterns()")
	}
	if strings.Join(rsp2, ``) != strings.Join(testPaths, ``)+strings.Join(testPatterns, ``) {
		t.Errorf("Erorr Patterns(), incorrest result")
	}

	rt.Reset()
	testEmpty(t, rt)
}
