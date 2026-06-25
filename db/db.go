package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "turso.tech/database/tursogo"
)

// DB is the global database connection pool
var DB *sql.DB

// InitDB initializes the database connection.
// If TURSO_DATABASE_URL is set, it connects to Turso/libSQL.
// Otherwise, it falls back to a local SQLite file (local.db).
func InitDB() error {
	var sqlDB *sql.DB
	var err error

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL != "" {
		// Remote Turso connection via libsql driver
		dsn := dbURL
		if authToken != "" {
			dsn = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
		}

		sqlDB, err = sql.Open("turso", dsn)
		if err != nil {
			return fmt.Errorf("failed to open database with turso driver: %w", err)
		}
	} else {
		// Fall back to local SQLite file
		dbPath := os.Getenv("LOCAL_DB_PATH")
		if dbPath == "" {
			dbPath = "local.db"
		}
		sqlDB, err = sql.Open("turso", dbPath)
		if err != nil {
			return fmt.Errorf("failed to initialize local turso connection: %w", err)
		}
	}

	// Verify connection
	if err = sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = sqlDB
	return nil
}

// CreateTables creates the necessary database tables if they do not exist.
func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id TEXT PRIMARY KEY,
			note TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			tags TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS note_vectors (
			id TEXT PRIMARY KEY,
			parent_id TEXT NOT NULL,
			chunk TEXT NOT NULL,
			embedding F32_BLOB(768),
			tags TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS documents (
			id TEXT PRIMARY KEY,
			file_path TEXT UNIQUE NOT NULL,
			title TEXT,
			summary TEXT,
			tags TEXT,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS document_vectors (
			id TEXT PRIMARY KEY,
			document_id TEXT NOT NULL,
			chunk_number INTEGER NOT NULL,
			embedding F32_BLOB(768),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query %q: %w", query, err)
		}
	}
	return nil
}
