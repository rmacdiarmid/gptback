package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
)

func init() {
	// Initialize test logger
	initTestLogger()

	// Initialize test database
	testDB, err := database.InitDB(":memory:")
	if err != nil {
		logger.DualLog.Fatalf("Failed to initialize test database: %v", err)
	}

	// Replace the global DB variable with the test database
	database.DB = testDB

	// Create the "tasks" table in the test database
	createTableSQL := `
	CREATE TABLE tasks (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = testDB.Exec(createTableSQL)
	if err != nil {
		logger.DualLog.Fatalf("Failed to create tasks table in test database: %v", err)
	}

	// Replace the global DB variable with the test database
	database.DB = testDB
}

func TestCreateTaskHandler(t *testing.T) {
	// Prepare the request payload
	task := map[string]interface{}{
		"title":       "Test Task",
		"description": "This is a test task for integration testing.",
	}
	payload, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to prepare request payload: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the CreateTaskHandler with the request and ResponseRecorder
	CreateTaskHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("CreateTaskHandler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Check the response body
	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if resp["Title"] != task["title"] || resp["Description"] != task["description"] {
		t.Errorf("CreateTaskHandler returned incorrect task: got %v, want %v", resp, task)
	}
}

func TestReadTaskHandler(t *testing.T) {
	// Create a task to be read
	task := database.Task{
		Title:       "Test Task",
		Description: "This is a test task for integration testing.",
	}
	taskID, err := database.CreateTask(task.Title, task.Description)
	if err != nil {
		t.Fatalf("Failed to create task for testing: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("GET", fmt.Sprintf("/tasks/%d", taskID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the ReadTaskHandler with the request and ResponseRecorder
	ReadTaskHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ReadTaskHandler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	var resp map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	if resp["ID"].(float64) != float64(taskID) || resp["Title"] != task.Title || resp["Description"] != task.Description {
		t.Errorf("ReadTaskHandler returned incorrect task: got %v, want %v", resp, task)
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	taskTitle := "Test Task"
	taskDescription := "This is a test task for integration testing."
	newTaskTitle := "Updated Test Task"
	newTaskDescription := "This is an updated test task for integration testing."

	// Create a task to update
	taskID, err := database.CreateTask(taskTitle, taskDescription)
	if err != nil {
		t.Fatalf("Failed to create task for testing: %v", err)
	}

	// Prepare the request body
	updatedTask := Task{
		ID:          taskID,
		Title:       newTaskTitle,
		Description: newTaskDescription,
	}
	reqBody, err := json.Marshal(updatedTask)
	if err != nil {
		t.Fatalf("Failed to marshal task JSON: %v", err)
	}

	// Create a request to update the task
	req, err := http.NewRequest("PUT", fmt.Sprintf("/tasks/%d", taskID), bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(UpdateTaskHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("UpdateTaskHandler returned wrong status code: got %v, want %v", status, http.StatusNoContent)
	}

	// Read the updated task
	updatedTaskFromDB, err := database.ReadTask(int(taskID))
	if err != nil {
		t.Fatalf("Failed to read updated task: %v", err)
	}

	// Check if the task was updated correctly
	if updatedTaskFromDB.Title != newTaskTitle || updatedTaskFromDB.Description != newTaskDescription {
		t.Errorf("UpdateTaskHandler did not update task correctly: got %+v, want %+v", updatedTaskFromDB, updatedTask)
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	// Create a task to be deleted
	task := database.Task{
		Title:       "Test Task",
		Description: "This is a test task for integration testing.",
	}
	taskID, err := database.CreateTask(task.Title, task.Description)
	if err != nil {
		t.Fatalf("Failed to create task for testing: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", taskID), nil)
	req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the DeleteTaskHandler with the request and ResponseRecorder
	DeleteTaskHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("DeleteTaskHandler returned wrong status code: got %v, want %v", status, http.StatusNoContent)
	}

	// Try to read the deleted task
	_, err = database.ReadTask(int(taskID))
	if err != sql.ErrNoRows {
		t.Errorf("DeleteTaskHandler did not delete task: task with ID %d still exists", taskID)
	}
}
