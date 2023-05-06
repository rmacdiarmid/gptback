package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
	"github.com/rmacdiarmid/gptback/pkg/storage"
)

var templates *template.Template

type TemplateData struct {
	Content template.HTML
	Data    interface{}
}

func init() {
	// Get the templates path from the environment variable
	templatesPath := os.Getenv("TEMPLATES_PATH")

	if templatesPath == "" {
		// If the environment variable is not set, use the default path
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Error getting the current working directory:", err)
		}
		// Check if the templates folder is in the current working directory
		if _, err := os.Stat(filepath.Join(cwd, "templates")); err == nil {
			templatesPath = filepath.Join(cwd, "templates/*.gohtml")
		} else {
			templatesPath = filepath.Join(cwd, "../templates/*.gohtml")
		}
	}

	// Create a FuncMap with the custom printData function
	var funcMap = template.FuncMap{
		"printData": printData,
	}

	// Initialize the global templates variable with the custom function map
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob(templatesPath))
}

func printData(data interface{}) string {
	return fmt.Sprintf("%+v", data)
}

func RenderTemplateWithData(w http.ResponseWriter, tmpl string, contentTemplateName string, data interface{}) {
	logger.DualLog.Println("Starting RenderTemplateWithData function...")
	defer logger.DualLog.Println("Exiting RenderTemplateWithData function.")

	logger.DualLog.Printf("Rendering template: %s", tmpl)

	// Log the loaded template names
	logger.DualLog.Println("Global templates variable contains the following templates:")
	for _, t := range templates.Templates() {
		logger.DualLog.Printf("- %s", t.Name())
	}

	// Execute the content template and store the output in a buffer
	var contentBuf bytes.Buffer
	err := templates.ExecuteTemplate(&contentBuf, contentTemplateName, data)
	if err != nil {
		logger.DualLog.Printf("Error executing content template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Uncomment the line when troubleshooting
	//logger.DualLog.Printf("Content template output: %s", contentBuf.String())

	// Create a TemplateData instance with the content template's output as a string
	templateData := TemplateData{
		Content: template.HTML(contentBuf.String()),
		Data:    data,
	}

	// Render the base template with the content template's output
	err = templates.ExecuteTemplate(w, tmpl, templateData)
	if err != nil {
		logger.DualLog.Printf("Error executing base template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// internal/handlers.go

func ArticlesHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("ArticlesHandler called")
	defer logger.DualLog.Println("ArticlesHandler exited")

	articles, err := database.GetArticles()
	if err != nil {
		logger.DualLog.Printf("Error fetching articles: %s", err.Error()) // Log the error with DualLog
		http.Error(w, "Error fetching articles", http.StatusInternalServerError)
		return
	}
	logger.DualLog.Printf("Fetched articles: %v", articles)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(articles)
}

//handle File storage
func HandleFile(fileStorage storage.FileStorage, w http.ResponseWriter, r *http.Request) {
	filePath := strings.TrimPrefix(r.URL.Path, "/static")
	file, err := fileStorage.GetFile(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	http.ServeContent(w, r, filePath, time.Now(), file)
}

func GetImageBaseURLHandler(w http.ResponseWriter, r *http.Request) {
	imageBaseURL := os.Getenv("IMAGE_BASE_URL")
	if imageBaseURL == "" {
		http.Error(w, "Image base URL not found", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, imageBaseURL)
}

