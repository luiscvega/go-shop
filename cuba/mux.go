package cuba

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func New() Mux {
	return Mux{make([]route, 0)}
}

type route struct {
	method     string
	pattern    string
	paramNames []string
	handler    ContextHandler
}

type Mux struct {
	routes []route
}

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

	context := &Context{w, r, make(map[string]string), r.URL.Path}

	err := m.serveContext(context)
	if err != nil {
		log.Println(err)
	}
}

func (m Mux) serveContext(c *Context) error {
	var err error

	for _, route := range m.routes {
		if route.method != "ALL" && c.R.Method != route.method {
			continue
		}

		// Check for exact matches (e.g. "/products/new" == "/products/new")
		if c.PathInfo == route.pattern {
			err = route.handler.serveContext(c)
			break
		}

		if route.pattern == "/" {
			return nil
		}

		// Check for captures (e.g. "/products/([^\/]+)" =~ "/products/123")
		// /\A\/(#{pattern})(\/|\z)/)
		matches := regexp.MustCompile(route.pattern).FindAllStringSubmatch(c.PathInfo, -1)
		fmt.Println(matches, route.method, route.pattern, c.PathInfo)

		if len(matches) > 0 {
			if len(matches[0]) == 2 {
				for i, match := range matches {
					c.Params[route.paramNames[i]] = match[1]
				}
			}

			err = route.handler.serveContext(c)
			break
		}
	}

	return err
}

func consume(names []string, pattern, pathInfo string) map[string]string {
	params := make(map[string]string, len(names))

	values := regexp.MustCompile(pattern).FindAllStringSubmatch(pathInfo, -1)

	for i, value := range values {
		params[names[i]] = value[1]
	}

	return params
}

var re = regexp.MustCompile(`:(\w+)`)

func getPatternAndParamNames(pattern string) (string, []string) {
	matches := re.FindAllStringSubmatch(pattern, -1)
	paramNames := make([]string, len(matches))

	if len(matches) > 0 {
		pattern = re.ReplaceAllLiteralString(pattern, "([^\\/]+)")

		for i, match := range matches {
			paramNames[i] = match[1]
		}
	}

	if pattern != "/" {
		pattern = "/" + pattern
	}

	return pattern, paramNames
}

func (m *Mux) On(pattern string, handler func(*Mux)) {
	newPattern, paramNames := getPatternAndParamNames(pattern)
	m.routes = append(m.routes, route{"ALL", newPattern, paramNames, MuxHandler(handler)})
}

func (m *Mux) Add(method, pattern string, handler func(*Context) error) {
	newPattern, paramNames := getPatternAndParamNames(pattern)

	m.routes = append(m.routes, route{method, newPattern, paramNames, Handler(handler)})
}

func (m *Mux) Get(pattern string, handler Handler) {
	m.Add("GET", pattern, handler)
}

func (m *Mux) Post(pattern string, handler Handler) {
	m.Add("POST", pattern, handler)
}

func (m *Mux) Put(pattern string, handler Handler) {
	m.Add("PUT", pattern, handler)
}

func (m *Mux) Delete(pattern string, handler Handler) {
	m.Add("DELETE", pattern, handler)
}

func (m Mux) Routes() []route {
	return m.routes
}
