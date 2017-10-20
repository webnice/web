package verify

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import ()

// Interface interface
type Interface interface {
	// Code Set code
	Code(int64) Interface

	// Message Set message
	Message(string) Interface

	// Add new error to slice
	Add(Error) Interface

	// Response as structure
	Response() *Response

	// Json Serialize and return object as json data
	Json() []byte
}

// Response object structure
type Response struct {
	Error struct {
		Code    int64   `json:"code"`    // Уникальный код ошибки из справочника ошибок
		Message string  `json:"message"` // Описание ошибки в локализации текущей сессии, сообщение готовое для вывода в интерфейсе пользователя
		Errors  []Error `json:"errors"`  // Массив объектов с описанием имён полей и ошибок в них
	} `json:"error"`
}

// Error object structure
type Error struct {
	Field      string `json:"field"`      // Название поля
	FieldValue string `json:"fieldValue"` // Переданное значение поля field
	Message    string `json:"message"`    // Описание ошибки
}
