package database

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

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
	FirstName    string
	LastName     string
	Gender       string
	DateOfBirth  string
	Email        string
	PasswordHash string
	RoleId       int64
}
