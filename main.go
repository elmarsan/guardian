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
	l := log.New(os.Stdout, "Guardian ", log.LstdFlags)

	checkEnv(l)

	userRepo, err := auth.NewSqliteUserRepository(l)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	base, _ := os.LookupEnv("BASE_PATH")

	storage, err := files.NewLocalStorage(base)
	if err != nil {
		log.Fatal(err)
	}

	// Auth handlers
	postLogin := handlers.NewPostLogin(l, "/login", userRepo)
	getLogin := handlers.NewGetLogin(l, "/login", "./templates/login.tmpl")

	// File handlers
	getFiles := handlers.NewGetFiles(l, "/files", "./templates/files.tmpl", storage)
	getDownloadFile := handlers.NewGetDownloadFile(l, "/files/download/{path}", storage)
	postUploadFile := handlers.NewPostUploadFile(l, "/files/upload/{filename}", storage)

	// Attach handler to router
	r.Handle(getLogin.Path, getLogin).Methods("GET")
	r.Handle(postLogin.Path, postLogin).Methods("POST")
	r.Handle(postUploadFile.Path, postUploadFile).Methods("POST")

	r.Handle(getFiles.Path, middlewares.Auth(getFiles)).Methods("GET")
	r.PathPrefix(getDownloadFile.Path).Handler(middlewares.Auth(getDownloadFile)).Methods("GET")
	r.Handle(postUploadFile.Path, postUploadFile).Methods("POST")

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
	err := godotenv.Load()
	if err != nil {
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
