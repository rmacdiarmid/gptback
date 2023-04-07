package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/handlers"
)

func main() {
	// Initialize the database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	// Close the database when the program exits
	defer db.Close()

	// Assign the database to the handlers package variable
	handlers.DB = db

	// Create the router and add the routes
	r := mux.NewRouter()

	// Task CRUD handlers
	r.HandleFunc("/tasks", handlers.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", handlers.ReadTaskHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", handlers.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", handlers.DeleteTaskHandler).Methods("DELETE")

	// Route handlers
	r.HandleFunc("/", handlers.IndexHandler)
	r.HandleFunc("/about", handlers.AboutHandler)
	r.HandleFunc("/contact", handlers.ContactHandler)
	r.HandleFunc("/activity", handlers.ActivityHandler)
	r.HandleFunc("/task_list", handlers.TaskListHandler) // <-- Add this line
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// Static file handling
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server
	fmt.Println("Starting server on :8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
