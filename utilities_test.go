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
			addr string
			host string
			port uint16
			err  error
		}{
			{addr: "localhost", host: "localhost", port: 0, err: nil},
			{addr: ":https", host: "", port: 443, err: nil},
			{addr: ":http", host: "", port: 80, err: nil},
			{addr: "localhost:80", host: "127.0.0.1", port: 80, err: nil},
			{addr: "localhost:443", host: "127.0.0.1", port: 443, err: nil},
			{addr: "", host: "", port: 0, err: nil},
			{addr: "Nmm4zmSnXh1ymGctOtgf6", host: "Nmm4zmSnXh1ymGctOtgf6", port: 0, err: nil},
			{addr: ":abc", host: "", port: 0, err: errors.New("lookup tcp/abc: nodename nor servname provided, or not known")},
		}
		n int
	)

	for n = range tests {
		cfg, err = parseAddress(tests[n].addr)
		if tests[n].err != nil && err == nil {
			t.Errorf(
				"функция parseAddress(%q), ошибка: \"%v\", ожидалось: \"%v\"",
				tests[n].addr, err, tests[n].err,
			)
		}
		if tests[n].host != cfg.Host {
			t.Errorf("функция parseAddress(%q), Host: %q, ожидалось: %q", tests[n].addr, cfg.Host, tests[n].host)
		}
		if tests[n].port != cfg.Port {
			t.Errorf("функция parseAddress(%q), Port: %d, ожидалось: %d", tests[n].addr, cfg.Port, tests[n].port)
		}
	}
}
