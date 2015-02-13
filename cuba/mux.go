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
	for _, route := range m.routes {
		if route.method != "ALL" && c.R.Method != route.method {
			continue
		}

		err := route.try(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (route route) try(c *Context) error {
	origPath := c.PathInfo
	defer func() {
		c.PathInfo = origPath
	}()

	// Check for exact matches (e.g. "/products/new" == "/products/new")
	if c.PathInfo == route.pattern {
		return route.handler.serveContext(c)
	}

	if route.pattern == "/" {
		return nil
	}

	// Check for captures (e.g. "/products/([^\/]+)" =~ "/products/123")
	pattern := `\A\/` + route.pattern + `(\/|\z)`
	matched, _ := regexp.MatchString(pattern, c.PathInfo)

	fmt.Println("PATTERN:", pattern, "PATH:", c.PathInfo, "MATCHED:", matched)

	if matched {
		matchData := regexp.MustCompile(pattern).FindAllStringSubmatch(c.PathInfo, -1)[0]
		captures := matchData[1 : len(matchData)-1]
		//lastMatch := matchData[len(matchData)]

		if len(captures) > 0 {
			for i, capture := range captures {
				c.Params[route.paramNames[i]] = capture
			}
		}

		return route.handler.serveContext(c)
	}

	return nil
}

type route struct {
	method     string
	pattern    string
	paramNames []string
	handler    ContextHandler
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
