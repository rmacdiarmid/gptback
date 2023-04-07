package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
)

var DB *sql.DB

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new template and parse the HTML file
	t, err := template.ParseFiles("templates/task_list.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the tasks data
	err = t.Execute(w, tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
