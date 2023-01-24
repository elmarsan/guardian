package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elmarsan/guardian/files"
)

// GetDownloadFile represents GET HTTP method handler for serving files.
type GetDownloadFile struct {
	l       *log.Logger
	Path    string
	storage files.Storage
}

// NewGetDownloadFile returns DownloadFile http handler.
func NewGetDownloadFile(l *log.Logger, path string, storage files.Storage) *GetDownloadFile {
	return &GetDownloadFile{
		l:       l,
		Path:    path,
		storage: storage,
	}
}

// ServeHTTP handles file download.
func (h *GetDownloadFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract file path from url
	fpath := "./" + r.URL.Path[len("/files/download/"):]

	h.l.Printf("%s - Downloading %s", h.Path, fpath)

	finfo, err := h.storage.Write(fpath, w)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers and return file
	w.Header().Set("Content-Disposition", "attachment; filename="+finfo.Name)
	w.Header().Set("Content-Type", finfo.Mime)
	w.Header().Set("Content-Length", strconv.FormatInt(finfo.Size, 10))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("ETag", finfo.Name)
}
