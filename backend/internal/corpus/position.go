package corpus

import (
	"fmt"
	"strconv"
	"strings"
)

// ApplyUCI joue un coup et renvoie le FEN de la position obtenue.
//
// Le coup n'est PAS vérifié comme légal, et c'est volontaire : les seuls coups
// proposés au client sont ceux que la base connaît pour la position courante,
// donc légaux par construction. Un moteur de génération de coups complet
// n'apporterait rien ici, sinon un second endroit où diverger de python-chess.
//
// Ce qui est vérifié, en revanche : la forme du coup, la présence d'une pièce
// au départ, et qu'elle appartient au trait. Le FEN vient du navigateur.
func ApplyUCI(fen, uci string) (string, error) {
	f := strings.Fields(fen)
	if len(f) < 4 {
		return "", errFEN
	}
	board, err := parsePlacement(f[0])
	if err != nil {
		return "", err
	}
	from, to, promo, err := parseUCI(uci)
	if err != nil {
		return "", err
	}

	piece := board[from]
	if piece == 0 {
		return "", fmt.Errorf("corpus: aucune pièce en %s", uci[:2])
	}
	whiteToMove := f[1] == "w"
	if isWhitePiece(piece) != whiteToMove {
		return "", fmt.Errorf("corpus: la pièce en %s n'est pas au trait", uci[:2])
	}

	captured := board[to] != 0
	lower := piece | 0x20 // minuscule : le type de pièce, couleur mise de côté

	// Prise en passant : un pion se déplace en diagonale vers une case vide.
	if lower == 'p' && from%8 != to%8 && board[to] == 0 {
		board[to-8+16*b2i(!whiteToMove)] = 0 // le pion capturé est sur la rangée de départ du preneur
		captured = true
	}

	// Roque : le roi saute deux colonnes, la tour l'accompagne.
	if lower == 'k' && abs(from%8-to%8) == 2 {
		rookFrom, rookTo := to+1, to-1 // petit roque
		if to%8 == 2 {                 // grand roque
			rookFrom, rookTo = to-2, to+1
		}
		board[rookTo], board[rookFrom] = board[rookFrom], 0
	}

	board[to], board[from] = piece, 0
	if promo != 0 {
		if whiteToMove {
			board[to] = promo &^ 0x20 // majuscule
		} else {
			board[to] = promo | 0x20
		}
	}

	next := []string{
		writePlacement(board),
		map[bool]string{true: "b", false: "w"}[whiteToMove],
		nextCastling(f[2], from, to),
		nextEnPassant(lower, from, to),
		nextHalfmove(f, lower == 'p' || captured),
		nextFullmove(f, whiteToMove),
	}
	return strings.Join(next, " "), nil
}

// nextCastling retire les droits qu'un coup fait perdre : le roi qui bouge, la
// tour qui bouge, et la tour qui se fait prendre sur sa case d'origine.
func nextCastling(rights string, from, to int) string {
	lost := map[int]string{
		4: "KQ", 60: "kq", // rois
		0: "Q", 7: "K", 56: "q", 63: "k", // tours
	}
	drop := lost[from] + lost[to]
	out := strings.Map(func(r rune) rune {
		if strings.ContainsRune(drop, r) {
			return -1
		}
		return r
	}, rights)
	if out == "" {
		return "-"
	}
	return out
}

// nextEnPassant n'annonce une case que sur une poussée de deux. On l'écrit
// systématiquement, sans regarder si un pion peut réellement prendre : c'est
// Hash() qui applique cette règle-là, et il le fait sur le FEN quel qu'il soit.
func nextEnPassant(lower byte, from, to int) string {
	if lower != 'p' || abs(from/8-to/8) != 2 {
		return "-"
	}
	mid := (from + to) / 2
	return string(rune('a'+mid%8)) + strconv.Itoa(mid/8+1)
}

func nextHalfmove(f []string, reset bool) string {
	if reset {
		return "0"
	}
	if len(f) < 5 {
		return "0"
	}
	n, err := strconv.Atoi(f[4])
	if err != nil {
		return "0"
	}
	return strconv.Itoa(n + 1)
}

func nextFullmove(f []string, whiteToMove bool) string {
	n := 1
	if len(f) >= 6 {
		if v, err := strconv.Atoi(f[5]); err == nil {
			n = v
		}
	}
	if !whiteToMove { // les Noirs viennent de jouer : le numéro avance
		n++
	}
	return strconv.Itoa(n)
}

func writePlacement(board [64]byte) string {
	var b strings.Builder
	for rank := 7; rank >= 0; rank-- {
		empty := 0
		for file := 0; file < 8; file++ {
			pc := board[rank*8+file]
			if pc == 0 {
				empty++
				continue
			}
			if empty > 0 {
				b.WriteString(strconv.Itoa(empty))
				empty = 0
			}
			b.WriteByte(pc)
		}
		if empty > 0 {
			b.WriteString(strconv.Itoa(empty))
		}
		if rank > 0 {
			b.WriteByte('/')
		}
	}
	return b.String()
}

func parseUCI(uci string) (from, to int, promo byte, err error) {
	if len(uci) != 4 && len(uci) != 5 {
		return 0, 0, 0, fmt.Errorf("corpus: coup UCI invalide %q", uci)
	}
	f, ok1 := parseSquare(uci[:2])
	t, ok2 := parseSquare(uci[2:4])
	if !ok1 || !ok2 {
		return 0, 0, 0, fmt.Errorf("corpus: coup UCI invalide %q", uci)
	}
	if len(uci) == 5 {
		switch uci[4] {
		case 'q', 'r', 'b', 'n':
			promo = uci[4]
		default:
			return 0, 0, 0, fmt.Errorf("corpus: promotion invalide %q", uci)
		}
	}
	return f, t, promo, nil
}

func isWhitePiece(p byte) bool { return p >= 'A' && p <= 'Z' }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
