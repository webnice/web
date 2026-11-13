package web

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Проверка нормализации списка доменов конфигурации.
func TestNewDomainFilter(t *testing.T) {
	const (
		tplHostsNotEqual  = "список доменов без порта: %v, ожидалось: %v"
		tplPortsNotEqual  = "список доменов с портом: %v, ожидалось: %v"
		tplFilterEmpty    = "фильтр пустого списка доменов не пуст"
		tplFilterNotEmpty = "фильтр не пустого списка доменов пуст"
	)

	var df *domainFilter

	// Пустой список доменов.
	df = newDomainFilter(nil)
	if !df.empty() {
		t.Error(tplFilterEmpty)
	}
	// Домена без порта, с портом, wildcard-шаблон и не корректный домен.
	df = newDomainFilter([]string{
		"example.com", "example.com:8443", "*.example.org", "домен.рф", "не*корректный..домен",
		"*.exa mple.com", "plain exa mple.com",
	})
	if df.empty() {
		t.Error(tplFilterNotEmpty)
	}
	if len(df.hosts) != 3 {
		t.Errorf(tplHostsNotEqual, df.hosts, "3 домена")
	}
	if len(df.hostPorts) != 1 {
		t.Errorf(tplPortsNotEqual, df.hostPorts, "1 домен")
	}
	// Проверка преобразования домена национального алфавита в punycode.
	if df.hosts[2] != "xn--d1acufc.xn--p1ai" {
		t.Errorf(tplHostsNotEqual, df.hosts[2], "xn--d1acufc.xn--p1ai")
	}
}

// Проверка сравнения домена с шаблоном.
func TestMatchDomainPattern(t *testing.T) {
	const (
		tplMatch          = "шаблон %q, домен %q: %t, ожидалось: %t"
		tplListMatchedNil = "функция matchDomainList(), пустой список совпал с доменом"
	)
	var tests = []struct {
		pattern string
		host    string
		ok      bool
	}{
		{pattern: "example.com", host: "example.com", ok: true},
		{pattern: "example.com", host: "other.com", ok: false},
		{pattern: "*.example.com", host: "www.example.com", ok: true},
		{pattern: "*.example.com", host: "a.b.example.com", ok: true},
		{pattern: "*.example.com", host: "example.com", ok: false},
		{pattern: "*.example.com", host: "badexample.com", ok: false},
	}

	for n := range tests {
		if ret := matchDomainPattern(tests[n].pattern, tests[n].host); ret != tests[n].ok {
			t.Errorf(tplMatch, tests[n].pattern, tests[n].host, ret, tests[n].ok)
		}
	}
	// Проверка списка без совпадений.
	if matchDomainList("example.com", nil) {
		t.Error(tplListMatchedNil)
	}
}

// Проверка фильтрации домена по заголовку Host запроса.
//
//goland:noinspection DuplicatedCode
func TestDomainFilterHandler(t *testing.T) {
	const (
		tplWrongCode     = "функция ServeHTTP(), код ответа: %d, ожидалось: %d"
		tplWrongBody     = "функция ServeHTTP(), тело ответа: %q, ожидалось: %q"
		tplWrongHeader   = "функция ServeHTTP(), заголовок Connection: %q, ожидалось: %q"
		tplNewRequest    = "функция NewRequest(), ошибка: %v"
		tplReadAll       = "функция ReadAll(), ошибка: %v"
		testBodyOK       = "ok"
		testHostDenied   = "denied.example.com"
		testHostAllowed  = "allowed.example.com"
		testHostSub      = "www.allowed.example.org"
		testHostWithPort = "allowed.example.com:8443"
		testHostIDN      = "домен.рф"
		testHeaderClose  = "close"
		testHostInvalid  = "exa mple.com"
	)
	var (
		df   *domainFilter
		rec  *httptest.ResponseRecorder
		req  *http.Request
		body []byte
		err  error
	)

	df = newDomainFilter([]string{"allowed.example.com", "*.allowed.example.org", "allowed.example.com:8443", "домен.рф"})
	var handler = df.handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testBodyOK))
	}))
	// Разрешённый домен.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostAllowed+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
	// Разрешённый wildcard-поддомен.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostSub+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
	// Разрешённый домен с портом.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostWithPort+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
	// Разрешённый домен национального алфавита.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostIDN+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
	// Не корректный домен запроса пропускается без блокировки.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostAllowed+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	req.Host = testHostInvalid
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf(tplWrongCode, rec.Code, http.StatusOK)
	}
	// Запрещённый домен: статус 434, закрытие соединения.
	rec = httptest.NewRecorder()
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testHostDenied+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	handler.ServeHTTP(rec, req)
	if rec.Code != Status.RequestedHostUnavailable.Code() {
		t.Errorf(tplWrongCode, rec.Code, Status.RequestedHostUnavailable.Code())
	}
	if rec.Header().Get(Header.Connection.String()) != testHeaderClose {
		t.Errorf(tplWrongHeader, rec.Header().Get(Header.Connection.String()), testHeaderClose)
	}
	if body, err = io.ReadAll(rec.Result().Body); err != nil {
		t.Fatalf(tplReadAll, err)
	}
	if string(body) != Status.RequestedHostUnavailable.String() {
		t.Errorf(tplWrongBody, string(body), Status.RequestedHostUnavailable.String())
	}
}

// Проверка фильтрации домена на этапе TLS handshake по SNI.
func TestDomainFilterSNI(t *testing.T) {
	const (
		tplErrorExpected = "функция sniFilter(), ошибка не вернулась"
		tplErrorNotExp   = "функция sniFilter(), ошибка: %v, ожидалось: %v"
	)
	var (
		df  *domainFilter
		err error
	)

	df = newDomainFilter([]string{"allowed.example.com", "*.allowed.example.org"})
	// Домен не указан, проверка пропускается.
	if _, err = df.sniFilter(&tls.ClientHelloInfo{}); err != nil {
		t.Errorf(tplErrorNotExp, err, nil)
	}
	// Разрешённый домен.
	if _, err = df.sniFilter(&tls.ClientHelloInfo{ServerName: "allowed.example.com"}); err != nil {
		t.Errorf(tplErrorNotExp, err, nil)
	}
	// Разрешённый wildcard-поддомен.
	if _, err = df.sniFilter(&tls.ClientHelloInfo{ServerName: "www.allowed.example.org"}); err != nil {
		t.Errorf(tplErrorNotExp, err, nil)
	}
	// Запрещённый домен.
	if _, err = df.sniFilter(&tls.ClientHelloInfo{ServerName: "denied.example.com"}); err == nil {
		t.Errorf(tplErrorExpected)
	}
	// Не корректный домен пропускается без блокировки.
	if _, err = df.sniFilter(&tls.ClientHelloInfo{ServerName: "exa mple.com"}); err != nil {
		t.Errorf(tplErrorNotExp, err, nil)
	}
}

// Проверка фильтрации домена в полном цикле: запуск сервера и HTTP запросы.
//
//goland:noinspection DuplicatedCode
func TestImpl_ListenAndServe_Domain(t *testing.T) {
	const (
		tplFatalError   = "функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v"
		tplWrongCode    = "функция Get(), код ответа: %d, ожидалось: %d"
		tplNewRequest   = "функция NewRequest(), ошибка: %v"
		tplDoError      = "функция Do(), ошибка: %v"
		testAddress     = "127.0.0.1:18451"
		testHostAllowed = "allowed.example.com"
		testHostDenied  = "denied.example.com"
	)
	var (
		err error
		web Interface
		cfg *Configuration
		req *http.Request
		rsp *http.Response
	)

	if cfg, err = parseAddress(testAddress); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	cfg.Domain = []string{"allowed.example.com"}
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeWithConfig(cfg)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	// Запрос на разрешённый домен: заголовок Host передаётся напрямую.
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testAddress+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	req.Host = testHostAllowed
	if rsp, err = http.DefaultClient.Do(req); err != nil {
		t.Fatalf(tplDoError, err)
	}
	_ = rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		t.Errorf(tplWrongCode, rsp.StatusCode, http.StatusOK)
	}
	// Запрос на запрещённый домен: заголовок Host передаётся напрямую.
	if req, err = http.NewRequest(http.MethodGet, testSchemeHTTP+testAddress+testRouteRoot, nil); err != nil {
		t.Fatalf(tplNewRequest, err)
	}
	req.Host = testHostDenied
	if rsp, err = http.DefaultClient.Do(req); err != nil {
		t.Fatalf(tplDoError, err)
	}
	_ = rsp.Body.Close()
	if rsp.StatusCode != Status.RequestedHostUnavailable.Code() {
		t.Errorf(tplWrongCode, rsp.StatusCode, Status.RequestedHostUnavailable.Code())
	}
}

// Проверка применения настроек TCP keep-alive к конфигурации слушателя.
func TestImpl_ApplyKeepAliveConfig(t *testing.T) {
	const (
		tplKeepAliveNoSet     = "функция applyKeepAliveConfig(), ListenConfig создан без настроек keep-alive"
		tplKeepAliveNotOff    = "функция applyKeepAliveConfig(), keep-alive не отключён"
		tplKeepAliveNotOn     = "функция applyKeepAliveConfig(), keep-alive не включён"
		testKeepAliveIdle     = time.Second * 10
		testKeepAliveInterval = time.Second * 5
		testKeepAliveCount    = 3
	)
	var (
		web  *impl
		conf *Configuration
	)

	web = New().(*impl)
	// Без настроек keep-alive слушатель не изменяется.
	conf = &Configuration{}
	web.applyKeepAliveConfig(conf)
	if conf.ListenConfig != nil {
		t.Error(tplKeepAliveNoSet)
	}
	// Отключение keep-alive.
	conf = &Configuration{KeepAliveDisable: true}
	web.applyKeepAliveConfig(conf)
	if conf.ListenConfig == nil || conf.ListenConfig.KeepAlive != -1 {
		t.Error(tplKeepAliveNotOff)
	}
	// Включение keep-alive с настройками.
	conf = &Configuration{KeepAliveEnable: true, KeepAliveIdle: testKeepAliveIdle, KeepAliveInterval: testKeepAliveInterval, KeepAliveCount: testKeepAliveCount}
	web.applyKeepAliveConfig(conf)
	if conf.ListenConfig == nil || !conf.ListenConfig.KeepAliveConfig.Enable {
		t.Error(tplKeepAliveNotOn)
	}
}
