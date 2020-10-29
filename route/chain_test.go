package route

import (
	"net/http"
	"testing"
)

type testHandler struct {
	C int64
}

func (h *testHandler) ServeHTTP(wr http.ResponseWriter, rq *http.Request) { h.C++ }
func (h *testHandler) Middleware(next http.Handler) http.Handler {
	var fn = func(wr http.ResponseWriter, rq *http.Request) {
		h.ServeHTTP(wr, rq)
		next.ServeHTTP(wr, rq)
	}
	return http.HandlerFunc(fn)
}

func TestChainHandlerServeHTTP(t *testing.T) {
	var (
		obj *ChainHandler
		h   *testHandler
		rq  *http.Request
	)

	obj = new(ChainHandler)
	h = new(testHandler)
	obj.chain = h
	obj.Endpoint = h
	rq, _ = http.NewRequest("", "", nil)
	obj.ServeHTTP(nil, rq)
	obj.ServeHTTP(nil, rq)
	if h.C != 2 {
		t.Errorf("Error (*ChainHandler) ServeHTTP()")
	}
}

func TestHandler(t *testing.T) {
	var (
		mws  *Middlewares
		h    *testHandler
		hndl http.Handler
		rq   *http.Request
	)

	mws = new(Middlewares)
	h = new(testHandler)
	hndl = mws.Handler(h)
	rq, _ = http.NewRequest("", "", nil)
	hndl.ServeHTTP(nil, rq)
	hndl.ServeHTTP(nil, rq)
	if h.C != 2 {
		t.Errorf("Error (Middlewares) Handler()")
	}
}

func TestHandlerFunc(t *testing.T) {
	var (
		mws  *Middlewares
		h    *testHandler
		hndl http.Handler
		rq   *http.Request
	)

	mws = new(Middlewares)
	h = new(testHandler)
	hndl = mws.HandlerFunc(h.ServeHTTP)
	rq, _ = http.NewRequest("", "", nil)
	hndl.ServeHTTP(nil, rq)
	hndl.ServeHTTP(nil, rq)
	if h.C != 2 {
		t.Errorf("Error (Middlewares) HandlerFunc()")
	}
}

func TestChainMiddlewares(t *testing.T) {
	var (
		mws  Middlewares
		h    *testHandler
		hndl http.Handler
		rq   *http.Request
	)

	h = new(testHandler)
	mws = make(Middlewares, 0)
	mws = append(mws, h.Middleware)
	mws = append(mws, h.Middleware)
	hndl = mws.Handler(h)
	rq, _ = http.NewRequest("", "", nil)
	hndl.ServeHTTP(nil, rq)
	hndl.ServeHTTP(nil, rq)
	if h.C != 6 {
		t.Errorf("Error chain(). Counter is %d expected 6", h.C)
	}
}
