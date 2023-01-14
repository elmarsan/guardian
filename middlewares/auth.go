package middlewares

import (
	"net/http"

	"github.com/elmarsan/guardian/jwt"
)

// Auth middleware checks if jwt token is stored in cookies.
// Checks the validity and redirect to file server root path or login page.
func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil || len(token.Value) == 0 {
			loginRedirect(w, r)
			return
		}

		err = jwt.ValidateToken(token.Value)
		if err != nil {
			loginRedirect(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

// loginRedirect redirects to login page.
func loginRedirect(w http.ResponseWriter, r *http.Request) {
	loginHandler := http.RedirectHandler("/login", http.StatusFound)
	loginHandler.ServeHTTP(w, r)
}
