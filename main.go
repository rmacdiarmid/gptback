package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/handlers"
	"github.com/spf13/viper"
)

func main() {
	// Load configuration from the config file
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.AddConfigPath(".")      // Path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatal("Error reading config file:", err)
	}

	// Read the database path from the config
	dbPath := viper.GetString("database.path")

	// Initialize the database
	_, err = database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Load log configuration
	logDir := viper.GetString("log.dir")
	logFile := viper.GetString("log.file")

	// Create log directory if it doesn't exist
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating log directory:", err)
	}

	// Create log file
	logPath := filepath.Join(logDir, logFile)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}
	defer f.Close()

	// Set log output to file
	log.SetOutput(f)

	// Write log message
	log.Println("Hello, world!")

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
	r.HandleFunc("/task_list", handlers.TaskListHandler)
	r.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

	// Static file handling
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Start the server
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
