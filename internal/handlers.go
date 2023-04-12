package internal

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rmacdiarmid/GPTSite/logger"
)

var templates = template.Must(template.ParseGlob("templates/*.gohtml"))

func LoadTemplates(pattern string) {
	logger.DualLog.Println("Starting LoadTemplates function...")
	defer logger.DualLog.Println("Exiting LoadTemplates function.")

	// Load templates
	var allFiles []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".gohtml") {
			allFiles = append(allFiles, path)
		}
		return nil
	})
	if err != nil {
		logger.DualLog.Printf("Error walking the path %q: %v\n", pattern, err)
		return
	}

	templates, err = template.ParseFiles(allFiles...)
	if err != nil {
		logger.DualLog.Printf("Error loading templates: %v\n", err)
		return
	}

	// Log the loaded template names
	for _, tmpl := range templates.Templates() {
		logger.DualLog.Printf("- %s", tmpl.Name())
	}
}

func RenderTemplateWithData(w http.ResponseWriter, tmpl string, data interface{}) {
	logger.DualLog.Println("Starting RenderTemplateWithData function...")
	defer logger.DualLog.Println("Exiting RenderTemplateWithData function.")

	logger.DualLog.Printf("Rendering template: %s", tmpl)

	// Log the loaded template names
	logger.DualLog.Println("Global templates variable contains the following templates:")
	for _, t := range templates.Templates() {
		logger.DualLog.Printf("- %s", t.Name())
	}

	// Render the template
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		logger.DualLog.Printf("Error executing template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
