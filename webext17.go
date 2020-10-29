// +build go1.7,!go1.8

package web



import (
	"net/http"
)

// Stop web server
func (wsv *web) Stop() Interface {
	if wsv.listener != nil {
		wsv.err = wsv.listener.Close()
	}
	return wsv
}

func (wsv *web) loadConfiguration() *http.Server {
	if wsv.route.Errors().RouteConfigurationError(nil) != nil {
		wsv.err = wsv.route.Errors().RouteConfigurationError(nil)
	}
	return &http.Server{
		Addr:           wsv.conf.HostPort,
		ReadTimeout:    wsv.conf.ReadTimeout,
		WriteTimeout:   wsv.conf.WriteTimeout,
		MaxHeaderBytes: wsv.conf.MaxHeaderBytes,
		Handler:        wsv.route,
	}
}
