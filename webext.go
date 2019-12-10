// +build go1.8

package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"net/http"
)

// Stop web server
func (wsv *web) Stop() Interface {
	var (
		ctx context.Context
		cfn context.CancelFunc
	)

	if wsv.server != nil {
		ctx = context.Background()
		if wsv.conf.ShutdownTimeout > 0 {
			ctx, cfn = context.WithTimeout(ctx, wsv.conf.ShutdownTimeout)
			defer cfn()
		}
		wsv.err = wsv.server.Shutdown(ctx)
	} else if wsv.listener != nil {
		wsv.err = wsv.listener.Close()
	}

	return wsv
}

func (wsv *web) loadConfiguration() *http.Server {
	if wsv.route.Errors().RouteConfigurationError(nil) != nil {
		wsv.err = wsv.route.Errors().RouteConfigurationError(nil)
	}
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
