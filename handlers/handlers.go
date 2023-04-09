package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/rmacdiarmid/GPTSite/logger"
)

var templates *template.Template

func parsePartialTemplates() *template.Template {
	tmpl := template.New("").Funcs(template.FuncMap{})
	partialDirs := []string{"./templates/partials"}
	for _, dir := range partialDirs {
		partials, err := filepath.Glob(filepath.Join(dir, "*.gohtml"))
		if err != nil {
			logger.DualLog.Fatalf("Error reading partial templates directory: %v", err)
		}

		for _, partial := range partials {
			_, err := tmpl.ParseFiles(partial)
			if err != nil {
				logger.DualLog.Printf("Error parsing partial template: %v", err)
			}
		}
	}
	return tmpl
}

func LoadTemplates() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("templates", "*.gohtml"))
	if err != nil {
		logger.DualLog.Fatalf("Error parsing templates: %s", err)
	}

	logger.DualLog.Println("Loaded templates:")
	for _, t := range templates.Templates() {
		logger.DualLog.Printf("- %s", t.Name())
	}
}

// ...

func RenderTemplate(w http.ResponseWriter, templateName string) error {
	t := templates.Lookup(templateName)
	if t == nil {
		logger.DualLog.Printf("Template not found: %s", templateName)
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return fmt.Errorf("Template not found: %s", templateName)
	}

	logger.DualLog.Printf("Rendering template: %s", templateName)
	err := t.Execute(w, nil)
	if err != nil {
		logger.DualLog.Printf("Error executing template: %v", err)
		return err
	}

	return nil
}

func RenderTemplateWithData(w http.ResponseWriter, templateName string, data interface{}) {
	t := templates.Lookup(templateName)
	if t == nil {
		logger.DualLog.Printf("Template not found: %s", templateName)
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	logger.DualLog.Printf("Rendering template: %s", templateName)
	err := t.Execute(w, data)
	if err != nil {
		logger.DualLog.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
