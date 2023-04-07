package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/database"
)

type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDB()
	defer db.Close()

	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := database.CreateTask(db, task.Title, task.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	task.ID = id

	newTask := database.ReadTask(db, int(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

func ReadTaskHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDB()
	defer db.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task := database.ReadTask(db, id)
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDB()
	defer db.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var task Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.ID = database.CreateTask(db, task.Title, task.Description)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database.UpdateTask(db, id, task.Title, task.Description)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	db := database.InitDB()
	defer db.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.DeleteTask(db, id)
	w.WriteHeader(http.StatusNoContent)
}
func ReadTask(db *sql.DB, id int) (*Task, error) {
	// Define the SQL statement to retrieve the task
	query := "SELECT id, title, description FROM tasks WHERE id = ?"

	// Execute the SQL statement with the given id
	row := db.QueryRow(query, id)

	// Create a new task to hold the retrieved data
	var task Task

	// Populate the task with data from the query result
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return nil, err
	}

	// Return the populated task
	return &task, nil
}
