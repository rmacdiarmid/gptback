package internal

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/logger"
)

var templates *template.Template

type TemplateData struct {
	Content template.HTML
	Data    interface{}
}

func init() {
	// Create a FuncMap with the custom printData function
	var funcMap = template.FuncMap{
		"printData": printData,
	}

	// Initialize the global templates variable with the custom function map
	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.gohtml"))
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
	logger.DualLog.Printf("Content template output: %s", contentBuf.String())

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
