package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rmacdiarmid/GPTSite/logger"
)

var DB *sql.DB

func InitDB(dbPath string) (*sql.DB, error) {
	logger.DualLog.Printf("Initializing database with path %s", dbPath)

	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.DualLog.Printf("Error opening database: %s", err.Error())
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		logger.DualLog.Printf("Error pinging database: %s", err.Error())
		return nil, err
	}

	// Create the articles table if it doesn't exist
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS articles (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        image TEXT NOT NULL,
        preview TEXT NOT NULL
    );
    `

	_, err = DB.Exec(createTableQuery)
	if err != nil {
		logger.DualLog.Printf("Error creating table: %s", err.Error())
		return nil, err
	}

	logger.DualLog.Printf("Database initialized successfully")
	return DB, nil
}

func ReadAllTasks() ([]Task, error) {
	logger.DualLog.Printf("Fetching all tasks")
	rows, err := DB.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		logger.DualLog.Printf("Error fetching tasks: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description)
		if err != nil {
			logger.DualLog.Printf("Error scanning task: %s", err.Error())
			return nil, err
		}
		tasks = append(tasks, task)

		logger.DualLog.Printf("Fetched task: %#v", task)
	}

	err = rows.Err()
	if err != nil {
		logger.DualLog.Printf("Error iterating through rows: %s", err.Error())
		return nil, err
	}

	logger.DualLog.Printf("Fetched tasks: %#v", tasks)
	return tasks, nil
}

func CreateTask(title string, description string) (int64, error) {
	logger.DualLog.Printf("Creating task with title: %s, description: %s", title, description)

	result, err := DB.Exec("INSERT INTO tasks(title, description) VALUES (?, ?)", title, description)
	if err != nil {
		logger.DualLog.Printf("Error creating task: %s", err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.DualLog.Printf("Error getting last insert id: %s", err.Error())
		return 0, err
	}

	logger.DualLog.Printf("Created task with ID: %d, title: %s, description: %s", id, title, description)
	return id, nil
}

func UpdateTask(db *sql.DB, id int, title string, description string) error {
	logger.DualLog.Printf("Updating task with ID: %d, title: %s, description: %s", id, title, description)

	stmt, err := db.Prepare("UPDATE tasks SET title = ?, description = ? WHERE id = ?")
	if err != nil {
		logger.DualLog.Printf("Error preparing statement: %s", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(title, description, id)
	if err != nil {
		logger.DualLog.Printf("Error executing statement: %s", err.Error())
		return err
	}

	logger.DualLog.Printf("Updated task with ID: %d, title: %s, description: %s", id, title, description)
	return nil
}

func DeleteTask(db *sql.DB, id int) error {
	logger.DualLog.Printf("Deleting task with ID: %d", id)

	stmt, err := db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		logger.DualLog.Printf("Error preparing statement: %s", err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		logger.DualLog.Printf("Error executing statement: %s", err.Error())
		return err
	}

	logger.DualLog.Printf("Deleted task with ID: %d", id)
	return nil
}

func ReadTask(id int) (Task, error) {
	logger.DualLog.Printf("Reading task with ID: %d", id)

	var task Task
	err := DB.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Title, &task.Description)
	if err != nil {
		logger.DualLog.Printf("Error reading task: %s", err.Error())
		return Task{}, err
	}

	logger.DualLog.Printf("Read task with ID: %d, title: %s, description: %s", id, task.Title, task.Description)
	return task, nil
}

func CreateArticle(title, image, preview string) (int64, error) {
	logger.DualLog.Printf("Creating article with title: %s, image: %s, preview: %s", title, image, preview)

	result, err := DB.Exec("INSERT INTO articles(title, image, preview) VALUES (?, ?, ?)", title, image, preview)
	if err != nil {
		logger.DualLog.Printf("Error creating article: %s", err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.DualLog.Printf("Error getting last insert id: %s", err.Error())
		return 0, err
	}

	logger.DualLog.Printf("Created article with ID: %d, title: %s, image: %s, preview: %s", id, title, image, preview)
	return id, nil
}

func ReadArticle(id int64) (Article, error) {
	logger.DualLog.Printf("Reading article with ID: %d", id)

	var article Article
	err := DB.QueryRow("SELECT id, title, image, preview FROM articles WHERE id = ?", id).Scan(&article.ID, &article.Title, &article.Image, &article.Preview)
	if err != nil {
		logger.DualLog.Printf("Error reading article: %s", err.Error())
		return Article{}, err
	}

	logger.DualLog.Printf("Read article with ID: %d, title: %s, image: %s, preview: %s", id, article.Title, article.Image, article.Preview)
	return article, nil
}

func DeleteArticle(id int64) error {
	stmt, err := DB.Prepare("DELETE FROM articles WHERE id = ?")
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

// database.go

// UpdateArticle updates an existing article with the given ID and returns the updated article
func UpdateArticle(id int64, title, image, preview string) (*Article, error) {
	// Replace this with your own implementation to update the article in the database
	stmt, err := DB.Prepare("UPDATE articles SET title=?, image=?, preview=? WHERE id=?")
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(title, image, preview, id)
	if err != nil {
		return nil, err
	}

	updatedArticle, err := ReadArticle(id)
	if err != nil {
		return nil, err
	}

	return &updatedArticle, nil

}

func GetArticles() ([]Article, error) {
	logger.DualLog.Printf("Fetching articles")

	rows, err := DB.Query("SELECT id, title, image, preview FROM articles")
	if err != nil {
		logger.DualLog.Printf("Error fetching articles: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.ID, &article.Title, &article.Image, &article.Preview)
		if err != nil {
			logger.DualLog.Printf("Error scanning article: %s", err.Error())
			return nil, err
		}
		articles = append(articles, article)

		logger.DualLog.Printf("Fetched article: %#v", article)
	}

	err = rows.Err()
	if err != nil {
		logger.DualLog.Printf("Error iterating through rows: %s", err.Error())
		return nil, err
	}

	logger.DualLog.Printf("Fetched articles: %#v", articles)
	return articles, nil
}

func InsertArticle(title, image, preview string) (int64, error) {
	logger.DualLog.Printf("Inserting article with title: %s, image: %s, preview: %s", title, image, preview)

	result, err := DB.Exec("INSERT INTO articles(title, image, preview) VALUES (?, ?, ?)", title, image, preview)
	if err != nil {
		logger.DualLog.Printf("Error inserting article: %s", err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.DualLog.Printf("Error getting last insert id: %s", err.Error())
		return 0, err
	}

	logger.DualLog.Printf("Inserted article with ID: %d, title: %s, image: %s, preview: %s", id, title, image, preview)
	return id, nil
}
