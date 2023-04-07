package handlers

import (
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the tasks data
	data := struct {
		Tasks []database.Task
	}{
		Tasks: tasks,
	}

	renderTemplateWithData(w, "about.html", data)
}

func renderTemplateWithData(w http.ResponseWriter, tmpl string, data interface{}) {
	t := templates.Lookup(tmpl)
	if t == nil {
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	err := t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
