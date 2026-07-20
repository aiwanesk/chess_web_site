package tactics

import (
	"strings"
)

// MirrorFEN returns the color-swapped, vertically-flipped equivalent of a FEN.
//
// This is the anonymisation core: the position it returns is strategically
// IDENTICAL (same evaluation, same solution mirrored) but the side to move is
// the OPPOSITE colour and the resulting FEN does not match the original game in
// any database. Because a tactic's side-to-move is always Alexandre's colour,
// applying this unconditionally means every published puzzle shows the opposite
// colour to the one he actually played — his rule: White → shown as Black.
//
// Transform (reflection across the horizontal axis + colour swap):
//   - piece ranks reversed (rank 8 ↔ rank 1), each piece letter case-swapped;
//   - side to move swapped (w ↔ b);
//   - castling rights case-swapped (files preserved, so K↔k / Q↔q);
//   - en-passant square rank-flipped (r → 9-r), file preserved;
//   - clocks preserved.
func MirrorFEN(fen string) string {
	parts := strings.Fields(fen)
	if len(parts) < 4 {
		return fen // malformed — leave untouched
	}

	// 1) Piece placement: reverse rank order, swap case within each rank.
	ranks := strings.Split(parts[0], "/")
	for i, j := 0, len(ranks)-1; i < j; i, j = i+1, j-1 {
		ranks[i], ranks[j] = ranks[j], ranks[i]
	}
	for i, r := range ranks {
		ranks[i] = swapCase(r)
	}
	placement := strings.Join(ranks, "/")

	// 2) Side to move.
	side := parts[1]
	if side == "w" {
		side = "b"
	} else {
		side = "w"
	}

	// 3) Castling rights.
	castling := parts[2]
	if castling != "-" {
		castling = canonicalCastling(swapCase(castling))
	}

	// 4) En passant target square.
	ep := parts[3]
	if ep != "-" && len(ep) == 2 {
		ep = string(ep[0]) + flipRankDigit(ep[1])
	}

	out := []string{placement, side, castling, ep}
	if len(parts) > 4 {
		out = append(out, parts[4:]...) // half/full move clocks unchanged
	}
	return strings.Join(out, " ")
}

// MirrorUCIMove mirrors a UCI move (e.g. "e2e4" → "e7e5", "a7a8q" → "a2a1q") so
// the solution matches the mirrored board. Files are preserved, ranks flipped.
func MirrorUCIMove(move string) string {
	if len(move) < 4 {
		return move
	}
	var b strings.Builder
	b.WriteByte(move[0])
	b.WriteString(flipRankDigit(move[1]))
	b.WriteByte(move[2])
	b.WriteString(flipRankDigit(move[3]))
	if len(move) > 4 {
		b.WriteString(move[4:]) // promotion piece (colour implied by side to move)
	}
	return b.String()
}

func swapCase(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r - 32
		case r >= 'A' && r <= 'Z':
			return r + 32
		default:
			return r
		}
	}, s)
}

// canonicalCastling reorders rights to the conventional KQkq order.
func canonicalCastling(c string) string {
	var b strings.Builder
	for _, want := range "KQkq" {
		if strings.ContainsRune(c, want) {
			b.WriteRune(want)
		}
	}
	if b.Len() == 0 {
		return "-"
	}
	return b.String()
}

// flipRankDigit maps a rank char '1'..'8' to '8'..'1' (n → 9-n).
func flipRankDigit(d byte) string {
	if d >= '1' && d <= '8' {
		return string(byte('0') + 9 - (d - '0'))
	}
	return string(d)
}
