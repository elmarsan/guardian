package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elmarsan/guardian/jwt"
	"github.com/elmarsan/guardian/repository"
)

// PostLogin represents POST HTTP method handler.
type PostLogin struct {
	l    *log.Logger
	Path string
	ur   repository.UserRepository
}

// NewPostLogin returns PostLogin http handler.
func NewPostLogin(l *log.Logger, path string, ur repository.UserRepository) *PostLogin {
	return &PostLogin{
		l:    l,
		Path: path,
		ur:   ur,
	}
}

// ServeHTTP handles login process.
func (h *PostLogin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Printf("%s - PostLogin", h.Path)

	r.ParseForm()

	// Validate form data
	creds, err := NewLoginCredentials(
		r.FormValue("username"),
		r.FormValue("password"),
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Check if user exist
	err = h.ur.ValidateCredentials(creds.Username, creds.Password)
	if err != nil {
		if err.Error() == repository.InvalidCredentialsErr {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
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

	// Once login success redirect to file list page
	fsRedirect(w, r)
	return
}

// LoginCredentials represents data required for login process.
type LoginCredentials struct {
	Password string
	Username string
}

// NewLoginCredentials function builds and validate LoginCredentials struct.
func NewLoginCredentials(username, password string) (*LoginCredentials, error) {
	if len(username) == 0 {
		return nil, fmt.Errorf("Invalid username")
	}

	if len(password) == 0 {
		return nil, fmt.Errorf("Invalid password")
	}

	return &LoginCredentials{
		Username: username,
		Password: password,
	}, nil
}

// fsRedirect function redirects to files page.
func fsRedirect(w http.ResponseWriter, r *http.Request) {
	loginHandler := http.RedirectHandler("/files", http.StatusFound)
	loginHandler.ServeHTTP(w, r)
}
