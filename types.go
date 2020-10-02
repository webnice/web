package web

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/webnice/web.v1/context/errors"
	"gopkg.in/webnice/web.v1/context/handlers"
	"gopkg.in/webnice/web.v1/route"
)

const (
	netTcp        = `tcp`
	netTcp4       = `tcp4`
	netTcp6       = `tcp6`
	netUnix       = `unix`
	netUnixPacket = `unixpacket`
	netSocket     = `socket`
	netSystemd    = `systemd`
)

// Interface is an interface
type Interface interface {
	// ListenAndServe listens on the TCP network address addr and then calls Serve on incoming connections
	ListenAndServe(addr string) Interface

	// ListenAndServeTLS listens on the TCP network address address with TLS and then calls Serve with handler
	// to handle requests on incoming connections
	ListenAndServeTLS(addr string, certFile string, keyFile string, tlsConfig *tls.Config) Interface

	// ListenAndServeWithConfig Fully configurable web server listens and then calls Serve on incoming connections
	ListenAndServeWithConfig(*Configuration) Interface

	// ListenAndServeTLSWithConfig Fully configurable web server listens and then calls Serve on incoming connections
	ListenAndServeTLSWithConfig(conf *Configuration, tlsConfig *tls.Config) Interface

	// NewListener Make new listener from web server configuration
	NewListener(conf *Configuration) (ret net.Listener, err error)

	// NewListenerTLS Make new listener with TLS from web server configuration
	NewListenerTLS(conf *Configuration, tlsConfig *tls.Config) (ret net.Listener, err error)

	// Serve accepts incoming connections on the Listener, creating a new service goroutine for each
	Serve(net.Listener) Interface

	// ServeTLS accepts incoming connections on the listener with TLS configuration, creating a new web server goroutine
	ServeTLS(ltn net.Listener, tlsConfig *tls.Config) Interface

	// Error Return last error of web server
	Error() error

	// Wait while web server is running
	Wait() Interface

	// Stop web server
	Stop() Interface

	// Route interface
	Route() route.Interface

	// Errors interface
	Errors() errors.Interface

	// Handlers interface
	Handlers() handlers.Interface
}

// Is an private implementation of web server
type web struct {
	isRun       atomic.Value    // The indicator of web server goroutine. =true-goroutine is started, =false-goroutine is stopped
	inCloseUp   chan bool       // The indicator of web server state, true in channel means we're in shutdown goroutine and web server
	doCloseDone sync.WaitGroup  // Wait while goroutine stopped
	conf        *Configuration  // The web server configuration
	listener    net.Listener    // The web server listener
	server      *http.Server    // The net/http web server object
	route       route.Interface // Routing settings interface
	err         error           // The last of error
}

// Configuration is a structure of web server configuration
type Configuration struct { // nolint: maligned
	// HostPort (readonly) Адрес составленный автоматически из Host:Port
	// Значение создаётся автоматически при инициализации конфигурации
	// Default value: ":http"
	HostPort string `yaml:"-" json:"-"`

	// Address Публичный адрес на котором сервер доступен извне
	// Например если сервер находится за прокси, тут указывается реальный адрес подключения к серверу
	// Default value: "" - make automatically
	Address string `yaml:"Address" json:"address"`

	// TODO Сделать ограничение по доменам
	// Domains Список всех доменов, на которые отвечает сервер
	// Если не пусто, то для всех других доменов будет ответ "Requested host unavailable"
	// Default value: [] - all domain
	//Domains []string `yaml:"Domains" json:"domains"`

	// TLSPublicKeyPEM Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA сертификаты
	// всех промежуточных центров сертификации, если ими подписан ключ
	TLSPublicKeyPEM string `yaml:"TLSPublicKeyPEM" json:"tls_public_key_pem"`

	// TLSPrivateKeyPEM Путь и имя файла содержащего приватный ключ в PEM формате
	TLSPrivateKeyPEM string `yaml:"TLSPrivateKeyPEM" json:"tls_private_key_pem"`

	// Host IP адрес или имя хоста на котором запускается web сервер,
	// можно указывать 0.0.0.0 для всех ip адресов
	// Default value: ""
	Host string `yaml:"Host" json:"host"`

	// Port tcp/ip порт занимаемый сервером
	// Default value: 80
	Port uint16 `yaml:"Port" json:"port"`

	// Socket Unix socket на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac
	// Default value: "" - unix socket is off
	Socket string `yaml:"Socket" json:"socket"`

	// Mode Режим работы, tcp, tcp4, tcp6, unix, unixpacket, socket, systemd
	// systemd - systemd is configures to open the port or unix socket and passes over the file descriptor the port
	// listener to the service. See systemd.socket(5) manual.
	// Default value: "tcp"
	Mode string `yaml:"Mode" json:"mode"`

	// ReadTimeout Время в наносекундах ожидания запроса включая ReadHeaderTimeout
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	ReadTimeout time.Duration `yaml:"ReadTimeout" json:"read_timeout"`

	// ReadHeaderTimeout Время в наносекундах ожидания заголовка запроса
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	ReadHeaderTimeout time.Duration `yaml:"ReadHeaderTimeout" json:"read_header_timeout"`

	// WriteTimeout Время в наносекундах ожидания выдачи ответа
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	WriteTimeout time.Duration `yaml:"WriteTimeout" json:"write_timeout"`

	// IdleTimeout is the maximum amount of time to wait for the next request when keep-alives are enabled
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	IdleTimeout time.Duration `yaml:"IdleTimeout" json:"idle_timeout"`

	// ShutdownTimeout is the maximum amount of time to wait for the server graceful shutdown
	// Если не указано или рано 0 - таймаута нет
	// Default value: 30s
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout" json:"shutdown_timeout"`

	// MaxHeaderBytes controls the maximum number of bytes the server will read parsing the request header's keys and values, including the request line
	// Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
	MaxHeaderBytes int `yaml:"MaxHeaderBytes" json:"max_header_bytes"`

	// TODO Сделать ограничение на максимальный размер тела запроса
	// MaxBodyBytes controls the maximum number of bytes the server will read request body
	// Default value: 0 - unlimited
	//MaxBodyBytes uint64 `yaml:"MaxBodyBytes" json:"max_body_bytes"`

	// KeepAliveDisable if is equal true, keep alive are disabled, if false - keep alive are enabled
	// Default value: false - keep alive are enabled
	KeepAliveDisable bool `yaml:"KeepAliveDisable" json:"keep_alive_disable"`

	// ProxyProtocol if is equal true, accepting the PROXY Protocol
	// The PROXY protocol enables web server to receive client connection information passed through proxy servers and
	// load balancers such as HAproxy and Amazon Elastic Load Balancer (ELB).
	// With the PROXY protocol, web server can learn the originating IP address from HTTP, SSL, HTTP/2, SPDY, WebSocket, and TCP.
	ProxyProtocol bool `yaml:"ProxyProtocol" json:"proxy_protocol"`
}
