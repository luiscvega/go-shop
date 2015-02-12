package cuba

import (
	"fmt"
	"log"
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

	err := m.serveContext(context)
	if err != nil {
		log.Println(err)
	}
}

var re = regexp.MustCompile(`:(\w+)`)

func (m mux) serveContext(c *Context) error {
	for _, route := range m.routes {
		if c.R.Method != route.method {
			continue
		}

		routePath := route.pattern

		if c.PathInfo == routePath {
			err := route.handler.serveContext(c)
			if err != nil {
				return err
			}

			break
		}

		names := match(&routePath)

		if len(names) > 0 {
			c.Params = consume(names, &routePath, c.PathInfo)
		}

		match := regexp.MustCompile(routePath).FindString(c.PathInfo)
		if len(match) > 0 {
			c.PathInfo = c.PathInfo[len(match):]
			fmt.Println(route)

			err := route.handler.serveContext(c)
			if err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func match(matcher *string) []string {
	matches := re.FindAllStringSubmatch(*matcher, -1)

	names := make([]string, len(matches))
	for i, match := range matches {
		names[i] = match[1]
	}

	return names
}

func consume(names []string, routePath *string, pathInfo string) map[string]string {
	params := make(map[string]string, len(names))

	*routePath = re.ReplaceAllLiteralString(*routePath, "([^\\/]+)")

	re := regexp.MustCompile(*routePath)
	values := re.FindAllStringSubmatch(pathInfo, -1)

	for i, value := range values {
		params[names[i]] = value[1]
	}

	return params
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
