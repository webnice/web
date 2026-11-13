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

//goland:noinspection HttpUrlsUsage
const (
	testRouteRoot    = "/"
	testResponseBody = "ok"
	testModeSocket   = "socket"
	testAddressEmpty = ""
	testSchemeHTTP   = "http://"
)

// Тестовый обработчик запросов: отвечает 200 OK на корневой маршрут.
func getTestHandlerFn(_ *testing.T) (ret http.Handler) {
	var mux = chi.NewMux()

	mux.Get(testRouteRoot, func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(testResponseBody)) })

	return mux
}

func TestNew_InvalidAddress(t *testing.T) {
	const (
		tplAddressNotChecked = "функция ListenAndServe(), не корректная проверка адреса"
		invalidAddress       = ":170000"
	)
	var web Interface

	web = New().
		ListenAndServe(invalidAddress)
	if web.Error() == nil {
		t.Error(tplAddressNotChecked)
	}
	web = New().
		ListenAndServeTLS(invalidAddress, "", "", nil)
	if web.Error() == nil {
		t.Error(tplAddressNotChecked)
	}
}

func TestNew_NoConfiguration(t *testing.T) {
	const (
		tplAddressNotChecked = "функция ListenAndServeWithConfig(), не корректная проверка конфигурации"
		tplWrongError        = "функция ListenAndServeWithConfig(), получена не корректная ошибка"
	)
	var web Interface

	web = New()
	web.ListenAndServeWithConfig(nil)
	defer web.Stop()
	if web.Error() == nil {
		t.Error(tplAddressNotChecked)
	}
	if !errors.Is(web.Error(), Errors().NoConfiguration.Bind()) {
		t.Error(tplWrongError)
	}
}

//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServe_AlreadyRunning(t *testing.T) {
	const (
		tplErrorNotExpected = "функция ListenAndServe(), ошибка: %v, ожидалось: %v"
		tplErrorExpected    = "функция ListenAndServe(), ошибка: %v, ожидалось: %v"
		tplWrongError       = "функция ListenAndServe(), получена не корректная ошибка"
		testAddress1        = ":18080"
		testAddress2        = ":18081"
	)
	var web Interface

	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServe(testAddress1)
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	// Повторный запуск сервера, сервер уже запущен.
	web.ListenAndServe(testAddress2)
	if web.Error() == nil {
		t.Errorf(tplErrorExpected, web.Error(), Errors().AlreadyRunning.Bind())
	}
	if !errors.Is(web.Error(), Errors().AlreadyRunning.Bind()) {
		t.Error(tplWrongError)
	}
	// Запуск сервера после очистки последней ошибки, сервер уже запущен.
	if !errors.Is(web.
		Clean().
		ListenAndServe(testAddress1).
		Error(), Errors().AlreadyRunning.Bind()) {
		t.Errorf(tplErrorExpected, web.Error(), Errors().AlreadyRunning.Bind())
	}
}

func TestImpl_NewListener_NoConfiguration(t *testing.T) {
	const tplErrorNotEqual = "функция NewListener(), ошибка: %v, ожидалось: %v"
	var web Interface

	web = New()
	if _, err := web.NewListener(nil); !errors.Is(err, Errors().NoConfiguration.Bind()) {
		t.Errorf(tplErrorNotEqual, err, Errors().NoConfiguration.Bind())
	}
}

func TestImpl_NewListenerTLS_NoConfiguration(t *testing.T) {
	const tplErrorNotEqual = "функция NewListenerTLS(), ошибка: %v, ожидалось: %v"
	var web Interface

	web = New()
	if _, err := web.NewListenerTLS(nil, nil); !errors.Is(err, Errors().NoConfiguration.Bind()) {
		t.Errorf(tplErrorNotEqual, err, Errors().NoConfiguration.Bind())
	}
}

// Проверка ограничения максимального размера тела запроса: превышение лимита вызывает ответ 413.
//
//goland:noinspection ALL
func TestImpl_ListenAndServe_MaxBodyBytes(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplNewRequest = "функция NewRequest(), ошибка: %v"
		tplDoError    = "функция Do(), ошибка: %v"
		tplWrongCode  = "функция Do(), код ответа: %d, ожидалось: %d"
		testAddress   = "127.0.0.1:18450"
		testLimit     = 4
		testBodyOk    = "1234"
		testBodyBig   = "12345"
		testResponse  = http.StatusRequestEntityTooLarge
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
		_, _ = w.Write([]byte(testResponseBody))
	})
	if cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	cfg.MaxBodyBytes = testLimit
	web = New().
		Handler(handler).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	// Тело запроса меньше лимита: код ответа 200.
	if request, err = http.NewRequest(http.MethodPost, testSchemeHTTP+testAddress+testRouteRoot, strings.NewReader(testBodyOk)); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	if response, err = http.DefaultClient.Do(request); err != nil {
		t.Fatalf(tplDoError, err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Errorf(tplWrongCode, response.StatusCode, http.StatusOK)
	}
	// Тело запроса больше лимита: код ответа 413.
	if request, err = http.NewRequest(http.MethodPost, testSchemeHTTP+testAddress+testRouteRoot, strings.NewReader(testBodyBig)); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	if response, err = http.DefaultClient.Do(request); err != nil {
		t.Fatalf(tplDoError, err)
	}
	_ = response.Body.Close()
	if response.StatusCode != testResponse {
		t.Errorf(tplWrongCode, response.StatusCode, testResponse)
	}
}

//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServe_Wait(t *testing.T) {
	const (
		tplErrorNotExpected = "функция ListenAndServe(), ошибка: %v, ожидалась: %v"
		tplWaitBroken       = "функция Wait() повреждена"
		testAddress1        = "127.0.0.1:1080"
		testAddress2        = ".test.socket"
		ticTimeout          = time.Second / 4
	)
	var (
		err  error
		tic  *time.Ticker
		cou  uint32
		web  Interface
		conf *Configuration
	)

	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServe(testAddress1)
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	// Остановка сервера по таймеру.
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
		t.Error(tplWaitBroken)
	}
	// Проверка ожидания завершения сервера на UNIX сокете.
	web = New().
		Handler(getTestHandlerFn(t))
	if conf, err = parseAddress(testAddressEmpty); err != nil {
		t.Fatalf(tplErrorNotExpected, err, nil)
	}
	conf.Mode = testModeSocket
	conf.Socket = testAddress2
	web.ListenAndServeWithConfig(conf)
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	// Остановка сервера по таймеру.
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
		t.Error(tplWaitBroken)
	}
}
