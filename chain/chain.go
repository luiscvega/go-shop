package chain

import (
	"net/http"
)

func New(final http.Handler) *Chain {
	return &Chain{final}
}

type Chain struct {
	Final http.Handler
}

type Middleware func(http.Handler) http.Handler

func (c *Chain) Use(middleware Middleware) {
	c.Final = middleware(c.Final)
}

func (c Chain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Final.ServeHTTP(w, r)
}
