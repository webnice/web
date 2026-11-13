/*
Пакет web — микро сборка для создания управляемого веб сервера на основе "net/http".

Библиотека основана на пакетах "net", "net/http" и "github.com/webnice/net": открытие адресов и
сокетов, TLS, прокси-протокол, сокеты службы systemd и корректное завершение работы выполняются
библиотекой "github.com/webnice/net"; обработка HTTP запросов выполняется стандартным сервером
"net/http"; в качестве мультиплексора запросов может использоваться любой обработчик, рекомендуется
"github.com/go-chi/chi".

# Быстрый старт

	// Создание мультиплексора запросов.
	mux := chi.NewMux()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("привет"))
	})

	// Запуск веб сервера с ожиданием завершения.
	srv := web.New().Handler(mux)
	if err := srv.ListenAndServe(":8080").Error(); err != nil {
		log.Fatalf("запуск веб сервера прерван ошибкой: %s", err)
	}

Адрес сервера указывается в любом формате, понятном пакету "net": ":8080", "127.0.0.1:8080",
"localhost:http", "[::1]:443", а так же путь к UNIX сокету файловой системы.

# Запуск сервера

Конструктор New() возвращает объект удовлетворяющий интерфейсу Interface. Запуск сервера выполняется
на адресе или сокете функциями ListenAndServe() и ListenAndServeTLS(), либо по готовой конфигурации
функциями ListenAndServeWithConfig() и ListenAndServeTLSWithConfig(). Для запуска сервера на готовом
объекте net.Listener используются функции Serve(), ServeWithId(), ServeTLS() и ServeTLSWithId():

	// Запуск сервера по конфигурации.
	srv := web.New().Handler(mux)
	srv.ListenAndServeWithConfig(&web.Configuration{
		Configuration: wnet.Configuration{Host: "0.0.0.0", Port: 8080, Mode: "tcp"},
		ReadTimeout:   time.Second * 5,
	})

	// Запуск сервера на готовом слушателе.
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("открытие адреса прервано ошибкой: %s", err)
	}
	if err = web.New().Handler(mux).Serve(listener).Error(); err != nil {
		log.Fatalf("запуск веб сервера прерван ошибкой: %s", err)
	}

Сервер поднимается на TCP адресе (tcp, tcp4, tcp6), на UNIX сокете (unix, unixpacket), на сокете
файловой системы и на сокетах, открытых службой systemd; сокеты systemd загружаются функциями
ListenersSystemdWithoutNames(), ListenersSystemdWithNames(), ListenersSystemdTLSWithoutNames() и
ListenersSystemdTLSWithNames(). Слушатели без запуска сервера создаются функциями NewListener() и
NewListenerTLS().

# Конфигурация

Структура Configuration содержит встроенную конфигурацию библиотеки "github.com/webnice/net"
(адрес, режим открытия сокета, TLS ключи, прокси-протокол) и настройки веб сервера: таймауты чтения
и записи, лимиты размера заголовка и тела запроса, ограничение списка доменов, настройки протоколов
HTTP/1 и HTTP/2, TLS, TCP keep-alive и обработки запросов методом OPTIONS. Поля конфигурации
документированы в файле types_configuration.go, пример YAML конфигурации приведён в том же файле.

# Ограничение списка доменов

Сервер отвечает только на домена из списка Domain конфигурации. Фильтрация выполняется дважды: на
этапе TLS handshake по SNI и по заголовку Host запроса, для не разрешённых доменов сервер отвечает
статусом 434 Requested Host Unavailable и закрывает соединение. Домены указываются в
человекочитаемом виде, сравнение выполняется в punycode, поддерживаются wildcard-шаблоны:

	srv.ListenAndServeWithConfig(&web.Configuration{
		Configuration: wnet.Configuration{Host: "0.0.0.0", Port: 8443, TLSPublicKeyPEM: crt, TLSPrivateKeyPEM: key},
		Domain: []string{"example.com", "*.example.org", "домен.рф", "example.com:8443"},
	})

Домен с портом сравнивается с портом, домен без порта - без порта. Wildcard-шаблон "*.example.com"
совпадает с любым количеством поддоменов, но не с самим доменом. Пустой список - ответ на любой домен.

# Настройка TLS

Базовая конфигурация TLS создаётся по ключам TLSPublicKeyPEM и TLSPrivateKeyPEM. Произвольная
конфигурация TLS передаётся полем TLSConfig: SNI и виртуальный хостинг (GetCertificate,
GetConfigForClient), версии протокола и шифры, session tickets, проверка сертификатов клиента и
другие настройки tls.Config. Требование клиентского сертификата (mTLS) настраивается полями
TLSClientAuth и TLSClientCAs. Приоритет источников конфигурации TLS: аргумент функции запуска,
затем поле TLSConfig, затем ключи PEM:

	srv.ListenAndServeTLSWithConfig(&web.Configuration{
		Configuration: wnet.Configuration{Host: "0.0.0.0", Port: 8443},
		TLSConfig:     &tls.Config{GetCertificate: certBySNI},
	}, nil)

TCP keep-alive для принятых соединений настраивается полями KeepAliveEnable, KeepAliveIdle,
KeepAliveInterval и KeepAliveCount, отключение keep-alive выполняется полем KeepAliveDisable.

# Протокол HTTP/2

Протокол HTTP/2 включён по умолчанию для TLS соединений, как в "net/http"; при согласовании ALPN
сервер отдаёт предпочтение протоколу HTTP/2. Незашифрованный HTTP/2 (h2c, prior knowledge)
включается полем HTTP2UnencryptedEnable конфигурации, отключение HTTP/2 выполняется полем
HTTP2Disable, максимальное количество одновременных потоков соединения ограничивается полем
HTTP2MaxConcurrentStreams:

	srv := web.New().Handler(mux)
	srv.ListenAndServeTLSWithConfig(&web.Configuration{
		Configuration: wnet.Configuration{
			Host:             "0.0.0.0",
			Port:             8443,
			TLSPublicKeyPEM:  "/etc/application/certificate.pub",
			TLSPrivateKeyPEM: "/etc/application/certificate.key",
		},
		HTTP2MaxConcurrentStreams: 100,
		HTTP2WriteByteTimeout:     time.Second * 10,
	}, nil)

Полный список настроек протокола HTTP/2: HTTP2MaxConcurrentStreams, HTTP2MaxDecoderHeaderTableSize,
HTTP2MaxEncoderHeaderTableSize, HTTP2MaxReadFrameSize, HTTP2MaxReceiveBufferPerConnection,
HTTP2MaxReceiveBufferPerStream, HTTP2SendPingTimeout, HTTP2PingTimeout и HTTP2WriteByteTimeout.

# Ограничение размера тела запроса

Максимальный размер тела запроса ограничивается полем MaxBodyBytes конфигурации. При превышении
лимита сервер отвечает ошибкой 413 Request Entity Too Large. Значение 0 отключает ограничение:

	srv.ListenAndServeWithConfig(&web.Configuration{
		Configuration: wnet.Configuration{Host: "0.0.0.0", Port: 8080},
		MaxBodyBytes:  1 << 20, // 1 MB
	})

# Завершение работы сервера

Функция Stop() отправляет сигнал прерывания работы сервера, завершение выполняется с учётом
значения ShutdownTimeout конфигурации. Функция Wait() блокирует вызвавшую её горутину до полного
завершения работы сервера; если сервер не запущен, функция завершается немедленно:

	srv := web.New().Handler(mux).ListenAndServe(":8080")
	// ...в другой горутине, например по сигналу операционной системы:
	srv.Stop()
	srv.Wait()

# Ошибки

Функция Error() возвращает последнюю ошибку пакета или библиотеки "github.com/webnice/net". Ошибка
библиотеки "github.com/webnice/net" приводится к справочнику ошибок пакета, поэтому сравнение ошибок
функцией errors.Is() выполняется в обе стороны. Справочник ошибок реализован на основе библиотеки
"github.com/webnice/dic": ошибка создаётся методом Bind(), сравнивается по якорю или по коду:

	srv.ListenAndServe(":8080")
	srv.ListenAndServe(":8081")
	if errors.Is(srv.Error(), srv.Errors().AlreadyRunning.Bind()) {
		log.Printf("сервер уже запущен")
	}
	srv.Clean() // очистка последней ошибки.

# Справочники

Пакет предоставляет справочники библиотеки "github.com/webnice/dic": MIME типов, HTTP методов
запросов, заголовков и статусов HTTP ответов:

	web.Mime.TextHtml.String()      // "text/html;charset=UTF-8"
	web.Method.Get.String()         // "GET"
	web.Header.ContentType.String() // "Content-Type"
	web.Status.NotFound.String()    // "Not Found"
	web.File.MimeByExtension(".png").String()       // "image/png"
	web.File.MimeByFilename("example.txt").String() // "text/plain"

Справочники дополняются пользовательскими данными через встраивание, например собственный справочник
заголовков:

	// myHeaders Структура справочника заголовков с пользовательскими заголовками.
	type myHeaders struct {
		dic.Headers
		XRequestID dic.IHeader // Заголовок X-Request-Id.
	}

	var myHeader = &myHeaders{
		Headers:    dic.Header(),
		XRequestID: dic.NewHeader("X-Request-Id"),
	}

Полное описание методов интерфейса пакета: https://pkg.go.dev/github.com/webnice/web/v3.
*/
package web
