package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// GetFiles represents GET HTTP method handler for serving files.
type GetFiles struct {
	path string
}

// NewGetFiles returns GetFiles http handler.
func NewGetFiles(path string) *GetFiles {
	return &GetFiles{
		path: path,
	}
}

// ServeHTTP handles file serving.
func (gf *GetFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files := []string{}

	filepath.Walk(gf.path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			files = append(files, f.Name())
		}

		return nil
	})

	tmpl := template.Must(template.ParseFiles("./templates/files.tmpl"))
	tmpl.Execute(w, files)
}
