package database

// ...package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Define the Task struct
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatalf("Error creating tasks table: %s", err)
	}

	return db
}

func CreateTask(db *sql.DB, title string, description string) int64 {
	result, err := db.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		log.Fatalf("Error creating task: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting task ID: %s", err)
	}

	return id
}

func ReadTask(db *sql.DB, id int) *Task {
	row := db.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		log.Fatalf("Error reading task: %s", err)
	}

	return &task
}

func UpdateTask(db *sql.DB, id int, title string, description string) {
	_, err := db.Exec("UPDATE tasks SET title = ?, description = ? WHERE id = ?", title, description, id)
	if err != nil {
		log.Fatalf("Error updating task: %s", err)
	}
}

func DeleteTask(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		log.Fatalf("Error deleting task: %s", err)
	}
}

// Your other database functions for CRUD operations...
