package cuba

import (
	"net/http"
	"regexp"
)

func New() mux {
	return mux{make([]route, 0)}
}

type route struct {
	method  string
	pattern string
	handler ContextHandler
}

type mux struct {
	routes []route
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	context := &Context{w, r, make(map[string]string), r.URL.Path}

	m.serveContext(context)
}

func (m mux) serveContext(c *Context) error {
	for _, route := range m.routes {
		if c.R.Method != route.method {
			continue
		}

		if c.PathInfo == route.pattern {
			route.handler.serveContext(c)
			break
		}

		re := regexp.MustCompile(`:(\w+)`)
		matches := re.FindAllStringSubmatch(route.pattern, -1)

		if len(matches) > 0 {
			route.pattern = re.ReplaceAllLiteralString(route.pattern, "([^\\/]+)")

			re = regexp.MustCompile(route.pattern)
			values := re.FindAllStringSubmatch(c.PathInfo, -1)

			for i, value := range values {
				c.Params[matches[i][1]] = value[1]
			}
		}

		match := regexp.MustCompile(route.pattern).FindString(c.PathInfo)
		if len(match) > 0 {
			c.PathInfo = c.PathInfo[len(match):]
			route.handler.serveContext(c)
			break
		}
	}

	return nil
}

func (m *mux) On(pattern string, nmux mux) {
	m.routes = append(m.routes, route{"GET", pattern, nmux})
}

func (m *mux) Add(method, pattern string, handler func(*Context) error) {
	m.routes = append(m.routes, route{method, pattern, Handler(handler)})
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
