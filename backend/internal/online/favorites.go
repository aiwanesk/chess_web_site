package online

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Les favoris sont la seule chose gardée sur le disque : une liste de pseudos
// à recharger d'un clic avant une ronde. Pas les parties — elles changent
// chaque semaine, et se retéléchargent en quelques secondes.

// MaxFavorites borne la liste : au-delà, c'est un annuaire, plus une sélection.
const MaxFavorites = 100

// ErrTooManyFavorites : la liste est pleine.
var ErrTooManyFavorites = errors.New("100 favoris au maximum : retire-en un d'abord")

const favSchema = `
CREATE TABLE IF NOT EXISTS prep_favorites (
	source     TEXT NOT NULL,
	username   TEXT NOT NULL COLLATE NOCASE,
	note       TEXT NOT NULL DEFAULT '',
	created_at INTEGER NOT NULL,
	PRIMARY KEY (source, username)
);`

// Favorite est un pseudo gardé.
type Favorite struct {
	Source   Source `json:"source"`
	Username string `json:"username"`
	Note     string `json:"note"`
}

// Favorites stocke les favoris dans la base SQLite du site (DB_PATH).
type Favorites struct{ db *sql.DB }

// OpenFavorites ouvre (et crée au besoin) la table des favoris.
func OpenFavorites(path string) (*Favorites, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(favSchema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Favorites{db: db}, nil
}

func (f *Favorites) Close() error { return f.db.Close() }

// List renvoie les favoris par ordre alphabétique.
func (f *Favorites) List() ([]Favorite, error) {
	rows, err := f.db.Query(`SELECT source, username, note FROM prep_favorites
		ORDER BY username COLLATE NOCASE, source`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Favorite{}
	for rows.Next() {
		var x Favorite
		if err := rows.Scan(&x.Source, &x.Username, &x.Note); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

// Add ajoute un favori, ou met à jour sa note s'il existe déjà.
func (f *Favorites) Add(x Favorite) error {
	x.Note = strings.TrimSpace(x.Note)
	if len([]rune(x.Note)) > 80 {
		x.Note = string([]rune(x.Note)[:80])
	}
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var exists int
	_ = tx.QueryRow(`SELECT COUNT(*) FROM prep_favorites WHERE source = ? AND username = ?`,
		x.Source, x.Username).Scan(&exists)
	if exists == 0 {
		var n int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM prep_favorites`).Scan(&n); err != nil {
			return err
		}
		if n >= MaxFavorites {
			return ErrTooManyFavorites
		}
	}
	if _, err := tx.Exec(`INSERT INTO prep_favorites (source, username, note, created_at)
		VALUES (?, ?, ?, ?)
		ON CONFLICT (source, username) DO UPDATE SET note = excluded.note`,
		x.Source, x.Username, x.Note, time.Now().Unix()); err != nil {
		return err
	}
	return tx.Commit()
}

// Remove retire un favori ; retirer un absent n'est pas une erreur.
func (f *Favorites) Remove(src Source, username string) error {
	_, err := f.db.Exec(`DELETE FROM prep_favorites WHERE source = ? AND username = ?`, src, username)
	return err
}
