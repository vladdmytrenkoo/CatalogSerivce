package migration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/googleapis/go-sql-spanner"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dbURI, migrationDir string) error {
	dsn := fmt.Sprintf("%s?autoConfigEmulator=true", dbURI)
	db, err := sql.Open("spanner", dsn)
	if err != nil {
		return fmt.Errorf("failed to open spanner db for migrations: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("spanner"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	log.Printf("Running manual migrations from %s", migrationDir)
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration dir: %w", err)
	}

	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		path := filepath.Join(migrationDir, f.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", f.Name(), err)
		}

		lines := strings.Split(string(content), "\n")
		var sqlLines []string
		isUpSection := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "-- +goose Up") {
				isUpSection = true
				continue
			}
			if strings.HasPrefix(trimmed, "-- +goose Down") {
				isUpSection = false
				break
			}
			if !isUpSection || strings.HasPrefix(trimmed, "--") || trimmed == "" {
				continue
			}
			sqlLines = append(sqlLines, line)
		}

		sqlText := strings.TrimSpace(strings.Join(sqlLines, "\n"))
		if sqlText == "" {
			continue
		}

		log.Printf("Executing migrations from %s", f.Name())

		statements := strings.Split(sqlText, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {

				errMsg := err.Error()
				if strings.Contains(errMsg, "already exists") || strings.Contains(errMsg, "Duplicate") {
					continue
				}
				return fmt.Errorf("failed to execute statement in %s: %w\nStatement: %s", f.Name(), err, stmt)
			}
			log.Printf("Successfully executed statement: %s...", strings.Split(stmt, "\n")[0])
		}
	}

	return nil
}
