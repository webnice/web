package web

import (
	"crypto/tls"
	"net"
	"testing"
)

// Функции в этом файле невозможно протестировать в нормальных условиях,
// требуется linux + systemd + служба запускаемая из под systemd, так как файловые дескрипторы процессов имеют
// строгую нумерацию порядка.
// Поэтому тестируется только то, что функции есть, инициализируются конструктором и вызываются
// соответствующими функциями.

func TestImpl_ListenersSystemdWithoutNames(t *testing.T) {
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdWithoutNames == nil {
		t.Fatalf("конструктор New() не установил функцию listenersSystemdWithoutNames()")
	}
	web.listenersSystemdWithoutNames = func() (_ []net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdWithoutNames(); inc == 0 {
		t.Fatalf("функция ListenersSystemdWithoutNames() не вызывает функцию listenersSystemdWithoutNames()")
	}
}

func TestImpl_ListenersSystemdWithNames(t *testing.T) {
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdWithNames == nil {
		t.Fatalf("конструктор New() не установил функцию listenersSystemdWithNames()")
	}
	web.listenersSystemdWithNames = func() (_ map[string][]net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdWithNames(); inc == 0 {
		t.Fatalf("функция ListenersSystemdWithNames() не вызывает функцию listenersSystemdWithNames()")
	}
}

func TestImpl_ListenersSystemdTLSWithoutNames(t *testing.T) {
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdTLSWithoutNames == nil {
		t.Fatalf("конструктор New() не установил функцию listenersSystemdTLSWithoutNames()")
	}
	web.listenersSystemdTLSWithoutNames = func(_ *tls.Config) (_ []net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdTLSWithoutNames(nil); inc == 0 {
		t.Fatalf("функция ListenersSystemdTLSWithoutNames() не вызывает функцию listenersSystemdTLSWithoutNames()")
	}
}

func TestImpl_ListenersSystemdTLSWithNames(t *testing.T) {
	var (
		web *impl
		inc int
	)

	web = New().(*impl)
	if web.listenersSystemdTLSWithNames == nil {
		t.Fatalf("конструктор New() не установил функцию listenersSystemdTLSWithNames()")
	}
	web.listenersSystemdTLSWithNames = func(_ *tls.Config) (_ map[string][]net.Listener, _ error) { inc++; return }
	if _, _ = web.ListenersSystemdTLSWithNames(nil); inc == 0 {
		t.Fatalf("функция ListenersSystemdTLSWithNames() не вызывает функцию listenersSystemdTLSWithNames()")
	}
}
