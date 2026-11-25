package web

import (
	"bytes"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

func TestEnsureHTTP2CipherSuites(t *testing.T) {
	const (
		tplCiphersNil      = "функция ensureHTTP2CipherSuites(nil), список шифров: %v, ожидалось: %v"
		tplCiphersEqual    = "функция ensureHTTP2CipherSuites(), список шифров: %v, ожидалось: %v"
		tplCiphersAppended = "функция ensureHTTP2CipherSuites(), список шифров: %v, ожидалось: %v + обязательные шифры HTTP/2"
	)
	var (
		ciphers []uint16
		updated []uint16
	)

	// Пустой список шифров не изменяется.
	if updated = ensureHTTP2CipherSuites(nil); len(updated) != 0 {
		t.Errorf(tplCiphersNil, updated, nil)
	}
	// Список шифров, содержащий шифр, обязательный для HTTP/2, не изменяется.
	ciphers = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}
	if updated = ensureHTTP2CipherSuites(ciphers); len(updated) != len(ciphers) {
		t.Errorf(tplCiphersEqual, updated, ciphers)
	}
	// Список шифров без шифров, обязательных для HTTP/2, дополняется.
	ciphers = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}
	updated = ensureHTTP2CipherSuites(ciphers)
	if len(updated) != len(ciphers)+len(http2RequiredCipherSuites) {
		t.Errorf(tplCiphersAppended, updated, ciphers)
	}
}

func TestIsTLSListener(t *testing.T) {
	const (
		tplNilListener = "функция isTLSListener(nil), вернулось: true, ожидалось: false"
		tplPlainIsTls  = "функция isTLSListener(), незащищённый слушатель определён как TLS"
		tplTlsNotFound = "функция isTLSListener(), TLS слушатель не определён как TLS"
		tplListenError = "функция net.Listen(), ошибка: %v, ожидалось: %v"
		testAddress    = "127.0.0.1:0"
	)
	var (
		err    error
		ltn    net.Listener
		tlsLtn net.Listener
	)

	if isTLSListener(nil) {
		t.Error(tplNilListener)
	}
	if ltn, err = net.Listen(netTcp, testAddress); err != nil {
		t.Fatalf(tplListenError, err, nil)
	}
	defer func() { _ = ltn.Close() }()
	if isTLSListener(ltn) {
		t.Error(tplPlainIsTls)
	}
	tlsLtn = tls.NewListener(ltn, new(tls.Config))
	defer func() { _ = tlsLtn.Close() }()
	if !isTLSListener(tlsLtn) {
		t.Error(tplTlsNotFound)
	}
}

func TestParseAddress(t *testing.T) {
	const (
		tplTestName        = "тест %q"
		tplParseError      = "тест %q, функция parseAddress(%q), ошибка: %q, ожидалось: %q"
		tplHost            = "тест %q, функция parseAddress(%q), Host: %q, ожидалось: %q"
		tplPort            = "тест %q, функция parseAddress(%q), Port: %d, ожидалось: %d"
		testHostIP         = "127.0.0.1"
		testAddressHTTP    = "127.0.0.1:80"
		testAddressHTTPS   = "127.0.0.1:443"
		testHostnameRandom = "Nmm4zmSnXh1ymGctOtgf6"
		testPortHTTP       = 80
		testPortHTTPS      = 443
	)
	var (
		err   error
		cfg   *Configuration
		tests = []struct {
			id   string
			addr string
			host string
			port uint16
			err  error
		}{
			{id: "01", addr: testHostIP, host: testHostIP, port: 0, err: nil},
			{id: "02", addr: ":https", host: "", port: testPortHTTPS, err: nil},
			{id: "03", addr: ":http", host: "", port: testPortHTTP, err: nil},
			{id: "04", addr: testAddressHTTP, host: testHostIP, port: testPortHTTP, err: nil},
			{id: "05", addr: testAddressHTTPS, host: testHostIP, port: testPortHTTPS, err: nil},
			{id: "06", addr: "", host: "", port: 0, err: nil},
			{id: "07", addr: testHostnameRandom, host: testHostnameRandom, port: 0, err: nil},
			{id: "08", addr: ":abc", host: "", port: 0, err: errors.New(
				"lookup tcp/abc: nodename nor servname provided, or not known",
			)},
		}
		n int
	)

	for n = range tests {
		t.Logf(tplTestName, tests[n].id)
		cfg, err = parseAddress(tests[n].addr)
		if tests[n].err != nil && err == nil {
			t.Errorf(tplParseError, tests[n].id, tests[n].addr, err, tests[n].err)
		}
		if cfg != nil && tests[n].host != cfg.Host {
			t.Errorf(tplHost, tests[n].id, tests[n].addr, cfg.Host, tests[n].host)
		}
		if cfg != nil && tests[n].port != cfg.Port {
			t.Errorf(tplPort, tests[n].id, tests[n].addr, cfg.Port, tests[n].port)
		}
	}
}

// Проверка безопасного вызова функции оповещения об отказе в обслуживании запроса.
func TestSafeCallRefused(t *testing.T) {
	const (
		tplArgsNotEqual   = "аргументы функции оповещения об отказе: %q, %q, ожидалось: %q, %q"
		tplPanicNotLogged = "функция safeCallRefused(), паника не записана в журнал"
		tplStackNotLogged = "функция safeCallRefused(), стек вызовов не записан в журнал"
		testHost          = "denied.example.com"
		testRemoteAddr    = "192.0.2.1:1234"
		testPanicText     = "тестовая паника"
		testStackMark     = "goroutine "
	)
	var (
		buf        bytes.Buffer
		host       string
		remoteAddr string
	)

	// Вызов функции без паники: аргументы передаются без изменений.
	safeCallRefused(func(h, r string) { host, remoteAddr = h, r }, testHost, testRemoteAddr)
	if host != testHost || remoteAddr != testRemoteAddr {
		t.Errorf(tplArgsNotEqual, host, remoteAddr, testHost, testRemoteAddr)
	}
	// Перенаправление журнала в буфер на время проверки перехвата паники.
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	// Паника внутри функции перехватывается, в журнал записываются паника и стек вызовов.
	safeCallRefused(func(_, _ string) { panic(testPanicText) }, testHost, testRemoteAddr)
	if !strings.Contains(buf.String(), testPanicText) {
		t.Error(tplPanicNotLogged)
	}
	if !strings.Contains(buf.String(), testStackMark) {
		t.Error(tplStackNotLogged)
	}
}

// Проверка загрузки стека вызовов.
func TestStackShort(t *testing.T) {
	const (
		tplStackEmpty = "функция stackShort(), стек вызовов пустой"
		tplStackMark  = "функция stackShort(), стек вызовов не содержит маркер горутины"
		testStackMark = "goroutine "
	)
	var stack []byte

	stack = stackShort()
	if len(stack) == 0 {
		t.Error(tplStackEmpty)
	}
	if !strings.Contains(string(stack), testStackMark) {
		t.Error(tplStackMark)
	}
}
