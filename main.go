package main

import (
	"log"
	"net/http"
	"time"

	"github.com/elmarsan/guardian/handlers"
	"github.com/elmarsan/guardian/middlewares"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	r.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	r.PathPrefix("/").Handler(middlewares.Auth(http.FileServer(http.Dir("./static"))))

	log.Print("Guardian listening on :8080...")
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
