package handlers

import (
	"net/http"
	"text/template"
)

// LoginTmpl represents GET HTTP method handler.
type LoginTmpl struct {
	tmpl string
}

// NewLoginTmpl returns GetLogin handler.
func NewLoginTmpl(tmpl string) *LoginTmpl {
	return &LoginTmpl{
		tmpl: tmpl,
	}
}

// ServeHTTP returns login page.
func (h *LoginTmpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(h.tmpl))
	tmpl.Execute(w, "")
}
