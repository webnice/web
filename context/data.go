package context

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "gopkg.in/webnice/web.v1/header"
import "gopkg.in/webnice/web.v1/mime"
import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// RegisterGlobalVerifyPlugin Register global external library of data verification
func RegisterGlobalVerifyPlugin(vp VerifyPlugin) { globalVerifyPlugin = vp }

// DataError Create formated error
func (ctx *impl) DataError(format string, a ...interface{}) (rsp []byte, err error) {
	err = fmt.Errorf(format, a...)
	if globalVerifyPlugin != nil {
		rsp = globalVerifyPlugin.Error400(err)
	}
	return
}

// Data Extracting from a request and decoding data to structure of obj
func (ctx *impl) Data(obj interface{}) (rsp []byte, err error) {
	var req *bytes.Buffer
	var written int64
	var ct string

	if ctx.Request == nil {
		rsp, err = ctx.DataError("net/http is nil, can't retrieve data")
		return
	}
	req = &bytes.Buffer{}
	// Получение запроса
	defer func() { _ = ctx.Request.Body.Close() }()
	if written, err = io.Copy(req, ctx.Request.Body); err != nil {
		rsp, err = ctx.DataError("Error read request data: %s", err)
		return
	} else if written < 2 {
		rsp, err = ctx.DataError("Request data is empty")
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
		rsp, err = ctx.DataError("Decoding data error: %s", err)
		return
	}
	// Верификация данных с использованием внешней библиотеки
	if globalVerifyPlugin != nil {
		rsp, err = globalVerifyPlugin.Verify(obj)
	}

	return
}
