package web

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// Клиентское приветствие протокола HTTP/2 (prior knowledge).
const testH2CPreface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"

// Пустой кадр SETTINGS протокола HTTP/2: длина 0, тип 4, флаги 0, поток 0.
const testH2CSettings = "\x00\x00\x00\x04\x00\x00\x00\x00\x00"

// Функция возвращает HTTP клиента с включённым согласованием протокола HTTP/2.
func getTestHTTP2Client(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			ForceAttemptHTTP2: true,
		},
	}
}

// Выполнение тестового HTTPS запроса к серверу, возврат мажорной версии протокола ответа.
func requestTestHTTP2Proto(address string) (proto int, err error) {
	const (
		tplGetError = "функция Get(), ошибка: %w"
		testTimeout = time.Second * 5
	)
	var resp *http.Response

	if resp, err = getTestHTTP2Client(testTimeout).Get(testSchemeHTTPS + address + testRouteRoot); err != nil {
		err = fmt.Errorf(tplGetError, err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	proto = resp.ProtoMajor

	return
}

// Функция отправляет клиентское приветствие протокола HTTP/2 (prior knowledge) и проверяет, что сервер
// отвечает кадром SETTINGS протокола HTTP/2. Если функция вернула false - сервер ответил не в протоколе HTTP/2.
func probeTestH2C(address string, timeout time.Duration) (ok bool, err error) {
	const (
		tplDialError = "функция DialTimeout(), ошибка: %w"
		tplWrite     = "функция Write(), ошибка: %w"
		tplRead      = "функция Read(), ошибка: %v, прочитано байт: %d"
		frameHeader  = 9
		frameType    = 0x04
	)
	var (
		conn   net.Conn
		n      int
		header = make([]byte, frameHeader)
	)

	if conn, err = net.DialTimeout(netTcp, address, timeout); err != nil {
		err = fmt.Errorf(tplDialError, err)
		return
	}
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(timeout))
	if _, err = conn.Write([]byte(testH2CPreface + testH2CSettings)); err != nil {
		err = fmt.Errorf(tplWrite, err)
		return
	}
	if n, err = io.ReadFull(conn, header); err != nil || n != len(header) {
		err = fmt.Errorf(tplRead, err, n)
		return
	}
	// Тип кадра находится в четвёртом байте заголовка кадра HTTP/2.
	ok = header[3] == frameType

	return
}

// Проверка согласования протокола HTTP/2 для TLS соединения.
//
//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServeTLS_HTTP2(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		tplWrongProto = "функция Get(), версия протокола: HTTP/%d, ожидалось: HTTP/%d"
		testAddress   = "127.0.0.1:18446"
		protoHTTP2    = 2
	)
	var (
		proto int
		err   error
		key   *tmpFile
		crt   *tmpFile
		web   Interface
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress, crt.Filename, key.Filename, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if proto, err = requestTestHTTP2Proto(testAddress); err != nil {
		t.Fatalf(tplRequest, err)
	}
	if proto != protoHTTP2 {
		t.Errorf(tplWrongProto, proto, protoHTTP2)
	}
}

// Проверка отключения протокола HTTP/2 для TLS соединения.
func TestImpl_ListenAndServeTLS_HTTP2Disabled(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeTLSWithConfig(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		tplWrongProto = "функция Get(), версия протокола: HTTP/%d, ожидалось: HTTP/%d"
		testAddress   = "127.0.0.1:18447"
		protoHTTP1    = 1
	)
	var (
		proto int
		err   error
		key   *tmpFile
		crt   *tmpFile
		web   Interface
		cfg   *Configuration
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	if cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	cfg.TLSPublicKeyPEM, cfg.TLSPrivateKeyPEM = crt.Filename, key.Filename
	cfg.HTTP2Disable = true
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLSWithConfig(cfg, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if proto, err = requestTestHTTP2Proto(testAddress); err != nil {
		t.Fatalf(tplRequest, err)
	}
	if proto != protoHTTP1 {
		t.Errorf(tplWrongProto, proto, protoHTTP1)
	}
}

// Проверка работы не зашифрованного протокола HTTP/2 (h2c, prior knowledge).
//
//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServe_H2C(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		tplNoH2C      = "сервер не ответил в протоколе HTTP/2"
		testAddress   = "127.0.0.1:18448"
		testTimeout   = time.Second * 5
	)
	var (
		ok  bool
		err error
		web Interface
		cfg *Configuration
	)

	if cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	cfg.HTTP2UnencryptedEnable = true
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if ok, err = probeTestH2C(testAddress, testTimeout); err != nil {
		t.Fatalf(tplRequest, err)
	}
	if !ok {
		t.Error(tplNoH2C)
	}
}

// Проверка отключения не зашифрованного протокола HTTP/2 (h2c) по умолчанию.
//
//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServe_H2CDisabled(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		tplH2COn      = "сервер ответил в протоколе HTTP/2 при отключённом h2c"
		testAddress   = "127.0.0.1:18449"
		testTimeout   = time.Second * 5
	)
	var (
		ok  bool
		err error
		web Interface
		cfg *Configuration
	)

	if cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if ok, err = probeTestH2C(testAddress, testTimeout); err != nil {
		t.Fatalf(tplRequest, err)
	}
	if ok {
		t.Error(tplH2COn)
	}
}
