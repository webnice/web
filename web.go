package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/context/errors"
import "gopkg.in/webnice/web.v1/context/handlers"
import (
	"fmt"
	"net"
	"os"

	"gopkg.in/webnice/web.v1/route"
)

// New is a constructor of new web server implementation
func New() Interface {
	var wsv = new(web)
	wsv.route = route.New()
	wsv.inCloseUp = make(chan bool, 1)
	wsv.isRun.Store(false)
	return wsv
}

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections
func (wsv *web) ListenAndServe(addr string) {
	var conf *Configuration
	if wsv.isRun.Load().(bool) {
		wsv.err = ErrAlreadyRunning()
		return
	}
	if conf, wsv.err = parseAddress(addr); wsv.err != nil {
		return
	}
	wsv.ListenAndServeWithConfig(conf)
}

// ListenAndServeWithConfig Fully configurable web server listens and then calls Serve on incoming connections
func (wsv *web) ListenAndServeWithConfig(conf *Configuration) {
	if wsv.isRun.Load().(bool) {
		wsv.err = ErrAlreadyRunning()
		return
	}
	if conf == nil {
		wsv.err = ErrNoConfiguration()
		return
	}
	defaultConfiguration(conf)
	wsv.conf = conf
	wsv.Listen()
}

// Listen Begin listen port and web server serve
func (wsv *web) Listen() {
	var ltn net.Listener
	if wsv.conf.Mode == "unix" || wsv.conf.Mode == "unixpacket" {
		ltn, wsv.err = net.Listen(wsv.conf.Mode, wsv.conf.Socket)
	} else {
		ltn, wsv.err = net.Listen(wsv.conf.Mode, wsv.conf.HostPort)
	}
	if wsv.err != nil {
		return
	}
	wsv.Serve(ltn)
}

// Serve accepts incoming connections on the listener, creating a new web server goroutine
func (wsv *web) Serve(ltn net.Listener) {
	var conf *Configuration
	if wsv.isRun.Load().(bool) {
		wsv.err = ErrAlreadyRunning()
		return
	}
	if wsv.conf == nil {
		conf, _ = parseAddress(ltn.Addr().String())
		defaultConfiguration(conf)
		wsv.conf = conf
	}
	wsv.listener = ltn
	wsv.isRun.Store(true)
	wsv.doCloseDone.Add(1)
	go wsv.run()
}

// Wait while web server is running
func (wsv *web) Wait() { wsv.doCloseDone.Wait() }

// Error Return last error of web server
func (wsv *web) Error() error { return wsv.err }

// Route interface
func (wsv *web) Route() route.Interface { return wsv.route }

// Errors interface
func (wsv *web) Errors() errors.Interface { return wsv.route.Errors() }

// Handlers interface
func (wsv *web) Handlers() handlers.Interface { return wsv.route.Handlers() }

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
	// Catch panic in the wrong configured web server
	defer func() {
		if e := recover(); e != nil {
			wsv.err = fmt.Errorf("Panic in the web sserver: %s", e.(error).Error())
		}
	}()

	// Configure net/http web server
	wsv.server = wsv.loadConfiguration()

	// Configure keep alives of web server
	if wsv.conf.KeepAliveDisable {
		wsv.server.SetKeepAlivesEnabled(false)
	}
	// Begin serve
	wsv.err = wsv.server.Serve(wsv.listener)
}
