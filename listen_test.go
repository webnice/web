// +build !race

package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestInvalidPort(t *testing.T) {
	const invalidAddress = `:170000`
	var wsv = New()
	wsv.ListenAndServe(invalidAddress)
	if wsv.Error() == nil {
		t.Errorf("Error ListenAndServe(), don't cheack listen address")
	}
}

func TestAlreadyRunningError(t *testing.T) {
	const (
		testAddress1 = `localhost:18080`
		testAddress2 = `localhost:18081`
	)
	var wsv = New()
	wsv.ListenAndServe(testAddress1)
	defer wsv.Stop()
	if wsv.Error() != nil {
		t.Errorf("Error ListenAndServe(), error: %s", wsv.Error().Error())
	}
	wsv.ListenAndServe(testAddress2)
	if wsv.Error() == nil {
		t.Errorf("Error ListenAndServe(), do not check already running")
	}
	if wsv.Error() != ErrAlreadyRunning() {
		t.Errorf("Error ListenAndServe(), incorrect error")
	}
}

func TestNoConfigurationError(t *testing.T) {
	var wsv = New()

	wsv.ListenAndServeWithConfig(nil)
	defer wsv.Stop()
	if wsv.Error() == nil {
		t.Errorf("Error ListenAndServe(), do not check configuration")
	}
	if wsv.Error() != ErrNoConfiguration() {
		t.Errorf("Error ListenAndServe(), incorrect error")
	}
}

func TestPortIsBusy(t *testing.T) {
	const testAddress1 = `localhost:18080`
	var w1, w2 Interface

	w1 = New()
	w1.ListenAndServe(testAddress1)
	defer w1.Stop()
	if w1.Error() != nil {
		t.Errorf("Error ListenAndServe(): %s", w1.Error().Error())
	}

	w2 = New()
	w2.ListenAndServe(testAddress1)
	defer w2.Stop()
	if w2.Error() == nil {
		t.Errorf("Error ListenAndServe(), error is nil, but port busy another program")
	}
}

func TestUnixSocket(t *testing.T) {
	const testAddress1 = `.test.socket`
	var err error
	var conf *Configuration
	var w1 Interface
	var fi os.FileInfo

	conf, _ = parseAddress("")
	conf.Mode = "socket"
	conf.Socket = testAddress1
	conf.KeepAliveDisable = true
	w1 = New()
	w1.ListenAndServeWithConfig(conf)
	if w1.Error() != nil {
		t.Errorf("Error listen unix socket: %s", w1.Error().Error())
	}

	if fi, err = os.Stat(testAddress1); err != nil {
		t.Errorf("Error check unix socket: %s", err)
	}
	if fi.Mode().Perm() != os.FileMode(0666).Perm() {
		t.Logf("Error umix socket Mode(): %v expected %v", fi.Mode(), os.FileMode(0666))
	}

	w1.Stop()
	if _, err = os.Stat(testAddress1); os.IsExist(err) {
		t.Errorf("Error delete unix socket after server stop")
	}
}

func TestServe(t *testing.T) {
	const (
		testAddress1 = `localhost:18080`
		testAddress2 = `127.0.0.1:18080`
	)
	var err error
	var ltn net.Listener
	var w1 = New()

	if ltn, err = net.Listen("tcp", testAddress1); err != nil {
		t.Errorf("Testing error, failed to open port '%s': %s", testAddress1, err.Error())
	}
	defer func() { _ = ltn.Close() }()

	w1.Serve(ltn)
	if w1.(*web).conf == nil {
		t.Errorf("Error, configuration is nil")
	}
	if w1.(*web).conf.Address != testAddress1 && w1.(*web).conf.Address != testAddress2 {
		t.Errorf("Error restore server address from net.Listener. Address is '%s' expected '%s'",
			w1.(*web).conf.Address,
			testAddress1,
		)
	}
}

func TestWait(t *testing.T) {
	const (
		testAddress1 = `localhost:18080`
		testAddress2 = `.test.socket`
	)
	var tic *time.Ticker
	var cou uint32
	var w1 Interface
	var conf *Configuration

	w1 = New()
	w1.ListenAndServe(testAddress1)
	if w1.Error() != nil {
		t.Errorf("Error starting web server: %s", w1.Error().Error())
	}
	go func() {
		tic = time.NewTicker(time.Second / 2)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				w1.Stop()
				break
			}
		}
	}()
	w1.Wait()
	if cou <= 4 {
		t.Errorf("Error Wait()")
	}

	w1 = New()
	conf, _ = parseAddress("")
	conf.Mode = "socket"
	conf.Socket = testAddress2
	conf.KeepAliveDisable = true
	w1.ListenAndServeWithConfig(conf)
	if w1.Error() != nil {
		t.Errorf("Error starting web server: %s", w1.Error().Error())
	}
	go func() {
		tic = time.NewTicker(time.Second / 2)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				w1.Stop()
				break
			}
		}
	}()
	w1.Wait()
	if cou <= 4 {
		t.Errorf("Error Wait()")
	}
}

func TestRunRouteConfigurationError(t *testing.T) {
	const (
		testAddress1    = `localhost:18080`
		testErrorString = `SvDJFQxV4Bscfn2tdP9bCr7CGnK7dYPJWrc6w5MJ`
	)
	var tic *time.Ticker
	var w1 = New()

	_ = w1.Errors().RouteConfigurationError(fmt.Errorf(testErrorString))
	w1.ListenAndServe(testAddress1)
	go func() {
		tic = time.NewTicker(time.Second)
		defer tic.Stop()
		<-tic.C
		w1.Stop()
	}()
	w1.Wait()
	if w1.Error() == nil || w1.Error().Error() != testErrorString {
		t.Errorf("Error starting web server with routing configuration error")
	}
	w1.Stop()
}
