package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global GORM database connection pool
var DB *gorm.DB

// InitDB initializes the database connection.
// If TURSO_DATABASE_URL is set, it connects to Turso/libSQL.
// Otherwise, it falls back to a local SQLite file (local.db).
func InitDB() error {
	var gormDB *gorm.DB
	var err error

	dbURL := os.Getenv("TURSO_DATABASE_URL")
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL != "" {
		// Remote Turso connection via libsql driver
		dsn := dbURL
		if authToken != "" {
			dsn = fmt.Sprintf("%s?authToken=%s", dbURL, authToken)
		}

		sqlDB, err := sql.Open("libsql", dsn)
		if err != nil {
			return fmt.Errorf("failed to open database with libsql driver: %w", err)
		}

		// Verify connection
		if err = sqlDB.Ping(); err != nil {
			sqlDB.Close()
			return fmt.Errorf("failed to ping libsql database: %w", err)
		}

		gormDB, err = gorm.Open(sqlite.Dialector{
			Conn: sqlDB,
		}, &gorm.Config{})
		if err != nil {
			sqlDB.Close()
			return fmt.Errorf("failed to initialize gorm with libsql connection: %w", err)
		}
	} else {
		// Fall back to local SQLite file
		dbPath := os.Getenv("LOCAL_DB_PATH")
		if dbPath == "" {
			dbPath = "local.db"
		}
		gormDB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to initialize local sqlite connection: %w", err)
		}
	}

	DB = gormDB
	return nil
}
