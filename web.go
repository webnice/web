package web

import (
	"github.com/labstack/echo/v4"
)

// New is a constructor of new web server implementation
func New() Interface {
	var wsv = new(web)
	wsv.route = echo.New()
	wsv.inCloseUp = make(chan bool, 1)
	wsv.isRun.Store(false)

	return wsv
}

// Error Return last error of web server
func (wsv *web) Error() error { return wsv.err }

// Route interface
func (wsv *web) Route() *echo.Echo { return wsv.route }

// Errors interface
//func (wsv *web) Errors() errors.Interface { return wsv.route.Errors() }

// Handlers interface
//func (wsv *web) Handlers() handlers.Interface { return wsv.route.Handlers() }
