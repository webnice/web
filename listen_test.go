package web

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func getTestHandlerFn(t *testing.T) (ret http.Handler) {
	var mux = chi.NewMux()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("ok")) })

	return mux
}

func TestNew_InvalidAddress(t *testing.T) {
	const invalidAddress = ":170000"
	var web Interface

	web = New().
		ListenAndServe(invalidAddress)
	if web.Error() == nil {
		t.Errorf("функция ListenAndServe(), не корректная проверка адреса")
	}
	web = New().
		ListenAndServeTLS(invalidAddress, "", "", nil)
	if web.Error() == nil {
		t.Errorf("функция ListenAndServe(), не корректная проверка адреса")
	}
}

func TestNew_NoConfiguration(t *testing.T) {
	var web = New()

	web.ListenAndServeWithConfig(nil)
	defer web.Stop()
	if web.Error() == nil {
		t.Errorf("функция ListenAndServe(), не корректная проверка адреса")
	}
	if !errors.Is(web.Error(), Errors().NoConfiguration()) {
		t.Errorf("функция ListenAndServe(), получена не корректная ошибка")
	}
}

func TestImpl_ListenAndServe_AlreadyRunning(t *testing.T) {
	const (
		testAddress1 = ":18080"
		testAddress2 = ":18081"
	)
	var web Interface

	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServe(testAddress1)
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	web.ListenAndServe(testAddress2)
	if web.Error() == nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", web.Error(), Errors().AlreadyRunning())
	}
	if !errors.Is(web.Error(), Errors().AlreadyRunning()) {
		t.Errorf("функция ListenAndServe(), не корректная ошибка")
	}
	if !errors.Is(web.
		Clean().                      // Очистка последней ошибки.
		ListenAndServe(testAddress1). // Запуск сервера, который уже запущен.
		Error(), Errors().AlreadyRunning()) {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", web.Error(), Errors().AlreadyRunning())
	}
}

func TestImpl_NewListener_NoConfiguration(t *testing.T) {
	var web Interface

	web = New()
	if _, err := web.NewListener(nil); !errors.Is(err, Errors().NoConfiguration()) {
		t.Errorf("функция NewListener(), ошибка: %v, ожидалось: %v", err, Errors().NoConfiguration())
	}
}

func TestImpl_NewListenerTLS_NoConfiguration(t *testing.T) {
	var web Interface

	web = New()
	if _, err := web.NewListenerTLS(nil, nil); !errors.Is(err, Errors().NoConfiguration()) {
		t.Errorf("функция NewListenerTLS(), ошибка: %v, ожидалось: %v", err, Errors().NoConfiguration())
	}
}

func TestImpl_ListenAndServe_Wait(t *testing.T) {
	const (
		testAddress1 = "127.0.0.1:1080"
		testAddress2 = ".test.socket"
		ticTimeout   = time.Second / 4
	)
	var (
		tic  *time.Ticker
		cou  uint32
		web  Interface
		conf *Configuration
	)

	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServe(testAddress1)
	if web.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалась: %v", web.Error(), nil)
	}
	go func() {
		tic = time.NewTicker(ticTimeout)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				web.Stop()
				break
			}
		}
	}()
	web.Wait()
	if cou <= 4 {
		t.Errorf("функция Wait() повреждена")
	}
	web = New().
		Handler(getTestHandlerFn(t))
	conf, _ = parseAddress("")
	conf.Mode = "socket"
	conf.Socket = testAddress2
	web.ListenAndServeWithConfig(conf)
	if web.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалась: %v", web.Error(), nil)
	}
	go func() {
		tic = time.NewTicker(ticTimeout)
		defer tic.Stop()
		for {
			<-tic.C
			if cou++; cou > 4 {
				web.Stop()
				break
			}
		}
	}()
	web.Wait()
	if cou <= 4 {
		t.Errorf("функция Wait() повреждена")
	}
}
