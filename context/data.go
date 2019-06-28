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
	"reflect"
	stddebug "runtime/debug"
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

func indirectValue(rv reflect.Value) reflect.Value {
	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	return rv
}

func indirectType(rt reflect.Type) reflect.Type {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	return rt
}

// Verify Верификация данных с использованием внешней библиотеки
func (ctx *impl) Verify(obj interface{}) (rsp []byte, err error) {
	var (
		rv   reflect.Value
		rt   reflect.Type
		item interface{}
	)

	if globalVerifyPlugin == nil {
		return
	}
	// При вызове reflect возможна паника
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic recovery:\n%v\n%s", e, string(stddebug.Stack()))
		}
	}()
	switch rt = indirectType(reflect.TypeOf(obj)); rt.Kind() {
	case reflect.Slice:
		rv = indirectValue(reflect.ValueOf(obj))
		for i := 0; i < rv.Len(); i++ {
			item = rv.Index(i).Interface()
			if rsp, err = globalVerifyPlugin.Verify(item); err != nil || len(rsp) > 0 {
				return
			}
		}
	default:
		rsp, err = globalVerifyPlugin.Verify(obj)
	}

	return
}

// Data Extracting from a request and decoding data to structure of obj
func (ctx *impl) Data(obj interface{}) (rsp []byte, err error) {
	var (
		req     *bytes.Buffer
		written int64
		ct      string
	)

	if ctx.Request == nil {
		rsp, err = ctx.DataError("net/http is nil, can't retrieve data")
		return
	}
	req = &bytes.Buffer{}
	// Получение запроса
	defer func() { _ = ctx.Request.Body.Close() }()
	if written, err = io.Copy(req, ctx.Request.Body); err != nil {
		rsp, err = ctx.DataError("reading data of request error: %s", err)
		return
	} else if written < 2 {
		rsp, err = ctx.DataError("request data is empty")
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
		err = fmt.Errorf("unknown content type: %q", ct)
		return
	}
	if err != nil {
		rsp, err = ctx.DataError("decoding data error: %s", err)
		return
	}
	// Верификация данных с использованием внешней библиотеки
	if rsp, err = ctx.Verify(obj); err != nil {
		rsp, err = ctx.DataError("verification of data error: %s", err)
		return
	}

	return
}
