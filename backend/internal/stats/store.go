// Package stats persists Alexandre's private puzzle-interaction counters in an
// embedded SQLite database (modernc.org/sqlite — pure Go, so it works with
// CGO_ENABLED=0 + distroless). Only aggregate counters are stored — no IPs, no
// personal data.
package stats

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Kinds of interaction events.
const (
	View    = "view"
	Attempt = "attempt"
	Solved  = "solved"
)

var kindColumn = map[string]string{View: "views", Attempt: "attempts", Solved: "solved"}

const schema = `
CREATE TABLE IF NOT EXISTS puzzle_stats (
	week      TEXT NOT NULL,
	puzzle_id TEXT NOT NULL,
	views     INTEGER NOT NULL DEFAULT 0,
	attempts  INTEGER NOT NULL DEFAULT 0,
	solved    INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (week, puzzle_id)
);`

// Store wraps the SQLite database.
type Store struct{ db *sql.DB }

// Open opens (creating if needed) the database at path and ensures the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite: serialise writes
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Record increments the counter for one interaction event. `kind` must be one
// of View/Attempt/Solved; the column name comes from a fixed whitelist (no
// injection), the values are parameterised.
func (s *Store) Record(week, puzzleID, kind string) error {
	col, ok := kindColumn[kind]
	if !ok {
		return fmt.Errorf("stats: unknown kind %q", kind)
	}
	q := fmt.Sprintf(`
		INSERT INTO puzzle_stats (week, puzzle_id, %s) VALUES (?, ?, 1)
		ON CONFLICT(week, puzzle_id) DO UPDATE SET %s = %s + 1`, col, col, col)
	_, err := s.db.Exec(q, week, puzzleID)
	return err
}

// Row is one aggregated puzzle line for the dashboard.
type Row struct {
	Week, PuzzleID          string
	Views, Attempts, Solved int
}

// Summary returns every puzzle's counters, newest week first.
func (s *Store) Summary() ([]Row, error) {
	rows, err := s.db.Query(`SELECT week, puzzle_id, views, attempts, solved
		FROM puzzle_stats ORDER BY week DESC, puzzle_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Row
	for rows.Next() {
		var r Row
		if err := rows.Scan(&r.Week, &r.PuzzleID, &r.Views, &r.Attempts, &r.Solved); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
