package cuba

import (
	"net/http"
	"html/template"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	Params map[string]string
}

func (c Context) Redirect(url string) {
	http.Redirect(c.W, c.R, url, http.StatusFound)
}

func (c Context) Render(view string, locals interface{}) {
	tmpl := template.Must(template.ParseFiles("views/layout.html", "views/"+view+".html"))
	tmpl.ExecuteTemplate(c.W, "layout.html", locals)
}
