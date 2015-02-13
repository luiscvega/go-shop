package cuba

type Handler func(*Context) error

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
