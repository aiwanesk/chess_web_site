package corpus

import (
	"fmt"
	"strings"
)

// Le sens inverse de SANToUCI, et la liste des coups légaux. Les exercices
// tirés des parties en ont besoin : le navigateur reçoit les coups légaux de la
// position (il n'a aucune règle à connaître pour accepter un clic), et les
// coups de Stockfish, qui parle UCI, s'affichent en notation lisible.

// LegalUCIs renvoie les coups légaux de la position, en UCI.
func LegalUCIs(fen string) ([]string, error) {
	p, err := parsePosition(fen)
	if err != nil {
		return nil, err
	}
	legal := p.legalMoves()
	out := make([]string, len(legal))
	for i, m := range legal {
		out[i] = uciOf(m)
	}
	return out, nil
}

// UCIToSAN écrit un coup UCI en notation algébrique anglaise, échec et mat
// compris. Le coup doit être légal dans la position.
func UCIToSAN(fen, uci string) (string, error) {
	p, err := parsePosition(fen)
	if err != nil {
		return "", err
	}
	legal := p.legalMoves()
	var m move
	found := false
	for _, l := range legal {
		if uciOf(l) == strings.ToLower(uci) {
			m, found = l, true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("corpus: %q n'est pas légal ici", uci)
	}

	piece := p.board[m.from] | 0x20
	capture := p.board[m.to] != 0 || (piece == 'p' && m.from%8 != m.to%8)
	var san string
	switch {
	case piece == 'k' && abs(m.from%8-m.to%8) == 2:
		san = "O-O"
		if m.to%8 == 2 {
			san = "O-O-O"
		}
	case piece == 'p':
		if capture {
			san = string(rune('a'+m.from%8)) + "x"
		}
		san += squareName(m.to)
		if m.promo != 0 {
			san += "=" + strings.ToUpper(string(m.promo))
		}
	default:
		san = strings.ToUpper(string(piece)) + disambiguation(p, legal, m)
		if capture {
			san += "x"
		}
		san += squareName(m.to)
	}

	after := p.apply(m)
	if k := after.kingSquare(after.white); k >= 0 && after.attacked(k, !after.white) {
		if len(after.legalMoves()) == 0 {
			san += "#"
		} else {
			san += "+"
		}
	}
	return san, nil
}

// disambiguation : la colonne, la rangée, ou les deux, quand une autre pièce du
// même type peut aller sur la même case. Les pièces clouées ne comptent pas —
// elles ne sont pas dans legal.
func disambiguation(p position, legal []move, m move) string {
	sameFile, sameRank, other := false, false, false
	for _, l := range legal {
		if l.to != m.to || l.from == m.from || p.board[l.from] != p.board[m.from] {
			continue
		}
		other = true
		if l.from%8 == m.from%8 {
			sameFile = true
		}
		if l.from/8 == m.from/8 {
			sameRank = true
		}
	}
	switch {
	case !other:
		return ""
	case !sameFile:
		return string(rune('a' + m.from%8))
	case !sameRank:
		return string(rune('1' + m.from/8))
	default:
		return squareName(m.from)
	}
}
