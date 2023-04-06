package handlers

import (
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

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := database.CreateTask(db, task.Title, task.Description)
	task.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
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

	task.ID = int64(id)
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
