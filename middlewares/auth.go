package middlewares

import (
	"net/http"

	"github.com/elmarsan/guardian/jwt"
)

// Auth middleware redirects to login in case cookie is missing or jwt is not valid.
func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		// If cookie is missing redirect to login
		if err != nil {
			loginRedirect(w, r)
			return
		}

		// If token is not valid redirect to login
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
