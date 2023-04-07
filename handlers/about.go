package handlers

import (
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/rmacdiarmid/GPTSite/database"
)

// Add these imports

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

func CreateAboutTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a Task struct
	var task database.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// Save the task to the database
	_, err = database.CreateTask(task.Title, task.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the tasks data into JSON
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
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
