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
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/rmacdiarmid/GPTSite/graphqlschema"
	"github.com/rmacdiarmid/GPTSite/internal"
	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
	"github.com/rmacdiarmid/GPTSite/pkg/storage"
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

	// Initialize the logger with the custom dual writer
	logger.InitLogger(f)

	// Modify the server starting code inside the main() function

}

func main() {
	var err error
	corsMiddleware := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	//fileStorage
	useS3 := viper.GetBool("storage.useS3")

	var fileStorage storage.FileStorage
	if useS3 {
		region := viper.GetString("storage.region")
		bucket := viper.GetString("storage.bucket")
		s3Storage, err := storage.NewS3FileStorage(region, bucket)
		if err != nil {
			logger.DualLog.Fatalf("Failed to initialize S3 file storage: %v", err)
		}
		fileStorage = s3Storage
	} else {
		basePath := viper.GetString("storage.basePath")
		fileStorage = &storage.LocalFileStorage{BasePath: basePath}
	}

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		internal.HandleFile(fileStorage, w, r)
	})

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

	// GraphQL schema
	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphqlschema.Query,
	})

	// GraphQL handler
	h := handler.New(&handler.Config{
		Schema:   &schema,
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
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "static/images/favicon.ico") })

	// Start the server
	logger.DualLog.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", corsMiddleware(r)); err != nil {
		logger.DualLog.Fatalf("Error starting server: %s", err)
	}
}
