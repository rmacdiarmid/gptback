package migration

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func ExecuteMigration(dbPath, migrationPath string) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			content, err := ioutil.ReadFile(filepath.Join(migrationPath, file.Name()))
			if err != nil {
				log.Fatalf("Failed to read SQL file: %v", err)
			}

			_, err = db.Exec(string(content))
			if err != nil {
				log.Fatalf("Failed to execute SQL query: %v", err)
			}

			fmt.Printf("Executed migration: %s\n", file.Name())
		}
	}
}
