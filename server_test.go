package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestImpl_MakeServer(t *testing.T) {
	const (
		testAddress              = "localhost:18080"
		tplMakeServerError       = "функция makeServer(%v), ошибка: %v, ожидалось: %v"
		tplMakeServerConfig      = "функция makeServer(%v), вернулась конфигурация, ожидалось: %v"
		tplMakeServerIsTls       = "функция makeServer(%v), isTls: %t, ожидалось: %t"
		tplMakeServerHTTP2Config = "функция makeServer(%v), http.Server: %v, ожидалась конфигурация"
		tplMakeServerNextProto   = "функция makeServer(%v), TLSNextProto не равно: %v"
	)
	var (
		web   *impl
		srv   *http.Server
		isTls bool
		err   error
	)

	web = New().(*impl)
	// Проверка реакции на не установленный обработчик запросов.
	srv, isTls = web.makeServer(nil)
	if !errors.Is(web.Error(), web.Errors().HandlerIsNotSet.Bind()) {
		t.Errorf(tplMakeServerError, nil, web.Error(), web.Errors().HandlerIsNotSet.Bind())
	}
	if srv != nil {
		t.Errorf(tplMakeServerConfig, nil, nil)
	}
	if isTls {
		t.Errorf(tplMakeServerIsTls, nil, isTls, false)
	}
	// Проверка реакции на отсутствие конфигурации.
	web.Clean().
		Handler(getTestHandlerFn(t))
	srv, isTls = web.makeServer(nil)
	if !errors.Is(web.Error(), web.Errors().NoConfiguration.Bind()) {
		t.Errorf(tplMakeServerError, nil, web.Error(), web.Errors().NoConfiguration.Bind())
	}
	// Проверка создания сервера по умолчанию без TLS.
	web.Clean()
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplMakeServerError, nil, err, nil)
	}
	srv, isTls = web.makeServer(nil)
	if web.Error() != nil {
		t.Errorf(tplMakeServerError, nil, web.Error(), nil)
	}
	if srv == nil {
		t.Errorf(tplMakeServerHTTP2Config, nil, srv)
	}
	if isTls {
		t.Errorf(tplMakeServerIsTls, nil, isTls, false)
	}
	if srv != nil && srv.TLSNextProto != nil {
		t.Errorf(tplMakeServerNextProto, nil, nil)
	}
}

func TestImpl_MakeServerTlsConfig(t *testing.T) {
	const (
		tplMakeServerError   = "функция makeServer(), ошибка: %v, ожидалась ошибка"
		tplMakeServerNoError = "функция makeServer(), ошибка: %v, ошибка не ожидалась"
		testPublicKey        = "111"
		testPrivateKey       = "222"
		testHostLocal        = "localhost"
		testPortLocal        = 18080
	)
	var (
		err error
		key *tmpFile
		crt *tmpFile
		web *impl
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(chi.NewMux())
	web.cfg = &Configuration{
		Host:             testHostLocal,
		Port:             testPortLocal,
		TLSPublicKeyPEM:  testPublicKey,
		TLSPrivateKeyPEM: testPrivateKey,
	}
	// Проверка реакции на не корректные ключи TLS.
	if _, _ = web.makeServer(nil); web.Error() == nil {
		t.Errorf(tplMakeServerError, err)
	}
	// Проверка создания сервера с корректными ключами TLS.
	web.Clean()
	web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM = crt.Filename, key.Filename
	if _, _ = web.makeServer(nil); web.Error() != nil {
		t.Errorf(tplMakeServerNoError, err)
	}
}

// Проверка настройки протоколов HTTP сервера по умолчанию и в соответствии с конфигурацией.
func TestImpl_MakeServerProtocols(t *testing.T) {
	const (
		testAddress            = "localhost:18080"
		tplMakeServerFatal     = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplProtocolDisabled    = "функция makeServer(), протокол HTTP/1 отключён"
		tplProtocolHTTP2       = "функция makeServer(), протокол HTTP/2 отключён"
		tplProtocolH2cEnabled  = "функция makeServer(), протокол h2c включён"
		tplHTTP2ConfigNotSet   = "функция makeServer(), конфигурация HTTP/2 не установлена"
		tplProtocolHTTP2On     = "функция makeServer(), протокол HTTP/2 не отключён"
		tplHTTP2ConfigSet      = "функция makeServer(), конфигурация HTTP/2 установлена"
		tplProtocolH2cDisabled = "функция makeServer(), протокол h2c не включён"
	)
	var (
		web *impl
		srv *http.Server
		err error
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplMakeServerFatal, err, nil)
	}
	// По умолчанию включены протоколы HTTP/1 и HTTP/2, h2c отключён.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if !srv.Protocols.HTTP1() {
		t.Error(tplProtocolDisabled)
	}
	if !srv.Protocols.HTTP2() {
		t.Error(tplProtocolHTTP2)
	}
	if srv.Protocols.UnencryptedHTTP2() {
		t.Error(tplProtocolH2cEnabled)
	}
	if srv.HTTP2 == nil {
		t.Error(tplHTTP2ConfigNotSet)
	}
	// Отключение протокола HTTP/2.
	web.Clean()
	web.cfg.HTTP2Disable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if srv.Protocols.HTTP2() {
		t.Error(tplProtocolHTTP2On)
	}
	if srv.HTTP2 != nil {
		t.Error(tplHTTP2ConfigSet)
	}
	// Включение протокола h2c.
	web.Clean()
	web.cfg.HTTP2Disable = false
	web.cfg.HTTP2UnencryptedEnable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if !srv.Protocols.UnencryptedHTTP2() {
		t.Error(tplProtocolH2cDisabled)
	}
	if srv.HTTP2 == nil {
		t.Error(tplHTTP2ConfigNotSet)
	}
}

// Проверка настройки конфигурации протокола HTTP/2.
func TestImpl_MakeServerHTTP2Config(t *testing.T) {
	const (
		testAddress            = "localhost:18080"
		tplMakeServerFatal     = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplHTTP2ConfigNotSet   = "функция makeServer(), конфигурация HTTP/2 не установлена"
		tplMaxConcurrentStream = "функция makeServer(), MaxConcurrentStreams: %d, ожидалось: %d"
		testStreams            = 100
	)
	var (
		web *impl
		srv *http.Server
		err error
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplMakeServerFatal, err, nil)
	}
	web.cfg.HTTP2MaxConcurrentStreams = testStreams
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if srv.HTTP2 == nil {
		t.Fatalf(tplHTTP2ConfigNotSet)
	}
	if srv.HTTP2.MaxConcurrentStreams != testStreams {
		t.Errorf(tplMaxConcurrentStream, srv.HTTP2.MaxConcurrentStreams, testStreams)
	}
}

// Проверка ограничения максимального размера тела запроса.
func TestImpl_MakeServerMaxBodyBytes(t *testing.T) {
	const (
		testAddress          = "localhost:18080"
		tplMakeServerFatal   = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplHandlerChanged    = "функция makeServer(), обработчик изменён без ограничения размера тела запроса"
		tplHandlerNotLimited = "функция makeServer(), обработчик не ограничен размером тела запроса"
		tplHandlerBadLimit   = "функция makeServer(), обработчик ограничен некорректным значением размера тела запроса"
	)
	var (
		web     *impl
		srv     *http.Server
		handler http.Handler
		err     error
	)

	handler = getTestHandlerFn(t)
	web = New().(*impl)
	web.Handler(handler)
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplMakeServerFatal, err, nil)
	}
	// Без ограничения размер тела запроса обработчик не оборачивается.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if srv.Handler != handler {
		t.Error(tplHandlerChanged)
	}
	// С ограничением обработчик оборачивается.
	web.cfg.MaxBodyBytes = 10
	srv, _ = web.makeServer(nil)
	if srv.Handler == handler {
		t.Error(tplHandlerNotLimited)
	}
	// Значение больше максимального int64 не применяется.
	web.cfg.MaxBodyBytes = 1 << 63
	srv, _ = web.makeServer(nil)
	if srv.Handler != handler {
		t.Error(tplHandlerBadLimit)
	}
}

// Проверка обёртки чтения тела запроса с ограничением размера.
func TestMaxBodyBytesReader(t *testing.T) {
	const (
		tplReadError  = "функция Read(), прочитано байт: %d, ошибка: %v, ожидалось: %d"
		tplCloseError = "функция Close(), ошибка: %v, ожидалось: %v"
		testContent   = "ok"
		testCount     = 2
	)
	var (
		mbr *maxBodyBytesReader
		buf []byte
	)

	mbr = &maxBodyBytesReader{w: httptest.NewRecorder(), r: io.NopCloser(strings.NewReader(testContent))}
	buf = make([]byte, testCount)
	if n, err := mbr.Read(buf); err != nil || n != testCount {
		t.Errorf(tplReadError, n, err, testCount)
	}
	if err := mbr.Close(); err != nil {
		t.Errorf(tplCloseError, err, nil)
	}
}

// Проверка согласования протокола ALPN при включённом и отключённом протоколе HTTP/2.
func TestImpl_MakeServerNextProtos(t *testing.T) {
	const (
		tplMakeServerFatal = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplAlpnH2Miss      = "функция makeServer(), протокол h2 отсутствует в списке ALPN: %q"
		tplAlpnH2Present   = "функция makeServer(), протокол h2 присутствует в списке ALPN: %q"
		tplAlpnHTTP1Miss   = "функция makeServer(), протокол http/1.1 отсутствует в списке ALPN: %q"
		protoH2            = "h2"
		protoHTTP1         = "http/1.1"
		testHostIP         = "127.0.0.1"
		testHostPort       = 18080
	)
	var (
		key *tmpFile
		crt *tmpFile
		web *impl
		srv *http.Server
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.cfg = &Configuration{
		Host:             testHostIP,
		Port:             testHostPort,
		TLSPublicKeyPEM:  crt.Filename,
		TLSPrivateKeyPEM: key.Filename,
	}
	// Протокол HTTP/2 включён: протокол h2 присутствует в списке ALPN.
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if !slices.Contains(srv.TLSConfig.NextProtos, protoH2) {
		t.Errorf(tplAlpnH2Miss, srv.TLSConfig.NextProtos)
	}
	// Протокол HTTP/2 отключён: протокол h2 отсутствует в списке ALPN.
	web.Clean()
	web.cfg.HTTP2Disable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplMakeServerFatal, web.Error(), nil)
	}
	if slices.Contains(srv.TLSConfig.NextProtos, protoH2) {
		t.Errorf(tplAlpnH2Present, srv.TLSConfig.NextProtos)
	}
	if !slices.Contains(srv.TLSConfig.NextProtos, protoHTTP1) {
		t.Errorf(tplAlpnHTTP1Miss, srv.TLSConfig.NextProtos)
	}
}

func TestImpl_ServeTLS(t *testing.T) {
	const (
		tplMakeServerError = "функция ServeTLS(), ошибка: %v, ожидалась ошибка"
		tplNewListener     = "функция NewListener(), ошибка: %v"
		testPublicKey      = "111"
		testPrivateKey     = "222"
		testHostLocal      = "localhost"
	)
	var (
		err      error
		key      *tmpFile
		crt      *tmpFile
		web      Interface
		listener net.Listener
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(chi.NewMux())
	web.(*impl).cfg = &Configuration{
		Host:             testHostLocal,
		Port:             18080,
		TLSPublicKeyPEM:  testPublicKey,
		TLSPrivateKeyPEM: testPrivateKey,
	}
	if listener, err = web.NewListener(web.(*impl).cfg); err != nil {
		t.Fatalf(tplNewListener, err)
	}
	defer func() { _ = listener.Close() }()
	// Запуск сервера в режиме TLS с не корректными ключами.
	if web.ServeTLS(listener, nil); web.Error() == nil {
		t.Errorf(tplMakeServerError, err)
	}
}

// Проверка настройки полей сервера по конфигурации: лимит заголовков, приоритизация, журнал ошибок,
// контексты и состояния соединений.
func TestImpl_MakeServerFields(t *testing.T) {
	const (
		tplFatalError          = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplFieldNotSet         = "функция makeServer(), поле %s не назначено"
		tplWrongMaxHeaderValue = "функция makeServer(), MaxHeaderValueCount: %d, ожидалось: %d"
		testAddress            = "localhost:18080"
		testHeaderCount        = 100
	)
	var (
		web *impl
		srv *http.Server
		err error
		csf func(net.Conn, http.ConnState)
		bsf func(net.Listener) context.Context
		ccf func(context.Context, net.Conn) context.Context
	)

	csf = func(net.Conn, http.ConnState) {}
	bsf = func(net.Listener) context.Context { return context.Background() }
	ccf = func(ctx context.Context, c net.Conn) context.Context { return ctx }
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.cfg.MaxHeaderValueCount = testHeaderCount
	web.cfg.DisableClientPriority = true
	web.cfg.ConnState = csf
	web.cfg.BaseContext = bsf
	web.cfg.ConnContext = ccf
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.MaxHeaderValueCount != testHeaderCount {
		t.Errorf(tplWrongMaxHeaderValue, srv.MaxHeaderValueCount, testHeaderCount)
	}
	if !srv.DisableClientPriority {
		t.Errorf(tplFieldNotSet, "DisableClientPriority")
	}
	if reflect.ValueOf(srv.ConnState).Pointer() != reflect.ValueOf(csf).Pointer() {
		t.Errorf(tplFieldNotSet, "ConnState")
	}
	if reflect.ValueOf(srv.BaseContext).Pointer() != reflect.ValueOf(bsf).Pointer() {
		t.Errorf(tplFieldNotSet, "BaseContext")
	}
	if reflect.ValueOf(srv.ConnContext).Pointer() != reflect.ValueOf(ccf).Pointer() {
		t.Errorf(tplFieldNotSet, "ConnContext")
	}
}

// Проверка настройки полной конфигурации протокола HTTP/2.
func TestImpl_MakeServerHTTP2FullConfig(t *testing.T) {
	const (
		tplFatalError = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplWrongHTTP2 = "функция makeServer(), конфигурация HTTP/2 не корректна: %+v"
		testAddress   = "localhost:18080"
	)
	var (
		web *impl
		srv *http.Server
		err error
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.cfg.HTTP2MaxDecoderHeaderTableSize = 1000
	web.cfg.HTTP2MaxEncoderHeaderTableSize = 2000
	web.cfg.HTTP2MaxReadFrameSize = 16384
	web.cfg.HTTP2MaxReceiveBufferPerConnection = 65536
	web.cfg.HTTP2MaxReceiveBufferPerStream = 32768
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.HTTP2 == nil || srv.HTTP2.MaxDecoderHeaderTableSize != 1000 ||
		srv.HTTP2.MaxEncoderHeaderTableSize != 2000 || srv.HTTP2.MaxReadFrameSize != 16384 ||
		srv.HTTP2.MaxReceiveBufferPerConnection != 65536 || srv.HTTP2.MaxReceiveBufferPerStream != 32768 {
		t.Errorf(tplWrongHTTP2, srv.HTTP2)
	}
}

// Проверка приоритетов выбора TLS конфигурации и применения настроек mTLS.
//
//goland:noinspection DuplicatedCode
func TestImpl_MakeServerTLSConfig(t *testing.T) {
	const (
		tplFatalError  = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplTlsNotEqual = "функция makeServer(), TLS конфигурация не корректна"
		testAddress    = "localhost:18080"
	)
	var (
		web       *impl
		srv       *http.Server
		isTls     bool
		err       error
		customTLS = new(tls.Config)
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	// TLS конфигурация из Configuration без ключей PEM.
	web.cfg.TLSConfig = customTLS
	srv, isTls = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if !isTls || srv.TLSConfig != customTLS {
		t.Errorf(tplTlsNotEqual)
	}
	// Переданная аргументом TLS конфигурация имеет приоритет над конфигурацией из Configuration.
	argTLS := new(tls.Config)
	srv, _ = web.makeServer(argTLS)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.TLSConfig != argTLS {
		t.Errorf(tplTlsNotEqual)
	}
	// Настройки mTLS применяются к TLS конфигурации.
	web.Clean()
	web.cfg.TLSClientAuth = tls.RequireAndVerifyClientCert
	web.cfg.TLSClientCAs = x509.NewCertPool()
	srv, _ = web.makeServer(argTLS)
	if srv.TLSConfig.ClientAuth != tls.RequireAndVerifyClientCert || srv.TLSConfig.ClientCAs == nil {
		t.Errorf(tplTlsNotEqual)
	}
}

// Проверка фильтрации доменов по SNI на этапе TLS handshake и обёртки обработчика запросов.
func TestImpl_MakeServerDomainFilter(t *testing.T) {
	const (
		tplFatalError      = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplTlsNotEqual     = "функция makeServer(), TLS конфигурация не корректна"
		tplErrorMissed     = "функция GetConfigForClient(), ошибка не вернулась"
		tplGetConfigError  = "функция GetConfigForClient(), ошибка: %v, ожидалось: %v"
		tplUserFuncNotCall = "функция GetConfigForClient(), пользовательская функция не вызвана"
		tplHandlerNotWrap  = "функция makeServer(), обработчик запросов не обёрнут фильтром доменов"
		testAddress        = "localhost:18080"
		testHostDenied     = "denied.example.com"
	)
	var (
		web *impl
		srv *http.Server
		err error
		key *tmpFile
		crt *tmpFile
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.cfg = &Configuration{
		Host:             "localhost",
		Port:             18080,
		TLSPublicKeyPEM:  crt.Filename,
		TLSPrivateKeyPEM: key.Filename,
		Domain:           []string{"allowed.example.com"},
	}
	// Без TLS ключей обработчик запросов обёрнут фильтром доменов.
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.cfg.Domain = []string{"allowed.example.com"}
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.Handler == web.handler {
		t.Error(tplHandlerNotWrap)
	}
	// С TLS ключами фильтр доменов назначается на этап TLS handshake.
	web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM = crt.Filename, key.Filename
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.TLSConfig.GetConfigForClient == nil {
		t.Fatalf(tplTlsNotEqual)
	}
	if _, err = srv.TLSConfig.GetConfigForClient(&tls.ClientHelloInfo{ServerName: testHostDenied}); err == nil {
		t.Errorf(tplErrorMissed)
	}
	if _, err = srv.TLSConfig.GetConfigForClient(&tls.ClientHelloInfo{ServerName: "allowed.example.com"}); err != nil {
		t.Errorf(tplGetConfigError, err, nil)
	}
	// Пользовательская функция GetConfigForClient сохраняется и вызывается после проверки домена.
	web.Clean()
	web.cfg.TLSConfig = &tls.Config{GetConfigForClient: func(*tls.ClientHelloInfo) (*tls.Config, error) {
		return nil, nil
	}}
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	var called bool
	web.cfg.TLSConfig.GetConfigForClient = func(*tls.ClientHelloInfo) (*tls.Config, error) {
		called = true
		return nil, nil
	}
	srv, _ = web.makeServer(nil)
	if _, err = srv.TLSConfig.GetConfigForClient(&tls.ClientHelloInfo{ServerName: "allowed.example.com"}); err != nil {
		t.Errorf(tplGetConfigError, err, nil)
	}
	if !called {
		t.Error(tplUserFuncNotCall)
	}
}

// Проверка включения полнодуплексного режима HTTP/1.
//
//goland:noinspection DuplicatedCode
func TestImpl_MakeServerFullDuplex(t *testing.T) {
	const (
		tplFatalError     = "функция makeServer(), ошибка: %v, ожидалось: %v"
		tplHandlerNotWrap = "функция makeServer(), обработчик запросов не обёрнут полнодуплексным режимом"
		tplNewRequest     = "функция NewRequest(), ошибка: %v"
		tplWrongCode      = "функция ServeHTTP(), код ответа: %d, ожидалось: %d"
		testAddress       = "localhost:18080"
	)
	var (
		web *impl
		srv *http.Server
		err error
		req *http.Request
	)

	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	if web.cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.cfg.HTTP1FullDuplexEnable = true
	srv, _ = web.makeServer(nil)
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if srv.Handler == web.handler {
		t.Error(tplHandlerNotWrap)
	}
	// Проверка работы обёрнутого обработчика.
	rec := httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	srv.Handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
}
