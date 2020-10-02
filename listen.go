package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"crypto/tls"
	"net"
	"os"
)

// Wait while web server is running
func (wsv *web) Wait() Interface { wsv.doCloseDone.Wait(); return wsv }

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections
func (wsv *web) ListenAndServe(addr string) Interface {
	var conf *Configuration

	if conf, wsv.err = parseAddress(addr); wsv.err != nil {
		return wsv
	}

	return wsv.ListenAndServeWithConfig(conf)
}

// ListenAndServeTLS listens on the TCP network address address with TLS and then calls Serve with handler
// to handle requests on incoming connections
func (wsv *web) ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface {
	var conf *Configuration

	if conf, wsv.err = parseAddress(addr); wsv.err != nil {
		return wsv
	}
	conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM = certFile, keyFile

	return wsv.ListenAndServeTLSWithConfig(conf, tlsConfig)
}

// ListenAndServeWithConfig Fully configurable web server listens and then calls Serve on incoming connections
func (wsv *web) ListenAndServeWithConfig(conf *Configuration) Interface {
	if conf == nil {
		wsv.err = ErrNoConfiguration()
		return wsv
	}
	wsv.conf = conf

	return wsv.Listen(nil)
}

// ListenAndServeTLSWithConfig Fully configurable web server listens and then calls Serve on incoming connections
func (wsv *web) ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface {
	if conf == nil {
		wsv.err = ErrNoConfiguration()
		return wsv
	}
	wsv.conf = conf
	if tlsConfig == nil {
		if tlsConfig, wsv.err = wsv.tlsConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); wsv.err != nil {
			return wsv
		}
	}

	return wsv.Listen(tlsConfig)
}

// NewListener Make new listener from web server configuration
func (wsv *web) NewListener(conf *Configuration) (ret net.Listener, err error) {
	var listeners []net.Listener

	defaultConfiguration(conf)
	switch conf.Mode {
	case netSystemd:
		if listeners, err = wsv.ListenersSystemdWithoutNames(); err != nil {
			return
		}
		if len(listeners) != 1 {
			err = ErrListenSystemdUnexpectedNumber()
			return
		}
		ret = listeners[0]
	case netUnix, netUnixPacket:
		_ = os.Remove(conf.Socket)
		ret, err = net.Listen(conf.Mode, conf.Socket)
		_ = os.Chmod(conf.Socket, os.FileMode(0666))
	default:
		ret, err = net.Listen(conf.Mode, conf.HostPort)
	}

	return
}

// NewListenerTLS Make new listener with TLS from web server configuration
func (wsv *web) NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, err error) {
	var l net.Listener

	if l, err = wsv.NewListener(conf); err != nil {
		return
	}
	if tlsConfig == nil {
		if tlsConfig, err = wsv.tlsConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); err != nil {
			return
		}
	}
	ret = tls.NewListener(l, tlsConfig)

	return
}

// Конфигурация TLS по умолчанию
func (wsv *web) tlsConfigDefault(tlsPublicFile string, tlsPrivateFile string) (ret *tls.Config, err error) {
	ret = &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: make([]tls.Certificate, 1),
	}
	if ret.Certificates[0], err = tls.LoadX509KeyPair(tlsPublicFile, tlsPrivateFile); err != nil {
		return
	}

	return
}

// Listen Begin listen port and web server serve
func (wsv *web) Listen(tlsConfig *tls.Config) Interface {
	var ltn net.Listener

	if wsv.isRun.Load().(bool) {
		wsv.err = ErrAlreadyRunning()
		return wsv
	}
	switch tlsConfig == nil {
	case true:
		ltn, wsv.err = wsv.NewListener(wsv.conf)
	case false:
		ltn, wsv.err = wsv.NewListenerTLS(wsv.conf, tlsConfig)
	}
	if wsv.err != nil {
		return wsv
	}

	return wsv.ServeTLS(ltn, tlsConfig)
}

// Serve accepts incoming connections on the listener, creating a new web server goroutine
func (wsv *web) Serve(ltn net.Listener) Interface { return wsv.ServeTLS(ltn, nil) }

// ServeTLS accepts incoming connections on the listener with TLS configuration, creating a new web server goroutine
func (wsv *web) ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface {
	var conf *Configuration

	// TODO: Реализовать поддержку PROXY Protocol через "gopkg.in/webnice/web.v1/proxyp", conf.ProxyProtocol

	if wsv.conf == nil {
		conf, _ = parseAddress(ltn.Addr().String())
		defaultConfiguration(conf)
		wsv.conf = conf
	}
	wsv.listener = ltn
	wsv.isRun.Store(true)
	wsv.doCloseDone.Add(1)
	go wsv.run(tlsConfig)

	return wsv
}

// Goroutine of the web server
func (wsv *web) run(tlsConfig *tls.Config) {
	defer wsv.doCloseDone.Done()
	defer wsv.isRun.Store(false)
	defer func() {
		if wsv.conf.Socket == "" {
			return
		}
		switch wsv.conf.Mode {
		case netSystemd:
			return
		case netUnix, netUnixPacket:
			_ = os.Remove(wsv.conf.Socket)
		}
	}()

	// Configure net/http web server
	if wsv.server = wsv.loadConfiguration(tlsConfig); wsv.err != nil {
		return
	}
	// Configure keepalive of web server
	if wsv.conf.KeepAliveDisable {
		wsv.server.SetKeepAlivesEnabled(false)
	}
	// Begin serve
	if wsv.conf.TLSPrivateKeyPEM == "" || wsv.conf.TLSPublicKeyPEM == "" {
		wsv.err = wsv.server.Serve(wsv.listener)
		return
	}
	wsv.err = wsv.server.ServeTLS(wsv.listener, wsv.conf.TLSPublicKeyPEM, wsv.conf.TLSPrivateKeyPEM)
}
