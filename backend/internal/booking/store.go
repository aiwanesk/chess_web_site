// Package booking persists lesson reservations in the shared SQLite volume.
// A booking is a date + a contiguous time range; overlapping ranges on the same
// day are rejected so a slot can't be double-booked.
package booking

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// ErrConflict means the requested slot overlaps an existing booking.
var ErrConflict = errors.New("créneau déjà réservé")

const schema = `
CREATE TABLE IF NOT EXISTS bookings (
	id         TEXT PRIMARY KEY,
	date       TEXT NOT NULL,
	start_min  INTEGER NOT NULL,
	end_min    INTEGER NOT NULL,
	name       TEXT NOT NULL,
	email      TEXT NOT NULL,
	price      INTEGER NOT NULL,
	created_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_bookings_date ON bookings(date);`

// Booking is one reserved lesson.
type Booking struct {
	ID       string
	Date     string // YYYY-MM-DD
	StartMin int    // minutes from midnight
	EndMin   int
	Name     string
	Email    string
	Price    int // CHF
}

// Store wraps the SQLite database.
type Store struct{ db *sql.DB }

// Open opens (creating if needed) the bookings database at path.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Create stores a booking, rejecting any overlap on the same day (ErrConflict).
func (s *Store) Create(b Booking) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var one int
	// Overlap iff NOT (existing ends at/before new start OR existing starts at/after new end).
	err = tx.QueryRow(`SELECT 1 FROM bookings
		WHERE date = ? AND NOT (end_min <= ? OR start_min >= ?) LIMIT 1`,
		b.Date, b.StartMin, b.EndMin).Scan(&one)
	if err == nil {
		return ErrConflict
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	id, err := randID()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO bookings
		(id, date, start_min, end_min, name, email, price, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, b.Date, b.StartMin, b.EndMin, b.Name, b.Email, b.Price, time.Now().Unix()); err != nil {
		return err
	}
	return tx.Commit()
}

// Slot is an occupied time range (minutes from midnight), no personal data.
type Slot struct{ StartMin, EndMin int }

// TakenOn returns the occupied slots for a given day, earliest first.
func (s *Store) TakenOn(date string) ([]Slot, error) {
	rows, err := s.db.Query(`SELECT start_min, end_min FROM bookings WHERE date = ? ORDER BY start_min`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Slot
	for rows.Next() {
		var sl Slot
		if err := rows.Scan(&sl.StartMin, &sl.EndMin); err != nil {
			return nil, err
		}
		out = append(out, sl)
	}
	return out, rows.Err()
}

// Upcoming returns bookings on or after fromDate (YYYY-MM-DD), soonest first.
func (s *Store) Upcoming(fromDate string, limit int) ([]Booking, error) {
	rows, err := s.db.Query(`SELECT id, date, start_min, end_min, name, email, price
		FROM bookings WHERE date >= ? ORDER BY date, start_min LIMIT ?`, fromDate, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Booking
	for rows.Next() {
		var b Booking
		if err := rows.Scan(&b.ID, &b.Date, &b.StartMin, &b.EndMin, &b.Name, &b.Email, &b.Price); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func randID() (string, error) {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
