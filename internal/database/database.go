// Package database provides database connection utilities for SQLite and LevelDB.
// It handles database creation, connection management, and path building for multi-shard databases.
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	_ "github.com/glebarez/go-sqlite"
	"github.com/syndtr/goleveldb/leveldb"
)

// Common errors for database operations
var (
	ErrInvalidPath      = errors.New("invalid path: at least 2 path components required")
	ErrDatabaseOpen     = errors.New("failed to open database")
	ErrDirectoryCreate  = errors.New("failed to create directory")
)

// CreateSQLiteDatabase creates and returns a SQLite database connection.
// Params should include at least the base directory and shard name.
// Optional third parameter specifies the database filename.
//
// Example usage:
//
//	db, err := CreateSQLiteDatabase("Database/", "UserData", "users.db")
//
// Returns a *sql.DB connection and any error encountered.
func CreateSQLiteDatabase(params ...string) (*sql.DB, error) {
	if len(params) < 2 {
		return nil, ErrInvalidPath
	}

	dir, err := buildPath(params...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDirectoryCreate, err)
	}

	log.Printf("Opening SQLite database at: %s", dir)
	db, err := sql.Open("sqlite", dir)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("%w: connection verification failed: %v", ErrDatabaseOpen, err)
	}

	return db, nil
}

// CreateLevelDBDatabase creates and returns a LevelDB database connection.
// Params should include at least the base directory and shard name.
// Optional third parameter specifies the database filename or subfolder.
//
// Example usage:
//
//	db, err := CreateLevelDBDatabase("Database/", "Common", "Shard_0.sqlite")
//
// Returns a *leveldb.DB connection and any error encountered.
func CreateLevelDBDatabase(params ...string) (*leveldb.DB, error) {
	if len(params) < 2 {
		return nil, ErrInvalidPath
	}

	dir, err := buildPath(params...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDirectoryCreate, err)
	}

	log.Printf("Opening LevelDB database at: %s", dir)
	db, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseOpen, err)
	}

	return db, nil
}

// buildPath constructs the database directory path from the given parameters.
// It ensures the parent directory exists, creating it if necessary.
// Returns the full path and any error encountered.
func buildPath(params ...string) (string, error) {
	if len(params) < 2 {
		return "", ErrInvalidPath
	}

	// Clean and join the path components
	cleanedParams := make([]string, len(params))
	for i, p := range params {
		cleanedParams[i] = strings.TrimSpace(p)
	}

	// Build the base directory (first two components)
	dir := path.Join(cleanedParams[0], cleanedParams[1])

	// Create the directory if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Printf("Created directory: %s", dir)
	}

	// Append additional path components (e.g., database filename)
	if len(params) > 2 {
		dir = path.Join(dir, cleanedParams[2])
	}

	return dir, nil
}

// CleanPath removes leading/trailing slashes and normalizes the path.
func CleanPath(p string) string {
	return path.Clean(strings.Trim(p, "/"))
}
