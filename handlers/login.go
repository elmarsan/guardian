package handlers

import (
	"net/http"
	"text/template"

	"github.com/elmarsan/guardian/jwt"
)

// AuthForm represents data required for login process.
type AuthForm struct {
	Password string
	Username string
}

var users = map[string]string{
	"ana":  "ana",
	"juan": "juan",
}

// PostLogin handles login process.
func PostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	creds := &AuthForm{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	if len(creds.Username) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid username"))
		return
	}

	if len(creds.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid password"))
		return
	}

	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid username or password"))
		return
	}

	token, err := jwt.NewToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})

	fsRedirect(w, r)
	return
}

// GetLogin handler returns login page.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/login.html"))
	tmpl.Execute(w, "")
}

// fsRedirect function redirects to files page.
func fsRedirect(w http.ResponseWriter, r *http.Request) {
	loginHandler := http.RedirectHandler("/", http.StatusFound)
	loginHandler.ServeHTTP(w, r)
}
