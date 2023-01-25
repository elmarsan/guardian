package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/elmarsan/guardian/auth"
	"github.com/elmarsan/guardian/files"
	"github.com/elmarsan/guardian/handlers"
	"github.com/elmarsan/guardian/middlewares"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	l := log.Default()

	checkEnv(l)

	// Create sqlite user repository
	userRepo, err := auth.NewSqliteUserRepository(l)
	if err != nil {
		log.Fatal(err)
	}

	base, _ := os.LookupEnv("BASE_PATH")

	// Create local storage
	storage, err := files.NewLocalStorage(base)
	if err != nil {
		log.Fatal(err)
	}

	// Router
	r := mux.NewRouter()
	r.Use(middlewares.Log, middlewares.Auth)

	// Auth handlers
	login := handlers.NewLogin(userRepo)
	loginTmpl := handlers.NewLoginTmpl("./templates/login.tmpl")

	r.Handle("/login", loginTmpl).Methods("GET")
	r.Handle("/login", login).Methods("POST")

	// File  handlers
	files := handlers.NewFiles(storage, "./templates/files.tmpl")
	downloadFile := handlers.NewDownloadFile(storage)
	uploadFile := handlers.NewUploadFile(storage)

	r.Handle("/files", files).Methods("GET")
	r.PathPrefix("/files/download/{path}").Handler(downloadFile).Methods("GET")
	r.Handle("/files/upload/{filename}", uploadFile).Methods("POST")

	s := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		l.Println("Listening on :8000...")

		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Println("Got signal:", sig)

	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func checkEnv(l *log.Logger) {
	if err := godotenv.Load(); err != nil {
		l.Fatal("Error loading .env file")
	}

	if _, ok := os.LookupEnv("JWT_KEY"); !ok {
		l.Fatal("Missing JWT_KEY env")
	}

	if _, ok := os.LookupEnv("DATABASE_URL"); !ok {
		l.Fatal("Missing DATABASE_URL env")
	}

	if _, ok := os.LookupEnv("BASE_PATH"); !ok {
		l.Fatal("Missing BASE_PATH env")
	}
}
