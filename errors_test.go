package web

import (
	"errors"
	"strings"
	"testing"

	wnet "github.com/webnice/net"
)

// Проверка эквивалентности ошибок справочника ошибкам библиотеки "github.com/webnice/net".
//
//goland:noinspection DuplicatedCode
func TestErrorsFromNet(t *testing.T) {
	const (
		tplErrorNotEqual                   = "ошибка %s не эквивалентна ошибке библиотеки net"
		testErrorNameAlreadyRunning        = "AlreadyRunning"
		testErrorNameNoConfiguration       = "NoConfiguration"
		testErrorNameListenSystemdPID      = "ListenSystemdPID"
		testErrorNameListenSystemdFDS      = "ListenSystemdFDS"
		testErrorNameSystemdNotFound       = "ListenSystemdNotFound"
		testErrorNameQuantityNotMatch      = "ListenSystemdQuantityNotMatch"
		testErrorNameTLSIsNil              = "TLSIsNil"
		testErrorNameServerHandlerIsNotSet = "ServerHandlerIsNotSet"
	)
	var err error

	if err = Errors().AlreadyRunning.Bind(); !errors.Is(err, wnet.Errors().AlreadyRunning()) {
		t.Errorf(tplErrorNotEqual, testErrorNameAlreadyRunning)
	}
	if err = Errors().NoConfiguration.Bind(); !errors.Is(err, wnet.Errors().NoConfiguration()) {
		t.Errorf(tplErrorNotEqual, testErrorNameNoConfiguration)
	}
	if err = Errors().ListenSystemdPID.Bind(); !errors.Is(err, wnet.Errors().ListenSystemdPID()) {
		t.Errorf(tplErrorNotEqual, testErrorNameListenSystemdPID)
	}
	if err = Errors().ListenSystemdFDS.Bind(); !errors.Is(err, wnet.Errors().ListenSystemdFDS()) {
		t.Errorf(tplErrorNotEqual, testErrorNameListenSystemdFDS)
	}
	if err = Errors().ListenSystemdNotFound.Bind(); !errors.Is(err, wnet.Errors().ListenSystemdNotFound()) {
		t.Errorf(tplErrorNotEqual, testErrorNameSystemdNotFound)
	}
	if err = Errors().ListenSystemdQuantityNotMatch.Bind(); !errors.Is(err, wnet.Errors().ListenSystemdQuantityNotMatch()) {
		t.Errorf(tplErrorNotEqual, testErrorNameQuantityNotMatch)
	}
	if err = Errors().TLSIsNil.Bind(); !errors.Is(err, wnet.Errors().TLSIsNil()) {
		t.Errorf(tplErrorNotEqual, testErrorNameTLSIsNil)
	}
	if err = Errors().ServerHandlerIsNotSet.Bind(); !errors.Is(err, wnet.Errors().ServerHandlerIsNotSet()) {
		t.Errorf(tplErrorNotEqual, testErrorNameServerHandlerIsNotSet)
	}
}

// Проверка внутренних ошибок пакета.
func TestErrorsInternal(t *testing.T) {
	const (
		tplErrorText   = "ошибка HandlerIsNotSet содержит не корректный текст: %q"
		tplErrorUnbind = "ошибка HandlerIsNotSet не найдена справочником через функцию Unbind()"
		tplErrorParse  = "ошибка HandlerIsNotSet не найдена справочником через функцию ParseError()"
		tplErrorCode   = "ошибка HandlerIsNotSet имеет не корректный код: %d"
		tplErrorCopies = "копии ошибки HandlerIsNotSet не эквивалентны"
	)
	var err error

	err = Errors().HandlerIsNotSet.Bind()
	if !strings.Contains(err.Error(), cHandlerIsNotSet) {
		t.Errorf(tplErrorText, err.Error())
	}
	// Сравнение ошибок через справочник.
	if !Errors().Unbind(err).IsEqual(Errors().HandlerIsNotSet) {
		t.Error(tplErrorUnbind)
	}
	if !Errors().ParseError(err, errSingleton).IsEqual(Errors().HandlerIsNotSet) {
		t.Error(tplErrorParse)
	}
	// Код ошибки.
	if Errors().HandlerIsNotSet.CodeU8().Get() != eHandlerIsNotSet {
		t.Errorf(tplErrorCode, Errors().HandlerIsNotSet.CodeU8().Get())
	}
	// Эквивалентность копий ошибки.
	if !errors.Is(Errors().HandlerIsNotSet.Bind(), Errors().HandlerIsNotSet.Bind()) {
		t.Error(tplErrorCopies)
	}
}
