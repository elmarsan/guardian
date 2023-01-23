package handlers

import (
	"net/http"
	"text/template"
)

// GetLogin representS GET HTTP method handler.
type GetLogin struct{}

// NewGetLogin returns GetLogin handler.
func NewGetLogin() *GetLogin {
	return &GetLogin{}
}

// ServeHTTP returns login page.
func (gl *GetLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/login.tmpl"))
	tmpl.Execute(w, "")
}
