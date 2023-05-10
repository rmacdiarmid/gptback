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
	//"github.com/rmacdiarmid/GPTSite/logger"
	"github.com/rmacdiarmid/GPTSite/pkg/database"
	"github.com/rmacdiarmid/gptback/logger"
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

func TestArticlesHandler(t *testing.T) {
	// Create an article to be retrieved
	article := database.Article{
		Title:   "Test Article",
		Image:   "test_image.jpg",
		Preview: "This is a test article for integration testing.",
		Text:    "This is the main text of the article it should be longer than the preview but I really don't want to type much more",
	}

	_, err := database.CreateArticle(article.Title, article.Image, article.Preview, article.Text)
	if err != nil {
		t.Fatalf("Failed to create article for testing: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the ArticlesHandler with the request and ResponseRecorder
	ArticlesHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ArticlesHandler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the response body
	var resp []map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// Find the article in the response
	found := false
	for _, a := range resp {
		if a["Title"] == article.Title && a["Image"] == article.Image && a["Preview"] == article.Preview && a["Preview"] == article.Text {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("ArticlesHandler did not return the test article: %v", article)
	}
}

func TestPrintData(t *testing.T) {
	// Define a test case
	testData := map[string]interface{}{
		"foo": 123,
		"bar": "hello",
	}

	// Call the printData function
	result := printData(testData)

	// Check that the result matches the expected output
	expected := "map[bar:hello foo:123]"
	if result != expected {
		t.Errorf("printData returned wrong result: got %v, want %v", result, expected)
	}
}

func TestRenderTemplateWithData(t *testing.T) {
	// Create a mock response writer
	recorder := httptest.NewRecorder()

	// Define a test data struct
	type Article struct {
		Image   string
		Title   string
		Preview string
		Text    string
	}

	testData := struct {
		Articles []Article
	}{
		Articles: []Article{
			{
				Image:   "https://example.com/image1.jpg",
				Title:   "Article 1",
				Preview: "This is a preview of article 1.",
				Text:    "This is a fake text about a fake thing, it's the main body",
			},
			{
				Image:   "https://example.com/image2.jpg",
				Title:   "Article 2",
				Preview: "This is a preview of article 2.",
				Text:    "This is a fake text about a fake thing, it's the main body for a second thing",
			},
		},
	}

	// Render the template with the test data
	RenderTemplateWithData(recorder, "base", "content", testData)

	// Check that the response writer contains the expected output
	expected := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="shortcut icon" type="image/x-icon" href="/favicon.ico" />
		<link rel="stylesheet" href="/static/css/main.css">
		<title>myFireGPT</title>
	</head>
	<body>
	  <header>
		<div >
			<nav class="container">
			   <div class="logo-container">
				  <img src="/static/images/logo.png" alt="Logo" class="logo">
				</div>
				<ul class="navbar">
					<li><a href="/">Home</a></li>
					<li><a href="/about">About</a></li>
					<li><a href="/contact">Contact</a></li>
					<li><a href="/task_list">Task List</a></li>
					<li><a href="/article-generator">Article Generator</a></li>
				</ul>
			</nav>
		</div>
	  </header>
	<div class="main-content">
	<section class="hero">
	<div class="hero-text-container">
	  <div class="search-container">
		<form class="myform">
		  <input type="text" placeholder="Search for articles...">
		  <button type="submit">Search</button>
		</form>
	  </div>
	</div>
  </section>
  <div class="article-container">
	<h2>Featured Articles</h2>
	<div class="articles">
	  {{range .Articles}}
	  <div class="article">
	  <div class="article-img-container">
		<img src="{{.Image}}" alt="Article Image">
	  </div>
		<h3 class="article-title">{{.Title}}</h3>
		<p class="article-preview">{{.Preview}}</p>
	  </div>
	  {{end}}
	</div>
  </div>
  
  <div id="task-list-container">
	
  </div>
	</div>
	  <footer>
		<div class="container">
			<p class="footer-text">&copy; MyPetGPT 2023. All rights reserved.</p>
		</div>
	  </footer>
	</body>
	</html>`
	if recorder.Body.String() != expected {
		t.Errorf("RenderTemplateWithData returned wrong result: got %v, want %v", recorder.Body.String(), expected)
	}

	// Render the template with the test data
	RenderTemplateWithData(recorder, "base", "indexContent", testData)

	// Print the actual output for debugging purposes
	fmt.Printf("Actual output: %v\n", recorder.Body.String())

	// Check that the response writer contains the expected output
	// ... (rest of the code)
}
