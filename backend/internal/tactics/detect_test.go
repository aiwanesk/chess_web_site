package tactics

import (
	"testing"

	"github.com/notnil/chess"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name             string
		best, sec, play  Line
		wantTactic, mate bool
		wantKind         string
	}{
		{
			name: "played mate", best: Line{Mate: 1, PV: []string{"d8h4"}}, sec: Line{CP: 50}, play: Line{Mate: 1},
			wantTactic: true, mate: true, wantKind: "played",
		},
		{
			name: "missed winning shot", best: Line{CP: 500, PV: []string{"e2e4"}}, sec: Line{CP: 40}, play: Line{CP: 20},
			wantTactic: true, mate: false, wantKind: "missed",
		},
		{
			name: "no only-move gap", best: Line{CP: 320}, sec: Line{CP: 300}, play: Line{CP: 320},
			wantTactic: false,
		},
		{
			name: "already winning (not a combination)", best: Line{CP: 800}, sec: Line{CP: 700}, play: Line{CP: 800},
			wantTactic: false,
		},
		{
			name: "decisive but 2nd best already wins", best: Line{CP: 900}, sec: Line{CP: 400}, play: Line{CP: 900},
			wantTactic: false, // createsWin fails (2nd best 400 > balancedCP)
		},
	}
	for _, c := range cases {
		gotTactic, gotKind, gotMate := classify(c.best, c.sec, c.play)
		if gotTactic != c.wantTactic || (c.wantTactic && (gotKind != c.wantKind || gotMate != c.mate)) {
			t.Fatalf("%s: got (tactic=%v kind=%q mate=%v), want (tactic=%v kind=%q mate=%v)",
				c.name, gotTactic, gotKind, gotMate, c.wantTactic, c.wantKind, c.mate)
		}
	}
}

// mockEval returns a canned mate at one specific position, and quiet evals
// everywhere else — so detection is verified without a real engine.
type mockEval struct{ tacticFEN, playedUCI string }

func (m mockEval) Lines(fen string, _ int) ([]Line, error) {
	if fen == m.tacticFEN {
		return []Line{{Mate: 1, PV: []string{m.playedUCI}}, {CP: 50}}, nil
	}
	return []Line{{CP: 20}, {CP: 10}}, nil
}
func (m mockEval) Move(fen, _ string) (Line, error) {
	if fen == m.tacticFEN {
		return Line{Mate: 1}, nil
	}
	return Line{CP: 15}, nil
}

func TestAnalyzeGameFindsPlayedMate(t *testing.T) {
	// Fool's mate: Alexandre is Black and plays Qh4#.
	setup := chess.NewGame()
	for _, m := range []string{"f3", "e5", "g4"} {
		if err := setup.MoveStr(m); err != nil {
			t.Fatal(err)
		}
	}
	tacticFEN := setup.Position().String()

	game := Game{
		PGN:     `[Event "?"]` + "\n\n" + `1. f3 e5 2. g4 Qh4# 0-1`,
		MyColor: "black",
		Source:  "lichess",
	}
	tactics, err := AnalyzeGame(mockEval{tacticFEN: tacticFEN, playedUCI: "d8h4"}, game)
	if err != nil {
		t.Fatal(err)
	}
	if len(tactics) != 1 {
		t.Fatalf("want exactly 1 tactic, got %d: %+v", len(tactics), tactics)
	}
	if tactics[0].Kind != "played" || !tactics[0].Mate {
		t.Fatalf("want played mate, got kind=%q mate=%v", tactics[0].Kind, tactics[0].Mate)
	}

	// End-to-end: the published puzzle must be mirrored (Black-to-move original
	// → White-to-move puzzle) with the solution mirrored too.
	puzzles := TopPuzzles(mockEval{tacticFEN: tacticFEN, playedUCI: "d8h4"}, []Game{game}, 10)
	if len(puzzles) != 1 {
		t.Fatalf("want 1 puzzle, got %d", len(puzzles))
	}
	if puzzles[0].SideToMove != "w" {
		t.Fatalf("original was Black to move → puzzle should be White to move, got %q", puzzles[0].SideToMove)
	}
	if puzzles[0].Solution[0] != "d1h5" { // d8h4 mirrored
		t.Fatalf("solution not mirrored: got %v", puzzles[0].Solution)
	}
}
