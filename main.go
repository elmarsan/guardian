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
	checkEnv()

	// Create sqlite user repository
	userRepo, err := auth.NewSqliteUserRepository()
	if err != nil {
		log.Fatal(err)
	}

	base, _ := os.LookupEnv("BASE_PATH")

	// Create local storage
	storage := files.NewLocalStorage(base)

	// Router
	r := mux.NewRouter()
	r.Use(middlewares.Log)

	// Auth handlers
	login := handlers.NewLogin(userRepo)
	loginTmpl := handlers.NewLoginTmpl("./templates/login.tmpl")

	r.Handle("/login", loginTmpl).Methods("GET")
	r.Handle("/login", login).Methods("POST")

	// File  handlers
	files := handlers.NewFiles(storage, "./templates/files.tmpl")
	downloadFile := handlers.NewDownloadFile(storage)
	uploadFile := handlers.NewUploadFile(storage)

	r.Handle("/files", middlewares.Auth(files)).Methods("GET")
	r.PathPrefix("/files/download/{path}").Handler(middlewares.Auth(downloadFile)).Methods("GET")
	r.Handle("/files/upload/{filename}", middlewares.Auth(uploadFile)).Methods("POST")

	s := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start server
	go func() {
		log.Println("Listening on :8000...")

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
	log.Println("Got signal:", sig)

	// Gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func checkEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if _, ok := os.LookupEnv("JWT_KEY"); !ok {
		log.Fatal("Missing JWT_KEY env")
	}

	if _, ok := os.LookupEnv("DATABASE_URL"); !ok {
		log.Fatal("Missing DATABASE_URL env")
	}

	if _, ok := os.LookupEnv("BASE_PATH"); !ok {
		log.Fatal("Missing BASE_PATH env")
	}
}
