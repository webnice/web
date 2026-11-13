package web

import (
	"crypto/tls"
	"net"
	"testing"
)

// Функции в этом файле невозможно протестировать в нормальных условиях,
// требуется Linux + systemd + служба запускаемая из-под systemd, так как файловые дескрипторы процессов имеют
// строгую нумерацию порядка.
// Поэтому тестируется только то, что функции есть, инициализируются конструктором и вызываются
// соответствующими функциями.

func TestImpl_ListenersSystemdWithoutNames(t *testing.T) {
	const (
		tplConstructor = "конструктор New() не установил функцию listenersSystemdWithoutNames()"
		tplNotCalling  = "функция ListenersSystemdWithoutNames() не вызывает функцию listenersSystemdWithoutNames()"
	)
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdWithoutNames == nil {
		t.Fatalf(tplConstructor)
	}
	web.listenersSystemdWithoutNames = func() (_ []net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdWithoutNames(); inc == 0 {
		t.Fatalf(tplNotCalling)
	}
}

func TestImpl_ListenersSystemdWithNames(t *testing.T) {
	const (
		tplConstructor = "конструктор New() не установил функцию listenersSystemdWithNames()"
		tplNotCalling  = "функция ListenersSystemdWithNames() не вызывает функцию listenersSystemdWithNames()"
	)
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdWithNames == nil {
		t.Fatalf(tplConstructor)
	}
	web.listenersSystemdWithNames = func() (_ map[string][]net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdWithNames(); inc == 0 {
		t.Fatalf(tplNotCalling)
	}
}

func TestImpl_ListenersSystemdTLSWithoutNames(t *testing.T) {
	const (
		tplConstructor = "конструктор New() не установил функцию listenersSystemdTLSWithoutNames()"
		tplNotCalling  = "функция ListenersSystemdTLSWithoutNames() не вызывает функцию listenersSystemdTLSWithoutNames()"
	)
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdTLSWithoutNames == nil {
		t.Fatalf(tplConstructor)
	}
	web.listenersSystemdTLSWithoutNames = func(_ *tls.Config) (_ []net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdTLSWithoutNames(nil); inc == 0 {
		t.Fatalf(tplNotCalling)
	}
}

func TestImpl_ListenersSystemdTLSWithNames(t *testing.T) {
	const (
		tplConstructor = "конструктор New() не установил функцию listenersSystemdTLSWithNames()"
		tplNotCalling  = "функция ListenersSystemdTLSWithNames() не вызывает функцию listenersSystemdTLSWithNames()"
	)
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdTLSWithNames == nil {
		t.Fatalf(tplConstructor)
	}
	web.listenersSystemdTLSWithNames = func(_ *tls.Config) (_ map[string][]net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdTLSWithNames(nil); inc == 0 {
		t.Fatalf(tplNotCalling)
	}
}
