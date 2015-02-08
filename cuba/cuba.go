package cuba

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

func New() mux {
	var (
		patterns = make([]string, 0)
		handlers = make(map[string]map[string]route)
	)

	return mux{patterns, handlers}
}

type route struct {
	captures []string
	handler  func(*Context)
}

type mux struct {
	patterns []string
	table    map[string]map[string]route
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{w, r, make(map[string]string)}

	if r.Method == "POST" {
		r.ParseForm()

		if r.FormValue("_method") == "PUT" {
			r.Method = "PUT"
		}

		if r.FormValue("_method") == "DELETE" {
			r.Method = "DELETE"
		}
	}

	routes, ok := m.table[r.Method]
	if !ok {
		return
	}

	for _, pattern := range m.patterns {
		if pattern == r.URL.Path {
			routes[r.URL.Path].handler(context)
			return
		}

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			route := routes[pattern]

			for i, name := range route.captures {
				context.Params[name] = matches[0][i+1]
			}

			route.handler(context)
			return
		}
	}
}

func (m *mux) Add(method, pattern string, handler func(*Context)) {
	re := regexp.MustCompile(":[a-zA-Z0-9]+")

	captures := make([]string, 0)
	for _, match := range re.FindAllStringSubmatch(pattern, -1) {
		captures = append(captures, strings.Replace(match[0], ":", "", 1))
	}

	pattern = re.ReplaceAllLiteralString(pattern, "([a-zA-Z0-9]+)")

	m.patterns = append(m.patterns, pattern)

	_, ok := m.table[method]
	if !ok {
		m.table[method] = make(map[string]route)
	}

	m.table[method][pattern] = route{captures, handler}
}

func (m *mux) Get(pattern string, handler func(*Context)) {
	m.Add("GET", pattern, handler)
}

func (m *mux) Post(pattern string, handler func(*Context)) {
	m.Add("POST", pattern, handler)
}

func (m *mux) Put(pattern string, handler func(*Context)) {
	m.Add("PUT", pattern, handler)
}

func (m *mux) Delete(pattern string, handler func(*Context)) {
	m.Add("DELETE", pattern, handler)
}

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string
}

func (c Context) Redirect(url string) {
	http.Redirect(c.W, c.R, url, http.StatusFound)
}

func (c Context) Render(view string, locals interface{}) {
	tmpl := template.Must(template.ParseFiles("views/" + view + ".html"))
	tmpl.ExecuteTemplate(c.W, view+".html", locals)
}
