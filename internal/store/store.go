// Package store owns Chronicle's SQLite metadata database.
//
// It is the only package that talks to SQLite directly. Callers open a Store,
// use it, and Close it. The schema is applied idempotently on open so that
// opening a fresh database initializes it.
package store

import (
	"database/sql"
	_ "embed"
	"fmt"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

// Store wraps the SQLite connection holding Chronicle metadata.
type Store struct {
	db *sql.DB
}

// Open opens (creating if necessary) the SQLite database at dbPath and applies
// the schema. The caller is responsible for ensuring the parent directory
// exists.
func Open(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// Enforce foreign keys; off by default in SQLite.
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the underlying database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
