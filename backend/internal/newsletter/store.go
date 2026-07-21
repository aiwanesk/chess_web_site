// Package newsletter persists the mailing list in the same SQLite volume as the
// stats DB. It is built for GDPR / Swiss nLPD compliance by design:
//
//   - Data minimisation: only the e-mail, a language preference, a status and
//     two random tokens are stored. No IP, no name, no tracking.
//   - Explicit consent via double opt-in: a subscriber stays "pending" (and is
//     never sent anything but the single confirmation mail) until they click the
//     confirmation link. Only "confirmed" rows ever receive announcements.
//   - Right to erasure: unsubscribing DELETES the row outright.
package newsletter

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	seededSentinel  = "__seeded__" // marks that the first-run backfill has happened
)

const schema = `
CREATE TABLE IF NOT EXISTS subscribers (
	email         TEXT PRIMARY KEY,
	lang          TEXT NOT NULL DEFAULT 'fr',
	status        TEXT NOT NULL DEFAULT 'pending',
	confirm_token TEXT NOT NULL,
	unsub_token   TEXT NOT NULL,
	created_at    INTEGER NOT NULL,
	confirmed_at  INTEGER
);
CREATE TABLE IF NOT EXISTS notified (
	content_id  TEXT PRIMARY KEY,
	notified_at INTEGER NOT NULL
);`

// Store wraps the SQLite database.
type Store struct{ db *sql.DB }

// Open opens (creating if needed) the newsletter database at path. A busy
// timeout lets it coexist with the stats handle on the same file.
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

// NormalizeEmail lower-cases and trims an address for storage/lookup.
func NormalizeEmail(e string) string { return strings.ToLower(strings.TrimSpace(e)) }

// Subscribe records (or refreshes) a pending subscription and returns the
// confirmation token to e-mail. If the address is already confirmed it returns
// alreadyConfirmed=true and no token (the caller must not leak this difference
// to the visitor — always show the same "check your inbox" message).
func (s *Store) Subscribe(email, lang string) (confirmToken string, alreadyConfirmed bool, err error) {
	email = NormalizeEmail(email)
	if lang != "en" {
		lang = "fr"
	}
	now := time.Now().Unix()

	var status, existingConfirm, existingUnsub string
	row := s.db.QueryRow(`SELECT status, confirm_token, unsub_token FROM subscribers WHERE email = ?`, email)
	switch err := row.Scan(&status, &existingConfirm, &existingUnsub); {
	case err == nil:
		if status == StatusConfirmed {
			return "", true, nil
		}
		// Pending → issue a fresh confirmation link, keep the unsub token.
		ct, e := randToken()
		if e != nil {
			return "", false, e
		}
		if _, e := s.db.Exec(`UPDATE subscribers SET lang = ?, confirm_token = ?, created_at = ? WHERE email = ?`,
			lang, ct, now, email); e != nil {
			return "", false, e
		}
		return ct, false, nil
	case errors.Is(err, sql.ErrNoRows):
		ct, e := randToken()
		if e != nil {
			return "", false, e
		}
		ut, e := randToken()
		if e != nil {
			return "", false, e
		}
		if _, e := s.db.Exec(`INSERT INTO subscribers (email, lang, status, confirm_token, unsub_token, created_at)
			VALUES (?, ?, ?, ?, ?, ?)`, email, lang, StatusPending, ct, ut, now); e != nil {
			return "", false, e
		}
		return ct, false, nil
	default:
		return "", false, err
	}
}

// Confirm marks a pending subscriber as confirmed (double opt-in). Returns the
// e-mail/lang on success so the caller can greet the right language.
func (s *Store) Confirm(token string) (email, lang string, ok bool, err error) {
	if token == "" {
		return "", "", false, nil
	}
	row := s.db.QueryRow(`SELECT email, lang FROM subscribers WHERE confirm_token = ? AND status = ?`,
		token, StatusPending)
	if err := row.Scan(&email, &lang); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", false, nil
		}
		return "", "", false, err
	}
	if _, err := s.db.Exec(`UPDATE subscribers SET status = ?, confirmed_at = ? WHERE confirm_token = ?`,
		StatusConfirmed, time.Now().Unix(), token); err != nil {
		return "", "", false, err
	}
	return email, lang, true, nil
}

// Unsubscribe deletes the subscriber for the given token (right to erasure).
// Returns ok=true if a row was removed (idempotent: unknown token → ok=false).
func (s *Store) Unsubscribe(token string) (bool, error) {
	if token == "" {
		return false, nil
	}
	res, err := s.db.Exec(`DELETE FROM subscribers WHERE unsub_token = ?`, token)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

// Recipient is a confirmed subscriber ready to receive an announcement.
type Recipient struct{ Email, Lang, UnsubToken string }

// ConfirmedRecipients returns every confirmed subscriber.
func (s *Store) ConfirmedRecipients() ([]Recipient, error) {
	rows, err := s.db.Query(`SELECT email, lang, unsub_token FROM subscribers WHERE status = ?`, StatusConfirmed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Recipient
	for rows.Next() {
		var r Recipient
		if err := rows.Scan(&r.Email, &r.Lang, &r.UnsubToken); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// --- announcement bookkeeping ----------------------------------------------

// Seeded reports whether the first-run backfill has already run. On a fresh DB
// it is false, letting the announcer mark all current content as "already
// notified" so deploying the feature never mass-mails existing content.
func (s *Store) Seeded() (bool, error) {
	var one int
	err := s.db.QueryRow(`SELECT 1 FROM notified WHERE content_id = ?`, seededSentinel).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// MarkNotified records that these content IDs have been announced (idempotent).
// The seeded sentinel is written alongside so future runs skip the backfill.
func (s *Store) MarkNotified(ids ...string) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now().Unix()
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO notified (content_id, notified_at) VALUES (?, ?)
		ON CONFLICT(content_id) DO NOTHING`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, id := range ids {
		if _, err := stmt.Exec(id, now); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// MarkSeeded writes the sentinel (called after the first-run backfill).
func (s *Store) MarkSeeded() error { return s.MarkNotified(seededSentinel) }

// IsNotified reports whether a content ID has already been announced.
func (s *Store) IsNotified(id string) (bool, error) {
	var one int
	err := s.db.QueryRow(`SELECT 1 FROM notified WHERE content_id = ?`, id).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// randToken returns a URL-safe, unguessable 144-bit token.
func randToken() (string, error) {
	b := make([]byte, 18)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
