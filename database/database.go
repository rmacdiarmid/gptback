package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		return nil, err
	}

	return DB, nil
}

type Article struct {
	ID      int64
	Title   string
	Image   string
	Preview string
}

func GetArticles() ([]Article, error) {
	rows, err := DB.Query("SELECT id, title, image, preview FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Image, &article.Preview)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	log.Printf("Fetched articles: %#v", articles)
	return articles, nil
}

type Task struct {
	ID          int64
	Title       string
	Description string
}

func ReadAllTasks() ([]Task, error) {
	rows, err := DB.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
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

	log.Printf("Fetched tasks: %#v", tasks)
	return tasks, nil
}

func CreateTask(title string, description string) (int64, error) {
	result, err := DB.Exec("INSERT INTO tasks(title, description) VALUES (?, ?)", title, description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	log.Printf("Created task with ID: %d", id)
	return id, nil
}

func UpdateTask(db *sql.DB, id int, title, description string) error {
	stmt, err := db.Prepare("UPDATE tasks SET title = ?, description = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, description, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(db *sql.DB, id int) error {
	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func ReadTask(db *sql.DB, id int) (Task, error) {
	var task Task
	err := db.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
