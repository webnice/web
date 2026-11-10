package web

import (
	"errors"
	"strings"
	"testing"

	wnet "github.com/webnice/net"
)

func TestErrorsFromNet(t *testing.T) {
	var err error

	if err = Errors().AlreadyRunning(); !errors.Is(err, wnet.Errors().AlreadyRunning()) {
		t.Errorf("ошибка AlreadyRunning() неправильно определена")
	}
	if err = Errors().NoConfiguration(); !errors.Is(err, wnet.Errors().NoConfiguration()) {
		t.Errorf("ошибка NoConfiguration() неправильно определена")
	}
	if err = Errors().ListenSystemdPID(); !errors.Is(err, wnet.Errors().ListenSystemdPID()) {
		t.Errorf("ошибка ListenSystemdPID() неправильно определена")
	}
	if err = Errors().ListenSystemdFDS(); !errors.Is(err, wnet.Errors().ListenSystemdFDS()) {
		t.Errorf("ошибка ListenSystemdFDS() неправильно определена")
	}
	if err = Errors().ListenSystemdNotFound(); !errors.Is(err, wnet.Errors().ListenSystemdNotFound()) {
		t.Errorf("ошибка ListenSystemdNotFound() неправильно определена")
	}
	if err = Errors().ListenSystemdQuantityNotMatch(); !errors.Is(err, wnet.Errors().ListenSystemdQuantityNotMatch()) {
		t.Errorf("ошибка ListenSystemdQuantityNotMatch() неправильно определена")
	}
	if err = Errors().TLSIsNil(); !errors.Is(err, wnet.Errors().TLSIsNil()) {
		t.Errorf("ошибка TLSIsNil() неправильно определена")
	}
	if err = Errors().ServerHandlerIsNotSet(); !errors.Is(err, wnet.Errors().ServerHandlerIsNotSet()) {
		t.Errorf("ошибка ServerHandlerIsNotSet() неправильно определена")
	}
	// Внутренние ошибки:
	err = Errors().HandlerIsNotSet()
	if !strings.Contains(err.Error(), "Не установлен обработчик запросов ВЕБ сервера.") {
		t.Errorf("ошибка HandlerIsNotSet() не определена")
	}
}
