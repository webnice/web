# web

Библиотека golang

Микро сборка для создания управляемого веб сервера на основе net/http.

[![Go Reference](https://pkg.go.dev/badge/github.com/webnice/web/v3.svg)](https://pkg.go.dev/github.com/webnice/web/v3)
[![Coverage Status](https://coveralls.io/repos/github/webnice/web/badge.svg?branch=master%0Av3)](https://coveralls.io/github/webnice/web?branch=master%0Av3)

## Описание

`web` — это обёртка над пакетами `net`, `net/http` и `github.com/webnice/net` для создания управляемого
веб сервера. Открытие адресов и сокетов, TLS, прокси-протокол, работа с сокетами службы systemd и корректное
завершение работы выполняются библиотекой `github.com/webnice/net`; обработка HTTP запросов выполняется
стандартным сервером `net/http`; в качестве мультиплексора запросов может использоваться любой обработчик,
рекомендуется `github.com/go-chi/chi`.

### Особенности

- Запуск сервера на `tcp`, `tcp4`, `tcp6`, `unix`, `unixpacket` сокетах, на UNIX сокете файловой системы
  и на сокетах, открытых службой systemd;
- TLS из коробки, протокол HTTP/2 включён по умолчанию для TLS соединений, поддержка h2c (незашифрованный
  HTTP/2, prior knowledge), настройка лимита потоков;
- Ограничение максимального размера тела запроса с ответом 413 Request Entity Too Large;
- Поддержка прокси-протокола версии 1 и 2, реализованного в HAProxy, nginx и других балансировщиках;
- Корректное завершение работы сервера с учётом значения ShutdownTimeout;
- Справочники MIME типов, HTTP методов, заголовков и статусов ответов;
- Справочник ошибок на основе библиотеки `github.com/webnice/dic`, сравнение ошибок функцией `errors.Is()`;
- Покрытие кода тестами 100%.

## Требования

- Go 1.26 или новее (используются `http.Protocols`, `http.HTTP2Config`, `errors.AsType`).

## Установка

```
go get -u github.com/webnice/web/v3
```

Подключение:

```
import "github.com/webnice/web/v3"
```

#### Зависимости

```
github.com/webnice/dic
github.com/webnice/net
github.com/go-chi/chi/v5      (используется в тестах)
github.com/google/uuid         (используется в тестах)
github.com/pires/go-proxyproto (косвенная зависимость)
```

## Примеры использования

### Простой веб сервер

```
mux := chi.NewMux()
mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("привет"))
})
if err := web.New().Handler(mux).ListenAndServe(":8080").Error(); err != nil {
	log.Fatalf("запуск веб сервера прерван ошибкой: %s", err)
}
```

Адрес можно указывать в любом формате, понятном пакету `net`: `:8080`, `127.0.0.1:8080`,
`localhost:http`, `[::1]:443`, а также путь к UNIX сокету.

### Веб сервер с конфигурацией

```
srv := web.New().Handler(mux)
srv.ListenAndServeWithConfig(&web.Configuration{
	Configuration: wnet.Configuration{
		Host: "0.0.0.0",
		Port: 8080,
		Mode: "tcp",
	},
	ReadTimeout:     time.Second * 5,
	WriteTimeout:    time.Second * 10,
	ShutdownTimeout: time.Second * 30,
	MaxBodyBytes:    1 << 20, // 1 MB, ответ 413 при превышении
})
```

Полный список полей конфигурации с описанием находится в документации пакета
(https://pkg.go.dev/github.com/webnice/web/v3#Configuration), пример YAML конфигурации приведён в исходном
коде файла types_configuration.go.

### TLS и HTTP/2

```
if err := web.New().Handler(mux).
	ListenAndServeTLS(":8443", "/etc/application/certificate.pub", "/etc/application/certificate.key", nil).
	Error(); err != nil {
	log.Fatalf("запуск веб сервера прерван ошибкой: %s", err)
}
```

Для TLS соединений протокол HTTP/2 включён по умолчанию, как в `net/http`. Отключение HTTP/2, поддержка
h2c и ограничение количества потоков настраиваются полями `HTTP2Disable`, `HTTP2UnencryptedEnable` и
`HTTP2MaxConcurrentStreams` конфигурации.

### Graceful shutdown

```
srv := web.New().Handler(mux).ListenAndServe(":8080")
// ...в другой горутине, например по сигналу операционной системы:
srv.Stop() // корректное завершение с учётом ShutdownTimeout
srv.Wait() // ожидание завершения работы сервера
```

### Работа с готовым слушателем

```
listener, err := net.Listen("tcp", ":8080")
if err != nil {
	log.Fatalf("открытие адреса прервано ошибкой: %s", err)
}
if err = web.New().Handler(mux).Serve(listener).Error(); err != nil {
	log.Fatalf("запуск веб сервера прерван ошибкой: %s", err)
}
```

Готовые TLS слушатели можно получить из функций `NewListenerTLS()`, `ListenersSystemdTLSWithoutNames()`,
`ListenersSystemdTLSWithNames()` — повторное TLS оборачивание выполняться не будет.

### Сравнение ошибок

```
srv := web.New().Handler(mux).ListenAndServe(":8080")
srv.ListenAndServe(":8081")
if errors.Is(srv.Error(), srv.Errors().AlreadyRunning.Bind()) {
	log.Printf("сервер уже запущен")
}
srv.Clean() // очистка последней ошибки
```

Ошибки библиотеки `github.com/webnice/net` приводятся к справочнику ошибок пакета, поэтому сравнение
функцией `errors.Is()` работает в обе стороны.

### Справочники

```
web.Mime.TextHtml.String()     // справочник MIME типов: "text/html;charset=UTF-8"
web.Method.Get.String()        // справочник HTTP методов запросов: "GET"
web.Header.ContentType.String() // справочник заголовков: "Content-Type"
web.Status.NotFound.String()   // справочник статусов HTTP ответов: "Not Found"
```

## Интерфейс пакета

Конструктор `New()` возвращает объект удовлетворяющий интерфейсу `web.Interface`:

- `ID` — уникальный идентификатор сервера;
- `Handler` — назначение обработчика запросов, обязательно до запуска сервера;
- `ListenAndServe`, `ListenAndServeWithConfig` — запуск сервера на адресе или сокете;
- `ListenAndServeTLS`, `ListenAndServeTLSWithConfig` — запуск сервера в режиме TLS;
- `ListenersSystemdWithoutNames`, `ListenersSystemdWithNames` — получение сокетов, переданных службой systemd;
- `ListenersSystemdTLSWithoutNames`, `ListenersSystemdTLSWithNames` — то же самое в режиме TLS;
- `NewListener`, `NewListenerTLS` — создание слушателя по конфигурации без запуска сервера;
- `Serve`, `ServeWithId`, `ServeTLS`, `ServeTLSWithId` — запуск сервера на готовом слушателе;
- `Stop`, `Wait` — прерывание и ожидание завершения работы сервера;
- `Clean` — очистка последней ошибки;
- `Errors` — справочник ошибок;
- `Error` — последняя ошибка веб сервера или библиотеки `github.com/webnice/net`.

Полное описание методов интерфейса доступно в документации пакета: https://pkg.go.dev/github.com/webnice/web/v3.

## Тестирование

```
# Стандартный прогон с покрытием.
go test -cover ./...

# Прогон через Makefile (агрегированный профиль покрытия в coverage.log).
make test
```

## Лицензия

© WEB NICE, распространяется под [лицензией MIT](LICENSE).
