package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// GetFiles represents GET HTTP method handler for serving files.
type GetFiles struct {
	l          *log.Logger
	Path       string
	staticPath string
}

// NewServeFiles returns GetFiles http handler.
func NewServeFiles(l *log.Logger, path string, staticPath string) *GetFiles {
	return &GetFiles{
		l:          l,
		Path:       path,
		staticPath: staticPath,
	}
}

type ServerFile struct {
	Path string
	Name string
}

// ServeHTTP handles file serving.
func (h *GetFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - GetFiles", h.Path)

	files := []ServerFile{}

	filepath.Walk(h.staticPath, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			f := ServerFile{
				Path: path,
				Name: f.Name(),
			}

			files = append(files, f)
		}

		return nil
	})

	tmpl := template.Must(template.ParseFiles("./templates/files.tmpl"))
	tmpl.Execute(w, files)
}
