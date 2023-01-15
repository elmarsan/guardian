package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/elmarsan/guardian/jwt"
)

type MockHandler struct{}

func (mh *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuth(t *testing.T) {
	handler := Auth(&MockHandler{})

	t.Run("should redirect when missing cookie token", func(t *testing.T) {
		os.Unsetenv("JWT_EXPIRATION_TIME")

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusFound {
			t.Errorf("Auth should redirect")
		}
	})

	t.Run("should redirect when cookie token has expired or is not valid", func(t *testing.T) {
		os.Setenv("JWT_EXPIRATION_TIME", "0")

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		token, _ := jwt.NewToken()
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})

		handler.ServeHTTP(rec, req)
		if rec.Result().StatusCode != http.StatusFound {
			t.Errorf("Auth should redirect")
		}
	})

	t.Run("should NOT redirect when cookie token is valid", func(t *testing.T) {
		os.Unsetenv("JWT_EXPIRATION_TIME")

		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()

		token, _ := jwt.NewToken()
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})

		handler.ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Error("Auth should redirect when token is valid")
		}
	})
}
