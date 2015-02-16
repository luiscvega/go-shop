package cuba

type MuxHandler func (*Mux)

func (mh MuxHandler) serveContext(context *Context) error {
	m := New()

	mh(&m)

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
