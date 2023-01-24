package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/elmarsan/guardian/files"
)

// GetFiles represents GET HTTP method handler for serving files.
type GetFiles struct {
	l       *log.Logger
	Path    string
	tmpl    string
	storage files.Storage
}

// NewGetFiles returns GetFiles http handler.
func NewGetFiles(l *log.Logger, path string, tmpl string, storage files.Storage) *GetFiles {
	return &GetFiles{
		l:       l,
		Path:    path,
		tmpl:    tmpl,
		storage: storage,
	}
}

// ServeHTTP handles file serving.
func (h *GetFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - GetFiles", h.Path)

	files, err := h.storage.GetAllInfo()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(h.tmpl))
	tmpl.Execute(w, files)
}
