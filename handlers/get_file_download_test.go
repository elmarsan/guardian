package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetDownloadFile(t *testing.T) {
	l := log.Default()

	t.Run("should download file", func(t *testing.T) {
		handler := NewGetDownloadFile(l, "/files/download/{path}")

		req := httptest.NewRequest("GET", "/files/download/../test_files/file.txt", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}

		if rec.Header().Get("Content-Disposition") != "attachment; filename=file.txt" {
			t.Errorf("Wrong Content-Disposition header")
		}
	})

	t.Run("should return not found when file does not exist", func(t *testing.T) {
		handler := NewGetDownloadFile(l, "/files/download/{path}")

		req := httptest.NewRequest("GET", "/files/download/../test_files/test.txt", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Errorf("StatusCode should be 404")
		}
	})
}
