package tactics

import "testing"

func TestMirrorStartPositionFlipsSideOnly(t *testing.T) {
	start := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	got := MirrorFEN(start)
	want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1"
	if got != want {
		t.Fatalf("mirror(start):\n got  %q\n want %q", got, want)
	}
}

func TestMirrorRoundTripIsIdentity(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
		"r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R b KQkq - 0 5",
		"8/8/8/4k3/8/3K4/4P3/8 w - - 0 1",
	}
	for _, f := range fens {
		if got := MirrorFEN(MirrorFEN(f)); got != f {
			t.Fatalf("round-trip changed the FEN:\n in  %q\n out %q", f, got)
		}
	}
}

func TestMirrorChangesFENAndFlipsSide(t *testing.T) {
	// An asymmetric position must map to a *different* FEN (so it can't be
	// matched back to the game) and to the opposite side to move.
	fen := "r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R b KQkq - 0 5"
	m := MirrorFEN(fen)
	if m == fen {
		t.Fatal("mirrored FEN must differ from the original (anonymisation)")
	}
	// side to move flipped b -> w
	if fields(m)[1] != "w" {
		t.Fatalf("expected side 'w' after mirror, got %q", fields(m)[1])
	}
}

func TestMirrorUCIMove(t *testing.T) {
	cases := map[string]string{
		"e2e4":  "e7e5",
		"a7a8q": "a2a1q",
		"e1g1":  "e8g8", // white kingside castle → black kingside castle
		"d4c5":  "d5c4",
	}
	for in, want := range cases {
		if got := MirrorUCIMove(in); got != want {
			t.Fatalf("MirrorUCIMove(%q) = %q, want %q", in, got, want)
		}
	}
}

func fields(s string) []string {
	out := []string{}
	cur := ""
	for _, r := range s {
		if r == ' ' {
			if cur != "" {
				out = append(out, cur)
				cur = ""
			}
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
