package status

const (
	// Continue Частично переданные на сервер данные проверены и сервер удовлетворён начальными данными. Клиент может продолжать отправку заголовков и данных
	Continue int = 100

	// SwitchingProtocols Switching Protocols. Сервер предлагает сменить протокол, перейти на более подходящий для указанного ресурса протокол. Список предлагаемых протоколов передаётся в заголовке Update Если клиент готов сменить протокол то новый запрос ожидается с указанием другого протокола
	SwitchingProtocols int = 101

	// Processing Processing (WebDAV). Запрос принят, на обработку запроса потребуется длительное время. Клиент, при получении такого ответа, должен сбросить таймер ожидания и ожидать следующего ответа в обычном режиме
	Processing int = 102

	// NameNotResolved Name Not Resolved. Ошибка при разрешении доменного имени, в связи с не верным или отсутствующих ip адресом DNS сервера
	NameNotResolved int = 105

	// Ok Успешный запрос. Результат запроса передаётся в заголовке или теле ответа
	Ok int = 200

	// Created Успешный запрос, в результате которого был создан новый ресурс. В ответе возвращается заголовок Location с указанием созданного ресурса. Так же возвращаются характеристики нового ресурса, Content-Type.
	Created int = 201

	// Accepted Запрос принят на обработку, но обработка не завершена. Клиенту нет необходимости ожидать окончания обработки запроса, так как процесс может быть весьма длительным
	Accepted int = 202

	// NonAuthoritativeInformation Non-Authoritative Information. Успешный запрос, но передаваемая в ответе информация взята из кэша или не является гарантированно достоверной (могла устареть)
	NonAuthoritativeInformation int = 203

	// NoContent No Content. Успешный запрос, в ответе передаются только заголовки без тела сообщения. Клиент не должен обновлять тело документа/ресурса, но может применить к нему полученные данные
	NoContent int = 204

	// ResetContent Reset Content. Сервер обязывает клиента сбросить введённые пользователем данные. Тела сообщения сервер клиенту не передаёт. Документ на стороне клиента обновлять не обязательно
	ResetContent int = 205

	// PartialContent Partial Content. Сервер удачно выполнил частичный GET запрос и возвращает только часть данных В заголовке Content-Range сервер указывает байтовые диапазоны содержимого
	PartialContent int = 206

	// MultiStatus Multi-Status (WebDAV). Сервер возвращает результат выполнения сразу нескольких независимых операций. Результат ожидается в виде XML тела ответа с объектом multistatus
	MultiStatus int = 207

	// AlreadyReported Already Reported (WebDAV). Код возвращается в случае запроса к ресурсу состоящему из коллекций данных повторяющихся между собой и фактически является указателем на то что данные ресурса А можно взять из ресурса Б так как они идентичны, либо данные ответа текущего запроса с итерацией совпадают с ответом предыдущих запросов с итерацией (ссылка). Так же сервер возвращает заголовок с ссылкой на данные на которые ссылается
	AlreadyReported int = 208

	// IMUsed IM Used. Заголовок A-IM от клиента был успешно принят и обработан. Сервер возвращает содержимое с учётом указанных параметров
	IMUsed int = 226

	// MultipleChoices Multiple Choices. Существует несколько вариантов предоставления ресурса по типу MIME Сервер возвращает список альтернатив для выбора на стороне клиента
	MultipleChoices int = 300

	// MovedPermanently Moved Permanently. Запрошенный ресурс был окончательно перенесён на новый URI Новый URI возвращается в заголовке Location
	MovedPermanently int = 301

	// FoundMovedTemporarily Found. Moved Temporarily. Запрошенный ресурс временно доступен по другому URI возвращаемому в заголовке Location
	FoundMovedTemporarily int = 302

	// SeeOther See Other. Ресурс по запрашиваемому URI необходимо запросить по адресу передаваемому в поле Location и только методом GET несмотря на первоначальный запрос иным методом
	SeeOther int = 303

	// NotModified Not Modified. Сервер отвечает кодом 304 если ресурс был запрошен методом GET с использованием заголовков If-Modified-Since или If-None-Match и документ не изменился с указанного момента. П этом ответ сервера не содержит тело ресурса
	NotModified int = 304

	// UseProxy Use Proxy. Запрос к запрашиваемому ресурсу должен осуществляться через прокси URI которого указывается в заголовки ответа в Location
	UseProxy int = 305

	// TemporaryRedirect Temporary Redirect. Запрошенный ресурс на короткое время доступен по другому URI указанному в заголовке ответа Location
	TemporaryRedirect int = 307

	// PermanentRedirect Permanent Redirect (experiemental). Текущий запрос и все будущие запросы необходимо выполнить по другому URI ресурса. Новый адрес ресурса возвращается в заголовке Location ответа сервера
	PermanentRedirect int = 308

	// BadRequest Bad Request. В запросе клиента присутствует синтаксическая ошибка или не указаны обязательные параметры запроса
	BadRequest int = 400

	// Unauthorized Unauthorized. Для доступа у ресурсу необходима аутентификация клиента. Заголовок ответа сервера будет содержать поле WWW-Authenticate с перечнем условий  аутентификации. Клиент может повторить запрос включив в новый запрос заголовок Authorization с требуемыми для аутентификации данными
	Unauthorized int = 401

	// PaymentRequired Payment Required. Для доступа к ресурсу необходимо произвести оплату
	PaymentRequired int = 402

	// Forbidden Forbidden. Сервер отказывается выполнять запрос из за ограничений в доступе клиента к указанному ресурсу
	Forbidden int = 403

	// NotFound Not Found. На сервере отсутствует запрашиваемый ресурс
	NotFound int = 404

	// MethodNotAllowed Method Not Allowed. Указанный клиентом метод нельзя применять к запрошенному ресурсу. В ответе сервера будет заголовок Allow с перечисленными через запятую методами запроса к ресурсу
	MethodNotAllowed int = 405

	// NotAcceptable Not Acceptable. Запрошенный ресурс не удовлетворяет переданным в заголовке характеристикам. Если запрос был не HEAD, то сервер в ответе вернёт список доступных характеристик запрашиваемого ресурса
	NotAcceptable int = 406

	// ProxyAuthenticationRequired Proxy Authentication Required. Для доступа у ресурсу необходима аутентификация клиента на прокси сервере. Заголовок ответа сервера будет содержать поле WWW-Authenticate с перечнем условий  аутентификации. Клиент может повторить запрос включив в новый запрос заголовок Authorization с требуемыми для аутентификации данными
	ProxyAuthenticationRequired int = 407

	// RequestTimeout Request Timeout. Время ожидания сервера окончания передачи данных клиентом истекло, запрос прерван
	RequestTimeout int = 408

	// Conflict Conflict. Запрос не может быть выполнен из за конфликта обращения к ресурсу, например ресурс заблокирован другим клиентом
	Conflict int = 409

	// Gone Gone. Ресурс ранее находившийся по указанному адресу удалён и не доступен, серверу не известно новый адрес ресурса
	Gone int = 410

	// LengthRequired Length Required. Для выполнения запроса к указанному ресурсу клиент должен передать Content-Length в заголовке запроса
	LengthRequired int = 411

	// PreconditionFailed Precondition Failed. Сервер не смог распознать ни один заголовок запроса
	PreconditionFailed int = 412

	// RequestEntityTooLarge Request Entity Too Large. Сервер отказывается выполнять запрос из за слишком большого тела запроса В случае если проблема временная, то возвращается заголовок Retry-After с указанием времени, по истечении которого можно повторить запрос
	RequestEntityTooLarge int = 413

	// RequestURLTooLong Request-URL Too Long. Сервер не может выполнить запрос к ресурсу из за слишком большого запрашиваемого URI
	RequestURLTooLong int = 414

	// UnsupportedMediaType Unsupported Media Type. Сервер отказывается работать с указанным типом данных передаваемого контента данным методом запроса
	UnsupportedMediaType int = 415

	// RequestedRangeNotSatisfiable Requested Range Not Satisfiable. Не корректный запрос с указанием Range. В поле Range заголовка запроса указан диапазон за пределами ресурса и отсутствует поле If-Range Если клиент передал байтовый диапазон, то сервер вернёт ресурса в заголовке Content-Range ответа
	RequestedRangeNotSatisfiable int = 416

	// ExpectationFailed Expectation Failed. Сервер не может удовлетворить переданный от клиента заголовок запроса Expect (ждать)
	ExpectationFailed int = 417

	// ImATeapot I'm a teapot (RFC 2324). Я чайник. Первоапрельская шутка от IETF :)
	ImATeapot int = 418

	// EnhanceYourCalm Enhance Your Calm (Twitter). Ответ аналогичен ответу 429 – слишком много запросов, призван заставить клиента отправлять меньшее количество запросов к ресурсу. В ответе сервера возвращается заголовок Retry-After с указанием времени через которое можно повторить запрос к серверу
	EnhanceYourCalm int = 420

	// UnprocessableEntity Unprocessable Entity (WebDAV).  Не обрабатываемая сущность. Запрос к серверу корректный и верный, но в теле запроса имеется логическая ошибка из за которой не возможно выполнить запрос к ресурсу
	UnprocessableEntity int = 422

	// Locked Locked (WebDAV). В результате запроса ресурс успешно заблокирован
	Locked int = 423

	// FailedDependency Failed Dependency (WebDAV). Успешность выполнения запроса зависит от не разрешимой зависимости (другого запроса) который ещё не выполнен
	FailedDependency int = 424

	// UnorderedCollection Unordered Collection. Запрос коллекции ресурсов в не корректном порядке или запрос к упорядоченному ресурсу в не корректном порядке
	UnorderedCollection int = 425

	// UpgradeRequired Upgrade Required. Клиенту необходимо обновить протокол запроса. В заголовке ответа сервера Upgrade и Connection возвращаются указания, которые должен выполнить клиент для успешности запроса
	UpgradeRequired int = 426

	// PreconditionRequired Precondition Required. Сервер требует от клиента выполнить запрос с указанием заголовков Range и If-Match
	PreconditionRequired int = 428

	// TooManyRequests Too Many Requests. Слишком много запросов клиента к ресурсу. В ответе сервера возвращается заголовок Retry-After с указанием времени через которое можно повторить запрос к серверу
	TooManyRequests int = 429

	// RequestHeaderFieldsTooLarge Request Header Fields Too Large. От клиента получено слишком много заголовков или длинна заголовка превысило допустимые размеры. Запрос прерван
	RequestHeaderFieldsTooLarge int = 431

	// RequestedHostUnavailable Requested Host Unavailable. Запрашиваемый ресурс не доступен
	RequestedHostUnavailable int = 434

	// NoResponse No Response (Nginx). Сервер отказывается или не может вернуть результат запроса к ресурсу и незамедлительно закрывает соединение
	NoResponse int = 444

	// RetryWith Retry With (Microsoft). В запросе к ресурсу не достаточно информации для его успешного выполнения. В заголовке ответа сервера передаётся Ms-Echo-Request с указанием необходимых полей
	RetryWith int = 449

	// BlockedByParentalControls Blocked by Parental Controls (Microsoft). Запрос заблокирован системой «родительский контроль»
	BlockedByParentalControls int = 450

	// UnavailableForLegalReasons Unavailable For Legal Reasons. Доступ к ресурсу закрыт по юридическим причинам или по требованиям органов государственной власти
	UnavailableForLegalReasons int = 451

	// UnrecoverableError Unrecoverable Error. Обработка запроса вызывает не обрабатываемые сбои в базе данных или её таблицах
	UnrecoverableError int = 456

	// ClientClosedRequest Client Closed Request (Nginx). Код введён nginx для записи в логи. Указывает на то что клиент закрыл соединение не попытавшись получить от сервера ответ
	ClientClosedRequest int = 499

	// InternalServerError Internal Server Error. Любая не описанная ошибка на стороне сервера
	InternalServerError int = 500

	// NotImplemented Not Implemented. Сервер не имеет возможностей для удовлетворения доступа к ресурсу или реализация обработки запроса ещё не закончена
	NotImplemented int = 501

	// BadGateway Bad Gateway. Сервер, выступающий в роли шлюза или прокси, получил не корректный ответ от вышестоящего сервера
	BadGateway int = 502

	// ServiceUnavailable Service Unavailable. По техническим причинам сервер не может выполнить запрос к ресурсу. В ответе сервер возвращает заголовок Retry-After с указанием времени через которое клиент может повторить запрос
	ServiceUnavailable int = 503

	// GatewayTimeout Gateway Timeout. Сервер, выступающий в роли шлюза или прокси, не дождался ответа от вышестоящего сервера
	GatewayTimeout int = 504

	// HTTPVersionNotSupported HTTP Version Not Supported. Сервер отказывается поддерживать указанную в запросе версию HTTP протокола
	HTTPVersionNotSupported int = 505

	// VariantAlsoNegotiates Variant Also Negotiates. Ошибка на стороне сервера связанная с циклической или рекурсивной задачей которая не может завершиться
	VariantAlsoNegotiates int = 506

	// InsufficientStorage Insufficient Storage (WebDAV). Не достаточно места или ресурсов пользователя для выполнения запросов. Код может возвращаться как по причинам физической не хватки ресурсов так и по причине текущих лимитов пользователя
	InsufficientStorage int = 507

	// LoopDetected Loop Detected (WebDAV). Бесконечный цикл. Запрос завершился в результате вызванного на сервере бесконечного цикла который не мог закончится
	LoopDetected int = 508

	// BandwidthLimitExceeded Bandwidth Limit Exceeded (Apache). Запрос прерван в связи с превышением со стороны клиента ограничений на скорость доступа к ресурсам
	BandwidthLimitExceeded int = 509

	// NotExtended Not Extended. На сервере отсутствует расширение которое необходимо для успешного выполнения запроса к ресурсу
	NotExtended int = 510

	// NetworkAuthenticationRequired Network Authentication Required. Запрос был прерван сервером посредником, прокси или шлюзом из за необходимости авторизации клиента на сервере посреднике до начала выполнения запросов
	NetworkAuthenticationRequired int = 511

	// NetworkReadTimeoutError Network Read Timeout Error. Ответ возвращается сервером который является прокси или шлюзом перед вышестоящим сервером и говорит о том что сервер не смог считать ответ на запрос к вышестоящему серверу
	NetworkReadTimeoutError int = 598

	// NetworkConnectTimeoutError Network Connect Timeout Error. Ответ возвращается сервером который является прокси или шлюзом перед вышестоящим сервером и говорит о том что сервер не смог установить связь с вышестоящим сервером (подключиться к серверу)
	NetworkConnectTimeoutError int = 599
)

var statusText = map[int]string{
	Continue:                      `Continue`,                        // Частично переданные на сервер данные проверены и сервер удовлетворён начальными данными. Клиент может продолжать отправку заголовков и данных
	SwitchingProtocols:            `Switching Protocols`,             // Сервер предлагает сменить протокол, перейти на более подходящий для указанного ресурса протокол. Список предлагаемых протоколов передаётся в заголовке Update Если клиент готов сменить протокол то новый запрос ожидается с указанием другого протокола
	Processing:                    `Processing`,                      // Запрос принят, на обработку запроса потребуется длительное время. Клиент, при получении такого ответа, должен сбросить таймер ожидания и ожидать следующего ответа в обычном режиме
	NameNotResolved:               `Name Not Resolved`,               // Ошибка при разрешении доменного имени, в связи с не верным или отсутствующих ip адресом DNS сервера
	Ok:                            `Ok`,                              // Успешный запрос. Результат запроса передаётся в заголовке или теле ответа
	Created:                       `Created`,                         // Успешный запрос, в результате которого был создан новый ресурс. В ответе возвращается заголовок Location с указанием созданного ресурса. Так же возвращаются характеристики нового ресурса, Content-Type.
	Accepted:                      `Accepted`,                        // Запрос принят на обработку, но обработка не завершена. Клиенту нет необходимости ожидать окончания обработки запроса, так как процесс может быть весьма длительным
	NonAuthoritativeInformation:   `Non-Authoritative Information`,   // Успешный запрос, но передаваемая в ответе информация взята из кэша или не является гарантированно достоверной (могла устареть)
	NoContent:                     `No Content`,                      // Успешный запрос, в ответе передаются только заголовки без тела сообщения. Клиент не должен обновлять тело документа/ресурса, но может применить к нему полученные данные
	ResetContent:                  `Reset Content`,                   // Сервер обязывает клиента сбросить введённые пользователем данные. Тела сообщения сервер клиенту не передаёт. Документ на стороне клиента обновлять не обязательно
	PartialContent:                `Partial Content`,                 // Сервер удачно выполнил частичный GET запрос и возвращает только часть данных В заголовке Content-Range сервер указывает байтовые диапазоны содержимого
	MultiStatus:                   `Multi-Status`,                    // Сервер возвращает результат выполнения сразу нескольких независимых операций. Результат ожидается в виде XML тела ответа с объектом multistatus
	AlreadyReported:               `Already Reported`,                // Код возвращается в случае запроса к ресурсу состоящему из коллекций данных повторяющихся между собой и фактически является указателем на то что данные ресурса А можно взять из ресурса Б так как они идентичны, либо данные ответа текущего запроса с итерацией совпадают с ответом предыдущих запросов с итерацией (ссылка). Так же сервер возвращает заголовок с ссылкой на данные на которые ссылается
	IMUsed:                        `IM Used`,                         // Заголовок A-IM от клиента был успешно принят и обработан. Сервер возвращает содержимое с учётом указанных параметров
	MultipleChoices:               `Multiple Choices`,                // Существует несколько вариантов предоставления ресурса по типу MIME Сервер возвращает список альтернатив для выбора на стороне клиента
	MovedPermanently:              `Moved Permanently`,               // Запрошенный ресурс был окончательно перенесён на новый URI Новый URI возвращается в заголовке Location
	FoundMovedTemporarily:         `Found. Moved Temporarily`,        // Запрошенный ресурс временно доступен по другому URI возвращаемому в заголовке Location
	SeeOther:                      `See Other`,                       // Ресурс по запрашиваемому URI необходимо запросить по адресу передаваемому в поле Location и только методом GET несмотря на первоначальный запрос иным методом
	NotModified:                   `Not Modified`,                    // Сервер отвечает кодом 304 если ресурс был запрошен методом GET с использованием заголовков If-Modified-Since или If-None-Match и документ не изменился с указанного момента. П этом ответ сервера не содержит тело ресурса
	UseProxy:                      `Use Proxy`,                       // Запрос к запрашиваемому ресурсу должен осуществляться через прокси URI которого указывается в заголовки ответа в Location
	TemporaryRedirect:             `Temporary Redirect`,              // Запрошенный ресурс на короткое время доступен по другому URI указанному в заголовке ответа Location
	PermanentRedirect:             `Permanent Redirect`,              // Текущий запрос и все будущие запросы необходимо выполнить по другому URI ресурса. Новый адрес ресурса возвращается в заголовке Location ответа сервера
	BadRequest:                    `Bad Request`,                     // В запросе клиента присутствует синтаксическая ошибка или не указаны обязательные параметры запроса
	Unauthorized:                  `Unauthorized`,                    // Для доступа у ресурсу необходима аутентификация клиента. Заголовок ответа сервера будет содержать поле WWW-Authenticate с перечнем условий  аутентификации. Клиент может повторить запрос включив в новый запрос заголовок Authorization с требуемыми для аутентификации данными
	PaymentRequired:               `Payment Required`,                // Для доступа к ресурсу необходимо произвести оплату
	Forbidden:                     `Forbidden`,                       // Сервер отказывается выполнять запрос из за ограничений в доступе клиента к указанному ресурсу
	NotFound:                      `Not Found`,                       // На сервере отсутствует запрашиваемый ресурс
	MethodNotAllowed:              `Method Not Allowed`,              // Указанный клиентом метод нельзя применять к запрошенному ресурсу. В ответе сервера будет заголовок Allow с перечисленными через запятую методами запроса к ресурсу
	NotAcceptable:                 `Not Acceptable`,                  // Запрошенный ресурс не удовлетворяет переданным в заголовке характеристикам. Если запрос был не HEAD, то сервер в ответе вернёт список доступных характеристик запрашиваемого ресурса
	ProxyAuthenticationRequired:   `Proxy Authentication Required`,   // Для доступа у ресурсу необходима аутентификация клиента на прокси сервере. Заголовок ответа сервера будет содержать поле WWW-Authenticate с перечнем условий  аутентификации. Клиент может повторить запрос включив в новый запрос заголовок Authorization с требуемыми для аутентификации данными
	RequestTimeout:                `Request Timeout`,                 // Время ожидания сервера окончания передачи данных клиентом истекло, запрос прерван
	Conflict:                      `Conflict`,                        // Запрос не может быть выполнен из за конфликта обращения к ресурсу, например ресурс заблокирован другим клиентом
	Gone:                          `Gone`,                            // Ресурс ранее находившийся по указанному адресу удалён и не доступен, серверу не известно новый адрес ресурса
	LengthRequired:                `Length Required`,                 // Для выполнения запроса к указанному ресурсу клиент должен передать Content-Length в заголовке запроса
	PreconditionFailed:            `Precondition Failed`,             // Сервер не смог распознать ни один заголовок запроса
	RequestEntityTooLarge:         `Request Entity Too Large`,        // Сервер отказывается выполнять запрос из за слишком большого тела запроса В случае если проблема временная, то возвращается заголовок Retry-After с указанием времени, по истечении которого можно повторить запрос
	RequestURLTooLong:             `Request-URL Too Long`,            // Сервер не может выполнить запрос к ресурсу из за слишком большого запрашиваемого URI
	UnsupportedMediaType:          `Unsupported Media Type`,          // Сервер отказывается работать с указанным типом данных передаваемого контента данным методом запроса
	RequestedRangeNotSatisfiable:  `Requested Range Not Satisfiable`, // Не корректный запрос с указанием Range. В поле Range заголовка запроса указан диапазон за пределами ресурса и отсутствует поле If-Range Если клиент передал байтовый диапазон, то сервер вернёт ресурса в заголовке Content-Range ответа
	ExpectationFailed:             `Expectation Failed`,              // Сервер не может удовлетворить переданный от клиента заголовок запроса Expect (ждать)
	ImATeapot:                     `I'm a teapot (RFC 2324)`,         // Я чайник. Первоапрельская шутка от IETF :)
	EnhanceYourCalm:               `Enhance Your Calm`,               // Ответ аналогичен ответу 429 – слишком много запросов, призван заставить клиента отправлять меньшее количество запросов к ресурсу. В ответе сервера возвращается заголовок Retry-After с указанием времени через которое можно повторить запрос к серверу
	UnprocessableEntity:           `Unprocessable Entity`,            // Не обрабатываемая сущность. Запрос к серверу корректный и верный, но в теле запроса имеется логическая ошибка из за которой не возможно выполнить запрос к ресурсу
	Locked:                        `Locked`,                          // В результате запроса ресурс успешно заблокирован
	FailedDependency:              `Failed Dependency`,               // Успешность выполнения запроса зависит от не разрешимой зависимости (другого запроса) который ещё не выполнен
	UnorderedCollection:           `Unordered Collection`,            // Запрос коллекции ресурсов в не корректном порядке или запрос к упорядоченному ресурсу в не корректном порядке
	UpgradeRequired:               `Upgrade Required`,                // Клиенту необходимо обновить протокол запроса. В заголовке ответа сервера Upgrade и Connection возвращаются указания, которые должен выполнить клиент для успешности запроса
	PreconditionRequired:          `Precondition Required`,           // Сервер требует от клиента выполнить запрос с указанием заголовков Range и If-Match
	TooManyRequests:               `Too Many Requests`,               // Слишком много запросов клиента к ресурсу. В ответе сервера возвращается заголовок Retry-After с указанием времени через которое можно повторить запрос к серверу
	RequestHeaderFieldsTooLarge:   `Request Header Fields Too Large`, // От клиента получено слишком много заголовков или длинна заголовка превысило допустимые размеры. Запрос прерван
	RequestedHostUnavailable:      `Requested Host Unavailable`,      // Запрашиваемый ресурс не доступен
	NoResponse:                    `No Response`,                     // Сервер отказывается или не может вернуть результат запроса к ресурсу и незамедлительно закрывает соединение
	RetryWith:                     `Retry With`,                      // В запросе к ресурсу не достаточно информации для его успешного выполнения. В заголовке ответа сервера передаётся Ms-Echo-Request с указанием необходимых полей
	BlockedByParentalControls:     `Blocked by Parental Controls`,    // Запрос заблокирован системой «родительский контроль»
	UnavailableForLegalReasons:    `Unavailable For Legal Reasons`,   // Доступ к ресурсу закрыт по юридическим причинам или по требованиям органов государственной власти
	UnrecoverableError:            `Unrecoverable Error`,             // Обработка запроса вызывает не обрабатываемые сбои в базе данных или её таблицах
	ClientClosedRequest:           `Client Closed Request`,           // Код введён nginx для записи в логи. Указывает на то что клиент закрыл соединение не попытавшись получить от сервера ответ
	InternalServerError:           `Internal Server Error`,           // Любая не описанная ошибка на стороне сервера
	NotImplemented:                `Not Implemented`,                 // Сервер не имеет возможностей для удовлетворения доступа к ресурсу или реализация обработки запроса ещё не закончена
	BadGateway:                    `Bad Gateway`,                     // Сервер, выступающий в роли шлюза или прокси, получил не корректный ответ от вышестоящего сервера
	ServiceUnavailable:            `Service Unavailable`,             // По техническим причинам сервер не может выполнить запрос к ресурсу. В ответе сервер возвращает заголовок Retry-After с указанием времени через которое клиент может повторить запрос
	GatewayTimeout:                `Gateway Timeout`,                 // Сервер, выступающий в роли шлюза или прокси, не дождался ответа от вышестоящего сервера
	HTTPVersionNotSupported:       `HTTP Version Not Supported`,      // Сервер отказывается поддерживать указанную в запросе версию HTTP протокола
	VariantAlsoNegotiates:         `Variant Also Negotiates`,         // Ошибка на стороне сервера связанная с циклической или рекурсивной задачей которая не может завершиться
	InsufficientStorage:           `Insufficient Storage`,            // Не достаточно места или ресурсов пользователя для выполнения запросов. Код может возвращаться как по причинам физической не хватки ресурсов так и по причине текущих лимитов пользователя
	LoopDetected:                  `Loop Detected`,                   // Бесконечный цикл. Запрос завершился в результате вызванного на сервере бесконечного цикла который не мог закончится
	BandwidthLimitExceeded:        `Bandwidth Limit Exceeded`,        // Запрос прерван в связи с превышением со стороны клиента ограничений на скорость доступа к ресурсам
	NotExtended:                   `Not Extended`,                    // На сервере отсутствует расширение которое необходимо для успешного выполнения запроса к ресурсу
	NetworkAuthenticationRequired: `Network Authentication Required`, // Запрос был прерван сервером посредником, прокси или шлюзом из за необходимости авторизации клиента на сервере посреднике до начала выполнения запросов
	NetworkReadTimeoutError:       `Network Read Timeout Error`,      // Ответ возвращается сервером который является прокси или шлюзом перед вышестоящим сервером и говорит о том что сервер не смог считать ответ на запрос к вышестоящему серверу
	NetworkConnectTimeoutError:    `Network Connect Timeout Error`,   // Ответ возвращается сервером который является прокси или шлюзом перед вышестоящим сервером и говорит о том что сервер не смог установить связь с вышестоящим сервером (подключиться к серверу)
}
