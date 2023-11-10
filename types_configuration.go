// Package web
package web

import "time"

// Configuration Структура конфигурации веб сервера.
type Configuration struct {
	// HostPort (readonly) Адрес составленный автоматически из Host:Port.
	// Значение создаётся автоматически при инициализации конфигурации.
	// Default value: ""
	HostPort string `yaml:"-" json:"-"`

	// Address Публичный адрес на котором сервер доступен извне.
	// Например, если сервер находится за прокси, тут указывается реальный адрес подключения к серверу.
	// Default value: "" - make automatically
	Address string `yaml:"Address" json:"address"`

	// TODO Сделать ограничение по доменам.
	// Domains Список всех доменов, на которые отвечает сервер.
	// Если не пусто, то для всех других доменов будет ответ "Requested host unavailable".
	// Default value: [] - all domain
	//Domains []string `yaml:"Domains" json:"domains"`

	// TLSPublicKeyPEM Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA сертификаты
	// всех промежуточных центров сертификации, если ими подписан ключ.
	TLSPublicKeyPEM string `yaml:"TLSPublicKeyPEM" json:"tls_public_key_pem"`

	// TLSPrivateKeyPEM Путь и имя файла содержащего приватный ключ в PEM формате.
	TLSPrivateKeyPEM string `yaml:"TLSPrivateKeyPEM" json:"tls_private_key_pem"`

	// Host IP адрес или имя хоста на котором запускается web сервер, можно указывать 0.0.0.0 для всех ip адресов.
	// Default value: "0.0.0.0"
	Host string `yaml:"Host" json:"host" default-value:"0.0.0.0"`

	// Port tcp/ip порт занимаемый сервером.
	// Default value: 80
	Port uint16 `yaml:"Port" json:"port" default-value:"80"`

	// Socket Unix socket на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac.
	// Default value: "" - unix socket is off
	Socket string `yaml:"Socket" json:"socket" default-value:"-"`

	// Mode Режим работы, tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
	// systemd - Открытие порта выполняется службой линукса systemd, веб серверу передаётся открытое соединение, через
	// файловый дескриптор, в который поступают входящие соединения.
	// Более подробно можно посмотреть в документации man systemd.socket(5).
	// Default value: "tcp"
	Mode string `yaml:"Mode" json:"mode" default-value:"tcp"`

	// ReadTimeout Время в наносекундах ожидания запроса включая ReadHeaderTimeout.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	ReadTimeout time.Duration `yaml:"ReadTimeout" json:"read_timeout"`

	// ReadHeaderTimeout Время в наносекундах ожидания заголовка запроса.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	ReadHeaderTimeout time.Duration `yaml:"ReadHeaderTimeout" json:"read_header_timeout"`

	// WriteTimeout Время в наносекундах ожидания выдачи ответа.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 0 - no timeout
	WriteTimeout time.Duration `yaml:"WriteTimeout" json:"write_timeout"`

	// IdleTimeout Максимальное время ожидания следующего входящего соединения для открытого сокета, дл его закрытия.
	// Используется при включённом keep-alives.
	// Если не указано или рано 0 - таймаута нет
	// Default value: 0 - no timeout
	IdleTimeout time.Duration `yaml:"IdleTimeout" json:"idle_timeout"`

	// ShutdownTimeout Максимальное время ожидания завершения работы веб сервера до начала принудительного обрыва
	// соединений и остановки процессов.
	// Если не указано или рано 0 - таймаута нет.
	// Default value: 30s
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout" json:"shutdown_timeout" default-value:"30s"`

	// MaxHeaderBytes Максимальный размер заголовка запроса.
	// Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
	MaxHeaderBytes int `yaml:"MaxHeaderBytes" json:"max_header_bytes" default-value:"1048576"`

	// TODO Сделать ограничение на максимальный размер тела запроса.
	// MaxBodyBytes Максимальный размер тела запроса.
	// Default value: 0 - unlimited
	//MaxBodyBytes uint64 `yaml:"MaxBodyBytes" json:"max_body_bytes"`

	// KeepAliveDisable Отключение режима "оставаться в живых" (keep alive).
	// Default value: false - keep alive are enabled
	KeepAliveDisable bool `yaml:"KeepAliveDisable" json:"keep_alive_disable"`

	// ProxyProtocol Включение прокси протокола.
	// Прокси протокол позволяет веб-серверу получать информацию о подключении клиента, передаваемую через
	// прокси-серверы и средства балансировки нагрузки, такие как HAProxy, Amazon Elastic Load Balancer (ELB) и другие.
	// С помощью прокси протокола веб-сервер может узнать IP-адрес клиента для HTTP, SSL, HTTP/2, SPDY, WebSocket, TCP
	// запросов приходящих от прокси сервера.
	// Default value: false
	ProxyProtocol bool `yaml:"ProxyProtocol" json:"proxy_protocol"`
}
