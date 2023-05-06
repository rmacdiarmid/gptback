package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/rmacdiarmid/gptback/config"
	"github.com/rmacdiarmid/gptback/graphqlschema"
	"github.com/rmacdiarmid/gptback/internal"
	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
	"github.com/rmacdiarmid/gptback/pkg/storage"
)

var templates *template.Template

func main() {
	var err error
	// Load configuration from the config file
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.DualLog.Fatalf("Error loading config: %v", err)
	}

	//Set Image Base Url
	os.Setenv("IMAGE_BASE_URL", cfg.Image.BaseURL)

	//Set Log Dir
	logDir := cfg.Log.Dir
	logFile := cfg.Log.File

	//Set Storage Dir
	useS3 := cfg.Storage.UseS3

	//Set Database Path
	dbPath := cfg.Database.Path

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

	// Initialize the logger with the custom dual writer
	logger.InitLogger(f)

	// Modify the server starting code inside the main() function

	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	//Initiate GraphQl
	graphqlschema.InitSchema()

	//fileStorage
	var fileStorage storage.FileStorage
	if useS3 {
		region := cfg.Storage.Region
		bucket := cfg.Storage.Bucket
		s3Storage, err := storage.NewS3FileStorage(region, bucket)
		if err != nil {
			logger.DualLog.Fatalf("Failed to initialize S3 file storage: %v", err)
		}
		fileStorage = s3Storage
	} else {
		basePath := cfg.Storage.BasePath
		fileStorage = &storage.LocalFileStorage{BasePath: basePath}
	}

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleFile(fileStorage, w, r)
	})

	// Use logger.DualLog instead of the previously used dualLog variable
	logger.DualLog.Println("Reading the database path from the config...")

	// Read the database path from the config
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

	// GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &graphqlschema.Schema, // Use the new schema here
		Pretty:   true,
		GraphiQL: true,
	})

	// Create the router and add the routes
	r := mux.NewRouter()

	// GraphQL Router
	r.Handle("/graphql", h)

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
	r.HandleFunc("/graphql/get-image-base-url", internal.GetImageBaseURLHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "static/images/favicon.ico") })

	// Start the server
	logger.DualLog.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		logger.DualLog.Fatalf("Error starting server: %s", err)
	}
}
