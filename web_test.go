package web

import (
	"errors"
	"net"
	"testing"

	"github.com/google/uuid"
)

func TestImpl_IdGenerate(t *testing.T) {
	const (
		tplErrorNotExpected = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplIdIsEmpty        = "функция ID(), ошибка, ожидалось не пустое значение"
		testHostLocal       = "localhost"
		testPortLocal       = 10180
		testModeTcp         = "tcp"
	)
	var web Interface

	web = New().
		Handler(getTestHandlerFn(t))
	web.ListenAndServeWithConfig(&Configuration{
		Host: testHostLocal,
		Port: testPortLocal,
		Mode: testModeTcp,
	})
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	if web.ID() == "" {
		t.Error(tplIdIsEmpty)
	}
}

func TestImpl_IdStatic(t *testing.T) {
	const (
		tplErrorNotExpected = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplIdNotEqual       = "функция ID(), вернулось: %q, ожидалось: %q"
		testHostLocal       = "localhost"
		testPortLocal       = 10180
		testModeTcp         = "tcp"
	)
	var (
		id  string
		web Interface
	)

	id = uuid.NewString()
	web = New().
		Handler(getTestHandlerFn(t))
	web.ListenAndServeWithConfig(&Configuration{
		ID:   id,
		Host: testHostLocal,
		Port: testPortLocal,
		Mode: testModeTcp,
	})
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	if web.ID() != id {
		t.Errorf(tplIdNotEqual, web.ID(), id)
	}
}

func TestImpl_ServeWithId(t *testing.T) {
	const (
		tplNewListenerError = "функция NewListener(), ошибка: %v, ожидалось: %v"
		tplServeError       = "функция ServeWithId(), ошибка: %v, ожидалось: %v"
		tplIdNotEqual       = "функция ID(), вернулось: %q, ожидалось: %q"
		testHostLocal       = "localhost"
		testPortLocal       = 18080
	)
	var (
		err      error
		id       string
		web      Interface
		listener net.Listener
	)

	id = uuid.NewString()
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.(*impl).cfg = &Configuration{
		Host: testHostLocal,
		Port: testPortLocal,
	}
	if listener, err = web.NewListener(web.(*impl).cfg); err != nil {
		t.Fatalf(tplNewListenerError, err, nil)
	}
	if web.ServeWithId(listener, id); web.Error() != nil {
		t.Errorf(tplServeError, web.Error(), nil)
	}
	defer web.Stop()
	if web.ID() != id {
		t.Errorf(tplIdNotEqual, web.ID(), id)
	}
}

// Запуск сервера через функцию Serve() без указания ID не должен изменять ID переданный в конфигурации.
func TestImpl_ServeConfiguredID(t *testing.T) {
	const (
		tplNewListenerError = "функция NewListener(), ошибка: %v, ожидалось: %v"
		tplServeError       = "функция Serve(), ошибка: %v, ожидалось: %v"
		tplIdNotEqual       = "функция ID(), вернулось: %q, ожидалось: %q"
		testHostLocal       = "localhost"
		testPort            = 18082
	)
	var (
		err      error
		id       string
		web      Interface
		listener net.Listener
	)

	id = uuid.NewString()
	web = New().(*impl)
	web.Handler(getTestHandlerFn(t))
	web.(*impl).cfg = &Configuration{
		ID:   id,
		Host: testHostLocal,
		Port: testPort,
	}
	if listener, err = web.NewListener(web.(*impl).cfg); err != nil {
		t.Fatalf(tplNewListenerError, err, nil)
	}
	web.Serve(listener)
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf(tplServeError, web.Error(), nil)
	}
	if web.ID() != id {
		t.Errorf(tplIdNotEqual, web.ID(), id)
	}
}

// Проверка получения ошибки библиотеки "github.com/webnice/net" через функцию Error() при отсутствии
// собственной ошибки пакета.
func TestImpl_ErrorFromNet(t *testing.T) {
	const tplErrorNotEqual = "функция Error(), ошибка: %v, ожидалось: %v"
	var web *impl

	web = New().(*impl)
	if web.Error() != nil {
		t.Fatalf(tplErrorNotEqual, web.Error(), nil)
	}
	// Формирование ошибки внутри библиотеки "github.com/webnice/net" без ошибки пакета.
	web.net.ListenAndServeWithConfig(nil)
	if !errors.Is(web.Error(), Errors().NoConfiguration.Bind()) {
		t.Errorf(tplErrorNotEqual, web.Error(), Errors().NoConfiguration.Bind())
	}
}
