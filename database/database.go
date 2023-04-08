package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Define the Task struct
type Task struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Article struct {
	ID      int
	Title   string
	Image   string
	Preview string
}

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Cannot open database: %v", err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT
    );`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatalf("Error creating tasks table: %s", err)
	}

	createArticlesTable := `CREATE TABLE IF NOT EXISTS articles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        image TEXT NOT NULL,
        preview TEXT NOT NULL
    );`

	_, err = DB.Exec(createArticlesTable)
	if err != nil {
		log.Fatalf("Error creating articles table: %s", err)
	}
}

func CreateTask(title, description string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO tasks(title, description) VALUES (?, ?)")
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

func ReadTask(id int) (*Task, error) {
	row := DB.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	} else if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTask(id int, title string, description string) error {
	_, err := DB.Exec("UPDATE tasks SET title = ?, description = ? WHERE id = ?", title, description, id)
	if err != nil {
		log.Printf("Error updating task: %s", err)
		return err
	}
	return nil
}

func DeleteTask(id int) error {
	_, err := DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func ReadAllTasks() ([]Task, error) {
	query := "SELECT id, title, description FROM tasks"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	tasks := []Task{}
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetArticles retrieves all articles from the database.
func GetArticles() ([]Article, error) {
	rows, err := DB.Query("SELECT id, title, image, preview FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.Title, &a.Image, &a.Preview)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Test data to check if the problem is with the database connection or not
	articles = append(articles, Article{ID: 1, Title: "Test article 1", Image: "https://via.placeholder.com/150", Preview: "This is a test article."})

	log.Printf("Fetched articles: %#v", articles)
	return articles, nil
}
