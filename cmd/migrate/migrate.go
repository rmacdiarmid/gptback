package main

import (
	"github.com/rmacdiarmid/gptback/config"
	"github.com/rmacdiarmid/gptback/migration"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	dbPath := cfg.Database.Path
	migrationPath := cfg.Migration.Path

	migration.ExecuteMigration(dbPath, migrationPath)
}
