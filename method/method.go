package method // import "gopkg.in/webnice/web.v1/method"

import "strings"
import "fmt"

const (
	// Get HTTP Method "GET"
	Get Method = 1 << iota

	// Post HTTP Method "POST"
	Post

	// Put HTTP Method "PUT"
	Put

	// Delete HTTP Method "DELETE"
	Delete

	// Connect HTTP Method "CONNECT"
	Connect

	// Head HTTP Method "HEAD"
	Head

	// Patch HTTP Method "PATCH" RFC 5789
	Patch

	// Options HTTP Method "OPTIONS"
	Options

	// Trace HTTP Method "TRACE"
	Trace

	// Stub internal Method "STUB"
	Stub

	// Any http method
	Any Method = Get | Post | Put | Delete | Connect | Head | Patch | Options | Trace
)

var methodMap = map[string]Method{
	`GET`:     Get,
	`POST`:    Post,
	`PUT`:     Put,
	`DELETE`:  Delete,
	`CONNECT`: Connect,
	`HEAD`:    Head,
	`PATCH`:   Patch,
	`OPTIONS`: Options,
	`TRACE`:   Trace,
}

// Method type
type Method int64

// All return all http methods
func All() (ret []Method) {
	var key string
	for key = range methodMap {
		ret = append(ret, methodMap[key])
	}
	return
}

// Parse string and return method
func Parse(mtd string) (ret Method, err error) {
	var ok bool
	if ret, ok = methodMap[strings.ToUpper(mtd)]; !ok {
		err = fmt.Errorf("Unknown request method %s", mtd)
		return
	}
	return
}

// String Convert to string
func (mt Method) String() string {
	switch mt {
	case Get:
		return `GET`
	case Post:
		return `POST`
	case Put:
		return `PUT`
	case Delete:
		return `DELETE`
	case Connect:
		return `CONNECT`
	case Head:
		return `HEAD`
	case Patch:
		return `PATCH`
	case Options:
		return `OPTIONS`
	case Trace:
		return `TRACE`
	}
	return ``
}

// Int64 Convert to int64
func (mt Method) Int64() int64 { return int64(mt) }
