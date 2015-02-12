package cuba

import (
	"log"
	"net/http"
)

type Handler func(*Context) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{w, r, make(map[string]string), r.URL.Path}

	err := h.serveContext(context)
	if err != nil {
		log.Println(err)
	}
}

func (h Handler) serveContext(context *Context) error {
	err := h(context)
	if err != nil {
		return err
	}

	return nil
}

type ContextHandler interface {
	serveContext(*Context) error
}
