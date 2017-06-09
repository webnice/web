// +build go1.8

package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"net/http"
)

// Stop web server
func (wsv *web) Stop() {
	var ctx context.Context
	if wsv.server != nil {
		ctx = context.Background()
		if wsv.conf.ShutdownTimeout > 0 {
			ctx, _ = context.WithTimeout(ctx, wsv.conf.ShutdownTimeout)
		}
		wsv.err = wsv.server.Shutdown(ctx)
	} else if wsv.listener != nil {
		wsv.err = wsv.listener.Close()
	}
}

func (wsv *web) loadConfiguration() *http.Server {
	return &http.Server{
		Addr:              wsv.conf.HostPort,
		IdleTimeout:       wsv.conf.IdleTimeout,
		ReadHeaderTimeout: wsv.conf.ReadHeaderTimeout,
		ReadTimeout:       wsv.conf.ReadTimeout,
		WriteTimeout:      wsv.conf.WriteTimeout,
		MaxHeaderBytes:    wsv.conf.MaxHeaderBytes,
		Handler:           wsv.route,
	}
}
