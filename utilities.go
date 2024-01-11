package web

import (
	"net"
	"strings"
)

// Разбор адреса, определение порта через net.LookupPort, в том числе портов заданных как ":http"
func parseAddress(addr string) (ret *Configuration, err error) {
	const bColon = ":"
	var (
		sp []string
		n  int
	)

	ret = new(Configuration)
	if sp = strings.Split(addr, bColon); len(sp) <= 1 {
		ret.Host = sp[0]
		return
	}
	if n, err = net.LookupPort(netTcp, strings.Join(sp[1:], bColon)); err != nil {
		return
	}
	ret.Host, ret.Port = sp[0], uint16(n)

	return
}
