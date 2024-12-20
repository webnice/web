package web

import (
	"errors"
	"testing"
)

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
