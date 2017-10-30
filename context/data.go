package context

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import "gopkg.in/webnice/web.v1/context/verify"
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Data Extracting from a request and decoding data to structure of obj
func (ctx *impl) Data(obj interface{}) (vfi verify.Interface, err error) {
	var req *bytes.Buffer
	var written int64
	var ct string

	if ctx.Request == nil {
		err = errors.New("net/http is nil, can't retrieve data")
		return
	}
	req = &bytes.Buffer{}
	vfi = verify.E4xx()

	// Получение запроса
	if written, err = io.Copy(req, ctx.Request.Body); err != nil {
		err = fmt.Errorf("Error read request data: %s", err.Error())
		return
	} else if written < 2 {
		vfi.Code(5)
		err = fmt.Errorf("Request data is empty")
		return
	}
	// Тип кодирования выбирается на основе Content-Type заголовка
	ct = ctx.Request.Header.Get(header.ContentType)
	switch {
	case strings.Contains(ct, mime.ApplicationJSON):
		err = json.NewDecoder(req).Decode(obj)
	case strings.Contains(ct, mime.TextXML), strings.Contains(ct, mime.ApplicationXML):
		err = xml.NewDecoder(req).Decode(obj)
	default:
		err = fmt.Errorf("Unknown content type: %q", ct)
		return
	}
	if err != nil {
		err = fmt.Errorf("Decoding data error: %s", err.Error())
		return
	}
	// Верификация полученных данных
	if vfi, err = verify.Verify(obj); err != nil {
		return
	}

	return
}
