package web

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	wnet "github.com/webnice/net"
)

func TestImpl_MakeServer(t *testing.T) {
	const testAddress = "localhost:18080"
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
	if srv != nil && srv.TLSNextProto != nil {
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

// Проверка настройки протоколов HTTP сервера по умолчанию и в соответствии с конфигурацией.
func TestImpl_MakeServerProtocols(t *testing.T) {
	const testAddress = "localhost:18080"
	var (
		web *impl
		srv *http.Server
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.cfg, _ = parseAddress(testAddress)
	// По умолчанию включены протоколы HTTP/1 и HTTP/2, h2c отключён.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if !srv.Protocols.HTTP1() {
		t.Errorf("функция makeServer(), протокол HTTP/1 отключён")
	}
	if !srv.Protocols.HTTP2() {
		t.Errorf("функция makeServer(), протокол HTTP/2 отключён")
	}
	if srv.Protocols.UnencryptedHTTP2() {
		t.Errorf("функция makeServer(), протокол h2c включён")
	}
	if srv.HTTP2 == nil {
		t.Errorf("функция makeServer(), конфигурация HTTP/2 не установлена")
	}
	// Отключение протокола HTTP/2.
	web.Clean()
	web.cfg.HTTP2Disable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if srv.Protocols.HTTP2() {
		t.Errorf("функция makeServer(), протокол HTTP/2 не отключён")
	}
	if srv.HTTP2 != nil {
		t.Errorf("функция makeServer(), конфигурация HTTP/2 установлена")
	}
	// Включение протокола h2c.
	web.Clean()
	web.cfg.HTTP2Disable = false
	web.cfg.HTTP2UnencryptedEnable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if !srv.Protocols.UnencryptedHTTP2() {
		t.Errorf("функция makeServer(), протокол h2c не включён")
	}
	if srv.HTTP2 == nil {
		t.Errorf("функция makeServer(), конфигурация HTTP/2 не установлена")
	}
}

// Проверка настройки конфигурации протокола HTTP/2.
func TestImpl_MakeServerHTTP2Config(t *testing.T) {
	const testAddress = "localhost:18080"
	var (
		web *impl
		srv *http.Server
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.cfg, _ = parseAddress(testAddress)
	web.cfg.HTTP2MaxConcurrentStreams = 100
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if srv.HTTP2 == nil {
		t.Fatalf("функция makeServer(), конфигурация HTTP/2 не установлена")
	}
	if srv.HTTP2.MaxConcurrentStreams != 100 {
		t.Errorf("функция makeServer(), MaxConcurrentStreams: %d, ожидалось: %d",
			srv.HTTP2.MaxConcurrentStreams, 100)
	}
}

// Проверка ограничения максимального размера тела запроса.
func TestImpl_MakeServerMaxBodyBytes(t *testing.T) {
	var (
		web     *impl
		srv     *http.Server
		handler http.Handler
	)

	handler = getTestHandlerFn(t)
	web = New().(*impl)
	web.Handler(handler)
	web.cfg, _ = parseAddress("localhost:18080")
	// Без ограничения размер тела запроса обработчик не оборачивается.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if srv.Handler != handler {
		t.Errorf("функция makeServer(), обработчик изменён без ограничения размера тела запроса")
	}
	// С ограничением обработчик оборачивается.
	web.cfg.MaxBodyBytes = 10
	srv, _ = web.makeServer(nil)
	if srv.Handler == handler {
		t.Errorf("функция makeServer(), обработчик не ограничен размером тела запроса")
	}
	// Значение больше максимального int64 не применяется.
	web.cfg.MaxBodyBytes = 1 << 63
	srv, _ = web.makeServer(nil)
	if srv.Handler != handler {
		t.Errorf("функция makeServer(), обработчик ограничен некорректным значением размера тела запроса")
	}
}

// Проверка обёртки чтения тела запроса с ограничением размера.
func TestMaxBodyBytesReader(t *testing.T) {
	var mbr *maxBodyBytesReader

	mbr = &maxBodyBytesReader{w: httptest.NewRecorder(), r: io.NopCloser(strings.NewReader("ok"))}
	buf := make([]byte, 2)
	if n, err := mbr.Read(buf); err != nil || n != 2 {
		t.Errorf("функция Read(), прочитано байт: %d, ошибка: %v, ожидалось: %d", n, err, 2)
	}
	if err := mbr.Close(); err != nil {
		t.Errorf("функция Close(), ошибка: %v, ожидалось: %v", err, nil)
	}
}

// Проверка согласования протокола ALPN при включённом и отключённом протоколе HTTP/2.
func TestImpl_MakeServerNextProtos(t *testing.T) {
	var (
		key *tmpFile
		crt *tmpFile
		web *impl
		srv *http.Server
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host:             "127.0.0.1",
			Port:             18080,
			TLSPublicKeyPEM:  crt.Filename,
			TLSPrivateKeyPEM: key.Filename,
		},
	}
	// Протокол HTTP/2 включён: протокол h2 присутствует в списке ALPN.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if !slices.Contains(srv.TLSConfig.NextProtos, "h2") {
		t.Errorf("функция makeServer(), протокол h2 отсутствует в списке ALPN: %q", srv.TLSConfig.NextProtos)
	}
	// Протокол HTTP/2 отключён: протокол h2 отсутствует в списке ALPN.
	web.Clean()
	web.cfg.HTTP2Disable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf("функция makeServer(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if slices.Contains(srv.TLSConfig.NextProtos, "h2") {
		t.Errorf("функция makeServer(), протокол h2 присутствует в списке ALPN: %q", srv.TLSConfig.NextProtos)
	}
	if !slices.Contains(srv.TLSConfig.NextProtos, "http/1.1") {
		t.Errorf("функция makeServer(), протокол http/1.1 отсутствует в списке ALPN: %q",
			srv.TLSConfig.NextProtos)
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
