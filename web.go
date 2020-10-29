package web

import (
	"github.com/webnice/web/v1/context/errors"
	"github.com/webnice/web/v1/context/handlers"
	"github.com/webnice/web/v1/route"
)

// New is a constructor of new web server implementation
func New() Interface {
	var wsv = new(web)
	wsv.route = route.New()
	wsv.inCloseUp = make(chan bool, 1)
	wsv.isRun.Store(false)

	return wsv
}

// Error Return last error of web server
func (wsv *web) Error() error { return wsv.err }

// Route interface
func (wsv *web) Route() route.Interface { return wsv.route }

// Errors interface
func (wsv *web) Errors() errors.Interface { return wsv.route.Errors() }

// Handlers interface
func (wsv *web) Handlers() handlers.Interface { return wsv.route.Handlers() }
