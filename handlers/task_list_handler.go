package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

var DB *sql.DB

func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	// Log message for starting TaskListHandler function
	logger.DualLog.Println("Starting TaskListHandler function...")

	// Get all tasks from the database
	tasks, err := database.ReadAllTasks()
	if err != nil {
		logger.DualLog.Printf("Error reading all tasks: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.DualLog.Println("All tasks read successfully.")

	// Create a new template and parse the HTML file
	t, err := template.ParseFiles("templates/task_list.html")
	if err != nil {
		logger.DualLog.Printf("Error parsing task_list.html file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.DualLog.Println("task_list.html file parsed successfully.")

	// Execute the template with the tasks data
	err = t.Execute(w, tasks)
	if err != nil {
		logger.DualLog.Printf("Error executing task list template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logger.DualLog.Println("Task list template executed successfully.")

	// Log message for successful completion of TaskListHandler function
	logger.DualLog.Println("TaskListHandler function completed successfully.")
}
