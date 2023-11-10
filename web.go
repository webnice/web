package web

import "net/http"

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New() Interface {
	var wbo = &web{
		inCloseUp: make(chan struct{}, 1),
	}
	wbo.isRun.Store(false)

	return wbo
}

// Handler Назначение обработчика запросов ВЕБ сервера.
// Обработчик необходимо назначить до запуска ВЕБ сервера.
func (wbo *web) Handler(handler http.Handler) Interface {
	wbo.handler = handler

	return wbo
}

// Error Функция возвращает последнюю ошибку веб сервера.
func (wbo *web) Error() error { return wbo.err }

// Errors Справочник ошибок.
func (wbo *web) Errors() *Error { return errSingleton }
