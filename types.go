package web

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync/atomic"
)

const (
	netTcp        = `tcp`
	netTcp4       = `tcp4`
	netTcp6       = `tcp6`
	netUnix       = `unix`
	netUnixPacket = `unixpacket`
	netSocket     = `socket`
	netSystemd    = `systemd`
)

// Interface Интерфейс пакета.
type Interface interface {
	// Handler Назначение обработчика запросов ВЕБ сервера.
	// Обработчик необходимо назначить до запуска ВЕБ сервера.
	Handler(handler http.Handler) Interface

	// ListenAndServe Открытие адреса или сокета без использования конфигурации веб сервера (конфигурация по
	// умолчанию), запуск веб сервера для обслуживания входящих соединений.
	ListenAndServe(addr string) Interface

	// ListenAndServeTLS Открытие адреса или сокета с использованием TLS, без использования конфигурации веб сервера
	// (конфигурация по умолчанию), запуск веб сервера для обслуживания входящих соединений.
	ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface

	// ListenAndServeWithConfig Настройка сервера с использованием переданной конфигурации, открытие адреса или сокета
	// на прослушивание, запуск веб сервера для обслуживания входящих соединений.
	ListenAndServeWithConfig(conf *Configuration) Interface

	// ListenAndServeTLSWithConfig Настройка сервера с использованием переданной конфигурации в режиме TLS, открытие
	// адреса или сокета на прослушивание, запуск веб сервера для обслуживания входящих соединений.
	ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface

	// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов переданных в процесс веб сервера из systemd.
	ListenersSystemdWithoutNames() (ret []net.Listener, err error)

	// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов переданных в процесс веб сервера
	// из systemd.
	ListenersSystemdWithNames() (ret map[string][]net.Listener, err error)

	// ListenersSystemdTLSWithoutNames Возвращает срез net.listener для TLS сокетов переданных в процесс веб сервера
	// из systemd.
	ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) (ret []net.Listener, err error)

	// ListenersSystemdTLSWithNames Возвращает карту срезов net.listener для TLS сокетов переданных в процесс веб сервера
	// из systemd.
	ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (ret map[string][]net.Listener, err error)

	// NewListener Создание нового слушателя соединений net.Listener на основе конфигурации веб сервера.
	NewListener(conf *Configuration) (ret net.Listener, err error)

	// NewListenerTLS Создание нового слушателя соединений net.Listener в режиме TLS, на основе конфигурации
	// веб сервера.
	NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, err error)

	// Serve Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener.
	Serve(net.Listener) Interface

	// ServeTLS Запуск веб сервера для входящих соединений на основе переданного слушателя net.Listener с
	// использованием TLS.
	ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface

	// Wait Блокируемая функция ожидания завершения веб сервера, если он запущен.
	// Если сервер не запущен, функция завершается немедленно.
	Wait() Interface

	// Stop Отправка сигнала прерывания работы веб сервера с учётом значения ShutdownTimeout.
	Stop() Interface

	// ОШИБКИ

	// Errors Справочник ошибок.
	Errors() *Error

	// Error Функция возвращает последнюю ошибку веб сервера.
	Error() error
}

// Объект сущности, реализующий интерфейс Interface.
type web struct {
	isRun       atomic.Value   // Состояние запуска веб сервера. =истина - запущен. =ложь - остановлен.
	inCloseUp   chan struct{}  // Канал начала завершения работы веб сервера.
	onCloseDone chan struct{}  // Канал сигнала окончания завершения работы веб сервера.
	conf        *Configuration // Конфигурация веб сервера.
	listener    net.Listener   // Слушатель сокета веб сервера.
	server      *http.Server   // Объект net/http веб сервера.
	handler     http.Handler   // Интерфейс обработчика запросов ВЕБ сервера.
	err         error          // Последняя ошибка веб сервера.
}
