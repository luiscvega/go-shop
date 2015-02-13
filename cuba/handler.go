package cuba

import (
	"fmt"
)

type MuxHandler func (*Mux)

func (mh MuxHandler) serveContext(context *Context) error {
	m := New()

	mh(&m)

	for _, routes := range m.Routes() {
		fmt.Println(routes)
	}

	err := m.serveContext(context)
	if err != nil {
		return err
	}

	return nil
}

type ContextHandler interface {
	serveContext(*Context) error
}

type Handler func(*Context) error

func (h Handler) serveContext(context *Context) error {
	err := h(context)
	if err != nil {
		return err
	}

	return nil
}
