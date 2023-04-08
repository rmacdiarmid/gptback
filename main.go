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
	database.InitDB("tasks.db")

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
	err := http.ListenAndServe(":8080", r) // Add the colon before the equal sign
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
