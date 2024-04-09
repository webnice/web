package web

import (
	"net"
	"testing"

	wnet "github.com/webnice/net"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestImpl_IdGenerate(t *testing.T) {
	var web Interface

	web = New().
		Handler(getTestHandlerFn(t))
	web.ListenAndServeWithConfig(&Configuration{
		Configuration: wnet.Configuration{
			Host: "localhost",
			Port: 10180,
			Mode: "tcp",
		},
	})
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if web.ID() == "" {
		t.Errorf("функция ID(), ошибка, ожидалось не пустое значение")
	}
}

func TestImpl_IdStatic(t *testing.T) {
	var (
		id  string
		web Interface
	)

	id = uuid.NewString()
	web = New().
		Handler(getTestHandlerFn(t))
	web.ListenAndServeWithConfig(&Configuration{
		Configuration: wnet.Configuration{
			ID:   id,
			Host: "localhost",
			Port: 10180,
			Mode: "tcp",
		},
	})
	defer web.Stop()
	if web.Error() != nil {
		t.Errorf("функция ListenAndServeWithConfig(), ошибка: %v, ожидалось: %v", web.Error(), nil)
	}
	if web.ID() != id {
		t.Errorf("функция ID(), вернулось: %q, ожидалось: %q", web.ID(), id)
	}
}

func TestImpl_ServeWithId(t *testing.T) {
	var (
		err      error
		id       string
		web      Interface
		listener net.Listener
	)

	id = uuid.NewString()
	web = New().(*impl)
	web.Handler(chi.NewMux())
	web.(*impl).cfg = &Configuration{
		Configuration: wnet.Configuration{
			Host: "localhost",
			Port: 18080,
		},
	}
	listener, err = web.NewListener(web.(*impl).cfg)
	if web.ServeWithId(listener, id); web.Error() != nil {
		t.Errorf("функция ServeWithId(), ошибка: %v, ожидалась: %v", err, nil)
	}
	defer web.Stop()
	if web.ID() != id {
		t.Errorf("функция ID(), вернулось: %q, ожидалось: %q", web.ID(), id)
	}
}
