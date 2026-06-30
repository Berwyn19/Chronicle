package store

import (
	"database/sql"
	"fmt"
)

// Event is one recorded AI interaction. CommitHash is empty when the
// interaction was not tied to a specific commit.
type Event struct {
	ID         string
	Prompt     string
	Model      string
	Timestamp  string
	PatchPath  string
	CommitHash string
}

// List returns all recorded events, newest first.
func (s *Store) List() ([]Event, error) {
	rows, err := s.db.Query(`
		SELECT id, prompt, model, timestamp, patch_path, commit_hash
		FROM events
		ORDER BY timestamp DESC`)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		// commit_hash is nullable; scan through a sql.NullString.
		var commit sql.NullString
		if err := rows.Scan(&e.ID, &e.Prompt, &e.Model, &e.Timestamp, &e.PatchPath, &commit); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		e.CommitHash = commit.String
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}
	return events, nil
}
