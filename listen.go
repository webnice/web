package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
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

// ListenAndServeWithConfig Fully configurable web server listens and then calls Serve on incoming connections
func (wsv *web) ListenAndServeWithConfig(conf *Configuration) Interface {
	if conf == nil {
		wsv.err = ErrNoConfiguration()
		return wsv
	}
	defaultConfiguration(conf)
	wsv.conf = conf
	return wsv.Listen()
}

// Listen Begin listen port and web server serve
func (wsv *web) Listen() Interface {
	var ltn net.Listener

	if wsv.isRun.Load().(bool) {
		wsv.err = ErrAlreadyRunning()
		return wsv
	}
	if wsv.conf.Mode == "unix" || wsv.conf.Mode == "unixpacket" {
		_ = os.Remove(wsv.conf.Socket)
		ltn, wsv.err = net.Listen(wsv.conf.Mode, wsv.conf.Socket)
		_ = os.Chmod(wsv.conf.Socket, os.FileMode(0666))
	} else {
		ltn, wsv.err = net.Listen(wsv.conf.Mode, wsv.conf.HostPort)
	}
	if wsv.err != nil {
		return wsv
	}
	return wsv.Serve(ltn)
}

// Serve accepts incoming connections on the listener, creating a new web server goroutine
func (wsv *web) Serve(ltn net.Listener) Interface {
	var conf *Configuration

	if wsv.conf == nil {
		conf, _ = parseAddress(ltn.Addr().String())
		defaultConfiguration(conf)
		wsv.conf = conf
	}
	wsv.listener = ltn
	wsv.isRun.Store(true)
	wsv.doCloseDone.Add(1)
	go wsv.run()
	return wsv
}

// Goroutine of the web server
func (wsv *web) run() {
	defer wsv.doCloseDone.Done()
	defer wsv.isRun.Store(false)
	defer func() {
		if wsv.conf.Socket == "" {
			return
		}
		if wsv.conf.Mode == "unix" || wsv.conf.Mode == "unixpacket" {
			_ = os.Remove(wsv.conf.Socket)
		}
	}()

	// Configure net/http web server
	wsv.server = wsv.loadConfiguration()
	if wsv.err != nil {
		return
	}

	// Configure keep alives of web server
	if wsv.conf.KeepAliveDisable {
		wsv.server.SetKeepAlivesEnabled(false)
	}
	// Begin serve
	wsv.err = wsv.server.Serve(wsv.listener)
}
