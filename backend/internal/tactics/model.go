package tactics

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

// Puzzle is one anonymised tactic ready to publish. Only the json-tagged fields
// are serialised into the published file: an already-mirrored position, its
// solution, and neutral metadata. Fields with `json:"-"` are private (never
// published) — they exist for Alexandre's own review/analytics.
//
// There is deliberately NO username, opponent, date-of-game, game URL or rating
// anywhere in this struct.
type Puzzle struct {
	ID         string   `json:"id"`         // stable hash of the mirrored FEN
	FEN        string   `json:"fen"`        // mirrored position; side to move = solver
	SideToMove string   `json:"sideToMove"` // "w" | "b" → "Les Blancs/Noirs jouent et gagnent"
	Solution   []string `json:"solution"`   // mirrored UCI moves, the winning line
	Mate       bool     `json:"mate,omitempty"`
	Sacrifice  bool     `json:"sacrifice,omitempty"`
	Beauty     int      `json:"beauty"` // ranking score

	// --- Private (never serialised) ---
	Kind   string `json:"-"` // "played" | "missed"
	Source string `json:"-"` // "chesscom" | "lichess"
}

// NewPuzzle builds an anonymised puzzle from a detected tactic given in the
// ORIGINAL orientation (side to move = Alexandre's colour). It applies the
// mandatory mirror so the published side is the opposite colour, and mirrors
// the solution line to match.
func NewPuzzle(originalFEN string, solutionUCI []string, kind, source string, mate, sacrifice bool, beauty int) Puzzle {
	fen := MirrorFEN(originalFEN)
	sol := make([]string, len(solutionUCI))
	for i, m := range solutionUCI {
		sol[i] = MirrorUCIMove(m)
	}
	side := "w"
	if f := strings.Fields(fen); len(f) > 1 {
		side = f[1]
	}
	sum := sha1.Sum([]byte(fen))
	return Puzzle{
		ID:         hex.EncodeToString(sum[:])[:12],
		FEN:        fen,
		SideToMove: side,
		Solution:   sol,
		Mate:       mate,
		Sacrifice:  sacrifice,
		Beauty:     beauty,
		Kind:       kind,
		Source:     source,
	}
}
