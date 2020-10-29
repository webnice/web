package context

import (
	stdContext "context"
	"fmt"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	var ctx Interface

	if ctx = New(); ctx == nil {
		t.Errorf("Error New(), returns nil")
	}
	if ctx.Route() == nil {
		t.Errorf("Error New(), Route() is nil")
	}
	if ctx.Errors() == nil {
		t.Errorf("Error New(), Errors() is nil")
	}
	ctx.Errors(ctx.Errors())
	if ctx.Handlers() == nil {
		t.Errorf("Error New(), Handlers() is nil")
	}
	ctx.Handlers(ctx.Handlers())
	if ctx = New(""); ctx != nil {
		t.Errorf("Error New(), incorrect response")
	}
	if constContextKey.String() == "" {
		t.Errorf("Error key, string is empty")
	}
}

func TestContext(t *testing.T) {
	var (
		ctxOrig Interface
		ctx     *impl
		c       stdContext.Context
	)

	ctxOrig = New()
	c = stdContext.Background()
	if ctx = context(c); ctx != nil {
		t.Errorf("Error func context(), returns not nil")
	}
	c = stdContext.WithValue(stdContext.Background(), constContextKey.String(), ctxOrig)
	if ctx = context(c); ctx == nil {
		t.Errorf("Error func context(), returns nil")
	}
	if ctx != ctxOrig.(*impl) {
		t.Errorf("Error func context(), returns unknown context object")
	}
}

func TestNewFromRequest(t *testing.T) {
	var (
		ctx, ctxOrig Interface
		rq           *http.Request
		c            stdContext.Context
	)

	ctxOrig = New()
	c = stdContext.WithValue(stdContext.Background(), constContextKey.String(), ctxOrig)
	rq, _ = http.NewRequest("", `http://www.google.com/search?q=foo&q=bar`, nil)
	rq = rq.WithContext(c)

	if ctx = New(rq); ctx == nil {
		t.Errorf("Error New(*http.Request), returns nil")
	}
	if ctx != ctxOrig {
		t.Errorf("Error New(*http.Request), can't find original context from net/http.Request context")
	}
}

func TestNewFromStdContext(t *testing.T) {
	var (
		ctx, ctxOrig Interface
		c            stdContext.Context
	)

	ctxOrig = New()
	// Empty standard context
	c = stdContext.Background()
	if ctx = New(c); c == nil {
		t.Errorf("Error func New(stdContext), returns nil")
	}
	if ctx == ctxOrig {
		t.Errorf("Error func New(stdContext), returns incorrect context object")
	}
	// Standard context with context object
	c = stdContext.WithValue(stdContext.Background(), constContextKey.String(), ctxOrig)
	if ctx = New(c); c == nil {
		t.Errorf("Error func New(stdContext), returns nil")
	}
	if ctx != ctxOrig {
		t.Errorf("Error func New(stdContext), returns incorrect context object")
	}
}

func TestNewFromInterface(t *testing.T) {
	var (
		ctx, ctxOrig                        Interface
		addrErrors, addrHandlers, addrRoute string
	)

	ctxOrig = New()
	ctxOrig.(*impl).errors = nil
	ctxOrig.(*impl).handlers = nil
	ctxOrig.(*impl).route = nil
	ctx = New(ctxOrig)
	if ctx == nil {
		t.Errorf("Error func New(Interface), returns nil")
	}
	if ctx != ctxOrig {
		t.Errorf("Error func New(Interface), returns different object")
	}
	if ctx.Errors() == nil {
		t.Errorf("Error func New(Interface), returns object not contains Errors() interface")
	}
	if ctx.Handlers() == nil {
		t.Errorf("Error func New(Interface), returns object not contains Handlers() interface")
	}
	if ctx.Route() == nil {
		t.Errorf("Error func New(Interface), returns object not contains Route() interface")
	}
	ctxOrig = New()
	addrErrors = fmt.Sprintf("%p", ctxOrig.Errors())
	addrHandlers = fmt.Sprintf("%p", ctxOrig.Handlers())
	addrRoute = fmt.Sprintf("%p", ctxOrig.Route())
	ctx = New(ctxOrig)
	if addrErrors != fmt.Sprintf("%p", ctx.Errors()) {
		t.Errorf("Error func New(Interface), returns different Errors() interface")
	}
	if addrHandlers != fmt.Sprintf("%p", ctx.Handlers()) {
		t.Errorf("Error func New(Interface), returns different Handlers() interface")
	}
	if addrRoute == fmt.Sprintf("%p", ctx.Route()) {
		t.Errorf("Error func New(Interface), returns same Route() interface, expected new Route()")
	}
}

func TestIsContext(t *testing.T) {
	var (
		ctx Interface
		rq  *http.Request
	)

	ctx = New()
	rq, _ = http.NewRequest("", `http://www.google.com/search?q=foo&q=bar`, nil)
	rq = rq.WithContext(stdContext.Background())
	if IsContext(rq) {
		t.Errorf("Error IsContext(*http.Request), returns true, expected false")
	}
	rq = rq.WithContext(stdContext.WithValue(stdContext.Background(), constContextKey.String(), ctx))
	if !IsContext(rq) {
		t.Errorf("Error IsContext(*http.Request), returns false, expected true")
	}
}

func TestNewRequest(t *testing.T) {
	const (
		testKey   = `A2BD00BB-4E19-4F77-B5BE-A3F863C17129`
		testValue = `a0dae22b3922d1ff50a4c4e91aa9b3f32b876dfa394a8affff30b76fa2aed41a69de1b2b2bcb8bdb96874c0149e5a75bfc6c2a86eda2995b17a216df49356516`
	)
	var (
		ctx Interface
		rq  *http.Request
		c   = stdContext.WithValue(stdContext.Background(), testKey, testValue)
	)

	rq, _ = http.NewRequest("", `http://www.google.com/search?q=foo&q=bar`, nil)
	rq = rq.WithContext(c)
	ctx = New()
	rq = ctx.NewRequest(rq)
	if rq.Context().Value(testKey).(string) != testValue {
		t.Errorf("Error NewRequest() context inheritance error")
	}
	if New(rq) != ctx {
		t.Errorf("Error New(*http.Request) is not contains context")
	}
}
