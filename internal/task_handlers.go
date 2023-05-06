package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/gptback/logger"
	"github.com/rmacdiarmid/gptback/pkg/database"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	logger.DualLog.Println("Starting CreateTaskHandler function...")

	var task database.Task
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
	w.WriteHeader(http.StatusCreated)
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

	task, err := database.ReadTask(id)
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

func UpdateTaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Starting UpdateTaskHandler function...")

	id, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		log.Printf("Error converting ID to integer: %v", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task Task
	err = json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateTask(database.DB, id, task.Title, task.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Println("UpdateTaskHandler function completed successfully.")
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
