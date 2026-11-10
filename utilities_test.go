package web

import (
	"crypto/tls"
	"errors"
	"net"
	"testing"
)

func TestEnsureHTTP2CipherSuites(t *testing.T) {
	var (
		ciphers []uint16
		updated []uint16
	)

	// Пустой список шифров не изменяется.
	if updated = ensureHTTP2CipherSuites(nil); len(updated) != 0 {
		t.Errorf("функция ensureHTTP2CipherSuites(nil), список шифров: %v, ожидалось: %v", updated, nil)
	}
	// Список шифров, содержащий шифр, обязательный для HTTP/2, не изменяется.
	ciphers = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}
	if updated = ensureHTTP2CipherSuites(ciphers); len(updated) != len(ciphers) {
		t.Errorf("функция ensureHTTP2CipherSuites(), список шифров: %v, ожидалось: %v", updated, ciphers)
	}
	// Список шифров без шифров, обязательных для HTTP/2, дополняется.
	ciphers = []uint16{tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384}
	updated = ensureHTTP2CipherSuites(ciphers)
	if len(updated) != len(ciphers)+len(http2RequiredCipherSuites) {
		t.Errorf("функция ensureHTTP2CipherSuites(), список шифров: %v, ожидалось: %v + обязательные шифры HTTP/2",
			updated, ciphers)
	}
}

func TestIsTLSListener(t *testing.T) {
	var (
		err    error
		ltn    net.Listener
		tlsLtn net.Listener
	)

	if isTLSListener(nil) {
		t.Errorf("функция isTLSListener(nil), вернулось: true, ожидалось: false")
	}
	if ltn, err = net.Listen(netTcp, "127.0.0.1:0"); err != nil {
		t.Fatalf("функция net.Listen(), ошибка: %v, ожидалось: %v", err, nil)
	}
	defer func() { _ = ltn.Close() }()
	if isTLSListener(ltn) {
		t.Errorf("функция isTLSListener(), незащищённый слушатель определён как TLS")
	}
	tlsLtn = tls.NewListener(ltn, new(tls.Config))
	defer func() { _ = tlsLtn.Close() }()
	if !isTLSListener(tlsLtn) {
		t.Errorf("функция isTLSListener(), TLS слушатель не определён как TLS")
	}
}

func TestParseAddress(t *testing.T) {
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
			{id: "01", addr: "127.0.0.1", host: "127.0.0.1", port: 0, err: nil},
			{id: "02", addr: ":https", host: "", port: 443, err: nil},
			{id: "03", addr: ":http", host: "", port: 80, err: nil},
			{id: "04", addr: "127.0.0.1:80", host: "127.0.0.1", port: 80, err: nil},
			{id: "05", addr: "127.0.0.1:443", host: "127.0.0.1", port: 443, err: nil},
			{id: "06", addr: "", host: "", port: 0, err: nil},
			{id: "07", addr: "Nmm4zmSnXh1ymGctOtgf6", host: "Nmm4zmSnXh1ymGctOtgf6", port: 0, err: nil},
			{id: "08", addr: ":abc", host: "", port: 0, err: errors.New(
				"lookup tcp/abc: nodename nor servname provided, or not known",
			)},
		}
		n int
	)

	for n = range tests {
		t.Logf("тест %q", tests[n].id)
		cfg, err = parseAddress(tests[n].addr)
		if tests[n].err != nil && err == nil {
			t.Errorf(
				"тест %q, функция parseAddress(%q), ошибка: \"%v\", ожидалось: \"%v\"",
				tests[n].id, tests[n].addr, err, tests[n].err,
			)
		}
		if cfg != nil && tests[n].host != cfg.Host {
			t.Errorf(
				"тест %q, функция parseAddress(%q), Host: %q, ожидалось: %q",
				tests[n].id, tests[n].addr, cfg.Host, tests[n].host,
			)
		}
		if cfg != nil && tests[n].port != cfg.Port {
			t.Errorf(
				"тест %q, функция parseAddress(%q), Port: %d, ожидалось: %d",
				tests[n].id, tests[n].addr, cfg.Port, tests[n].port,
			)
		}
	}
}
