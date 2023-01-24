package handlers

import (
	"log"
	"net/http"
	"text/template"
)

// GetLogin representS GET HTTP method handler.
type GetLogin struct {
	l    *log.Logger
	Path string
	tmpl string
}

// NewGetLogin returns GetLogin handler.
func NewGetLogin(l *log.Logger, path string, tmpl string) *GetLogin {
	return &GetLogin{
		l:    l,
		Path: path,
		tmpl: tmpl,
	}
}

// ServeHTTP returns login page.
func (h *GetLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - GetLogin", h.Path)

	tmpl := template.Must(template.ParseFiles(h.tmpl))
	tmpl.Execute(w, "")
}
