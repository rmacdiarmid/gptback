package database

import (
	_ "github.com/mattn/go-sqlite3"
)

type Article struct {
	ID      int64
	Title   string
	Image   string
	Preview string
}

type Task struct {
	ID          int64
	Title       string
	Description string
}
