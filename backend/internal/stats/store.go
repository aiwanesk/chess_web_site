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
);
CREATE TABLE IF NOT EXISTS traffic (
	day   TEXT NOT NULL,
	kind  TEXT NOT NULL,   -- 'human' | 'bot'
	count INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (day, kind)
);
CREATE TABLE IF NOT EXISTS geo (
	country TEXT PRIMARY KEY, -- 2-letter code (offline lookup, no IP stored)
	count   INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS visitors (
	day TEXT NOT NULL,
	fp  TEXT NOT NULL,        -- salted daily hash of IP+UA (non-reversible, no IP)
	PRIMARY KEY (day, fp)
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

// RecordPageview records one page view. Bots are counted in `traffic` only;
// humans also count toward the page's views, the day's unique visitors (by a
// non-reversible daily fingerprint — no IP stored) and the country tally.
func (s *Store) RecordPageview(path, country string, isBot bool, fp string) error {
	if len(path) > 200 { // defensive cap against pathological URLs
		path = path[:200]
	}
	day := time.Now().UTC().Format("2006-01-02")
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	bump := func(q string, args ...any) error { _, e := tx.Exec(q, args...); return e }

	kind := "human"
	if isBot {
		kind = "bot"
	}
	if err := bump(`INSERT INTO traffic (day, kind, count) VALUES (?, ?, 1)
		ON CONFLICT(day, kind) DO UPDATE SET count = count + 1`, day, kind); err != nil {
		return err
	}
	if !isBot {
		if err := bump(`INSERT INTO pageviews (day, path, count) VALUES (?, ?, 1)
			ON CONFLICT(day, path) DO UPDATE SET count = count + 1`, day, path); err != nil {
			return err
		}
		if country != "" {
			if err := bump(`INSERT INTO geo (country, count) VALUES (?, 1)
				ON CONFLICT(country) DO UPDATE SET count = count + 1`, country); err != nil {
				return err
			}
		}
		if fp != "" {
			if err := bump(`INSERT INTO visitors (day, fp) VALUES (?, ?) ON CONFLICT DO NOTHING`, day, fp); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// CountryRow is one country's human view tally.
type CountryRow struct {
	Country string
	Count   int
}

// TopCountries returns the busiest visitor countries (all time).
func (s *Store) TopCountries(limit int) ([]CountryRow, error) {
	rows, err := s.db.Query(`SELECT country, count FROM geo ORDER BY count DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CountryRow
	for rows.Next() {
		var r CountryRow
		if err := rows.Scan(&r.Country, &r.Count); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// TrafficRow is a day's human vs bot split.
type TrafficRow struct {
	Day                 string
	Human, Bot, Uniques int
}

// Traffic returns per-day human/bot totals + unique human visitors, recent first.
func (s *Store) Traffic(days int) ([]TrafficRow, error) {
	rows, err := s.db.Query(`
		SELECT d.day,
		  COALESCE(SUM(CASE WHEN t.kind='human' THEN t.count END), 0),
		  COALESCE(SUM(CASE WHEN t.kind='bot'   THEN t.count END), 0),
		  COALESCE((SELECT COUNT(*) FROM visitors v WHERE v.day = d.day), 0)
		FROM (SELECT DISTINCT day FROM traffic) d
		LEFT JOIN traffic t ON t.day = d.day
		GROUP BY d.day ORDER BY d.day DESC LIMIT ?`, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TrafficRow
	for rows.Next() {
		var r TrafficRow
		if err := rows.Scan(&r.Day, &r.Human, &r.Bot, &r.Uniques); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
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
