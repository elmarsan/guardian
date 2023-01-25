package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
	l := log.Default()

	t.Run("should return http status ok with file.tmpl", func(t *testing.T) {
		handler := NewFiles(l, "/files", "static", "../templates/files.tmpl")

		// Create folder that serves files
		os.Mkdir("static", 0777)

		req := httptest.NewRequest("GET", handler.Path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}
	})

	t.Run("should return http internal server error when static path does not exist", func(t *testing.T) {
		handler := NewFiles(l, "/files", "non-existing-static-path", "../templates/files.tmpl")

		req := httptest.NewRequest("GET", handler.Path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("StatusCode should be 500")
		}
	})
}
