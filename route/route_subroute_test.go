package route

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/webnice/web/v3/context"
	"github.com/webnice/web/v3/method"
	"github.com/webnice/web/v3/status"
)

var countTestSubroute int64

func testSubrouteFunc(wr http.ResponseWriter, rq *http.Request) {
	const keyIpv4 = `address`
	var (
		ctx   context.Interface
		ipsrc string
	)

	ctx = context.New(rq)
	ipsrc = ctx.Route().Params().Get(keyIpv4)
	if ipsrc == "188.42.231.207" {
		countTestSubroute++
	}
	wr.WriteHeader(status.Ok)
}

func testSubrouteRoute(r Interface) {
	r.Subroute("/v1.0", func(r Interface) {
		r.Subroute("/info", func(r Interface) {
			r.Get("/:address", testSubrouteFunc)
		})
	})
}

func TestSubroute(t *testing.T) {
	var (
		err error
		w   Interface
		srv *httptest.Server
		cou int64
	)

	w = New()
	w.Subroute("/api", testSubrouteRoute)
	srv = httptest.NewServer(w)
	countTestSubroute = 0
	for cou = 0; cou < 10000; cou++ {
		_, _, err = testRequest(t, method.Get.String(), srv.URL+"/api/v1.0/info/188.42.231.207", &bytes.Buffer{})
		if err != nil {
			t.Fatalf("Error httptest get %s: %s", srv.URL, err.Error())
		}
	}
	if countTestSubroute != 10000 {
		t.Fatalf("Error TestSubroute, count %d", countTestSubroute)
	}
}
