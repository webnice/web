// +build go1.7,!go1.8

package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
)

// Stop web server
func (wsv *web) Stop() {
	if wsv.listener != nil {
		wsv.err = wsv.listener.Close()
	}
}

func (wsv *web) loadConfiguration() *http.Server {
	return &http.Server{
		Addr:           wsv.conf.HostPort,
		ReadTimeout:    wsv.conf.ReadTimeout,
		WriteTimeout:   wsv.conf.WriteTimeout,
		MaxHeaderBytes: wsv.conf.MaxHeaderBytes,
		Handler:        wsv.route,
	}
}
