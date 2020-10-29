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
func (wsv *web) ListenLoadFilesFdWithNames() (ret []*os.File, err error) {
	const (
		listenFdBegin  = 3
		listenPID      = `LISTEN_PID`
		listenFds      = `LISTEN_FDS`
		listenFdNames  = `LISTEN_FDNAMES`
		listenFdPrefix = `LISTEN_FD_`
		errPIDTpl      = `getting pid from environment %q, error: %s`
		errFDSTpl      = `getting FD id from environment %q, error: %s`
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

// ListenersSystemdWithoutNames returns a net.Listener for each matching socket type passed to this process
func (wsv *web) ListenersSystemdWithoutNames() (ret []net.Listener, err error) {
	var (
		file  *os.File
		files []*os.File
		n     int
		pc    net.Listener
	)

	files, err = wsv.ListenLoadFilesFdWithNames()
	ret = make([]net.Listener, len(files))
	for n, file = range files {
		if pc, err = net.FileListener(file); err == nil {
			ret[n], err = pc, file.Close()
		}
	}

	return
}

// ListenersSystemdWithNames maps a listener name to a set of net.Listener instances
func (wsv *web) ListenersSystemdWithNames() (ret map[string][]net.Listener, err error) {
	var (
		file    *os.File
		files   []*os.File
		current []net.Listener
		pc      net.Listener
		ok      bool
	)

	files, err = wsv.ListenLoadFilesFdWithNames()
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

// ListenersSystemdTLSWithoutNames returns a net.listener for each matching TCP socket type passed to this process
func (wsv *web) ListenersSystemdTLSWithoutNames(tlsConfig *tls.Config) (ret []net.Listener, err error) {
	var (
		listeners []net.Listener
		l         net.Listener
		n         int
	)

	if listeners, err = wsv.ListenersSystemdWithoutNames(); listeners == nil || err != nil {
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

// ListenersSystemdTLSWithNames maps a listener name to a net.Listener with the associated TLS configuration
func (wsv *web) ListenersSystemdTLSWithNames(tlsConfig *tls.Config) (ret map[string][]net.Listener, err error) {
	var (
		listeners map[string][]net.Listener
		ll        []net.Listener
		l         net.Listener
		n         int
	)

	if listeners, err = wsv.ListenersSystemdWithNames(); listeners == nil || err != nil {
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
