package web

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/idna"
)

const (
	tplDomainAbort = "запрашиваемый домен %q не разрешён конфигурацией сервера"
	domainWildcard = "*."
)

// Профиль преобразования доменов в punycode: строгий режим с проверкой меток и символов.
var domainIDNA = idna.New(idna.StrictDomainName(true), idna.ValidateLabels(true), idna.MapForLookup())

// Фильтр разрешённых доменов: нормализованные списки доменов конфигурации.
type domainFilter struct {
	hosts     []string // Домены без порта.
	hostPorts []string // Домены с портом.
}

// Создание фильтра доменов: домены конфигурации преобразуются в punycode и разделяются на домены с портом
// и без порта. Не корректные записи конфигурации пропускаются.
func newDomainFilter(domains []string) (ret *domainFilter) {
	var (
		n   int
		hst string
		prt string
		err error
	)

	ret = &domainFilter{
		hosts:     make([]string, 0, len(domains)),
		hostPorts: make([]string, 0, len(domains)),
	}
	for n = range domains {
		// Разделение домена и порта, порт сохраняется в записи, если он указан.
		if hst, prt, err = net.SplitHostPort(domains[n]); err != nil {
			hst, prt = domains[n], ""
		}
		// Преобразование домена из человекочитаемого вида в технический вид punycode.
		// Wildcard-шаблон преобразуется по частям: звёздочка не является корректным символом домена.
		if strings.HasPrefix(hst, domainWildcard) {
			if hst, err = domainIDNA.ToASCII(hst[2:]); err != nil {
				// Не корректный домен конфигурации пропускается.
				continue
			}
			hst = domainWildcard + hst
		} else if hst, err = domainIDNA.ToASCII(hst); err != nil {
			// Не корректный домен конфигурации пропускается.
			continue
		}
		if prt != "" {
			ret.hostPorts = append(ret.hostPorts, net.JoinHostPort(hst, prt))
			continue
		}
		ret.hosts = append(ret.hosts, hst)
	}

	return
}

// Проверка того, что фильтр не содержит ни одного домена и ограничение не применяется.
func (df *domainFilter) empty() bool { return len(df.hosts) == 0 && len(df.hostPorts) == 0 }

// Проверка домена без учёта порта.
func (df *domainFilter) matchByHost(host string) bool { return matchDomainList(host, df.hosts) }

// Проверка домена с учётом порта.
func (df *domainFilter) matchByHostPort(host, port string) bool {
	return matchDomainList(net.JoinHostPort(host, port), df.hostPorts)
}

// Проверка домена по списку шаблонов: точное сравнение или wildcard-шаблон.
func matchDomainList(host string, patterns []string) bool {
	var n int

	for n = range patterns {
		if matchDomainPattern(patterns[n], host) {
			return true
		}
	}

	return false
}

// Сравнение домена с шаблоном: точное сравнение или wildcard-шаблон "*.домен", который совпадает с любым
// количеством поддоменов, но не с самим доменом.
func matchDomainPattern(pattern, host string) bool {
	var suffix string

	if !strings.HasPrefix(pattern, domainWildcard) {
		return pattern == host
	}
	// Сравнение по суффиксу ".домен".
	suffix = pattern[1:]
	if !strings.HasSuffix(host, suffix) {
		return false
	}
	// Домен без поддоменов не совпадает с шаблоном: host "example.com", suffix ".example.com".
	return len(host) > len(suffix)
}

// Фильтрация домена на этапе TLS handshake по SNI. Вернётся ошибка, если домен не разрешён конфигурацией.
func (df *domainFilter) sniFilter(hello *tls.ClientHelloInfo) (*tls.Config, error) {
	var (
		err  error
		host string
	)

	if hello.ServerName == "" {
		// Домен не указан, проверка будет выполнена по заголовку Host после расшифровки соединения.
		return nil, nil
	}
	// Преобразование домена из человекочитаемого вида в технический вид punycode.
	if host, err = domainIDNA.ToASCII(hello.ServerName); err != nil {
		// Не корректный домен пропускается, ограничение применяется максимально дружественно.
		return nil, nil
	}
	if df.matchByHost(host) {
		return nil, nil
	}

	return nil, fmt.Errorf(tplDomainAbort, hello.ServerName)
}

// Обёртка обработчика запросов с фильтрацией домена по заголовку Host. Для не разрешённых доменов сервер
// отвечает статусом 434 Requested Host Unavailable и закрывает соединение.
func (df *domainFilter) handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err  error
			host string
			prt  string
		)

		// Разделение домена и порта заголовка Host.
		if host, prt, err = net.SplitHostPort(r.Host); err != nil {
			host, prt = r.Host, ""
		}
		// Преобразование домена из человекочитаемого вида в технический вид punycode.
		if host, err = domainIDNA.ToASCII(host); err != nil {
			// Не корректный домен пропускается, ограничение применяется максимально дружественно.
			h.ServeHTTP(w, r)
			return
		}
		if df.matchByHost(host) || df.matchByHostPort(host, prt) {
			h.ServeHTTP(w, r)
			return
		}
		// Домен не разрешён конфигурацией.
		w.Header().Set(Header.Connection.String(), "close")
		w.WriteHeader(Status.RequestedHostUnavailable.Code())
		_, _ = w.Write([]byte(Status.RequestedHostUnavailable.String()))
	})
}
