package web

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	wnet "github.com/webnice/net"
)

const testSchemeHTTPS = "https://"

// Тестовый ECDSA секретный ключ в PEM формате.
const testKeyEcdsa = `-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDCJGaDfRSjCg2zYopdy
M7SqBKeIpcEriH3GWTtwy3hlQSiloiyGOk25Ekpt/Ha04PahZANiAARu/6BxP3/t
kYuOdvDeAKD9fsC2m3pLEOzM+ZY8phS7qg4CTFT7Yej8UTaEX1WSd4Sq5F/zmLto
BE2ulX63u0MqdUd/GU6XIpn31kDt2MVqKgprixw7Ow3zIH47KDvdwa0=
-----END PRIVATE KEY-----`

// Тестовый ECDSA публичный ключ (сертификат) в PEM формате.
const testCrtEcdsa = `-----BEGIN CERTIFICATE-----
MIICHjCCAaSgAwIBAgIUCUXV3GKvC6699j1lRVwby9bvbVQwCgYIKoZIzj0EAwIw
RTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGElu
dGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMzExMjcxNjI4MDFaGA8yMTIzMTEw
MzE2MjgwMVowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAf
BgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDB2MBAGByqGSM49AgEGBSuB
BAAiA2IABG7/oHE/f+2Ri4528N4AoP1+wLabeksQ7Mz5ljymFLuqDgJMVPth6PxR
NoRfVZJ3hKrkX/OYu2gETa6Vfre7Qyp1R38ZTpcimffWQO3YxWoqCmuLHDs7DfMg
fjsoO93BraNTMFEwHQYDVR0OBBYEFB+2RT/hFPvsibVkD5YixFyDSMAmMB8GA1Ud
IwQYMBaAFB+2RT/hFPvsibVkD5YixFyDSMAmMA8GA1UdEwEB/wQFMAMBAf8wCgYI
KoZIzj0EAwIDaAAwZQIwUDdaraaLyrL2+Lmj0xTvPI4+zUJ2qVPcVMgzKiElDCk+
75dcyLpKakC/CIgcIzhOAjEAouUFwLdcMQClcToO2jH0n3KWIJXv/y/X3yLN/t+W
4R/MDRGuSjdIm80Rtr3DbnB7
-----END CERTIFICATE-----`

type tmpFile struct {
	Filename string
}

// Clean Удаление временного файла.
func (tfo *tmpFile) Clean() { _ = os.RemoveAll(tfo.Filename) }

// Создание временного файла с содержимым content.
func newTmpFile(content []byte) (ret *tmpFile) {
	var (
		err error
		fh  *os.File
	)

	ret = new(tmpFile)
	if fh, err = os.CreateTemp(os.TempDir(), ""); err != nil {
		ret = nil
		return
	}
	defer func() { _ = fh.Close() }()
	ret.Filename = fh.Name()
	_, _ = fh.Write(content)

	return
}

func TestImpl_ListenAndServeTLS_Ok(t *testing.T) {
	const (
		tplListenError = "функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v"
		tplStopError   = "функция Stop(), ошибка: %v, ожидалось: %v"
		testAddress    = "127.0.0.1:8088"
	)
	var (
		err error
		key *tmpFile
		crt *tmpFile
		web Interface
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress, crt.Filename, key.Filename, nil)
	if err = web.Error(); err != nil {
		t.Errorf(tplListenError, err, nil)
	}
	if err = web.
		Stop().
		Error(); err != nil {
		t.Errorf(tplStopError, err, nil)
	}
}

func TestImpl_ListenAndServeTLS_Port(t *testing.T) {
	const (
		tplAddressNotChecked = "функция ListenAndServeTLS(), не корректная проверка адреса"
		invalidAddress       = ":170000"
	)
	var (
		key *tmpFile
		crt *tmpFile
		web Interface
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		ListenAndServeTLS(invalidAddress, crt.Filename, key.Filename, nil)
	if web.Error() == nil {
		t.Error(tplAddressNotChecked)
	}
}

func TestImpl_ListenAndServeTLSWithConfig(t *testing.T) {
	const (
		tplNoErrorExpected = "функция ListenAndServeTLSWithConfig(), не корректная проверка конфигурации"
		tplWrongError      = "функция ListenAndServeTLSWithConfig(), получена не корректная ошибка"
		tplErrorExpected   = "функция ListenAndServeTLSWithConfig(), ошибка: %v, ожидалась ошибка"
		testHostLocal      = "localhost"
		testPort           = 8088
	)
	var (
		err error
		web Interface
		cfg *Configuration
	)

	web = New()
	web.ListenAndServeTLSWithConfig(nil, nil)
	defer web.Stop()
	if web.Error() == nil {
		t.Error(tplNoErrorExpected)
	}
	if !errors.Is(web.Error(), Errors().NoConfiguration.Bind()) {
		t.Error(tplWrongError)
	}
	// Запуск сервера с конфигурацией, не содержащей TLS ключей.
	cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host: testHostLocal,
			Port: testPort,
		},
	}
	if web = web.
		ListenAndServeTLSWithConfig(cfg, nil); web.Error() == nil {
		t.Errorf(tplErrorExpected, err)
	}
}

func TestImpl_NewListenerTLS(t *testing.T) {
	const (
		tplAddressFormat = "%s:%d"
		tplBroken        = "функция NewListenerTLS(), функция повреждена"
		tplListenerError = "функция NewListenerTLS(), ошибка: %v, ожидалось: %v"
		tplServeBroken   = "функция Serve(), функция повреждена"
		tplStopError     = "функция Stop(), ошибка: %v, ожидалось: %v"
		testHost         = "127.0.0.1"
		testPort1        = 8088
		testPort2        = 8188
	)
	var (
		err error
		web Interface
		key *tmpFile
		crt *tmpFile
		cfg *Configuration
		lst net.Listener
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	if cfg, err = parseAddress(fmt.Sprintf(tplAddressFormat, testHost, testPort1)); err != nil {
		t.Fatalf(tplListenerError, err, nil)
	}
	defaultConfiguration(cfg)
	web = New().
		Handler(getTestHandlerFn(t))
	// Создание слушателя без TLS ключей в конфигурации.
	if lst, err = web.
		NewListenerTLS(cfg, nil); err == nil {
		t.Errorf(tplBroken)
	}
	// Создание слушателя с корректными TLS ключами и запуск сервера.
	cfg.TLSPrivateKeyPEM, cfg.TLSPublicKeyPEM, cfg.Port = key.Filename, crt.Filename, testPort2
	if lst, err = web.
		NewListenerTLS(cfg, nil); err != nil {
		t.Errorf(tplListenerError, err, nil)
	}
	web = web.Serve(lst)
	if err = web.Error(); err != nil {
		t.Errorf(tplServeBroken)
	}
	if err = web.Stop().
		Error(); err != nil {
		t.Errorf(tplStopError, err, nil)
	}
}

// Функция возвращает HTTP клиента для тестирования TLS серверов, доверяющего любому сертификату сервера.
func getTestHTTPSClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// Выполнение тестового HTTPS запроса к серверу, проверка кода ответа.
func requestTestHTTPS(address string) (err error) {
	const (
		tplGetError  = "функция Get(), ошибка: %w"
		tplWrongCode = "функция Get(), код ответа: %d, ожидалось: %d"
		testHTTPTOut = time.Second * 5
	)
	var resp *http.Response

	if resp, err = getTestHTTPSClient(testHTTPTOut).Get(testSchemeHTTPS + address + testRouteRoot); err != nil {
		err = fmt.Errorf(tplGetError, err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(tplWrongCode, resp.StatusCode, http.StatusOK)
	}

	return
}

// Проверка полного цикла HTTPS запроса к серверу запущенному функцией ListenAndServeTLS().
func TestImpl_ListenAndServeTLS_Request(t *testing.T) {
	const (
		tplFatalError = "функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		testAddress   = "127.0.0.1:18443"
	)
	var (
		key *tmpFile
		crt *tmpFile
		web Interface
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress, crt.Filename, key.Filename, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplFatalError, web.Error(), nil)
	}
	if err := requestTestHTTPS(testAddress); err != nil {
		t.Errorf(tplRequest, err)
	}
}

// Запуск сервера в режиме TLS через функцию Serve() с незащищённым слушателем: оборачивание слушателя в TLS
// выполняется функцией ServeTLSWithId().
func TestImpl_ServeTLSPlainListener_Request(t *testing.T) {
	const (
		tplFatalError = "функция NewListener(), ошибка: %v, ожидалось: %v"
		tplServeError = "функция Serve(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		testAddress   = "127.0.0.1:18444"
		testHostIP    = "127.0.0.1"
	)
	var (
		err      error
		key      *tmpFile
		crt      *tmpFile
		web      Interface
		listener net.Listener
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New()
	web.Handler(getTestHandlerFn(t))
	web.(*impl).cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host:             testHostIP,
			Port:             18444,
			TLSPublicKeyPEM:  crt.Filename,
			TLSPrivateKeyPEM: key.Filename,
		},
	}
	if listener, err = web.NewListener(web.(*impl).cfg); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.Serve(listener)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplServeError, web.Error(), nil)
	}
	if err = requestTestHTTPS(testAddress); err != nil {
		t.Errorf(tplRequest, err)
	}
}

// Запуск сервера через функцию Serve() со слушателем, уже находящимся в режиме TLS: повторное оборачивание
// слушателя в TLS выполняться не должно.
func TestImpl_NewListenerTLS_Request(t *testing.T) {
	const (
		tplFatalError = "функция NewListenerTLS(), ошибка: %v, ожидалось: %v"
		tplServeError = "функция Serve(), ошибка: %v, ожидалось: %v"
		tplRequest    = "%v"
		testAddress   = "127.0.0.1:18445"
		testHostIP    = "127.0.0.1"
	)
	var (
		err      error
		key      *tmpFile
		crt      *tmpFile
		web      Interface
		listener net.Listener
		cfg      *Configuration
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().Handler(getTestHandlerFn(t))
	cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host:             testHostIP,
			Port:             18445,
			TLSPublicKeyPEM:  crt.Filename,
			TLSPrivateKeyPEM: key.Filename,
		},
	}
	if listener, err = web.NewListenerTLS(cfg, nil); err != nil {
		t.Fatalf(tplFatalError, err, nil)
	}
	web.Serve(listener)
	defer web.Stop()
	if web.Error() != nil {
		t.Fatalf(tplServeError, web.Error(), nil)
	}
	if err = requestTestHTTPS(testAddress); err != nil {
		t.Errorf(tplRequest, err)
	}
}

//goland:noinspection DuplicatedCode
func TestImpl_NewListenerTLS_AlreadyRunning(t *testing.T) {
	const (
		tplErrorNotExpected = "функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v"
		tplErrorExpected    = "функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v"
		tplWrongError       = "функция ListenAndServeTLS(), получена не корректная ошибка"
		testAddress1        = "127.0.0.1:18080"
		testAddress2        = "127.0.0.1:18081"
	)
	var (
		web Interface
		key *tmpFile
		crt *tmpFile
	)

	key, crt = newTmpFile([]byte(testKeyEcdsa)), newTmpFile([]byte(testCrtEcdsa))
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress1, crt.Filename, key.Filename, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf(tplErrorNotExpected, web.Error(), nil)
	}
	// Повторный запуск сервера, сервер уже запущен.
	web.ListenAndServeTLS(testAddress2, crt.Filename, key.Filename, nil)
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
