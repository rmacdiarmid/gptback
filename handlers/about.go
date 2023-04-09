package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the AboutHandler function...")
	defer logger.DualLog.Println("Exiting the AboutHandler function.")

	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error reading all tasks from database: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template with the tasks data
	data := struct {
		Tasks []database.Task
	}{
		Tasks: tasks,
	}

	renderTemplateWithData(w, "about.gohtml", data)
}

func CreateAboutTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the CreateAboutTaskHandler function...")
	defer logger.DualLog.Println("Exiting the CreateAboutTaskHandler function.")

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error reading request body: %v", err))
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Unmarshal the JSON data into a Task struct
	var task database.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error unmarshalling JSON: %v", err))
		http.Error(w, "Error unmarshalling JSON", http.StatusBadRequest)
		return
	}

	// Save the task to the database
	_, err = database.CreateTask(task.Title, task.Description)
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error creating task in database: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.DualLog.Println("Task created successfully")
	w.WriteHeader(http.StatusCreated)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting the GetTasksHandler function...")
	defer logger.DualLog.Println("Exiting the GetTasksHandler function.")

	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error reading all tasks from database: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal the tasks data into JSON
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error marshalling tasks data into JSON: %v", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type header and write the JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	logger.DualLog.Println("Tasks data successfully sent as JSON response")
}

func renderTemplateWithData(w http.ResponseWriter, tmpl string, data interface{}) {
	logger.DualLog.Println("Starting the renderTemplateWithData function...")
	defer logger.DualLog.Println("Exiting the renderTemplateWithData function.")

	tmplObj := templates.Lookup(tmpl)
	if tmplObj == nil {
		logger.DualLog.Println(fmt.Sprintf("Error: template '%v' not found.", tmpl))
		http.Error(w, "Template not found.", http.StatusInternalServerError)
		return
	}

	err := tmplObj.Execute(w, data)
	if err != nil {
		logger.DualLog.Println(fmt.Sprintf("Error executing template '%v': %v", tmpl, err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.DualLog.Println(fmt.Sprintf("Template '%v' executed successfully", tmpl))
}
