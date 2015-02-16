package cuba

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	fmt.Println("Starting tests...")
}

func TestBoot(t *testing.T) {
	mux := New()

	mux.Get("/", func(c *Context) error {
		c.W.Write([]byte("+1"))
		return nil
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "+1" {
		t.Error("invalid body")
	}
}

func TestCaptures(t *testing.T) {
	mux := New()
	mux.Get("posts/:id", func(c *Context) error {
		c.W.Write([]byte(c.Params["id"]))
		return nil
	})

	req, _ := http.NewRequest("GET", "/posts/123", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "123" {
		t.Error("invalid captures")
	}
}

func TestMultipleCaptures(t *testing.T) {
	mux := New()
	mux.Get("posts/:id/name/:name", func(c *Context) error {
		c.W.Write([]byte(c.Params["id"] + "-" + c.Params["name"]))
		return nil
	})

	req, _ := http.NewRequest("GET", "/posts/123/name/watch", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "123-watch" {
		t.Error("invalid captures")
	}
}

func TestNestRoot(t *testing.T) {
	mux := New()

	mux.On("posts", func(mux *Mux) {
		mux.Get("/", func(c *Context) error {
			c.W.Write([]byte("+1"))
			return nil
		})
	})

	req, _ := http.NewRequest("GET", "/posts", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "+1" {
		t.Error("invalid body:", w.Body)
	}
}

func TestNestCaptures(t *testing.T) {
	mux := New()

	mux.On("posts", func(mux *Mux) {
		mux.Get(":id", func(c *Context) error {
			c.W.Write([]byte(c.Params["id"]))
			return nil
		})
	})

	req, _ := http.NewRequest("GET", "/posts/123", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "123" {
		t.Error("invalid body:", w.Body)
	}
}

func TestDoubleNest(t *testing.T) {
	mux := New()

	mux.On("posts", func(mux *Mux) {
		mux.On(":id", func(mux *Mux) {
			mux.Get("/", func(c *Context) error {
				c.W.Write([]byte(c.Params["id"]))
				return nil
			})
		})
	})

	req, _ := http.NewRequest("GET", "/posts/124", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "124" {
		t.Error("invalid body:", w.Body)
	}
}

func TestSkipRoute(t *testing.T) {
	mux := New()

	mux.On("posts", func(mux *Mux) {
		mux.Get("skipthis", func(c *Context) error {
			c.W.Write([]byte("i should have been skipped"))
			return nil
		})

		mux.Get("gethere", func(c *Context) error {
			c.W.Write([]byte("+1"))
			return nil
		})
	})

	req, _ := http.NewRequest("GET", "/posts/gethere", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Body.String() != "+1" {
		t.Error("invalid body:", w.Body)
	}
}
