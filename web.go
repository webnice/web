package web

import (
	"crypto/tls"
	"net"
	"net/http"

	wnet "github.com/webnice/net"
)

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New() Interface {
	var ret *impl

	ret = &impl{net: wnet.New()}
	// Функции systemd невозможно протестировать в обычных условиях, поэтому они выделены в поля объекта.
	ret.listenersSystemdWithoutNames = ret.net.ListenersSystemdWithoutNames
	ret.listenersSystemdWithNames = ret.net.ListenersSystemdWithNames
	ret.listenersSystemdTLSWithoutNames = ret.net.ListenersSystemdTLSWithoutNames
	ret.listenersSystemdTLSWithNames = ret.net.ListenersSystemdTLSWithNames

	return ret
}

// ID Уникальный идентификатор сервера.
// Если идентификатор не был указан в конфигурации, создаётся при запуске.
func (web *impl) ID() string { return web.net.ID() }

// Handler Назначение обработчика запросов ВЕБ сервера.
// Обработчик необходимо назначить до запуска ВЕБ сервера.
func (web *impl) Handler(handler http.Handler) Interface { web.handler = handler; return web }

// Clean Очистка последней ошибки.
func (web *impl) Clean() Interface { web.err = nil; web.net.Clean(); return web }

// Errors Справочник ошибок.
func (web *impl) Errors() *Error { return errSingleton }

// Error Функция возвращает последнюю ошибку веб сервера или библиотеки "github.com/webnice/net", на которой
// основан ВЕБ сервер.
// Ошибка библиотеки "github.com/webnice/net" приводится к справочнику ошибок пакета для сравнения ошибок
// функцией errors.Is().
func (web *impl) Error() (ret error) {
	switch {
	case web.err != nil:
		ret = web.err
	case web.net.Error() != nil:
		// Приведение ошибки библиотеки "github.com/webnice/net" к справочнику ошибок пакета.
		ret = Errors().ParseError(web.net.Error(), errSingleton)
	}

	return
}

// Serve Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener.
func (web *impl) Serve(ltn net.Listener) Interface { return web.ServeTLS(ltn, nil) }

// ServeWithId Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// указанием ID сервера.
func (web *impl) ServeWithId(ltn net.Listener, id string) Interface {
	return web.ServeTLSWithId(ltn, nil, id)
}

// ServeTLS Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// использованием TLS.
func (web *impl) ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface {
	return web.ServeTLSWithId(ltn, tlsConfig, "")
}

// ServeTLSWithId Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
// использованием TLS и указанием ID сервера.
func (web *impl) ServeTLSWithId(ltn net.Listener, tlsConfig *tls.Config, id string) Interface {
	var (
		listener net.Listener
		isTls    bool
	)

	// Создание конфигурации по умолчанию на основе адреса слушателя.
	if web.cfg == nil {
		web.cfg, web.err = parseAddress(ltn.Addr().String())
		defaultConfiguration(web.cfg)
	}
	// Назначение идентификатора сервера, если он передан.
	if id != "" {
		web.cfg.ID = id
	}
	if web.server, isTls = web.makeServer(tlsConfig); web.Error() != nil {
		return web
	}
	// Оборачивание слушателя в TLS режим, если слушатель ещё не обёрнут, например, слушатель полученный из
	// функций ListenersSystemdTLSWithoutNames() или ListenersSystemdTLSWithNames().
	// Конфигурация TLS берётся из объекта сервера, там находится переданная или созданная по умолчанию
	// конфигурация.
	if listener = ltn; isTls && !isTLSListener(ltn) {
		listener = tls.NewListener(ltn, web.server.TLSConfig)
	}
	web.net.Handler(web.server.Serve)
	web.net.ServeWithId(listener, web.cfg.ID)

	return web
}

// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
// Если сервер не запущен, функция завершается немедленно.
func (web *impl) Wait() Interface { web.net.Wait(); return web }

// Stop Отправка сигнала прерывания работы веб сервера с учётом значения ShutdownTimeout.
func (web *impl) Stop() Interface { web.net.Stop(); return web }
