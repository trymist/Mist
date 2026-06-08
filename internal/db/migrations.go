package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func dbMigrate(db *sql.DB) error {
	// create migrations table if doesn't exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY
		);
		`)
	if err != nil {
		return err
	}
	rows, err := db.Query(`SELECT version FROM schema_migrations`)
	if err != nil {
		return err
	}
	defer rows.Close()
	appliedMigrations := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		appliedMigrations[version] = true
	}
	// TODO: should be embedded in the binary in prod releases
	migrationsDir := "internal/db/migrations"
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if filepath.Ext(filename) != ".sql" {
			continue
		}
		if appliedMigrations[filename] {
			continue
		}
		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return err
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", filename, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, filename); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil

}
