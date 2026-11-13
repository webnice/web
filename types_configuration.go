package web

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"net/http"
	"time"

	wnet "github.com/webnice/net"
)

// Configuration Структура конфигурации веб сервера.
type Configuration struct {
	// Описание сетевого или локального доступа к серверу.
	wnet.Configuration `yaml:"Network" json:"network"`

	// Domain Список доменов, на которые отвечает сервер. Пустой список - ответ на любой домен.
	// Домены указываются в человекочитаемом виде, сравнение выполняется в punycode, поэтому домены
	// национальных алфавитов указываются как есть, например "домен.рф".
	// Домен может быть указан с портом (например "example.com:8443") - тогда сравнение выполняется с портом,
	// без порта - сравнение выполняется без порта.
	// Wildcard-шаблон "*.example.com" совпадает с любым количеством поддоменов, но не с самим доменом.
	// Фильтрация выполняется дважды: на этапе TLS handshake по SNI и по заголовку Host запроса.
	// Для доменов не входящих в список сервер отвечает статусом 434 Requested Host Unavailable
	// и закрывает соединение.
	// Default value: [] - все домена
	Domain []string `yaml:"Domain" json:"domain"`

	// ReadTimeout Время ожидания запроса включая ReadHeaderTimeout.
	// Если не указано или равно 0 - тайм-аута нет.
	// Default value: 0 - no timeout
	ReadTimeout time.Duration `yaml:"ReadTimeout" json:"read_timeout"`

	// ReadHeaderTimeout Время ожидания заголовка запроса.
	// Если не указано или равно 0 - тайм-аута нет.
	// Default value: 0 - no timeout
	ReadHeaderTimeout time.Duration `yaml:"ReadHeaderTimeout" json:"read_header_timeout"`

	// WriteTimeout Время ожидания выдачи ответа.
	// Если не указано или равно 0 - тайм-аута нет.
	// Default value: 0 - no timeout
	WriteTimeout time.Duration `yaml:"WriteTimeout" json:"write_timeout"`

	// IdleTimeout Максимальное время ожидания следующего входящего соединения для открытого сокета, до его закрытия.
	// Используется при включённом keep-alives.
	// Если не указано или равно 0 - тайм-аута нет
	// Default value: 0 - no timeout
	IdleTimeout time.Duration `yaml:"IdleTimeout" json:"idle_timeout"`

	// ShutdownTimeout Максимальное время ожидания завершения работы веб сервера до начала принудительного обрыва
	// соединений и остановки процессов.
	// Если не указано или равно 0 - тайм-аута нет.
	// Default value: 30s
	ShutdownTimeout time.Duration `yaml:"ShutdownTimeout" json:"shutdown_timeout" default-value:"30s"`

	// MaxHeaderBytes Максимальный размер заголовка запроса.
	// Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
	MaxHeaderBytes int `yaml:"MaxHeaderBytes" json:"max_header_bytes" default-value:"1048576"`

	// MaxBodyBytes Максимальный размер тела запроса в байтах.
	// Если не указано или равно 0 - размер тела запроса не ограничен.
	// При превышении размера сервер отвечает ошибкой 413 Request Entity Too Large.
	// Default value: 0 - unlimited
	MaxBodyBytes uint64 `yaml:"MaxBodyBytes" json:"max_body_bytes"`

	// KeepAliveDisable Отключение режима "оставаться в живых" (keep alive).
	// Default value: false - keep alive are enabled.
	KeepAliveDisable bool `yaml:"KeepAliveDisable" json:"keep_alive_disable"`

	// DisableGeneralOptionsHandler Отключение обработки HTTP запросов методом OPTIONS.
	// Если установлено значение "истина", на все запросы методом OPTIONS сервер отвечает кодом 200 и передаёт в
	// заголовке Content-Length значение 0.
	// Если установлено значение "ложь", запросы передаются контроллеру и обрабатываются обычным образом.
	// Default value: false - запросы передаются контроллеру.
	DisableGeneralOptionsHandler bool `yaml:"DisableGeneralOptionsHandler" json:"disable_general_options_handler"`

	// HTTP2Disable Отключение протокола HTTP/2 для TLS соединений.
	// Если установлено значение "ложь", при TLS подключении сервер согласует протокол HTTP/2 по умолчанию,
	// как это делает "net/http".
	// Default value: false - HTTP/2 включён.
	HTTP2Disable bool `yaml:"HTTP2Disable" json:"http2_disable"`

	// HTTP2UnencryptedEnable Включение протокола HTTP/2 без шифрования (h2c, prior knowledge).
	// Протокол применяется только для не зашифрованных соединений.
	// Default value: false - h2c отключён.
	HTTP2UnencryptedEnable bool `yaml:"HTTP2UnencryptedEnable" json:"http2_unencrypted_enable"`

	// HTTP2MaxConcurrentStreams Максимальное количество одновременных потоков в одном HTTP/2 соединении.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxConcurrentStreams int `yaml:"HTTP2MaxConcurrentStreams" json:"http2_max_concurrent_streams"`

	// HTTP2MaxDecoderHeaderTableSize Максимальный размер таблицы декодирования заголовков HTTP/2.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxDecoderHeaderTableSize int `yaml:"HTTP2MaxDecoderHeaderTableSize" json:"http2_max_decoder_header_table_size"`

	// HTTP2MaxEncoderHeaderTableSize Максимальный размер таблицы кодирования заголовков HTTP/2.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxEncoderHeaderTableSize int `yaml:"HTTP2MaxEncoderHeaderTableSize" json:"http2_max_encoder_header_table_size"`

	// HTTP2MaxReadFrameSize Максимальный размер читаемого кадра HTTP/2, от 16 KiB до 16 MiB.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxReadFrameSize int `yaml:"HTTP2MaxReadFrameSize" json:"http2_max_read_frame_size"`

	// HTTP2MaxReceiveBufferPerConnection Максимальный размер буфера приёма данных HTTP/2 соединения, не менее 64 KiB.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxReceiveBufferPerConnection int `yaml:"HTTP2MaxReceiveBufferPerConnection" json:"http2_max_receive_buffer_per_connection"`

	// HTTP2MaxReceiveBufferPerStream Максимальный размер буфера приёма данных потока HTTP/2.
	// Если не указано или равно 0 - используется значение по умолчанию.
	// Default value: 0 - значение по умолчанию
	HTTP2MaxReceiveBufferPerStream int `yaml:"HTTP2MaxReceiveBufferPerStream" json:"http2_max_receive_buffer_per_stream"`

	// HTTP2SendPingTimeout Интервал проверки соединения HTTP/2 кадром PING при отсутствии принятых кадров.
	// Если не указано или равно 0 - проверка не выполняется.
	// Default value: 0 - проверка не выполняется
	HTTP2SendPingTimeout time.Duration `yaml:"HTTP2SendPingTimeout" json:"http2_send_ping_timeout"`

	// HTTP2PingTimeout Время ожидания ответа на кадр PING, после которого соединение HTTP/2 закрывается.
	// Если не указано или равно 0 - используется значение по умолчанию 15 секунд.
	// Default value: 0 - значение по умолчанию
	HTTP2PingTimeout time.Duration `yaml:"HTTP2PingTimeout" json:"http2_ping_timeout"`

	// HTTP2WriteByteTimeout Время ожидания записи данных HTTP/2, после которого соединение закрывается.
	// Защита от медленных клиентов.
	// Если не указано или равно 0 - тайм-аута нет.
	// Default value: 0 - no timeout
	HTTP2WriteByteTimeout time.Duration `yaml:"HTTP2WriteByteTimeout" json:"http2_write_byte_timeout"`

	// MaxHeaderValueCount Максимальное количество значений одного заголовка запроса.
	// Защита от запросов с чрезмерным количеством значений заголовков.
	// Если не указано или равно 0 - используется значение по умолчанию 500.
	// Default value: 0 - значение по умолчанию
	MaxHeaderValueCount int `yaml:"MaxHeaderValueCount" json:"max_header_value_count"`

	// DisableClientPriority Отключение приоритизации запросов HTTP/2 клиента.
	// Default value: false - приоритизация включена.
	DisableClientPriority bool `yaml:"DisableClientPriority" json:"disable_client_priority"`

	// HTTP1FullDuplexEnable Включение полнодуплексного режима HTTP/1: чтение тела запроса одновременно
	// с записью ответа. Применяется только к HTTP/1 соединениям.
	// Default value: false - полнодуплексный режим отключён.
	HTTP1FullDuplexEnable bool `yaml:"HTTP1FullDuplexEnable" json:"http1_full_duplex_enable"`

	// KeepAliveEnable Включение TCP keep-alive проб для принятых соединений.
	// Если не указано и остальные поля keep-alive не указаны - используется поведение операционной системы.
	// Default value: false
	KeepAliveEnable bool `yaml:"KeepAliveEnable" json:"keep_alive_enable"`

	// KeepAliveIdle Время бездействия соединения до первой keep-alive пробы.
	// Default value: 0 - значение по умолчанию операционной системы
	KeepAliveIdle time.Duration `yaml:"KeepAliveIdle" json:"keep_alive_idle"`

	// KeepAliveInterval Интервал между keep-alive пробами.
	// Default value: 0 - значение по умолчанию операционной системы
	KeepAliveInterval time.Duration `yaml:"KeepAliveInterval" json:"keep_alive_interval"`

	// KeepAliveCount Максимальное количество не отвеченных keep-alive проб до закрытия соединения.
	// Default value: 0 - значение по умолчанию операционной системы
	KeepAliveCount int `yaml:"KeepAliveCount" json:"keep_alive_count"`

	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////// //
	// ПЕРЕМЕННЫЕ КОНФИГУРАЦИИ КОТОРЫЕ НАСТРАИВАЮТ ИЗ ПРОГРАММНОГО КОДА И НЕ ЧИТАЮТСЯ ИЗ КОНФИГУРАЦИИ YAML ИЛИ JSON   //
	// ////////////////////////////////////////////////////////////////////////////////////////////////////////////// //

	// ErrorLog Журнал внутренних ошибок сервера. Если не указан, используется стандартный журнал log.
	// Default value: nil
	ErrorLog *log.Logger `yaml:"-" json:"-"`

	// ConnState Функция, вызываемая при изменении состояния соединения. Используется для метрик
	// и тонкой настройки принятых соединений.
	// Default value: nil
	ConnState func(net.Conn, http.ConnState) `yaml:"-" json:"-"`

	// BaseContext Функция создания базового контекста сервера.
	// Default value: nil
	BaseContext func(net.Listener) context.Context `yaml:"-" json:"-"`

	// ConnContext Функция изменения контекста принятого соединения.
	// Default value: nil
	ConnContext func(context.Context, net.Conn) context.Context `yaml:"-" json:"-"`

	// TLSConfig Базовая конфигурация TLS соединения: SNI, версии протокола, шифры, кривые, mTLS,
	// session tickets, проверка сертификатов клиента и другие настройки tls.Config.
	// Если не указана, конфигурация TLS создаётся по ключам TLSPublicKeyPEM и TLSPrivateKeyPEM.
	// Default value: nil
	TLSConfig *tls.Config `yaml:"-" json:"-"`

	// TLSClientAuth Требование и проверка сертификата клиента (mTLS).
	// Применяется если конфигурация TLSConfig не указана или если её поле ClientAuth не изменено.
	// Default value: 0 - клиентский сертификат не запрашивается.
	TLSClientAuth tls.ClientAuthType `yaml:"-" json:"-"`

	// TLSClientCAs Пул сертификатов удостоверяющих центров, которыми подписываются клиентские
	// сертификаты при проверке mTLS.
	// Default value: nil
	TLSClientCAs *x509.CertPool `yaml:"-" json:"-"`
}

/**

   Пример конфигурации YAML:

      ## Описание сетевого или локального доступа к серверу.
      Network:
            ## ID Уникальный идентификатор сервера, может быть любым уникальным строковым значением.
            ## Если при запуске сервера значение не указано, создаётся уникальное временное значение, меняющееся
            ## при каждом запуске.
            ## Default value: ""
            ID string `yaml:"ID" json:"id"`

            ## Публичный адрес по которому сервер доступен извне.
            ## Например, если сервер находится за прокси, тут указывается реальный адрес подключения к серверу.
            ## Default value: ""
            Address: !!str "http://localhost/"

            ## IP адрес или имя хоста на котором поднимается сервер, можно указывать 0.0.0.0 для всех ip адресов.
            ## Default value: "0.0.0.0".
            #Host: !!str "[2a03:e2c0:a32::2]"
            #Host: !!str "example.hostname.local"
            Host: !!str "0.0.0.0"

            ## Tcp/ip порт занимаемый сервером.
            ## Default value: 80
            Port: !!int 80

            ## Юникс сокет на котором поднимается сервер, только для unix-like операционных систем Linux, Unix, Mac.
            ## Default value: ""
            Socket: !!str "run/example.sock"

            ## Файловые разрешения доступа к юникс-сокету.
            ## Значение задаётся в восьмеричной системе и не должно превышать 32 бита.
            ## Default value: "0666"
            SocketMode: !!str "0666"

            ## Режим открытия сокета, возможные значения: tcp, tcp4, tcp6, unix, unixpacket, socket, systemd.
            ## udp, udp4, udp6 - Сервер поднимается на указанном Host:Port;
            ## tcp, tcp4, tcp6 - Сервер поднимается на указанном Host:Port;
            ## unix, unixpacket - Сервер поднимается на указанном unix/unixpacket;
            ## socket  - Сервер поднимается на socket, только для unix-like операционных систем. Параметры Host:Port
            ##           игнорируются, используется только путь к сокету;
            ## systemd - Порт или сокет открывает systemd и передаёт слушателя порта через файловый дескриптор сервису,
            ##           запущенному от пользователя без права открытия привилегированных портов. Максимально удобный
            ##           способ при использовании правильного безопасно настроенного Linux сервера.
            ##           Более подробно можно посмотреть в документации man systemd.socket(5);
            ## Default value: "tcp"
            Mode: !!str "tcp"

            ## Путь и имя файла содержащего публичный ключ (сертификат) в PEM формате, включая CA
            ## сертификаты всех промежуточных центров сертификации, если ими подписан ключ.
            ## Применяется только для TCP соединений, для UDP не используется.
            ## Default value: ""
            TLSPublicKeyPEM: !!str "/etc/application/certificate.pub"

            ## Путь и имя файла содержащего секретный/приватный ключ в PEM формате.
            ## Применяется только для TCP соединений, для UDP не используется.
            ## Default value: ""
            TLSPrivateKeyPEM: !!str "/etc/application/certificate.key"

            ## Включение прокси-протокола.
            ## Прокси-протокол позволяет серверу получать информацию о подключении клиента, передаваемую через
            ## прокси-серверы и средства балансировки нагрузки, такие как Nginx, HAProxy, Amazon Elastic Load
            ## Balancer (ELB) и многие другие.
            ## Поддерживаются запросы с реализацией прокси протокола версий 1 и 2.
            ## PROXY protocol: https://www.haproxy.org/download/2.3/doc/proxy-protocol.txt.
            ## Default value: false
            ProxyProtocol: !!bool false

            ## Максимальное время ожидания получения данных о клиенте через прокси-протокол.
            ## Время ожидание используется только при включённом прокси-протоколе.
            ## Default value: 0s - no timeout
            ProxyProtocolReadHeaderTimeout: 0s

      ## Список доменов, на которые отвечает сервер. Домены указываются в человекочитаемом виде,
      ## сравнение выполняется в punycode. Домен с портом сравнивается с портом, без порта - без порта.
      ## Wildcard-шаблон "*.example.com" совпадает с любым количеством поддоменов.
      ## Фильтрация выполняется дважды: на этапе TLS handshake по SNI и по заголовку Host запроса.
      ## Для не разрешённых доменов сервер отвечает статусом 434 Requested Host Unavailable.
      ## Default value: [] - все домена.
      Domain:
      - !!str "localhost"
      - !!str "localhost:8080"
      - !!str "example.domain.tld"
      - !!str "*.example.domain.tld"
      - !!str "домен.рф"

      ## Время ожидания запроса включая ReadHeaderTimeout.
      ## Если не указано или равно 0 - тайм-аута нет.
      ## Default value: 0 - no timeout
      ReadTimeout: 0s

      ## Время ожидания заголовка запроса.
      ## Если не указано или равно 0 - тайм-аута нет.
      ## Default value: 0 - no timeout
      ReadHeaderTimeout: 0s

      ## Время ожидания выдачи ответа.
      ## Если не указано или равно 0 - тайм-аута нет.
      ## Default value: 0 - no timeout
      WriteTimeout: 0s

      ## Максимальное время ожидания следующего входящего соединения для открытого сокета, до его закрытия.
      ## Используется при включённом keep-alives.
      ## Если не указано или равно 0 - тайм-аута нет
      ## Default value: 0 - no timeout
      IdleTimeout: 0s

      ## Максимальное время ожидания завершения работы веб сервера до начала принудительного обрыва
      ## соединений и остановки процессов.
      ## Если не указано или равно 0 - тайм-аута нет.
      ## Default value: 30s
      ShutdownTimeout: 30s

      ## Максимальный размер заголовка запроса.
      ## Default value: 1 MB (from net/http/DefaultMaxHeaderBytes)
      MaxHeaderBytes: !!int 1048576

      ## Максимальный размер тела запроса в байтах.
      ## Если не указано или равно 0 - размер тела запроса не ограничен.
      ## При превышении размера сервер отвечает ошибкой 413 Request Entity Too Large.
      ## Default value: 0 - без ограничений.
      MaxBodyBytes: !!int 0

      ## Отключение режима "оставаться в живых" (keep alive).
      ## Default value: false - оставаться в живых включён.
      KeepAliveDisable: !!bool false

      ## Отключение обработки HTTP запросов методом OPTIONS.
      ## Если установлено значение "истина", на все запросы методом OPTIONS сервер отвечает кодом 200 и передаёт в
      ## заголовке Content-Length значение 0.
      ## Если установлено значение "ложь", запросы передаются контроллеру и обрабатываются обычным образом.
      ## Default value: false - запросы передаются контроллеру.
      DisableGeneralOptionsHandler: !!bool false

      ## Отключение протокола HTTP/2 для TLS соединений.
      ## Default value: false - HTTP/2 включён.
      HTTP2Disable: !!bool false

      ## Включение протокола HTTP/2 без шифрования (h2c, prior knowledge).
      ## Протокол применяется только для не зашифрованных соединений.
      ## Default value: false - h2c отключён.
      HTTP2UnencryptedEnable: !!bool false

      ## Максимальное количество одновременных потоков в одном HTTP/2 соединении.
      ## Если не указано или равно 0 - используется значение по умолчанию.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxConcurrentStreams: !!int 0

      ## Максимальный размер таблицы декодирования заголовков HTTP/2.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxDecoderHeaderTableSize: !!int 0

      ## Максимальный размер таблицы кодирования заголовков HTTP/2.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxEncoderHeaderTableSize: !!int 0

      ## Максимальный размер читаемого кадра HTTP/2, от 16 KiB до 16 MiB.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxReadFrameSize: !!int 0

      ## Максимальный размер буфера приёма данных HTTP/2 соединения, не менее 64 KiB.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxReceiveBufferPerConnection: !!int 0

      ## Максимальный размер буфера приёма данных потока HTTP/2.
      ## Default value: 0 - значение по умолчанию
      HTTP2MaxReceiveBufferPerStream: !!int 0

      ## Интервал проверки соединения HTTP/2 кадром PING при отсутствии принятых кадров.
      ## Default value: 0s - проверка не выполняется
      HTTP2SendPingTimeout: 0s

      ## Время ожидания ответа на кадр PING, после которого соединение HTTP/2 закрывается.
      ## Default value: 0s - значение по умолчанию
      HTTP2PingTimeout: 0s

      ## Время ожидания записи данных HTTP/2, после которого соединение закрывается.
      ## Default value: 0s - таймаута нет
      HTTP2WriteByteTimeout: 0s

      ## Максимальное количество значений одного заголовка запроса.
      ## Default value: 0 - значение по умолчанию 500
      MaxHeaderValueCount: !!int 0

      ## Отключение приоритизации запросов HTTP/2 клиента.
      ## Default value: false - приоритизация включена.
      DisableClientPriority: !!bool false

      ## Включение полнодуплексного режима HTTP/1.
      ## Default value: false - полнодуплексный режим отключён.
      HTTP1FullDuplexEnable: !!bool false

      ## Включение TCP keep-alive проб для принятых соединений.
      ## Default value: false
      KeepAliveEnable: !!bool false

      ## Время бездействия соединения до первой keep-alive пробы.
      ## Default value: 0s - значение по умолчанию операционной системы
      KeepAliveIdle: 0s

      ## Интервал между keep-alive пробами.
      ## Default value: 0s - значение по умолчанию операционной системы
      KeepAliveInterval: 0s

      ## Максимальное количество не отвеченных keep-alive проб до закрытия соединения.
      ## Default value: 0 - значение по умолчанию операционной системы
      KeepAliveCount: !!int 0

**/
