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

		//// Check for exact matches (e.g. "/products/new" == "/products/new")
		//if c.PathInfo == route.pattern {
		//return route.handler.serveContext(c)
		//}

		//if route.pattern == "/" {
		//return nil
		//}

		// Check for captures (e.g. "/products/([^\/]+)" =~ "/products/123")
		re := regexp.MustCompile(`\A\/` + route.pattern + `(\/|\z)`)

		matched := re.MatchString(c.PathInfo)
		fmt.Println(re, c.PathInfo, matched)

		if matched {
			matchData := re.FindAllStringSubmatch(c.PathInfo, -1)[0]
			path := matchData[0]

			var captures []string
			if len(matchData) > 2 {
				// There are captures
				captures = matchData[1 : len(matchData)-1]

				for i, name := range route.paramNames {
					c.Params[name] = captures[i]
				}
			}

			fmt.Println("BEFORE:", c.PathInfo)
			fmt.Println("MATCHDATA:", matchData)
			c.PathInfo = matchData[len(matchData)-1] + c.PathInfo[len(path):]
			fmt.Println("AFTER:", c.PathInfo)

			return route.handler.serveContext(c)
		}

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
