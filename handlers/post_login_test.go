package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type userMockRepository struct {
	err error
}

func (ur *userMockRepository) ValidateCredentials(user string, hash string) error {
	return ur.err
}

func TestLoginTest(t *testing.T) {
	l := log.Default()

	t.Run("should login when credentials are valid and redirec to file path", func(t *testing.T) {
		ur := &userMockRepository{
			err: nil,
		}

		handler := NewPostLogin(l, "/login", ur)

		form := url.Values{}
		form.Add("username", "ana")
		form.Add("password", "secret_hash")

		req := httptest.NewRequest("POST", handler.Path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusFound {
			t.Errorf("StatusCode should be 200")
		}
	})

	t.Run("should return not found when wrong credentials", func(t *testing.T) {
		ur := &userMockRepository{
			err: fmt.Errorf("Invalid credentials"),
		}

		handler := NewPostLogin(l, "/login", ur)

		form := url.Values{}
		form.Add("username", "ana")
		form.Add("password", "secret_hash")

		req := httptest.NewRequest("POST", handler.Path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Errorf("StatusCode should be 404")
		}
	})

	t.Run("should return bad request when form is not valid", func(t *testing.T) {
		ur := &userMockRepository{
			err: nil,
		}

		handler := NewPostLogin(l, "/login", ur)

		form := url.Values{}
		form.Add("username", "ana")
		form.Add("password", "")

		req := httptest.NewRequest("POST", handler.Path, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("StatusCode should be 400")
		}
	})
}
