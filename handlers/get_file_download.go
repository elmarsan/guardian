package handlers

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
)

// GetDownloadFile represents GET HTTP method handler for serving files.
type GetDownloadFile struct {
	l    *log.Logger
	Path string
}

// NewGetDownloadFile returns DownloadFile http handler.
func NewGetDownloadFile(l *log.Logger, path string) *GetDownloadFile {
	return &GetDownloadFile{
		l:    l,
		Path: path,
	}
}

// ServeHTTP handles file download.
func (h *GetDownloadFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract file path from url
	fpath := r.URL.Path[len("/files/download/"):]

	h.l.Printf("%s - Downloading %s", h.Path, fpath)

	file, err := os.Open(fpath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Figure out Content-type header from file extension
	ext := path.Ext(file.Name())
	ct := mime.TypeByExtension(ext)

	// Set headers and return file
	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
	w.Header().Set("Content-Type", ct)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("ETag", fileInfo.Name())
	w.WriteHeader(http.StatusOK)
	io.Copy(w, file)
}
