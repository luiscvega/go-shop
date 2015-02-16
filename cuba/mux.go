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

		// Special case for root patterns and root paths
		if route.pattern == "/" && c.PathInfo == "/" {
			return route.handler.serveContext(c)
		}

		// Check for captures (e.g. "/products/([^\/]+)" =~ "/products/123")
		re := regexp.MustCompile(`\A\/` + route.pattern + `(\/|\z)`)

		matched := re.MatchString(c.PathInfo)

		if matched {
			matchData := re.FindAllStringSubmatch(c.PathInfo, -1)[0]
			path := matchData[0]

			fmt.Println("====================================")
			fmt.Println("MATCHED:")
			fmt.Println("      PATH:", c.PathInfo)
			fmt.Println("   PATTERN:", re)

			if len(matchData) > 2 {
				captures := matchData[1 : len(matchData)-1]

				for i, name := range route.paramNames {
					c.Params[name] = captures[i]
				}
			}

			c.PathInfo = "/" + c.PathInfo[len(path):]

			fmt.Println("  NEW PATH:", c.PathInfo)
			fmt.Println("====================================")
			fmt.Println()

			return route.handler.serveContext(c)
		} else {
			fmt.Println("====================================")
			fmt.Println("NO MATCH:")
			fmt.Println("      PATH:", c.PathInfo)
			fmt.Println("   PATTERN:", re)
			fmt.Println("====================================")
			fmt.Println()
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
