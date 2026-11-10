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

// getTestHTTP2Client Функция возвращает HTTP клиента с включённым согласованием протокола HTTP/2.
func getTestHTTP2Client(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			ForceAttemptHTTP2: true,
		},
	}
}

// requestTestHTTP2Proto Выполнение тестового HTTPS запроса к серверу, возврат мажорной версии протокола ответа.
func requestTestHTTP2Proto(address string) (proto int, err error) {
	var resp *http.Response

	if resp, err = getTestHTTP2Client(time.Second * 5).Get("https://" + address + "/"); err != nil {
		err = fmt.Errorf("функция Get(), ошибка: %w", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	proto = resp.ProtoMajor

	return
}

// probeTestH2C Функция отправляет клиентское приветствие протокола HTTP/2 (prior knowledge) и проверяет, что
// сервер отвечает кадром SETTINGS протокола HTTP/2. Если функция вернула false - сервер ответил не в протоколе
// HTTP/2.
func probeTestH2C(address string, timeout time.Duration) (ok bool, err error) {
	const (
		preface  = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
		settings = "\x00\x00\x00\x04\x00\x00\x00\x00\x00"
	)
	var (
		conn   net.Conn
		n      int
		header = make([]byte, 9)
	)

	if conn, err = net.DialTimeout(netTcp, address, timeout); err != nil {
		err = fmt.Errorf("функция DialTimeout(), ошибка: %w", err)
		return
	}
	defer func() { _ = conn.Close() }()
	_ = conn.SetDeadline(time.Now().Add(timeout))
	if _, err = conn.Write([]byte(preface + settings)); err != nil {
		err = fmt.Errorf("функция Write(), ошибка: %w", err)
		return
	}
	if n, err = io.ReadFull(conn, header); err != nil || n != len(header) {
		err = fmt.Errorf("функция Read(), ошибка: %v, прочитано байт: %d", err, n)
		return
	}
	// Тип кадра находится в четвёртом байте заголовка кадра HTTP/2.
	ok = header[3] == 0x04

	return
}

// Проверка согласования протокола HTTP/2 для TLS соединения.
func TestImpl_ListenAndServeTLS_HTTP2(t *testing.T) {
	const testAddress = "127.0.0.1:18446"
	var (
		proto int
		err   error
		key   *tmpFile
		crt   *tmpFile
		web   Interface
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress, crt.Filename, key.Filename, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf("функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if proto, err = requestTestHTTP2Proto(testAddress); err != nil {
		t.Fatalf("%v", err)
	}
	if proto != 2 {
		t.Errorf("функция Get(), версия протокола: HTTP/%d, ожидалось: HTTP/%d", proto, 2)
	}
}

// Проверка отключения протокола HTTP/2 для TLS соединения.
func TestImpl_ListenAndServeTLS_HTTP2Disabled(t *testing.T) {
	const testAddress = "127.0.0.1:18447"
	var (
		proto int
		err   error
		key   *tmpFile
		crt   *tmpFile
		web   Interface
		cfg   *Configuration
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	cfg, _ = parseAddress(testAddress)
	cfg.TLSPublicKeyPEM, cfg.TLSPrivateKeyPEM = crt.Filename, key.Filename
	cfg.HTTP2Disable = true
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLSWithConfig(cfg, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf("функция ListenAndServeTLSWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if proto, err = requestTestHTTP2Proto(testAddress); err != nil {
		t.Fatalf("%v", err)
	}
	if proto != 1 {
		t.Errorf("функция Get(), версия протокола: HTTP/%d, ожидалось: HTTP/%d", proto, 1)
	}
}

// Проверка работы не зашифрованного протокола HTTP/2 (h2c, prior knowledge).
func TestImpl_ListenAndServe_H2C(t *testing.T) {
	const testAddress = "127.0.0.1:18448"
	var (
		ok  bool
		err error
		web Interface
		cfg *Configuration
	)

	cfg, _ = parseAddress(testAddress)
	cfg.HTTP2UnencryptedEnable = true
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if ok, err = probeTestH2C(testAddress, time.Second*5); err != nil {
		t.Fatalf("%v", err)
	}
	if !ok {
		t.Errorf("сервер не ответил в протоколе HTTP/2")
	}
}

// Проверка отключения не зашифрованного протокола HTTP/2 (h2c) по умолчанию.
func TestImpl_ListenAndServe_H2CDisabled(t *testing.T) {
	const testAddress = "127.0.0.1:18449"
	var (
		ok  bool
		err error
		web Interface
		cfg *Configuration
	)

	cfg, _ = parseAddress(testAddress)
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if ok, err = probeTestH2C(testAddress, time.Second*5); err != nil {
		t.Fatalf("%v", err)
	}
	if ok {
		t.Errorf("сервер ответил в протоколе HTTP/2 при отключённом h2c")
	}
}
