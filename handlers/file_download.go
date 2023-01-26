package handlers

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/elmarsan/guardian/files"
)

// DownloadFile represents GET HTTP method handler for serving files.
type DownloadFile struct {
	storage files.Storage
}

// NewDownloadFile returns DownloadFile http handler.
func NewDownloadFile(storage files.Storage) *DownloadFile {
	return &DownloadFile{
		storage: storage,
	}
}

// ServeHTTP handles file download.
func (h *DownloadFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract file path from url
	fpath := r.URL.Path[len("/files/download/"):]

	b := bytes.NewBuffer([]byte{})
	finfo, err := h.storage.Write(fpath, b)
	if err != nil {
		if err.Error() == files.FileNotFoundErr {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set headers and return file
	w.Header().Set("Content-Disposition", "attachment; filename="+finfo.Name)
	w.Header().Set("Content-Type", finfo.Mime)
	w.Header().Set("Content-Length", strconv.FormatInt(finfo.Size, 10))
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("ETag", finfo.Name)
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, b)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
