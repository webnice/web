package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Наполнение конфигурации значениями по умолчанию
// Проверка значений
func defaultConfiguration(conf *Configuration) {
	if conf.Port == 0 {
		conf.Port = 80
	}
	// Check mode
	switch strings.ToLower(conf.Mode) {
	case netTcp, netTcp4, netTcp6, netUnix, netUnixPacket:
		conf.Mode = strings.ToLower(conf.Mode)
	case netSocket:
		conf.Mode = netUnix
	default:
		conf.Mode = netTcp
	}
	if conf.Mode == netUnix && conf.Socket == "" || conf.Mode == netUnixPacket && conf.Socket == "" {
		conf.Mode = netTcp
	}
	// Check MaxHeaderBytes
	if conf.MaxHeaderBytes == 0 {
		conf.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	}
	// unix socket modes
	switch conf.Mode {
	case netUnix, netUnixPacket:
		conf.HostPort = fmt.Sprintf("%s:%s", conf.Mode, conf.Socket)
	default:
		conf.HostPort = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	}
	// Public address
	if conf.Address == "" && conf.Mode == netTcp {
		if conf.Port == 80 {
			conf.Address = conf.Host
		} else {
			conf.Address = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		}
	}
	// ShutdownTimeout
	if conf.ShutdownTimeout <= 0 {
		conf.ShutdownTimeout = shutdownTimeout
	}
}

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных как ":http"
func parseAddress(addr string) (ret *Configuration, err error) {
	const addrSep = ":"
	var (
		sp []string
		n  int
		e  error
	)

	ret = new(Configuration)
	defer defaultConfiguration(ret)
	sp = strings.Split(addr, addrSep)
	if len(sp) <= 1 {
		ret.Host = sp[0]
		return
	}
	n, e = net.LookupPort(netTcp, strings.Join(sp[1:], addrSep))
	if err = e; err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint32(n)

	return
}
