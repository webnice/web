package web

import (
	"net"
	"reflect"
	"strings"
	"time"
)

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных как ":http".
func parseAddress(addr string) (ret *Configuration, err error) {
	const bColon = ":"
	var (
		sp []string
		n  int
	)

	addr = strings.TrimSpace(addr)
	ret, sp = new(Configuration), make([]string, 2)
	if sp[0], sp[1], err = net.SplitHostPort(addr); err != nil {
		ret.Host, err = addr, nil
		return
	}
	if n, err = net.LookupPort(netTcp, strings.Join(sp[1:], bColon)); err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint16(n)
	switch sp, err = net.LookupHost(sp[0]); err {
	case nil:
		ret.Host = sp[0]
	default:
		err = nil
	}

	return
}

// isTLSListener Проверка того, что слушатель соединений уже обёрнут в TLS режим.
// Тип слушателя TLS из пакета "crypto/tls" не экспортируется, поэтому проверка выполняется по пакету типа.
func isTLSListener(ltn net.Listener) (ok bool) {
	const pkgTls = "crypto/tls"
	var tp = reflect.TypeOf(ltn)

	if tp == nil {
		return
	}
	for tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	ok = tp.PkgPath() == pkgTls

	return
}

// Заполнение конфигурации значениями по умолчанию.
func defaultConfiguration(cfg *Configuration) {
	const defaultShutdownTimeout = time.Duration(30) * time.Second

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = defaultShutdownTimeout
	}
}
