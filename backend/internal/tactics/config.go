// Package tactics builds the weekly, fully-anonymised tactics puzzles from
// Alexandre's own games. Usernames are read from the environment and NEVER
// written to any output (see Puzzle + mirror.go). The published JSON contains
// only a mirrored position + solution — no username, opponent, date or game URL.
package tactics

import "os"

// Config holds the private inputs. Sourced from env (.env locally / GitHub
// Secrets in CI) — never hard-coded, never committed.
type Config struct {
	ChessComUser string // e.g. from CHESSCOM_USERNAME
	LichessUser  string // e.g. from LICHESS_USERNAME
	LichessToken string // optional, higher API limits
}

func LoadConfig() Config {
	return Config{
		ChessComUser: os.Getenv("CHESSCOM_USERNAME"),
		LichessUser:  os.Getenv("LICHESS_USERNAME"),
		LichessToken: os.Getenv("LICHESS_TOKEN"),
	}
}
