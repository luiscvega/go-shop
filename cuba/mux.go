package cuba

import (
	"net/http"
	"regexp"
)

func New() mux {
	return mux{make(map[string][]route)}
}

type mux struct {
	table map[string][]route
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{w, r, make(map[string]string)}

	m.serveContext(context)
}

func (m mux) serveContext(c *Context) error {
	routes, ok := m.table[c.R.Method]
	if !ok {
		return nil
	}

	for _, route := range routes {
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

	// Initialize method
	_, ok := m.table[method]
	if !ok {
		m.table[method] = make([]route, 0)
	}

	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	captures := make([]string, 0, len(matches))
	for _, match := range matches {
		captures = append(captures, match[1])

		pattern = re.ReplaceAllLiteralString(pattern, "([^\\/]+)")

	}

	m.table[method] = append(m.table[method], route{pattern, captures, nmux})
}

func (m *mux) Add(method, pattern string, handler func(*Context) error) {
	// Initialize method
	_, ok := m.table[method]
	if !ok {
		m.table[method] = make([]route, 0)
	}

	re := regexp.MustCompile(`:(\w+)`)
	matches := re.FindAllStringSubmatch(pattern, -1)

	captures := make([]string, 0, len(matches))
	for _, match := range matches {
		captures = append(captures, match[1])

		pattern = re.ReplaceAllLiteralString(pattern, "([^\\/]+)")

	}

	m.table[method] = append(m.table[method], route{pattern, captures, Handler(handler)})
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

func (m mux) Table() map[string][]route {
	return m.table
}
