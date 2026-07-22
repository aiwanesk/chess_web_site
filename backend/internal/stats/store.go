// Package stats persists Alexandre's private puzzle-interaction counters in an
// embedded SQLite database (modernc.org/sqlite — pure Go, so it works with
// CGO_ENABLED=0 + distroless). Only aggregate counters are stored — no IPs, no
// personal data.
package stats

import (
	"database/sql"
	"fmt"
	"time"

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
);
CREATE TABLE IF NOT EXISTS pageviews (
	day   TEXT NOT NULL,
	path  TEXT NOT NULL,
	count INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (day, path)
);`

// Store wraps the SQLite database.
type Store struct{ db *sql.DB }

// Open opens (creating if needed) the database at path and ensures the schema.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite: serialise writes (busy_timeout lets the
	// newsletter handle on the same file wait rather than error on contention).
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

// --- privacy-first page-view analytics ------------------------------------
// Only aggregate counts per (day, path) are stored — no IP, no user-agent, no
// cookie, no per-visitor record. Pure self-hosted, RGPD-clean.

// RecordPageview increments today's counter for a page path.
func (s *Store) RecordPageview(path string) error {
	day := time.Now().UTC().Format("2006-01-02")
	_, err := s.db.Exec(`
		INSERT INTO pageviews (day, path, count) VALUES (?, ?, 1)
		ON CONFLICT(day, path) DO UPDATE SET count = count + 1`, day, path)
	return err
}

// PageRow is one path's total views.
type PageRow struct {
	Path  string
	Count int
}

// TopPages returns the most-viewed paths (all time), busiest first.
func (s *Store) TopPages(limit int) ([]PageRow, error) {
	rows, err := s.db.Query(`SELECT path, SUM(count) AS c FROM pageviews
		GROUP BY path ORDER BY c DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []PageRow
	for rows.Next() {
		var r PageRow
		if err := rows.Scan(&r.Path, &r.Count); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// DayRow is one day's total views.
type DayRow struct {
	Day   string
	Count int
}

// DailyViews returns per-day totals, most recent first.
func (s *Store) DailyViews(days int) ([]DayRow, error) {
	rows, err := s.db.Query(`SELECT day, SUM(count) AS c FROM pageviews
		GROUP BY day ORDER BY day DESC LIMIT ?`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DayRow
	for rows.Next() {
		var r DayRow
		if err := rows.Scan(&r.Day, &r.Count); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
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
