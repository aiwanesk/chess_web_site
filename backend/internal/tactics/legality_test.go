package tactics

import (
	"testing"

	"github.com/notnil/chess"
)

// legalMoveCount parses a FEN and returns the number of legal moves, failing if
// the FEN is not a legal/parseable position.
func legalMoveCount(t *testing.T, fen string) int {
	t.Helper()
	opt, err := chess.FEN(fen)
	if err != nil {
		t.Fatalf("illegal FEN %q: %v", fen, err)
	}
	return len(chess.NewGame(opt).ValidMoves())
}

// TestMirrorPreservesLegalityAcrossAGame plays a real game (Ruy Lopez, incl.
// castling) and, at every position, checks that the mirrored FEN is a LEGAL
// position with exactly the same number of legal moves. Mirroring is a bijection
// on legal moves, so any illegal/corrupt mirror would change the count or fail
// to parse.
func TestMirrorPreservesLegalityAcrossAGame(t *testing.T) {
	moves := []string{
		"e4", "e5", "Nf3", "Nc6", "Bb5", "a6", "Ba4", "Nf6", "O-O", "Be7",
		"Re1", "b5", "Bb3", "d6", "c3", "O-O", "h3", "Na5", "Bc2", "c5",
		"d4", "Qc7", "Nbd2", "cxd4", "cxd4", "Nc6",
	}
	g := chess.NewGame()
	fens := []string{g.Position().String()}
	for _, m := range moves {
		if err := g.MoveStr(m); err != nil {
			t.Fatalf("illegal setup move %q: %v", m, err)
		}
		fens = append(fens, g.Position().String())
	}

	for _, fen := range fens {
		orig := legalMoveCount(t, fen)
		mirrored := MirrorFEN(fen)
		mir := legalMoveCount(t, mirrored) // fails here if mirror is illegal
		if orig != mir {
			t.Fatalf("legal-move count changed under mirror (%d → %d)\n  fen:    %s\n  mirror: %s",
				orig, mir, fen, mirrored)
		}
	}
}

// TestMirrorLegalityTrickyPositions covers en-passant and partial castling
// rights explicitly — the cases most likely to break a naive mirror.
func TestMirrorLegalityTrickyPositions(t *testing.T) {
	cases := []string{
		// White to move, en-passant target d6.
		"rnbqkbnr/1pp1pppp/p7/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
		// Black to move, en-passant target e3.
		"rnbqkbnr/pppp1ppp/8/8/4pP2/8/PPPPP1PP/RNBQKBNR b KQkq f3 0 3",
		// Partial castling rights (white king+queen side gone, black intact).
		"r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w Kq - 0 1",
	}
	for _, fen := range cases {
		orig := legalMoveCount(t, fen)
		mir := legalMoveCount(t, MirrorFEN(fen))
		if orig != mir {
			t.Fatalf("count changed for tricky pos %q: %d → %d (mirror %q)", fen, orig, mir, MirrorFEN(fen))
		}
	}
}

// TestMirroredMoveIsLegalInMirroredPosition proves the solution stays valid:
// a move legal in P maps to a move legal in mirror(P).
func TestMirroredMoveIsLegalInMirroredPosition(t *testing.T) {
	fen := "r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 0 1"
	move := "b5c6" // Bxc6, a legal capture for White

	opt, err := chess.FEN(MirrorFEN(fen))
	if err != nil {
		t.Fatalf("mirrored FEN illegal: %v", err)
	}
	mg := chess.NewGame(opt)
	uci := chess.UCINotation{}
	if _, err := uci.Decode(mg.Position(), MirrorUCIMove(move)); err != nil {
		t.Fatalf("mirrored move %q illegal in mirrored position: %v", MirrorUCIMove(move), err)
	}
}
