package cuba

import (
	"net/http"
)

func New() App {
	mux := http.NewServeMux()
	return App{Mux: mux}
}

type App struct {
	Mux *http.ServeMux
}

func (app App) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	app.Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.ParseForm()

			if r.FormValue("_method") == "PUT" {
				r.Method = "PUT"
			}

			if r.FormValue("_method") == "DELETE" {
				r.Method = "DELETE"
			}
		}

		handler(w, r)
	})
}
