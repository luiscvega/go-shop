package cuba

import (
	"net/http"
	"regexp"
)

func New() mux {
	var patterns []string
	handlers := make(map[string]cubaHandler)

	return mux{patterns, handlers}
}

type cubaHandler struct {
	names []string
	http.Handler
}

type mux struct {
	patterns []string
	handlers map[string]cubaHandler
}

func (m mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}

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
			m.handlers[r.URL.Path].ServeHTTP(w, r)
			return
		}

		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(r.URL.Path, -1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			r.Form["params"] = append(r.Form["params"], matches[0][1])
			m.handlers[pattern].ServeHTTP(w, r)
			return
		}
	}
}

func (m *mux) Add(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	pattern, names := prepareHandler(pattern)
	m.patterns = append(m.patterns, pattern)
	m.handlers[pattern] = cubaHandler{names, http.HandlerFunc(handler)}
}

func prepareHandler(pattern string) (string, []string) {
	re := regexp.MustCompile(":[a-zA-Z0-9]+")
	pattern = re.ReplaceAllLiteralString(pattern, "([a-zA-Z0-9]+)")
	return pattern, []string{}
}
