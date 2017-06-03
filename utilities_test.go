package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net/http"
	"testing"
)

func TestDefaultConfiguration(t *testing.T) {
	const (
		_TestHost    = "5BjSMzDCuHJWHKx2JqbW.Dd5Vr8ytnD968dDrNc3s"
		_TestAddress = "https://abc.pw"
		_TestSocket  = "/var/run/test.socket"
	)
	var conf = new(Configuration)

	defaultConfiguration(conf)
	if conf.Port != 80 {
		t.Errorf("Error configuration defaults: Port is '80' expected '%d'", conf.Port)
	}
	if conf.Mode != "tcp" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "tcp")
	}
	if conf.MaxHeaderBytes != http.DefaultMaxHeaderBytes {
		t.Errorf("Error configuration defaults: MaxHeaderBytes is '%d' expected '%d'", conf.MaxHeaderBytes, http.DefaultMaxHeaderBytes)
	}
	if conf.ShutdownTimeout != _ShutdownTimeout {
		t.Errorf("Error configuration defaults: ShutdownTimeout is '%s' expected '%s'", conf.ShutdownTimeout.String(), _ShutdownTimeout.String())
	}

	conf.Host = _TestHost
	conf.Mode = "socket"
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "tcp")
	}
	if conf.Address != _TestHost {
		t.Errorf("Error configuration defaults: Address is '%s' expected '%s'", conf.Address, _TestHost)
	}

	conf = new(Configuration)
	conf.Host = _TestHost
	conf.Port = 1234
	conf.Mode = _TestHost
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "tcp")
	}
	if conf.Address != fmt.Sprintf("%s:%d", _TestHost, conf.Port) {
		t.Errorf("Error configuration defaults: Address is '%s' expected '%s'", conf.Address, fmt.Sprintf("%s:%d", _TestHost, conf.Port))
	}

	conf = new(Configuration)
	conf.Address = _TestAddress
	conf.Host = _TestHost
	conf.Port = 3210
	conf.Mode = "unixpacket"
	defaultConfiguration(conf)
	if conf.Mode != "tcp" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "tcp")
	}
	if conf.Address != _TestAddress {
		t.Errorf("Error configuration defaults: Address is '%s' expected '%s'", conf.Address, _TestAddress)
	}

	conf = new(Configuration)
	conf.Host = _TestHost
	conf.Mode = "socket"
	conf.Socket = _TestSocket
	defaultConfiguration(conf)
	if conf.Mode != "unix" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "unix")
	}
	if conf.Address != "" {
		t.Errorf("Error configuration defaults: Address is '%s' expected '%s'", conf.Address, _TestAddress)
	}
	if conf.Mode != "unix" {
		t.Errorf("Error configuration defaults: Mode is '%s' expected '%s'", conf.Mode, "unix")
	}
	if conf.HostPort != fmt.Sprintf("unix:%s", _TestSocket) {
		t.Errorf("Error configuration defaults: HostPort is '%s' expected '%s'", conf.HostPort, fmt.Sprintf("unix:%s", _TestSocket))
	}
}

func TestParseAddress(t *testing.T) {
	var err error
	var conf *Configuration

	conf, err = parseAddress("")
	if conf == nil && err == nil {
		t.Errorf("Error parseAddress(), configuration is nil")
	}

	conf, err = parseAddress(":https")
	if conf.Port != 443 || err != nil {
		t.Errorf("Error parseAddress(): Port is '%d' expected '443', error is '%v'", conf.Port, err)
	}
	conf, err = parseAddress(":http")
	if conf.Port != 80 || err != nil {
		t.Errorf("Error parseAddress(): Port is '%d' expected '80', error is '%v'", conf.Port, err)
	}
	conf, err = parseAddress("abcd:9080")
	if conf.Port != 9080 || err != nil {
		t.Errorf("Error parseAddress(): Port is '%d' expected '9080', error is '%v'", conf.Port, err)
	}
	if conf.Host != "abcd" {
		t.Errorf("Error parseAddress(): Host is '%s' expected 'abcd'", conf.Host)
	}

	// Check error
	_, err = parseAddress("abcd:abcd")
	if err == nil {
		t.Errorf("Error parseAddress()")
	}
}
