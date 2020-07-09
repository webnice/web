// +build go1.8

package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"crypto/tls"
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

func (wsv *web) loadConfiguration(tlsConfig *tls.Config) (srv *http.Server) {
	if wsv.route.Errors().RouteConfigurationError(nil) != nil {
		wsv.err = wsv.route.Errors().RouteConfigurationError(nil)
	}
	srv = &http.Server{
		Addr:              wsv.conf.HostPort,
		IdleTimeout:       wsv.conf.IdleTimeout,
		ReadHeaderTimeout: wsv.conf.ReadHeaderTimeout,
		ReadTimeout:       wsv.conf.ReadTimeout,
		WriteTimeout:      wsv.conf.WriteTimeout,
		MaxHeaderBytes:    wsv.conf.MaxHeaderBytes,
		Handler:           wsv.route,
	}
	if wsv.conf.TLSPrivateKeyPEM == "" || wsv.conf.TLSPublicKeyPEM == "" {
		return
	}
	// TLS конфигурация по умолчанию
	if tlsConfig == nil {
		tlsConfig, wsv.err = wsv.tlsConfigDefault(wsv.conf.TLSPublicKeyPEM, wsv.conf.TLSPrivateKeyPEM)
		if wsv.err != nil {
			return
		}
	}
	srv.TLSConfig, srv.TLSNextProto = tlsConfig, make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	return
}
