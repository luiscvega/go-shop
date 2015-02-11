package cuba

import (
	"net/http"
	"regexp"
)

func New() mux {
	return mux{make([]route, 0)}
}

type route struct {
	method   string
	pattern  string
	captures []string
	handler  ContextHandler
}

type mux struct {
	routes []route
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{w, r, make(map[string]string)}

	m.serveContext(context)
}

func (m mux) serveContext(c *Context) error {
	if c.R.URL.Path == "/favicon.ico" {
		return nil
	}

	for _, route := range m.routes {
		if route.method != c.R.Method {
			continue
		}

		if route.pattern == c.R.URL.Path {
			route.handler.serveContext(c)
			break
		}

		if route.pattern == "/" {
			return nil
		}

		re := regexp.MustCompile(route.pattern)
		matches := re.FindAllStringSubmatch(c.R.URL.Path, -1)

		if len(matches) > 0 {
			for i, name := range route.captures {
				c.Params[name] = matches[0][i+1]
			}

			route.handler.serveContext(c)
			break
		}
	}

	return nil
}

func (m *mux) On(pattern string, nmux mux) {
	method := "GET"

	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	captures := make([]string, 0, len(matches))
	for _, match := range matches {
		captures = append(captures, match[1])

		pattern = re.ReplaceAllLiteralString(pattern, "([^\\/]+)")

	}

	m.routes = append(m.routes, route{method, pattern, captures, nmux})
}

func (m *mux) Add(method, pattern string, handler func(*Context) error) {
	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	captures := make([]string, 0, len(matches))
	for _, match := range matches {
		captures = append(captures, match[1])

		pattern = re.ReplaceAllLiteralString(pattern, "([^\\/]+)")

	}

	m.routes = append(m.routes, route{method, pattern, captures, Handler(handler)})
}

func (m *mux) Get(pattern string, handler Handler) {
	m.Add("GET", pattern, handler)
}

func (m *mux) Post(pattern string, handler Handler) {
	m.Add("POST", pattern, handler)
}

func (m *mux) Put(pattern string, handler Handler) {
	m.Add("PUT", pattern, handler)
}

func (m *mux) Delete(pattern string, handler Handler) {
	m.Add("DELETE", pattern, handler)
}

func (m mux) Routes() []route {
	return m.routes
}
