//go:build go1.8
// +build go1.8

package web

import (
	"context"
	"crypto/tls"
	"net/http"
)

// Stop Отправка сигнала прерывания работы веб сервера с учётом значения ShutdownTimeout.
func (wbo *web) Stop() Interface {
	var (
		ctx context.Context
		cfn context.CancelFunc
	)

	if wbo.server != nil {
		ctx = context.Background()
		if wbo.conf.ShutdownTimeout > 0 {
			ctx, cfn = context.WithTimeout(ctx, wbo.conf.ShutdownTimeout)
			defer cfn()
		}
		wbo.err = wbo.server.Shutdown(ctx)
	} else if wbo.listener != nil {
		wbo.err = wbo.listener.Close()
	}

	return wbo
}

func (wbo *web) loadConfiguration(tlsConfig *tls.Config) (srv *http.Server) {
	srv = &http.Server{
		Addr:              wbo.conf.HostPort,
		Handler:           wbo.handler,
		ReadTimeout:       wbo.conf.ReadTimeout,
		ReadHeaderTimeout: wbo.conf.ReadHeaderTimeout,
		WriteTimeout:      wbo.conf.WriteTimeout,
		IdleTimeout:       wbo.conf.IdleTimeout,
		MaxHeaderBytes:    wbo.conf.MaxHeaderBytes,
	}
	if wbo.conf.TLSPrivateKeyPEM == "" || wbo.conf.TLSPublicKeyPEM == "" {
		return
	}
	// TLS конфигурация по умолчанию
	if tlsConfig == nil {
		tlsConfig, wbo.err = wbo.tlsConfigDefault(wbo.conf.TLSPublicKeyPEM, wbo.conf.TLSPrivateKeyPEM)
		if wbo.err != nil {
			return
		}
	}
	srv.TLSConfig, srv.TLSNextProto = tlsConfig, make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

	return
}
