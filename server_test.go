package web

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	wnet "github.com/webnice/net"
)

func TestImpl_MakeServer(t *testing.T) {
	const testAddress = `localhost:18080`
	var (
		web   *impl
		srv   *http.Server
		isTls bool
	)

	web = New().(*impl)
	srv, isTls = web.makeServer(nil)
	if !errors.Is(web.Error(), web.Errors().HandlerIsNotSet()) {
		t.Errorf("функция makeServer(%v), ошибка: %v, ожидалось: %v", nil, web.Error(), web.Errors().HandlerIsNotSet())
	}
	if srv != nil {
		t.Errorf("функция makeServer(%v), вернулась конфигурация, ожидалось: %v", nil, nil)
	}
	if isTls {
		t.Errorf("функция makeServer(%v), isTls: %t, ожидалось: %t", nil, isTls, false)
	}
	web.Clean().
		Handler(getTestHandlerFn(t))
	srv, isTls = web.makeServer(nil)
	if !errors.Is(web.Error(), web.Errors().NoConfiguration()) {
		t.Errorf("функция makeServer(%v), ошибка: %v, ожидалось: %v", nil, web.Error(), web.Errors().NoConfiguration())
	}
	web.Clean()
	web.cfg, _ = parseAddress(testAddress)
	srv, isTls = web.makeServer(nil)
	if web.Error() != nil {
		t.Errorf("функция makeServer(%v), ошибка: %v, ожидалось: %v", nil, web.Error(), nil)
	}
	if srv == nil {
		t.Errorf("функция makeServer(%v), http.Server: %v, ожидалась конфигурация", nil, srv)
	}
	if isTls {
		t.Errorf("функция makeServer(%v), isTls: %t, ожидалось: %t", nil, isTls, false)
	}
	if srv.TLSNextProto != nil {
		t.Errorf("функция makeServer(%v), TLSNextProto не равно: %v", nil, nil)
	}
}

func TestImpl_MakeServerTlsConfig(t *testing.T) {
	var (
		err error
		key *tmpFile
		crt *tmpFile
		web *impl
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(chi.NewMux())
	web.cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host:             "localhost",
			Port:             18080,
			TLSPublicKeyPEM:  "111",
			TLSPrivateKeyPEM: "222",
		},
	}
	if _, _ = web.makeServer(nil); web.Error() == nil {
		t.Errorf("функция makeServer(), ошибка: %v, ожидалась ошибка", err)
	}
	web.Clean()
	web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM = crt.Filename, key.Filename
	if _, _ = web.makeServer(nil); web.Error() != nil {
		t.Errorf("функция makeServer(), ошибка: %v, ошибка не ожидалась", err)
	}
}

func TestImpl_ServeTLS(t *testing.T) {
	var (
		err       error
		key       *tmpFile
		crt       *tmpFile
		web       Interface
		listener  net.Listener
		tlsConfig *tls.Config
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(chi.NewMux())
	web.(*impl).cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host:             "localhost",
			Port:             18080,
			TLSPublicKeyPEM:  "111",
			TLSPrivateKeyPEM: "222",
		},
	}
	listener, err = web.NewListener(web.(*impl).cfg)
	//tlsConfig, err = wnet.New().NewTLSConfigDefault(web.(*impl).cfg.TLSPublicKeyPEM, web.(*impl).cfg.TLSPrivateKeyPEM)
	if web.ServeTLS(listener, tlsConfig); web.Error() == nil {
		t.Errorf("функция makeServer(), ошибка: %v, ожидалась ошибка", err)
	}
}
