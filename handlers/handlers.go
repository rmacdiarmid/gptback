package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("templates", "*"))
	if err != nil {
		log.Fatalf("Error parsing templates: %s", err)
	}

	log.Println("Loaded templates:")
	for _, t := range templates.Templates() {
		log.Printf("- %s", t.Name())
	}
}

func renderTemplate(w http.ResponseWriter, templateName string) {
	t := templates.Lookup(templateName)
	if t == nil {
		log.Printf("Template not found: %s", templateName)
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering template: %s", templateName)
	err := t.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func RenderTemplateWithData(w http.ResponseWriter, templateName string, data interface{}) {
	t := templates.Lookup(templateName)
	if t == nil {
		log.Printf("Template not found: %s", templateName)
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	log.Printf("Rendering template: %s", templateName)
	err := t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
