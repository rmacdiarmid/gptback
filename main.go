package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/internal"
	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
	"github.com/spf13/viper"
)

var templates *template.Template

func init() {

	// Load configuration from the config file
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.AddConfigPath(".")      // Path to look for the config file in
	err := viper.ReadInConfig()   // Find and read the config file
	if err != nil {               // Handle errors reading the config file
		log.Fatal("Error reading config file:", err)
	}

	// Load log configuration
	logDir := viper.GetString("log.dir")
	logFile := viper.GetString("log.file")

	// Append a timestamp to the log file name
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileWithTimestamp := fmt.Sprintf("%s_%s", timestamp, logFile)

	// Create log directory if it doesn't exist
	err = os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating log directory:", err)
	}

	// Create log file
	logPath := filepath.Join(logDir, logFileWithTimestamp)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	//funcMap := template.FuncMap{"ExecTemplate": handlers.ExecTemplate}
	//templates = template.Must(template.New("").Funcs(funcMap).ParseFiles("templates/base.gohtml", "templates/index.gohtml"))

	// Initialize the logger with the custom dual writer
	logger.InitLogger(f)

	// Load templates
	internal.LoadTemplates("templates/*.gohtml")
}

func main() {
	var err error

	// Use logger.DualLog instead of the previously used dualLog variable
	logger.DualLog.Println("Reading the database path from the config...")

	// Read the database path from the config
	dbPath := viper.GetString("database.path")
	logger.DualLog.Println("Read database path from config successfully")

	// Initialize the database
	logger.DualLog.Println("Initializing database...")

	_, err = database.InitDB(dbPath)
	if err != nil {
		logger.DualLog.Fatalf("Failed to initialize database: %v", err)
	}
	logger.DualLog.Println("Database initialized successfully")

	// Load environment variables from .env file
	logger.DualLog.Println("Loading environmental variables...")

	err = internal.LoadEnvFile(".env")
	if err != nil {
		logger.DualLog.Fatalf("Failed to load environment variables from .env file: %s", err)
	}
	logger.DualLog.Println("Environmental variables loaded successfully")

	// Add the logger usage that was removed from the handlers package
	logger.DualLog.Println("Internal handlers package initialized")

	// Create the router and add the routes
	r := mux.NewRouter()

	// Task CRUD handlers
	r.HandleFunc("/tasks", internal.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", internal.ReadTaskHandler).Methods("GET")
	r.HandleFunc("/tasks/{id}", internal.UpdateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", internal.DeleteTaskHandler).Methods("DELETE")

	// Route handlers
	r.HandleFunc("/", internal.IndexHandler)
	r.HandleFunc("/about", internal.AboutHandler)
	r.HandleFunc("/contact", internal.ContactHandler)
	//r.HandleFunc("/activity", handlers.ActivityHandler)
	r.HandleFunc("/task_list", internal.TaskListHandler)
	r.HandleFunc("/success", internal.SuccessHandler)

	// New routes for generating and accepting articles
	r.HandleFunc("/generate-article", internal.GenerateArticleHandler)
	r.HandleFunc("/accept-article", internal.AcceptArticleHandler)
	r.HandleFunc("/article-generator", internal.ArticleGeneratorHandler)

	r.NotFoundHandler = http.HandlerFunc(internal.NotFoundHandler)

	// Static file handling
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "static/images/favicon.ico") })

	// Start the server
	logger.DualLog.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.DualLog.Fatalf("Error starting server: %s", err)
	}
}
