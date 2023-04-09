package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/database"
	"github.com/rmacdiarmid/GPTSite/logger"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting CreateTaskHandler function...")

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		logger.DualLog.Printf("Error decoding JSON request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := database.CreateTask(task.Title, task.Description)
	if err != nil {
		logger.DualLog.Printf("Error creating task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	task.ID = id
	json.NewEncoder(w).Encode(task)

	logger.DualLog.Println("CreateTaskHandler function completed successfully.")
}

func ReadTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting ReadTaskHandler function...")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.DualLog.Printf("Error converting ID to integer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := database.ReadTask(database.DB, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.DualLog.Println("Task not found.")
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			logger.DualLog.Printf("Error reading task: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(task)

	logger.DualLog.Println("ReadTaskHandler function completed successfully.")
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting UpdateTaskHandler function...")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.DualLog.Printf("Error converting ID to integer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		logger.DualLog.Printf("Error decoding JSON request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateTask(database.DB, id, updatedTask.Title, updatedTask.Description)
	if err != nil {
		logger.DualLog.Printf("Error updating task: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	logger.DualLog.Println("UpdateTaskHandler function completed successfully.")
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting DeleteTaskHandler function...")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logger.DualLog.Printf("Error converting ID to integer: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteTask(database.DB, id)
	if err != nil {
		logger.DualLog.Printf("Error deleting task with ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.DualLog.Printf("Task with ID %d deleted successfully.", id)
	w.WriteHeader(http.StatusNoContent)
}
