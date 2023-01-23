package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/gorilla/mux"
)

// ServeFiles represents GET HTTP method handler for serving files.
type ServeFiles struct {
	path string
}

// NewServeFiles returns GetFiles http handler.
func NewServeFiles(path string) *ServeFiles {
	return &ServeFiles{
		path: path,
	}
}

type ServerFile struct {
	Path string
	Name string
}

// ServeHTTP handles file serving.
func (gf *ServeFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files := []ServerFile{}

	filepath.Walk(gf.path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			re := regexp.MustCompile(`^[^\/]+`)
			fpath := re.ReplaceAllString(path, "")

			f := ServerFile{
				Path: fpath,
				Name: f.Name(),
			}

			files = append(files, f)
		}

		return nil
	})

	tmpl := template.Must(template.ParseFiles("./templates/files.tmpl"))
	tmpl.Execute(w, files)
}

// DownloadFiles represents GET HTTP method handler for serving files.
type DownloadFiles struct {
	mux  *mux.Router
	path string
}

// NewDownloadFiles returns DownloadFiles http handler.
func NewDownloadFiles(mux *mux.Router, path string) *DownloadFiles {
	return &DownloadFiles{
		mux:  mux,
		path: path,
	}
}

// ServeHTTP handles file download.
func (pf *DownloadFiles) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	path, ok := params["path"]
	if !ok {
		w.Write([]byte("Missing file"))
		w.WriteHeader(http.StatusBadRequest)
	}

	fpath := fmt.Sprintf("./%s/%s", pf.path, path)
	f, err := os.Open(fpath)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Length", fmt.Sprint(stat.Size()))
	io.Copy(w, f)
}
