package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/elmarsan/guardian/handlers"
	"github.com/elmarsan/guardian/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	checkEnv()

	userRepo, err := repository.NewSqliteUserRepository()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	base, _ := os.LookupEnv("BASE_PATH")

	postLoginHandler := handlers.NewPostLogin(userRepo)
	getLoginHandler := handlers.NewGetLogin()
	getFilesHandler := handlers.NewServeFiles(base)
	downloadFilesHandler := handlers.NewDownloadFiles(r, base)

	r.HandleFunc("/login", getLoginHandler.ServeHTTP).Methods("GET")
	r.HandleFunc("/login", postLoginHandler.ServeHTTP).Methods("POST")
	r.HandleFunc("/files", getFilesHandler.ServeHTTP).Methods("GET")
	r.HandleFunc("/files/download/{path}", downloadFilesHandler.ServeHTTP).Methods("GET")

	s := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start server
	go func() {
		log.Print("Guardian listening on :8000...")

		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func checkEnv() {
	err := godotenv.Load()
	if err != nil {
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
