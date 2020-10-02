package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// ListenSystemdSocket Загрузка сокета открытого systemd
func (wsv *web) ListenSystemdSocket() (ret net.Listener, err error) {
	var listeners []net.Listener

	if listeners, err = wsv.ListenersSystemdWithoutNames(); err != nil {
		return
	}
	if len(listeners) != 1 {
		err = ErrListenSystemdUnexpectedNumber()
		return
	}
	ret = listeners[0]

	return
}

// ListenLoadFilesFdWithNames Загрузка файловых дескрипторов на основе переменных окружения
func (wsv *web) ListenLoadFilesFdWithNames(unsetEnvAll bool) (ret []*os.File, err error) {
	const (
		listenFdBegin  = 3
		listenPID      = `LISTEN_PID`
		listenFds      = `LISTEN_FDS`
		listenFdNames  = `LISTEN_FDNAMES`
		listenFdPrefix = `LISTEN_FD_`
	)
	var (
		pID        int
		nFds       int
		names      []string
		name       string
		fd, offset int
	)

	if unsetEnvAll {
		defer func() { _ = os.Unsetenv(listenPID) }()
		defer func() { _ = os.Unsetenv(listenFds) }()
		defer func() { _ = os.Unsetenv(listenFdNames) }()
	}
	if pID, err = strconv.Atoi(os.Getenv(listenPID)); err != nil {
		err = fmt.Errorf("get pid from environment %q error: %s", listenPID, err)
		return
	}
	if nFds, err = strconv.Atoi(os.Getenv(listenFds)); err != nil {
		err = fmt.Errorf("get FDs from environment %q error: %s", listenFds, err)
		return
	}
	if pID != os.Getpid() || nFds == 0 {
		err = ErrListenSystemdPID()
		return
	}
	names = strings.Split(os.Getenv(listenFdNames), ":")
	ret = make([]*os.File, 0, nFds)
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

// ListenersSystemdWithoutNames Experimental function
func (wsv *web) ListenersSystemdWithoutNames() (ret []net.Listener, err error) {
	var (
		file  *os.File
		files []*os.File
		n     int
		pc    net.Listener
	)

	files, err = wsv.ListenLoadFilesFdWithNames(true)
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

	files, err = wsv.ListenLoadFilesFdWithNames(true)
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
