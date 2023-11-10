package web

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// ListenLoadFilesFdWithNames Загрузка файловых дескрипторов на основе переменных окружения
func (wbo *web) ListenLoadFilesFdWithNames() (ret []*os.File, err error) {
	const (
		listenFdBegin  = 3
		listenPID      = `LISTEN_PID`
		listenFds      = `LISTEN_FDS`
		listenFdNames  = `LISTEN_FDNAMES`
		listenFdPrefix = `LISTEN_FD_`
		errPIDTpl      = `получение PID из переменной окружения %q прервано ошибкой: %s`
		errFDSTpl      = `получение файлового дескриптора из переменной окружения %q прервано ошибкой: %s`
	)
	var (
		pID        int
		nFds       int
		names      []string
		name       string
		fd, offset int
	)

	if pID, err = strconv.Atoi(os.Getenv(listenPID)); err != nil {
		err = fmt.Errorf(errPIDTpl, listenPID, err)
		return
	}
	if nFds, err = strconv.Atoi(os.Getenv(listenFds)); err != nil {
		err = fmt.Errorf(errFDSTpl, listenFds, err)
		return
	}
	if pID != os.Getpid() {
		err = ErrListenSystemdPID()
		return
	}
	if nFds == 0 {
		err = ErrListenSystemdFDS()
		return
	}
	ret, names = make([]*os.File, 0, nFds), strings.Split(os.Getenv(listenFdNames), ":")
	for fd = listenFdBegin; fd < listenFdBegin+nFds; fd++ {
		syscall.CloseOnExec(fd)
		name = listenFdPrefix + strconv.Itoa(fd)
		if offset = fd - listenFdBegin; offset < len(names) && len(names[offset]) > 0 {
			name = names[offset]
		}
		ret = append(ret, os.NewFile(uintptr(fd), name))
	}

	return
}

// ListenersSystemdWithoutNames Возвращает срез net.Listener сокетов переданных в процесс веб сервера из systemd.
func (wbo *web) ListenersSystemdWithoutNames() (ret []net.Listener, err error) {
	var (
		file  *os.File
		files []*os.File
		n     int
		pc    net.Listener
	)

	files, err = wbo.ListenLoadFilesFdWithNames()
	ret = make([]net.Listener, len(files))
	for n, file = range files {
		if pc, err = net.FileListener(file); err == nil {
			ret[n], err = pc, file.Close()
		}
	}

	return
}

// ListenersSystemdWithNames Возвращает карту срезов net.Listener сокетов переданных в процесс веб сервера
// из systemd.
func (wbo *web) ListenersSystemdWithNames() (ret map[string][]net.Listener, err error) {
	var (
		file    *os.File
		files   []*os.File
		current []net.Listener
		pc      net.Listener
		ok      bool
	)

	files, err = wbo.ListenLoadFilesFdWithNames()
	ret = make(map[string][]net.Listener)
	for _, file = range files {
		if pc, err = net.FileListener(file); err == nil {
			if current, ok = ret[file.Name()]; !ok {
				ret[file.Name()] = []net.Listener{pc}
			} else {
				ret[file.Name()] = append(current, pc)
			}
			err = file.Close()
		}
	}

	return
}

// ListenersSystemdTLSWithoutNames Возвращает срез net.listener для TLS сокетов переданных в процесс веб сервера
// из systemd.
func (wbo *web) ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) (ret []net.Listener, err error) {
	var (
		listeners []net.Listener
		l         net.Listener
		n         int
	)

	if listeners, err = wbo.ListenersSystemdWithoutNames(); listeners == nil || err != nil {
		return
	}
	if tlsConfig == nil {
		err = ErrTLSIsNil()
		return
	}
	for n, l = range listeners {
		listeners[n] = tls.NewListener(l, tlsConfig)
	}

	return
}

// ListenersSystemdTLSWithNames Возвращает карту срезов net.listener для TLS сокетов переданных в процесс веб сервера
// из systemd.
func (wbo *web) ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (ret map[string][]net.Listener, err error) {
	var (
		listeners map[string][]net.Listener
		ll        []net.Listener
		l         net.Listener
		n         int
	)

	if listeners, err = wbo.ListenersSystemdWithNames(); listeners == nil || err != nil {
		return
	}
	if tlsConfig == nil {
		err = ErrTLSIsNil()
		return
	}
	for _, ll = range listeners {
		for n, l = range ll {
			ll[n] = tls.NewListener(l, tlsConfig)
		}
	}

	return
}
