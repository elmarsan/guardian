package handlers

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLogin(t *testing.T) {
	l := log.Default()

	t.Run("should return http status ok with login.tmpl", func(t *testing.T) {
		handler := NewGetLogin(l, "/login", "../templates/login.tmpl")

		req := httptest.NewRequest("GET", handler.Path, nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Errorf("StatusCode should be 200")
		}
	})
}
