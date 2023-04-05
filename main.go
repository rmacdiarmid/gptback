package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/handlers"
)

func main() {
	r := mux.NewRouter()

	// Route handlers
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
	r.HandleFunc("/contact", handlers.ContactHandler)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// Static file handling
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Starting server on :8080...")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
