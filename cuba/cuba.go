package cuba

import (
	"html/template"
	"net/http"
	"regexp"
)

func New() mux {
	table := make(map[string][]route)
	return mux{table}
}

type route struct {
	pattern  string
	captures []string
	handler  func(*Context)
}

type mux struct {
	Table map[string][]route
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

	routes, ok := m.Table[r.Method]
	if !ok {
		return
	}

	for _, route := range routes {
		if route.pattern == r.URL.Path {
			route.handler(context)
			return
		}

		re := regexp.MustCompile(route.pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			for i, name := range route.captures {
				context.Params[name] = matches[0][i+1]
			}

			route.handler(context)
			return
		}
	}
}

func (m *mux) Add(method, path string, handler func(*Context)) {
	// Initialize method
	_, ok := m.Table[method]
	if !ok {
		m.Table[method] = make([]route, 0)
	}

	re := regexp.MustCompile(":([a-zA-Z0-9]+)")
	matches := re.FindAllStringSubmatch(path, -1)

	for _, match := range matches {
		captures := make([]string, 0, len(matches))
		captures = append(captures, match[1])

		pattern := re.ReplaceAllLiteralString(path, "([a-zA-Z0-9]+)")

		m.Table[method] = append(m.Table[method], route{pattern, captures, handler})

		return
	}

	m.Table[method] = append(m.Table[method], route{path, nil, handler})
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
