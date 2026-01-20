// Package database provides database connection and initialization utilities.
package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

// Config holds database configuration.
type Config struct {
	Path           string
	MigrationsPath string
}

// Open opens a SQLite database connection and runs migrations.
func Open(cfg Config) (*sql.DB, error) {
	// Ensure the directory exists
	dir := filepath.Dir(cfg.Path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite3", cfg.Path+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations if path is provided
	if cfg.MigrationsPath != "" {
		if err := runMigrations(db, cfg.MigrationsPath); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return db, nil
}

// runMigrations runs all pending database migrations.
func runMigrations(db *sql.DB, migrationsPath string) error {
	goose.SetBaseFS(nil)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// OpenInMemory opens an in-memory SQLite database and runs migrations.
// Useful for testing.
func OpenInMemory(migrationsPath string) (*sql.DB, error) {
	return Open(Config{
		Path:           ":memory:",
		MigrationsPath: migrationsPath,
	})
}

// OpenTemp opens a temporary SQLite database file and runs migrations.
// Returns the database connection and a cleanup function.
// Useful for E2E testing with isolated databases.
func OpenTemp(migrationsPath string) (*sql.DB, func(), error) {
	tmpFile, err := os.CreateTemp("", "powerpro-test-*.db")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create temp db file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	db, err := Open(Config{
		Path:           tmpPath,
		MigrationsPath: migrationsPath,
	})
	if err != nil {
		os.Remove(tmpPath)
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
		os.Remove(tmpPath)
	}

	return db, cleanup, nil
}
