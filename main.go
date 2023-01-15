package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/elmarsan/guardian/handlers"
	"github.com/elmarsan/guardian/middlewares"
	"github.com/elmarsan/guardian/repository"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	checkEnv()

	userRepo, err := repository.NewSqliteUserRepository()
	if err != nil {
		log.Fatal(err)
	}

	postLoginHandler := handlers.NewPostLogin(userRepo)
	getLoginHandler := handlers.NewGetLogin()

	r := mux.NewRouter()
	r.HandleFunc("/login", getLoginHandler.ServeHTTP).Methods("GET")
	r.HandleFunc("/login", postLoginHandler.ServeHTTP).Methods("POST")
	r.PathPrefix("/").Handler(middlewares.Auth(http.FileServer(http.Dir("./static"))))

	log.Print("Guardian listening on :8080...")
	s := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// start server
	go func() {
		err = s.ListenAndServe()
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
	_, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		log.Fatal("Missing JWT_KEY env")
	}

	_, ok = os.LookupEnv("DATABASE_URL")
	if !ok {
		log.Fatal("Missing DATABASE_URL env")
	}
}
