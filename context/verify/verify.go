package verify

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"encoding/json"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// Code Set code
func (e *Response) Code(c int64) Interface { e.Error.Code = c; return e }

// Message Set message
func (e *Response) Message(m string) Interface { e.Error.Message = m; return e }

// Add new error to slice
func (e *Response) Add(err Error) Interface {
	e.Error.Errors = append(e.Error.Errors, err)
	return e
}

// Response as structure
func (e *Response) Response() *Response { return e }

// Json Serialize and return object as json data
func (e *Response) Json() (ret []byte) {
	var err error
	if ret, err = json.Marshal(e); err != nil {
		//e.Error = fmt.Errorf("Error marshal json: %s", err.Error())
		return
	}
	return
}

// Verify Проверяет структуру разобранную из json на ошибки описанные в теге validate
// В случае наличия ошибок возвращает Interface с данными составленными на основе ошибок верификации
func Verify(data interface{}) (ret Interface, err error) {
	var vrf *validator.Validate
	var i int
	var terr error
	var vErr validator.ValidationErrors

	ret = E4xx()
	vrf = validator.New()
	terr = vrf.Struct(data)
	if terr == nil {
		return
	}

	vErr = terr.(validator.ValidationErrors)
	for i = range vErr {
		ret.Add(Error{
			Field:      vErr[i].StructField(),
			FieldValue: fmt.Sprintf("%v", vErr[i].Value()),
			Message:    vErr[i].ActualTag(),
		})
	}
	err = fmt.Errorf("Found verification error")

	return
}
