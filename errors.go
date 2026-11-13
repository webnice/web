package web

import "github.com/webnice/dic"

// Коды ошибок.
const (
	eAlreadyRunning uint8 = iota + 1
	eNoConfiguration
	eListenSystemdPID
	eListenSystemdFDS
	eListenSystemdNotFound
	eListenSystemdQuantityNotMatch
	eTLSIsNil
	eServerHandlerIsNotSet
	eHandlerIsNotSet
)

// Текстовые значения кодов ошибок на основном языке приложения.
// Ошибки дублируют ошибки библиотеки "github.com/webnice/net", текстовые значения должны совпадать с ними.
const (
	cAlreadyRunning                = "Сервер уже запущен."
	cNoConfiguration               = "Конфигурация сервера отсутствует либо равна nil."
	cListenSystemdPID              = "Переменная окружения LISTEN_PID пустая, либо содержит неверное значение."
	cListenSystemdFDS              = "Переменная окружения LISTEN_FDS пустая, либо содержит неверное значение."
	cListenSystemdNotFound         = "Получение сокета systemd по имени, имя не найдено."
	cListenSystemdQuantityNotMatch = "Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES."
	cTLSIsNil                      = "Конфигурация TLS сервера пустая."
	cServerHandlerIsNotSet         = "Не установлен обработчик основной функции TCP сервера."
	cHandlerIsNotSet               = "Не установлен обработчик запросов ВЕБ сервера."
)

// Error Структура справочника ошибок.
type Error struct {
	dic.Errors

	// AlreadyRunning Сервер уже запущен.
	AlreadyRunning dic.IError

	// NoConfiguration Конфигурация сервера отсутствует либо равна nil.
	NoConfiguration dic.IError

	// ListenSystemdPID Переменная окружения LISTEN_PID пустая, либо содержит неверное значение.
	ListenSystemdPID dic.IError

	// ListenSystemdFDS Переменная окружения LISTEN_FDS пустая, либо содержит неверное значение.
	ListenSystemdFDS dic.IError

	// ListenSystemdNotFound Получение сокета systemd по имени, имя не найдено.
	ListenSystemdNotFound dic.IError

	// ListenSystemdQuantityNotMatch Полученное количество LISTEN_FDS не соответствует переданному LISTEN_FDNAMES.
	ListenSystemdQuantityNotMatch dic.IError

	// TLSIsNil Конфигурация TLS сервера пустая.
	TLSIsNil dic.IError

	// ServerHandlerIsNotSet Не установлен обработчик основной функции TCP сервера.
	ServerHandlerIsNotSet dic.IError

	// HandlerIsNotSet Не установлен обработчик запросов ВЕБ сервера.
	HandlerIsNotSet dic.IError
}

// Ошибки справочника создаются один раз, с фиксированным адресом якоря, ошибки могут сравниваться по
// содержимому, по указателю и т.д.
var errSingleton = &Error{
	Errors:                        dic.Error(),
	AlreadyRunning:                dic.NewError(cAlreadyRunning).CodeU8().Set(eAlreadyRunning),
	NoConfiguration:               dic.NewError(cNoConfiguration).CodeU8().Set(eNoConfiguration),
	ListenSystemdPID:              dic.NewError(cListenSystemdPID).CodeU8().Set(eListenSystemdPID),
	ListenSystemdFDS:              dic.NewError(cListenSystemdFDS).CodeU8().Set(eListenSystemdFDS),
	ListenSystemdNotFound:         dic.NewError(cListenSystemdNotFound).CodeU8().Set(eListenSystemdNotFound),
	ListenSystemdQuantityNotMatch: dic.NewError(cListenSystemdQuantityNotMatch).CodeU8().Set(eListenSystemdQuantityNotMatch),
	TLSIsNil:                      dic.NewError(cTLSIsNil).CodeU8().Set(eTLSIsNil),
	ServerHandlerIsNotSet:         dic.NewError(cServerHandlerIsNotSet).CodeU8().Set(eServerHandlerIsNotSet),
	HandlerIsNotSet:               dic.NewError(cHandlerIsNotSet).CodeU8().Set(eHandlerIsNotSet),
}

// Errors Справочник ошибок.
func Errors() *Error { return errSingleton }
