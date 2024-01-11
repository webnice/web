package web

import (
	"errors"
	"net"
	"os"
	"testing"

	wnet "github.com/webnice/net"
)

type tmpFile struct {
	Filename string
}

func (tfo *tmpFile) Clean() { _ = os.RemoveAll(tfo.Filename) }

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

func getKeyEcdsa() []byte {
	return []byte(`-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDCJGaDfRSjCg2zYopdy
M7SqBKeIpcEriH3GWTtwy3hlQSiloiyGOk25Ekpt/Ha04PahZANiAARu/6BxP3/t
kYuOdvDeAKD9fsC2m3pLEOzM+ZY8phS7qg4CTFT7Yej8UTaEX1WSd4Sq5F/zmLto
BE2ulX63u0MqdUd/GU6XIpn31kDt2MVqKgprixw7Ow3zIH47KDvdwa0=
-----END PRIVATE KEY-----`)
}

func getCrtEcdsa() []byte {
	return []byte(`-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`)
}

func TestImpl_ListenAndServeTLS_Ok(t *testing.T) {
	const testAddress = `localhost:8088`
	var (
		err error
		key *tmpFile
		crt *tmpFile
		web Interface
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress, crt.Filename, key.Filename, nil)
	if err = web.Error(); err != nil {
		t.Errorf("функция ListenAndServeTLS(), ошибка: %v, ожидалось: %v", err, nil)
	}
	if err = web.
		Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
}

func TestImpl_ListenAndServeTLS_Port(t *testing.T) {
	const invalidAddress = `:170000`
	var (
		key *tmpFile
		crt *tmpFile
		web Interface
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		ListenAndServeTLS(invalidAddress, crt.Filename, key.Filename, nil)
	if web.Error() == nil {
		t.Errorf("функция ListenAndServeTLS(), не корректная проверка адреса")
	}
}

func TestImpl_ListenAndServeTLSWithConfig(t *testing.T) {
	var (
		err error
		web Interface
		cfg *Configuration
	)

	web = New()
	web.ListenAndServeTLSWithConfig(nil, nil)
	defer web.Stop()
	if web.Error() == nil {
		t.Errorf("функция ListenAndServeTLSWithConfig(), не корректная проверка адреса")
	}
	if !errors.Is(web.Error(), Errors().NoConfiguration()) {
		t.Errorf("функция ListenAndServeTLSWithConfig(), получена не корректная ошибка")
	}
	cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host: "localhost",
			Port: 8088,
		},
	}
	if web = web.
		ListenAndServeTLSWithConfig(cfg, nil); web.Error() == nil {
		t.Errorf("функция ListenAndServeTLSWithConfig(), ошибка: %v, ожидалась ошибка", err)
	}
}

func TestImpl_NewListenerTLS(t *testing.T) {
	var (
		err error
		web Interface
		key *tmpFile
		crt *tmpFile
		cfg *Configuration
		lst net.Listener
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	cfg = new(Configuration)

	cfg.Host, cfg.Port = "localhost", 8088
	//defaultConfiguration(cfg)

	web = New().
		Handler(getTestHandlerFn(t))
	if lst, err = web.
		NewListenerTLS(cfg, nil); err == nil {
		t.Errorf("функция NewListenerTLS(), функция повреждена")
	}
	cfg.TLSPrivateKeyPEM, cfg.TLSPublicKeyPEM, cfg.Port = key.Filename, crt.Filename, 8188
	if lst, err = web.
		NewListenerTLS(cfg, nil); err != nil {
		t.Errorf("функция NewListenerTLS(), ошибка: %v, ожидалось: %v", err, nil)
	}
	web = web.Serve(lst)
	if err = web.Error(); err != nil {
		t.Errorf("функция Serve(), функция повреждена")
	}
	if err = web.Stop().
		Error(); err != nil {
		t.Errorf("функция Stop(), ошибка: %v, ожидалось: %v", err, nil)
	}
}

func TestImpl_NewListenerTLS_AlreadyRunning(t *testing.T) {
	const (
		testAddress1 = `localhost:18080`
		testAddress2 = `localhost:18081`
	)
	var (
		web Interface
		key *tmpFile
		crt *tmpFile
	)

	key, crt = newTmpFile(getKeyEcdsa()), newTmpFile(getCrtEcdsa())
	defer func() { key.Clean(); crt.Clean() }()
	web = New().
		Handler(getTestHandlerFn(t)).
		ListenAndServeTLS(testAddress1, crt.Filename, key.Filename, nil)
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf("функция ListenAndServe(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	web.ListenAndServeTLS(testAddress2, crt.Filename, key.Filename, nil)
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
