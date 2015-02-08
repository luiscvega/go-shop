package cuba

import (
	"html/template"
	"net/http"
	"regexp"
)

func New() mux {
	return mux{make(map[string][]route), nil}
}

type Handler func(*Context)

type route struct {
	pattern  string
	captures []string
	handler  Handler
}

type mux struct {
	table map[string][]route
	final http.Handler
}

func (m *mux) Use(mw func(http.Handler) http.Handler) {
	if m.final == nil {
		m.final = http.HandlerFunc(m.run)
	}

	m.final = mw(m.final)
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.final.ServeHTTP(w, r)
}

func (m mux) run(w http.ResponseWriter, r *http.Request) {
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

	var handler Handler
	for _, route := range routes {
		if route.pattern == r.URL.Path {
			handler = route.handler
			break
		}

		re := regexp.MustCompile(route.pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(matches) > 0 {
			for i, name := range route.captures {
				context.Params[name] = matches[0][i+1]
			}

			handler = route.handler
			break
		}
	}

	if handler != nil {
		handler(context)
	}
}

func (m *mux) Add(method, pattern string, handler func(*Context)) {
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

	m.table[method] = append(m.table[method], route{pattern, captures, handler})
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

func (m mux) Table() map[string][]route {
	return m.table
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
