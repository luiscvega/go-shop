package cuba

import (
	"net/http"
	"regexp"
	"strings"
)

func New() mux {
	var patterns []string
	handlers := make(map[string]route)

	return mux{patterns, handlers}
}

type route struct {
	names   []string
	handler func(*Context)
}

type mux struct {
	patterns []string
	routes   map[string]route
}

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		if r.FormValue("_method") == "PUT" {
			r.Method = "PUT"
		}

		if r.FormValue("_method") == "DELETE" {
			r.Method = "DELETE"
		}
	}

	params := make(map[string]string)

	c := &Context{w, r, params}

	for _, pattern := range m.patterns {
		if pattern == r.URL.Path {
			m.routes[r.URL.Path].handler(c)
			return
		}

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			route := m.routes[pattern]

			for i, name := range route.names {
				c.Params[name] = matches[0][i+1]
			}

			route.handler(c)
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
