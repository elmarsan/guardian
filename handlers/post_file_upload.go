package handlers

import (
	"log"
	"net/http"

	"github.com/elmarsan/guardian/files"
	"github.com/gorilla/mux"
)

// PostUploadFile represents POST HTTP method handler for uploading files.
type PostUploadFile struct {
	l       *log.Logger
	Path    string
	storage files.Storage
}

// NewPostUploadFile returns PostUploadFile http handler.
func NewPostUploadFile(l *log.Logger, path string, storage files.Storage) *PostUploadFile {
	return &PostUploadFile{
		l:       l,
		Path:    path,
		storage: storage,
	}
}

// ServeHTTP handles file download.
func (h *PostUploadFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - UploadFile", h.Path)

	vars := mux.Vars(r)
	fp := vars["filename"]

	err := h.storage.Save(fp, r.Body)
	if err != nil {
		h.l.Printf("Unable to save file, error: %s", err.Error())
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}
