package cuba

import (
	"net/http"
	"regexp"
	"strings"
)

func New() mux {
	var (
		patterns = make([]string, 0)
		handlers = make(map[string]route)
		params   = make(map[string]string)
		context  = &Context{Params: params}
	)

	return mux{context, patterns, handlers}
}

type route struct {
	names   []string
	handler func(*Context)
}

type mux struct {
	Context *Context

	patterns []string
	routes   map[string]route
}

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Context.W = w
	m.Context.R = r

	if r.Method == "POST" {
		r.ParseForm()

		if r.FormValue("_method") == "PUT" {
			r.Method = "PUT"
		}

		if r.FormValue("_method") == "DELETE" {
			r.Method = "DELETE"
		}
	}

	for _, pattern := range m.patterns {
		if pattern == r.URL.Path {
			m.routes[r.URL.Path].handler(m.Context)
			return
		}

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			route := m.routes[pattern]

			for i, name := range route.names {
				m.Context.Params[name] = matches[0][i+1]
			}

			route.handler(m.Context)
			return
		}
	}
}

func (m *mux) Add(pattern string, handler func(*Context)) {
	pattern, names := prepareHandler(pattern)

	m.patterns = append(m.patterns, pattern)
	m.routes[pattern] = route{names, handler}
}

func prepareHandler(pattern string) (string, []string) {
	re := regexp.MustCompile(":[a-zA-Z0-9]+")

	names := make([]string, 0)
	for _, match := range re.FindAllStringSubmatch(pattern, -1) {
		names = append(names, strings.Replace(match[0], ":", "", 1))
	}

	pattern = re.ReplaceAllLiteralString(pattern, "([a-zA-Z0-9]+)")

	return pattern, names
}
