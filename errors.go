package web

// Все ошибки определены как константы.
const (
	cAlreadyRunning                = `Веб сервер уже запущен.`
	cNoConfiguration               = `Конфигурация ВЕБ сервера отсутствует либо равна nil.`
	cListenSystemdPID              = `Переменная окружения LISTEN_PID пустая либо содержит не верное значение.`
	cListenSystemdFDS              = `Переменная окружения LISTEN_FDS пустая либо содержит не верное значение.`
	cListenSystemdUnexpectedNumber = `Неожиданное количество сокетов активации FDS.`
	cListenSystemdNotFound         = `Получение сокета systemd по имени, имя не найдено.`
	cTLSIsNil                      = `Конфигурация TLS веб сервера пустая.`
	cHandlerIsNotSet               = `Не установлен обработчик ВЕБ запросов.`
)

// Константы указываются в объектах в качестве фиксированного адреса на протяжении всего времени работы приложения.
// Ошибка с ошибкой могут сравниваться по содержимому, по адресу и т.д.
var (
	errSingleton                     = &Error{}
	errAlreadyRunning                = err(cAlreadyRunning)
	errNoConfiguration               = err(cNoConfiguration)
	errListenSystemdPID              = err(cListenSystemdPID)
	errListenSystemdFDS              = err(cListenSystemdFDS)
	errListenSystemdUnexpectedNumber = err(cListenSystemdUnexpectedNumber)
	errListenSystemdNotFound         = err(cListenSystemdNotFound)
	errTLSIsNil                      = err(cTLSIsNil)
	errHandlerIsNotSet               = err(cHandlerIsNotSet)
)

type (
	// Error object of package
	Error struct{}
	err   string
)

// Error The error built-in interface implementation
func (e err) Error() string { return string(e) }

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }

// ОШИБКИ.

// ErrAlreadyRunning Веб сервер уже запущен.
func ErrAlreadyRunning() error { return &errAlreadyRunning }

// ErrNoConfiguration Конфигурация ВЕБ сервера отсутствует либо равна nil.
func ErrNoConfiguration() error { return &errNoConfiguration }

// ErrListenSystemdPID Переменная окружения LISTEN_PID пустая либо содержит не верное значение.
func ErrListenSystemdPID() error { return &errListenSystemdPID }

// ErrListenSystemdFDS Переменная окружения LISTEN_FDS пустая либо содержит не верное значение.
func ErrListenSystemdFDS() error { return &errListenSystemdFDS }

// ErrListenSystemdUnexpectedNumber Неожиданное количество сокетов активации FDS.
func ErrListenSystemdUnexpectedNumber() error { return &errListenSystemdUnexpectedNumber }

// ErrListenSystemdNotFound Получение сокета systemd по имени, имя не найдено.
func ErrListenSystemdNotFound() error { return &errListenSystemdNotFound }

// ErrTLSIsNil Конфигурация TLS веб сервера пустая.
func ErrTLSIsNil() error { return &errTLSIsNil }

// ErrHandlerIsNotSet Не установлен обработчик ВЕБ запросов.
func ErrHandlerIsNotSet() error { return &errHandlerIsNotSet }
