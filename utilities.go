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
	case "tcp", "tcp4", "tcp6", "unix", "unixpacket":
		conf.Mode = strings.ToLower(conf.Mode)
	case "socket":
		conf.Mode = strings.ToLower("unix")
	default:
		conf.Mode = "tcp"
	}
	if conf.Mode == "unix" && conf.Socket == "" || conf.Mode == "unixpacket" && conf.Socket == "" {
		conf.Mode = "tcp"
	}
	// Check MaxHeaderBytes
	if conf.MaxHeaderBytes == 0 {
		conf.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	}
	// unix socket modes
	if conf.Mode == "unix" || conf.Mode == "unixpacket" {
		conf.HostPort = fmt.Sprintf("%s:%s", conf.Mode, conf.Socket)
	} else {
		conf.HostPort = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	}
	// Public address
	if conf.Address == "" && conf.Mode == "tcp" {
		if conf.Port == 80 {
			conf.Address = fmt.Sprintf("%s", conf.Host)
		} else {
			conf.Address = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		}
	}
	// ShutdownTimeout
	if conf.ShutdownTimeout == 0 {
		conf.ShutdownTimeout = _ShutdownTimeout
	}
}

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных как ":http"
func parseAddress(addr string) (ret *Configuration, err error) {
	var sp []string
	ret = new(Configuration)
	defer defaultConfiguration(ret)
	sp = strings.Split(addr, ":")
	if len(sp) <= 1 {
		ret.Host = sp[0]
		return
	}
	var n, e = net.LookupPort("tcp", strings.Join(sp[1:], ":"))
	if err = e; err != nil {
		return
	}
	ret.Host = sp[0]
	ret.Port = uint32(n)
	return
}
