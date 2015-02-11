package cuba

import (
	"net/http"
)

type Handler func(*Context) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{w, r, make(map[string]string)}

	h.serveContext(context)
}

func (h Handler) serveContext(context *Context) error {
	h(context)
	return nil
}

type ContextHandler interface {
	serveContext(*Context) error
}
