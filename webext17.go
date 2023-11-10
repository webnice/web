//go:build go1.7 && !go1.8
// +build go1.7,!go1.8

package web

import "net/http"

// Stop web server
func (wbo *web) Stop() Interface {
	if wbo.listener != nil {
		wbo.err = wbo.listener.Close()
	}
	return wbo
}

func (wbo *web) loadConfiguration() *http.Server {
	if wbo.route.Errors().RouteConfigurationError(nil) != nil {
		wbo.err = wbo.route.Errors().RouteConfigurationError(nil)
	}
	return &http.Server{
		Addr:           wbo.conf.HostPort,
		ReadTimeout:    wbo.conf.ReadTimeout,
		WriteTimeout:   wbo.conf.WriteTimeout,
		MaxHeaderBytes: wbo.conf.MaxHeaderBytes,
		Handler:        wbo.route,
	}
}
