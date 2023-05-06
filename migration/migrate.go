package main

import (
	"github.com/go-delve/delve/pkg/config"
	"github.com/rmacdiarmid/gptback/migration"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbPath := cfg.Database.Path
	migrationPath := "./path/to/your_migration_file.sql"

	migration.ExecuteMigration(dbPath, migrationPath)
}