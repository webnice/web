package web

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func getTestHandlerFn(_ *testing.T) (ret http.Handler) {
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

// Проверка ограничения максимального размера тела запроса: превышение лимита вызывает ответ 413.
//
//goland:noinspection ALL
func TestImpl_ListenAndServe_MaxBodyBytes(t *testing.T) {
	const (
		testAddress  = "127.0.0.1:18450"
		testLimit    = 4
		testBodyOk   = "1234"
		testBodyBig  = "12345"
		testResponse = http.StatusRequestEntityTooLarge
	)
	var (
		err      error
		web      Interface
		cfg      *Configuration
		handler  http.Handler
		request  *http.Request
		response *http.Response
	)

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, er := io.Copy(io.Discard, r.Body); er != nil {
			// Превышение лимита размера тела запроса, сервер отвечает ошибкой 413.
			return
		}
		_, _ = w.Write([]byte("ok"))
	})
	cfg, _ = parseAddress(testAddress)
	cfg.MaxBodyBytes = testLimit
	web = New().
		Handler(handler).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	// Тело запроса меньше лимита: код ответа 200.
	if request, err = http.NewRequest(http.MethodPost, "http://"+testAddress+"/", strings.NewReader(testBodyOk)); err != nil {
		t.Fatalf("функция NewRequest(), ошибка: %v", err)
	}
	if response, err = http.DefaultClient.Do(request); err != nil {
		t.Fatalf("функция Do(), ошибка: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf("функция Do(), код ответа: %d, ожидалось: %d", response.StatusCode, http.StatusOK)
	}
	// Тело запроса больше лимита: код ответа 413.
	if request, err = http.NewRequest(http.MethodPost, "http://"+testAddress+"/", strings.NewReader(testBodyBig)); err != nil {
		t.Fatalf("функция NewRequest(), ошибка: %v", err)
	}
	if response, err = http.DefaultClient.Do(request); err != nil {
		t.Fatalf("функция Do(), ошибка: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != testResponse {
		t.Errorf("функция Do(), код ответа: %d, ожидалось: %d", response.StatusCode, testResponse)
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
