package web

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
)

// Обёртка тела запроса, ограничивающая его размер и отвечающая ошибкой 413 Request Entity Too Large
// при превышении лимита.
type maxBodyBytesReader struct {
	w http.ResponseWriter
	r io.ReadCloser
}

// Read Чтение тела запроса.
func (mbr *maxBodyBytesReader) Read(p []byte) (n int, err error) {
	if n, err = mbr.r.Read(p); err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			// Лимит размера тела запроса превышен. Ответ 413 возможен если обработчик ещё не начал ответ.
			http.Error(mbr.w, http.StatusText(http.StatusRequestEntityTooLarge), http.StatusRequestEntityTooLarge)
		}
	}

	return
}

// Close Закрытие тела запроса.
func (mbr *maxBodyBytesReader) Close() error { return mbr.r.Close() }

// Ограничение максимального размера тела запроса в n байт с ответом ошибкой 413 Request Entity Too Large
// при превышении лимита.
func maxBodyBytesHandler(h http.Handler, n int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var r2 *http.Request

		// Ограничение размера тела запроса без изменения исходного запроса.
		r2 = new(http.Request)
		*r2 = *r
		r2.Body = &maxBodyBytesReader{w: w, r: http.MaxBytesReader(w, r.Body, n)}
		h.ServeHTTP(w, r2)
	})
}

// Включение полнодуплексного режима HTTP/1 для соединения: чтение тела запроса одновременно с записью ответа.
func http1FullDuplexHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ошибка ожидается для HTTP/2 соединений: полнодуплексный режим доступен только в HTTP/1.
		_ = http.NewResponseController(w).EnableFullDuplex()
		h.ServeHTTP(w, r)
	})
}

// Применение настроек из конфигурации к стандартному веб серверу.
func (web *impl) makeServer(tlsConfig *tls.Config) (ret *http.Server, isTLS bool) {
	const (
		http11 = "http/1.1"
		proto2 = "h2"
	)
	var (
		df      *domainFilter
		handler http.Handler
	)

	if web.handler == nil {
		web.err = Errors().HandlerIsNotSet.Bind()
		return
	}
	if web.cfg == nil {
		web.err = Errors().NoConfiguration.Bind()
		return
	}
	// Фильтр разрешённых доменов, применяется к TLS handshake и к обработке запросов.
	if df = newDomainFilter(web.cfg.Domain); df.empty() {
		df = nil
	}
	// Цепочка обработки запросов: ограничение доменов, ограничение размера тела, полнодуплексный режим.
	handler = web.handler
	if df != nil {
		handler = df.handler(handler)
	}
	if n := int64(web.cfg.MaxBodyBytes); n > 0 {
		handler = maxBodyBytesHandler(handler, n)
	}
	if web.cfg.HTTP1FullDuplexEnable {
		handler = http1FullDuplexHandler(handler)
	}
	ret = &http.Server{
		Handler:                      handler,
		Addr:                         web.cfg.HostPort(),
		DisableGeneralOptionsHandler: web.cfg.DisableGeneralOptionsHandler,
		ReadTimeout:                  web.cfg.ReadTimeout,
		ReadHeaderTimeout:            web.cfg.ReadHeaderTimeout,
		WriteTimeout:                 web.cfg.WriteTimeout,
		IdleTimeout:                  web.cfg.IdleTimeout,
		MaxHeaderBytes:               web.cfg.MaxHeaderBytes,
		MaxHeaderValueCount:          web.cfg.MaxHeaderValueCount,
		DisableClientPriority:        web.cfg.DisableClientPriority,
		ConnState:                    web.cfg.ConnState,
		ErrorLog:                     web.cfg.ErrorLog,
		BaseContext:                  web.cfg.BaseContext,
		ConnContext:                  web.cfg.ConnContext,
	}
	// Настройка протоколов HTTP. HTTP/1 всегда включён.
	ret.Protocols = new(http.Protocols)
	ret.Protocols.SetHTTP1(true)
	if !web.cfg.HTTP2Disable {
		ret.Protocols.SetHTTP2(true)
	}
	if web.cfg.HTTP2UnencryptedEnable {
		ret.Protocols.SetUnencryptedHTTP2(true)
	}
	// Настройка протокола HTTP/2.
	if ret.Protocols.HTTP2() || ret.Protocols.UnencryptedHTTP2() {
		ret.HTTP2 = &http.HTTP2Config{
			MaxConcurrentStreams:          web.cfg.HTTP2MaxConcurrentStreams,
			MaxDecoderHeaderTableSize:     web.cfg.HTTP2MaxDecoderHeaderTableSize,
			MaxEncoderHeaderTableSize:     web.cfg.HTTP2MaxEncoderHeaderTableSize,
			MaxReadFrameSize:              web.cfg.HTTP2MaxReadFrameSize,
			MaxReceiveBufferPerConnection: web.cfg.HTTP2MaxReceiveBufferPerConnection,
			MaxReceiveBufferPerStream:     web.cfg.HTTP2MaxReceiveBufferPerStream,
			SendPingTimeout:               web.cfg.HTTP2SendPingTimeout,
			PingTimeout:                   web.cfg.HTTP2PingTimeout,
			WriteByteTimeout:              web.cfg.HTTP2WriteByteTimeout,
		}
	}
	// Конфигурация TLS по убыванию приоритета: переданная аргументом, затем конфигурация из Configuration,
	// затем конфигурация, созданная по ключам PEM. Если ни один источник не доступен, сервер работает
	// без шифрования.
	if tlsConfig == nil {
		tlsConfig = web.cfg.TLSConfig
	}
	if tlsConfig == nil {
		if web.cfg.TLSPrivateKeyPEM == "" || web.cfg.TLSPublicKeyPEM == "" {
			isTLS = false
			return
		}
		if tlsConfig, web.err = web.net.NewTLSConfigDefault(web.cfg.TLSPublicKeyPEM, web.cfg.TLSPrivateKeyPEM); web.err != nil {
			return
		}
	}
	// Требование и проверка клиентского сертификата (mTLS).
	// Явно заданный пользователем ClientAuth конфигурации TLS не переопределяется.
	if web.cfg.TLSClientAuth != 0 && tlsConfig.ClientAuth == 0 {
		tlsConfig.ClientAuth = web.cfg.TLSClientAuth
	}
	if web.cfg.TLSClientCAs != nil {
		tlsConfig.ClientCAs = web.cfg.TLSClientCAs
	}
	// Фильтрация доменов на этапе TLS handshake по SNI.
	if df != nil {
		web.wrapSNIConfig(tlsConfig, df)
	}
	// Порядок протоколов определяет предпочтение сервера при согласовании ALPN.
	if ret.Protocols.HTTP2() {
		tlsConfig.NextProtos = append(tlsConfig.NextProtos, proto2)
		// Сервер с включённым протоколом HTTP/2 не запускается, если в списке шифров отсутствует хотя бы
		// один шифр, обязательный для HTTP/2.
		tlsConfig.CipherSuites = ensureHTTP2CipherSuites(tlsConfig.CipherSuites)
	}
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, http11)
	ret.TLSConfig, isTLS = tlsConfig, true

	return
}

// Обёртка функции GetConfigForClient конфигурации TLS: фильтрация домена по SNI с сохранением
// пользовательской функции, если она назначена.
func (web *impl) wrapSNIConfig(tlsConfig *tls.Config, df *domainFilter) {
	var userGetConfigForClient = tlsConfig.GetConfigForClient

	tlsConfig.GetConfigForClient = func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
		var err error
		if _, err = df.sniFilter(hello); err != nil {
			return nil, err
		}
		if userGetConfigForClient != nil {
			return userGetConfigForClient(hello)
		}

		return nil, nil
	}
}
