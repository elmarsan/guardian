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

// NewGetFiles returns GetFiles http handler.
func NewGetFiles(l *log.Logger, path string, staticPath string) *GetFiles {
	return &GetFiles{
		l:          l,
		Path:       path,
		staticPath: staticPath,
	}
}

// ServerFile represents file stored in server.
type ServerFile struct {
	Path string
	Name string
}

// ServeHTTP handles file serving.
func (h *GetFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - GetFiles", h.Path)

	files := []ServerFile{}

	err := filepath.Walk(h.staticPath, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !f.IsDir() {
			f := ServerFile{
				Path: path,
				Name: f.Name(),
			}

			files = append(files, f)
		}

		return nil
	})

	if err != nil {
		h.l.Printf("ERROR: Static path not found %s", h.staticPath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("./templates/files.tmpl"))
	tmpl.Execute(w, files)
}
