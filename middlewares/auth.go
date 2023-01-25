package middlewares

import (
	"fmt"
	"net/http"

	"github.com/elmarsan/guardian/jwt"
)

// Auth middleware checks if jwt token is stored in cookies.
// Checks the validity and redirect to file server root path or login page.
func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginPath := r.RequestURI == "/login"
		fmt.Println(loginPath)

		token, err := r.Cookie("token")
		fmt.Printf("%s\n", err.Error())
		fmt.Printf("token %s\n", token)

		// If cookie is missing and is no login path redirect to login,
		if err != nil && loginPath == false {
			fmt.Println("Not found cookie")
			loginRedirect(w, r)
			return
		}

		err = jwt.ValidateToken(token.Value)
		fmt.Printf("%s %s %s", r.RequestURI, token, err.Error())
		if err != nil {
			fmt.Println("Invalid token")
			// loginRedirect(w, r)
			return
		}

		// fmt.Printf("step 1 %s", token)

		// // If cookie is empty string redirect to login.
		// if len(token.Value) == 0 && !loginPath {
		// 	loginRedirect(w, r)
		// 	return
		// }

		// fmt.Println("step 2")

		// err = jwt.ValidateToken(token.Value)
		// fmt.Printf("%s %s %s", r.RequestURI, token, err.Error())
		// if err != nil && !loginPath {
		// 	loginRedirect(w, r)
		// 	return
		// }

		// fmt.Println("LA MIAUUUUU")

		// handler.ServeHTTP(w, r)
	})
}

// loginRedirect redirects to login page.
func loginRedirect(w http.ResponseWriter, r *http.Request) {
	loginHandler := http.RedirectHandler("/login", http.StatusFound)
	loginHandler.ServeHTTP(w, r)
}

// filesRedirect redirects to files page.
func filesRedirect(w http.ResponseWriter, r *http.Request) {
	loginHandler := http.RedirectHandler("/files/", http.StatusFound)
	loginHandler.ServeHTTP(w, r)
}
