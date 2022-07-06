package pprof

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/webnice/web/v3/middleware/nocache"
	"github.com/webnice/web/v3/route"
	"github.com/webnice/web/v3/status"
)

// Handler Middleware to profiling
func Handler() http.Handler {
	var rou = route.New()

	rou.Use(nocache.Handler)
	rou.Get("/", func(wr http.ResponseWriter, rq *http.Request) {
		http.Redirect(wr, rq, rq.RequestURI+"/pprof/", status.MovedPermanently)
	})
	rou.HandleFunc("/pprof", func(wr http.ResponseWriter, rq *http.Request) {
		http.Redirect(wr, rq, rq.RequestURI+"/", status.MovedPermanently)
	})
	rou.HandleFunc("/pprof/", pprof.Index)
	rou.HandleFunc("/pprof/cmdline", pprof.Cmdline)
	rou.HandleFunc("/pprof/profile", pprof.Profile)
	rou.HandleFunc("/pprof/symbol", pprof.Symbol)
	rou.HandleFunc("/pprof/trace", pprof.Trace)
	rou.Handle("/pprof/allocs", pprof.Handler("allocs"))
	rou.Handle("/pprof/block", pprof.Handler("block"))
	rou.Handle("/pprof/goroutine", pprof.Handler("goroutine"))
	rou.Handle("/pprof/heap", pprof.Handler("heap"))
	rou.Handle("/pprof/mutex", pprof.Handler("mutex"))
	rou.Handle("/pprof/threadcreate", pprof.Handler("threadcreate"))
	rou.HandleFunc("/vars", func(wr http.ResponseWriter, rq *http.Request) {
		var first = true

		wr.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = fmt.Fprintf(wr, "{\n") // nolint: errcheck
		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				_, _ = fmt.Fprintf(wr, ",\n") // nolint: errcheck
			}
			first = false
			_, _ = fmt.Fprintf(wr, "%q: %s", kv.Key, kv.Value) // nolint: errcheck
		})
		_, _ = fmt.Fprintf(wr, "\n}\n") // nolint: errcheck
	})

	return rou
}
