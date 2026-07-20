package tactics

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNewPuzzleMirrorsAndFlipsSide(t *testing.T) {
	// White to move in the original → published as Black to move.
	p := NewPuzzle("8/8/8/4k3/8/3K4/4P3/8 w - - 0 1", []string{"e2e4"}, "played", "lichess", false, false, 10)
	if p.SideToMove != "b" {
		t.Fatalf("side to move should be flipped to 'b', got %q", p.SideToMove)
	}
	if p.Solution[0] != "e7e5" {
		t.Fatalf("solution not mirrored: got %q", p.Solution[0])
	}
	if p.ID == "" {
		t.Fatal("puzzle id must be set")
	}
}

// The published JSON must never contain identifying data: no username, no
// platform source, no played/missed marker.
func TestPublishedJSONCarriesNoIdentifiers(t *testing.T) {
	p := NewPuzzle("r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R b KQkq - 0 5",
		[]string{"e2e4"}, "missed", "chesscom", true, true, 42)

	raw, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	s := string(raw)

	// Private metadata must not serialise into the published file.
	for _, forbidden := range []string{"chesscom", "missed", "Source", "Kind"} {
		if strings.Contains(s, forbidden) {
			t.Fatalf("published JSON leaked private field %q: %s", forbidden, s)
		}
	}

	// Real usernames (read from env, never hard-coded here so this test file is
	// safe to commit to the public repo) must never appear in the output.
	for _, key := range []string{"CHESSCOM_USERNAME", "LICHESS_USERNAME"} {
		if u := os.Getenv(key); u != "" && strings.Contains(strings.ToLower(s), strings.ToLower(u)) {
			t.Fatalf("published JSON leaked the username from %s", key)
		}
	}
}
