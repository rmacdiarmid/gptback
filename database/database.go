package database

// ...package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Define the Task struct
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "tasks.db")
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %s", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("Error creating tasks table: %s", err)
	}

	return db, nil
}

func CreateTask(title, description string) (int64, error) {
	db, err := sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatalf("Error opening the database: %s", err)
		return 0, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tasks(title, description) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(title, description)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func ReadTask(db *sql.DB, id int) (*Task, error) {
	row := db.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(db *sql.DB, id int, title string, description string) error {
	_, err := db.Exec("UPDATE tasks SET title = ?, description = ? WHERE id = ?", title, description, id)
	if err != nil {
		log.Printf("Error updating task: %s", err)
		return err
	}
	return nil
}

//func DeleteTask(db *sql.DB, id int) {
//	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
//	if err != nil {
//		log.Fatalf("Error deleting task: %s", err)
//	}
//}

func DeleteTask(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id=$1", id)
	return err
}

// ReadAllTasks retrieves all tasks from the database
func ReadAllTasks() ([]Task, error) {
	// Initialize the database connection
	db, err := InitDB()
	if err != nil {
		return nil, fmt.Errorf("Error initializing database: %w", err)
	}
	defer db.Close()

	// Define the SQL query to retrieve all tasks
	query := "SELECT id, title, description FROM tasks"

	// Execute the query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice to hold the tasks
	tasks := []Task{}

	// Iterate over the result rows and add each task to the slice
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for any errors encountered while iterating over the rows
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Return the slice of tasks
	return tasks, nil
}
