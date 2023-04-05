package handlers

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	t := templates.Lookup(tmpl)
	if t == nil {
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
