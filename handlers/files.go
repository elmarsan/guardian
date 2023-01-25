package handlers

import (
	"net/http"
	"text/template"

	"github.com/elmarsan/guardian/files"
)

// Files represents GET HTTP method handler for serving files.
type Files struct {
	tmpl    string
	storage files.Storage
}

// NewFiles returns Files http handler.
func NewFiles(storage files.Storage, tmpl string) *Files {
	return &Files{
		storage: storage,
		tmpl:    tmpl,
	}
}

// ServeHTTP handles file serving.
func (h *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := h.storage.GetAllInfo()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(h.tmpl))
	tmpl.Execute(w, files)
}
