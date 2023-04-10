package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/rmacdiarmid/GPTSite/logger"
)

var templates *template.Template

func render(tmplName string, data interface{}) (template.HTML, error) {
	var buf bytes.Buffer
	t := templates.Lookup(tmplName)
	if t == nil {
		return "", fmt.Errorf("template not found: %s", tmplName)
	}
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func LoadTemplates() {
	logger.DualLog.Println("Starting LoadTemplates function...")
	defer logger.DualLog.Println("Exiting LoadTemplates function.")

	// Create a new template
	templates = template.New("").Funcs(template.FuncMap{
		"render": render,
	})

	// Read all files in the templates directory
	templateFiles, err := ioutil.ReadDir("templates")
	if err != nil {
		logger.DualLog.Printf("Error reading templates directory: %v", err)
		return
	}

	// Iterate through the template files and add them to the templates object
	for _, file := range templateFiles {
		filePath := filepath.Join("templates", file.Name())
		_, err := templates.ParseFiles(filePath)
		if err != nil {
			logger.DualLog.Printf("Error parsing template file %s: %v", file.Name(), err)
			return
		}
		logger.DualLog.Printf("- %s", file.Name())
	}
}

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

func ExecTemplate(templateName string, data interface{}) (string, error) {
	var buf bytes.Buffer
	t := templates.Lookup(templateName)
	if t == nil {
		return "", fmt.Errorf("template not found: %s", templateName)
	}
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func RenderTemplateWithData(w http.ResponseWriter, templateName string, data interface{}) {
	t := templates.Lookup("base.gohtml")
	if t == nil {
		logger.DualLog.Printf("Template not found: base.gohtml")
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	logger.DualLog.Printf("Rendering template: base.gohtml with content: %s", templateName)
	err := t.Execute(w, struct {
		Template string
		Data     interface{}
	}{
		Template: templateName,
		Data:     data,
	})

	if err != nil {
		logger.DualLog.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
