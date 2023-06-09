package database

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database interface {
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Article struct {
	ID      int64
	Title   string
	Image   string
	Preview string
	Text    string
}

type Task struct {
	ID          int64
	Title       string
	Description string
}

type FrontendLog struct {
	ID        int64
	Message   string
	Timestamp time.Time
}

type User struct {
	UserId       int64
	Email        string
	PasswordHash string
}

type NewUser struct {
	Email                string
	Password             string
	PasswordConfirmation string
}
