package migrations

import (
	"embed"
	"fmt"
	"sort"

	"github.com/OvsienkoValeriya/GophKeeper/internal/logger"
	"github.com/jmoiron/sqlx"
)

//go:embed *.sql
var migrationFiles embed.FS

// Run runs all SQL migrations from the migrations folder in alphabetical order
func Run(db *sqlx.DB) error {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name for correct order of execution
	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) > 4 && entry.Name()[len(entry.Name())-4:] == ".sql" {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		content, err := migrationFiles.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		logger.Sugar.Infof("Applying migration: %s", filename)
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	logger.Sugar.Infof("All migrations applied successfully")
	return nil
}
