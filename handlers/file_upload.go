package handlers

import (
	"net/http"

	"github.com/elmarsan/guardian/files"
	"github.com/gorilla/mux"
)

// UploadFile represents POST HTTP method handler for uploading files.
type UploadFile struct {
	storage files.Storage
}

// NewUploadFile returns UploadFile http handler.
func NewUploadFile(storage files.Storage) *UploadFile {
	return &UploadFile{
		storage: storage,
	}
}

// ServeHTTP handles file download.
func (h *UploadFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fp := vars["filename"]

	err := h.storage.Save(fp, r.Body)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
}
