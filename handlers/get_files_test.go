package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetFiles(t *testing.T) {
	l := log.Default()

	t.Run("should return http status ok with file.tmpl", func(t *testing.T) {
		handler := NewGetFiles(l, "/files", "static")

		// Create folder that serves files
		os.Mkdir("static", 0777)
		// Go to root path for avoid err while loading templates
		os.Chdir("../")

		req := httptest.NewRequest("GET", handler.Path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}

		os.Remove("static")
	})

	t.Run("should return http internal server error when static path does not exist", func(t *testing.T) {
		handler := NewGetFiles(l, "/files", "non-existing-static-path")

		req := httptest.NewRequest("GET", handler.Path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("StatusCode should be 200")
		}
	})
}
