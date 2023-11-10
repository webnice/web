package web

import (
	"crypto/tls"
	"net"
	"os"
	"path"
)

// ListenAndServe Открытие адреса или сокета без использования конфигурации веб сервера (конфигурация по
// умолчанию), запуск веб сервера для обслуживания входящих соединений.
func (wbo *web) ListenAndServe(addr string) Interface {
	var conf *Configuration

	if conf, wbo.err = parseAddress(addr); wbo.err != nil {
		return wbo
	}

	return wbo.ListenAndServeWithConfig(conf)
}

// ListenAndServeTLS Открытие адреса или сокета с использованием TLS, без использования конфигурации веб сервера
// (конфигурация по умолчанию), запуск веб сервера для обслуживания входящих соединений.
func (wbo *web) ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface {
	var conf *Configuration

	if conf, wbo.err = parseAddress(addr); wbo.err != nil {
		return wbo
	}
	conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM = certFile, keyFile

	return wbo.ListenAndServeTLSWithConfig(conf, tlsConfig)
}

// ListenAndServeWithConfig Настройка сервера с использованием переданной конфигурации, открытие адреса или сокета
// на прослушивание, запуск веб сервера для обслуживания входящих соединений.
func (wbo *web) ListenAndServeWithConfig(conf *Configuration) Interface {
	if conf == nil {
		wbo.err = ErrNoConfiguration()
		return wbo
	}
	wbo.conf = conf

	return wbo.Listen(nil)
}

// ListenAndServeTLSWithConfig Настройка сервера с использованием переданной конфигурации в режиме TLS, открытие
// адреса или сокета на прослушивание, запуск веб сервера для обслуживания входящих соединений.
func (wbo *web) ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface {
	if conf == nil {
		wbo.err = ErrNoConfiguration()
		return wbo
	}
	wbo.conf = conf
	if tlsConfig == nil {
		if tlsConfig, wbo.err = wbo.tlsConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); wbo.err != nil {
			return wbo
		}
	}

	return wbo.Listen(tlsConfig)
}

// NewListener Создание нового слушателя соединений net.Listener на основе конфигурации веб сервера.
func (wbo *web) NewListener(conf *Configuration) (ret net.Listener, err error) {
	var (
		lstWithNames map[string][]net.Listener
		listeners    []net.Listener
		ok           bool
	)

	defaultConfiguration(conf)
	switch conf.Mode {
	case netSystemd:
		if conf.Socket != "" {
			// Имена сокетов указаны
			if lstWithNames, err = wbo.ListenersSystemdWithNames(); err != nil {
				return
			}
			// Выбор сокета по имени
			if listeners, ok = lstWithNames[path.Base(conf.Socket)]; !ok {
				err = ErrListenSystemdNotFound()
				return
			}
		} else {
			// Имена сокетов не указаны
			if listeners, err = wbo.ListenersSystemdWithoutNames(); err != nil {
				return
			}
		}
		if len(listeners) == 0 {
			err = ErrListenSystemdUnexpectedNumber()
			return
		}
		ret = listeners[0]
	case netUnix, netUnixPacket:
		_ = os.Remove(conf.Socket)
		ret, err = net.Listen(conf.Mode, conf.Socket)
		_ = os.Chmod(conf.Socket, os.FileMode(0666))
	default:
		ret, err = net.Listen(conf.Mode, conf.HostPort)
	}

	return
}

// NewListenerTLS Создание нового слушателя соединений net.Listener в режиме TLS, на основе конфигурации
// веб сервера.
func (wbo *web) NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, err error) {
	var lst net.Listener

	if lst, err = wbo.NewListener(conf); err != nil {
		return
	}
	if tlsConfig == nil {
		if tlsConfig, err = wbo.tlsConfigDefault(conf.TLSPublicKeyPEM, conf.TLSPrivateKeyPEM); err != nil {
			return
		}
	}
	if tlsConfig == nil {
		err = ErrTLSIsNil()
		return
	}
	ret = tls.NewListener(lst, tlsConfig)

	return
}

// Конфигурация TLS по умолчанию.
func (wbo *web) tlsConfigDefault(tlsPublicFile string, tlsPrivateFile string) (ret *tls.Config, err error) {
	ret = &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		Certificates: make([]tls.Certificate, 1),
	}
	if ret.Certificates[0], err = tls.LoadX509KeyPair(tlsPublicFile, tlsPrivateFile); err != nil {
		return
	}

	return
}

// Listen Запуск прослушивания входящих соединений и веб сервера.
func (wbo *web) Listen(tlsConfig *tls.Config) Interface {
	var ltn net.Listener

	if wbo.isRun.Load().(bool) {
		wbo.err = ErrAlreadyRunning()
		return wbo
	}
	switch tlsConfig == nil {
	case true:
		ltn, wbo.err = wbo.NewListener(wbo.conf)
	case false:
		ltn, wbo.err = wbo.NewListenerTLS(wbo.conf, tlsConfig)
	}
	if wbo.err != nil {
		return wbo
	}

	return wbo.ServeTLS(ltn, tlsConfig)
}

// Serve Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener.
func (wbo *web) Serve(ltn net.Listener) Interface { return wbo.ServeTLS(ltn, nil) }

// ServeTLS Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// использованием TLS.
func (wbo *web) ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface {
	var (
		conf *Configuration
		onUp chan struct{}
	)

	// TODO: Сделать поддержку PROXY Protocol через "github.com/webnice/web/v3/proxyp", conf.ProxyProtocol

	if wbo.conf == nil {
		conf, _ = parseAddress(ltn.Addr().String())
		defaultConfiguration(conf)
		wbo.conf = conf
	}
	wbo.listener, wbo.onCloseDone = ltn, make(chan struct{})
	wbo.isRun.Store(true)
	onUp = make(chan struct{})
	go wbo.run(onUp, tlsConfig)
	onEnd(onUp)

	return wbo
}

// Процесс веб сервера.
func (wbo *web) run(onUp chan struct{}, tlsConfig *tls.Config) {
	defer func() {
		wbo.isRun.Store(false)
		wbo.onCloseDone <- struct{}{}
	}()
	defer func() {
		if wbo.conf.Socket == "" {
			return
		}
		switch wbo.conf.Mode {
		case netSystemd:
			return
		case netUnix, netUnixPacket:
			_ = os.Remove(wbo.conf.Socket)
		}
	}()
	// Обеспечение синхронного запуска потока.
	onUp <- struct{}{}
	// Присвоение конфигурации веб серверу.
	if wbo.server = wbo.loadConfiguration(tlsConfig); wbo.err != nil {
		return
	}
	// Проверка наличия обработчика запросов ВЕБ сервера.
	if wbo.server.Handler == nil || wbo.handler == nil {
		wbo.err = ErrHandlerIsNotSet()
		return
	}
	// Конфигурация "оставаться в живых".
	if wbo.conf.KeepAliveDisable {
		wbo.server.SetKeepAlivesEnabled(false)
	}
	// Запуск веб сервера.
	if wbo.conf.TLSPrivateKeyPEM == "" || wbo.conf.TLSPublicKeyPEM == "" {
		wbo.err = wbo.server.Serve(wbo.listener)
		return
	}
	wbo.err = wbo.server.ServeTLS(wbo.listener, wbo.conf.TLSPublicKeyPEM, wbo.conf.TLSPrivateKeyPEM)
}

// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
// Если сервер не запущен, функция завершается немедленно.
func (wbo *web) Wait() Interface {
	if !wbo.isRun.Load().(bool) {
		return wbo
	}
	onEnd(wbo.onCloseDone)

	return wbo
}
