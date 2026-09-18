package corpus

import (
	"fmt"
	"strings"
)

// Conversion SAN → UCI, et la génération de coups légaux qu'elle suppose.
//
// J'avais évité ce morceau aussi longtemps que possible : tant que les coups
// venaient de la base, ils étaient légaux par construction et ApplyUCI
// suffisait. Ingérer du PGN change la donne — un fichier TWIC n'écrit que du
// SAN, et « Cbd2 » ne devient « b1d2 » qu'en sachant quels cavaliers peuvent
// réellement aller en d2.
//
// Le risque d'un générateur de coups est qu'il soit presque juste. D'où la
// validation : corpus.db stocke pour chaque arête LE SAN ET L'UCI, ce qui donne
// des millions de cas de référence produits par python-chess. Le test rejoue le
// graphe et vérifie que chaque SAN retombe sur l'UCI attendu.

type position struct {
	board    [64]byte
	white    bool
	castling string
	ep       int // -1 si aucune
}

func parsePosition(fen string) (position, error) {
	f := strings.Fields(fen)
	if len(f) < 4 {
		return position{}, errFEN
	}
	b, err := parsePlacement(f[0])
	if err != nil {
		return position{}, err
	}
	p := position{board: b, white: f[1] == "w", castling: f[2], ep: -1}
	if !p.white && f[1] != "b" {
		return position{}, errFEN
	}
	if sq, ok := parseSquare(f[3]); ok {
		p.ep = sq
	}
	return p, nil
}

func (p *position) mine(sq int) bool {
	c := p.board[sq]
	return c != 0 && isWhitePiece(c) == p.white
}

func (p *position) theirs(sq int) bool {
	c := p.board[sq]
	return c != 0 && isWhitePiece(c) != p.white
}

var (
	knightSteps = [8][2]int{{1, 2}, {2, 1}, {2, -1}, {1, -2}, {-1, -2}, {-2, -1}, {-2, 1}, {-1, 2}}
	kingSteps   = [8][2]int{{0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1}, {-1, 0}, {-1, 1}}
	bishopDirs  = [4][2]int{{1, 1}, {1, -1}, {-1, -1}, {-1, 1}}
	rookDirs    = [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
)

func onBoard(file, rank int) bool { return file >= 0 && file < 8 && rank >= 0 && rank < 8 }

type move struct {
	from, to int
	promo    byte
}

// pseudoMoves génère les coups sans vérifier que le roi reste en sécurité :
// c'est legalMoves qui filtre ensuite. Séparer les deux évite d'écrire deux
// fois la logique d'attaque.
func (p *position) pseudoMoves() []move {
	var out []move
	for sq := 0; sq < 64; sq++ {
		if !p.mine(sq) {
			continue
		}
		file, rank := sq%8, sq/8
		switch p.board[sq] | 0x20 {
		case 'p':
			out = p.pawnMoves(out, sq, file, rank)
		case 'n':
			for _, s := range knightSteps {
				f, r := file+s[0], rank+s[1]
				if onBoard(f, r) && !p.mine(r*8+f) {
					out = append(out, move{sq, r*8 + f, 0})
				}
			}
		case 'k':
			for _, s := range kingSteps {
				f, r := file+s[0], rank+s[1]
				if onBoard(f, r) && !p.mine(r*8+f) {
					out = append(out, move{sq, r*8 + f, 0})
				}
			}
			out = p.castleMoves(out, sq)
		case 'b':
			out = p.slide(out, sq, file, rank, bishopDirs[:])
		case 'r':
			out = p.slide(out, sq, file, rank, rookDirs[:])
		case 'q':
			out = p.slide(out, sq, file, rank, bishopDirs[:])
			out = p.slide(out, sq, file, rank, rookDirs[:])
		}
	}
	return out
}

func (p *position) slide(out []move, sq, file, rank int, dirs [][2]int) []move {
	for _, d := range dirs {
		for i := 1; ; i++ {
			f, r := file+d[0]*i, rank+d[1]*i
			if !onBoard(f, r) {
				break
			}
			to := r*8 + f
			if p.mine(to) {
				break
			}
			out = append(out, move{sq, to, 0})
			if p.board[to] != 0 {
				break
			}
		}
	}
	return out
}

func (p *position) pawnMoves(out []move, sq, file, rank int) []move {
	dir, start, last := 1, 1, 7
	if !p.white {
		dir, start, last = -1, 6, 0
	}
	push := func(to int) {
		if to/8 == last {
			for _, pr := range []byte{'q', 'r', 'b', 'n'} {
				out = append(out, move{sq, to, pr})
			}
			return
		}
		out = append(out, move{sq, to, 0})
	}
	if one := sq + 8*dir; onBoard(file, rank+dir) && p.board[one] == 0 {
		push(one)
		if rank == start {
			if two := sq + 16*dir; p.board[two] == 0 {
				out = append(out, move{sq, two, 0})
			}
		}
	}
	for _, df := range []int{-1, 1} {
		f, r := file+df, rank+dir
		if !onBoard(f, r) {
			continue
		}
		to := r*8 + f
		if p.theirs(to) || to == p.ep {
			push(to)
		}
	}
	return out
}

// castleMoves : les cases traversées doivent être vides, et ni la case de
// départ ni celles franchies ne doivent être attaquées. Le roque est le seul
// coup où une case INTERMÉDIAIRE compte — l'oublier laisse roquer à travers un
// échec, et ce genre de bug ne se voit que sur une partie sur mille.
func (p *position) castleMoves(out []move, sq int) []move {
	if p.white && sq != 4 || !p.white && sq != 60 {
		return out
	}
	short, long := byte('K'), byte('Q')
	if !p.white {
		short, long = 'k', 'q'
	}
	if strings.IndexByte(p.castling, short) >= 0 &&
		p.board[sq+1] == 0 && p.board[sq+2] == 0 &&
		!p.attacked(sq, !p.white) && !p.attacked(sq+1, !p.white) {
		out = append(out, move{sq, sq + 2, 0})
	}
	if strings.IndexByte(p.castling, long) >= 0 &&
		p.board[sq-1] == 0 && p.board[sq-2] == 0 && p.board[sq-3] == 0 &&
		!p.attacked(sq, !p.white) && !p.attacked(sq-1, !p.white) {
		out = append(out, move{sq, sq - 2, 0})
	}
	return out
}

// attacked : la case est-elle attaquée par le camp `byWhite` ?
func (p *position) attacked(sq int, byWhite bool) bool {
	file, rank := sq%8, sq/8
	pawn, knight, king, bishop, rook, queen := byte('p'), byte('n'), byte('k'), byte('b'), byte('r'), byte('q')
	if byWhite {
		pawn, knight, king, bishop, rook, queen = 'P', 'N', 'K', 'B', 'R', 'Q'
	}
	// Un pion attaque en diagonale, vers l'avant de SON camp : on regarde donc
	// depuis la case vers l'arrière de l'attaquant.
	pdir := -1
	if byWhite {
		pdir = 1
	}
	for _, df := range []int{-1, 1} {
		f, r := file+df, rank-pdir
		if onBoard(f, r) && p.board[r*8+f] == pawn {
			return true
		}
	}
	for _, s := range knightSteps {
		f, r := file+s[0], rank+s[1]
		if onBoard(f, r) && p.board[r*8+f] == knight {
			return true
		}
	}
	for _, s := range kingSteps {
		f, r := file+s[0], rank+s[1]
		if onBoard(f, r) && p.board[r*8+f] == king {
			return true
		}
	}
	scan := func(dirs [][2]int, a, b byte) bool {
		for _, d := range dirs {
			for i := 1; ; i++ {
				f, r := file+d[0]*i, rank+d[1]*i
				if !onBoard(f, r) {
					break
				}
				c := p.board[r*8+f]
				if c == 0 {
					continue
				}
				// Premier obstacle sur CETTE direction : il attaque ou il bloque,
				// dans les deux cas on passe à la direction suivante. Un `return`
				// ici déclarerait la case sûre dès que la première diagonale est
				// bouchée — et laisserait bouger une pièce clouée.
				if c == a || c == b {
					return true
				}
				break
			}
		}
		return false
	}
	return scan(bishopDirs[:], bishop, queen) || scan(rookDirs[:], rook, queen)
}

func (p *position) kingSquare(white bool) int {
	k := byte('k')
	if white {
		k = 'K'
	}
	for sq := 0; sq < 64; sq++ {
		if p.board[sq] == k {
			return sq
		}
	}
	return -1
}

// apply joue le coup sur une copie et rend la position obtenue, côté au trait
// inversé. Volontairement minimal : seul ce qu'il faut pour tester la légalité.
func (p position) apply(m move) position {
	piece := p.board[m.from]
	lower := piece | 0x20
	if lower == 'p' && m.from%8 != m.to%8 && p.board[m.to] == 0 {
		if p.white {
			p.board[m.to-8] = 0
		} else {
			p.board[m.to+8] = 0
		}
	}
	if lower == 'k' && abs(m.from%8-m.to%8) == 2 {
		rookFrom, rookTo := m.to+1, m.to-1
		if m.to%8 == 2 {
			rookFrom, rookTo = m.to-2, m.to+1
		}
		p.board[rookTo], p.board[rookFrom] = p.board[rookFrom], 0
	}
	p.board[m.to], p.board[m.from] = piece, 0
	if m.promo != 0 {
		if p.white {
			p.board[m.to] = m.promo &^ 0x20
		} else {
			p.board[m.to] = m.promo | 0x20
		}
	}
	p.white = !p.white
	return p
}

func (p *position) legalMoves() []move {
	var out []move
	for _, m := range p.pseudoMoves() {
		after := p.apply(m)
		if k := after.kingSquare(p.white); k >= 0 && !after.attacked(k, !p.white) {
			out = append(out, m)
		}
	}
	return out
}

// SANToUCI traduit un coup écrit en notation algébrique anglaise.
//
// Les cas qui font échouer une implémentation naïve : la désambiguïsation
// partielle (« Cbd2 », « C1d2 », « Cb1d2 »), le roque écrit avec des zéros,
// la promotion, et surtout le fait qu'une pièce clouée ne compte PAS comme
// candidate — c'est ce qui rend « Cd2 » non ambigu alors que deux cavaliers
// peuvent y aller sur le papier.
func SANToUCI(fen, san string) (string, error) {
	p, err := parsePosition(fen)
	if err != nil {
		return "", err
	}
	clean := strings.TrimRight(san, "+#!?")
	clean = strings.ReplaceAll(clean, "0", "O")
	legal := p.legalMoves()

	if clean == "O-O" || clean == "O-O-O" {
		want := 6 // colonne g
		if clean == "O-O-O" {
			want = 2 // colonne c
		}
		for _, m := range legal {
			if p.board[m.from]|0x20 == 'k' && abs(m.from%8-m.to%8) == 2 && m.to%8 == want {
				return uciOf(m), nil
			}
		}
		return "", fmt.Errorf("corpus: roque impossible (%s)", san)
	}

	var promo byte
	if i := strings.IndexByte(clean, '='); i >= 0 {
		if i+1 >= len(clean) {
			return "", fmt.Errorf("corpus: promotion incomplète (%s)", san)
		}
		promo = clean[i+1] | 0x20
		clean = clean[:i]
	}
	clean = strings.ReplaceAll(clean, "x", "")
	if len(clean) < 2 {
		return "", fmt.Errorf("corpus: coup illisible (%s)", san)
	}

	dest, ok := parseSquare(clean[len(clean)-2:])
	if !ok {
		return "", fmt.Errorf("corpus: case d'arrivée illisible (%s)", san)
	}
	rest := clean[:len(clean)-2]

	kind := byte('p')
	if len(rest) > 0 && rest[0] >= 'A' && rest[0] <= 'Z' {
		kind = rest[0] | 0x20
		rest = rest[1:]
	}
	// Ce qui reste est la désambiguïsation : une colonne, une rangée, ou les deux.
	var wantFile, wantRank = -1, -1
	for i := 0; i < len(rest); i++ {
		switch c := rest[i]; {
		case c >= 'a' && c <= 'h':
			wantFile = int(c - 'a')
		case c >= '1' && c <= '8':
			wantRank = int(c - '1')
		default:
			return "", fmt.Errorf("corpus: désambiguïsation illisible (%s)", san)
		}
	}

	var found []move
	for _, m := range legal {
		if m.to != dest || p.board[m.from]|0x20 != kind || m.promo != promo {
			continue
		}
		if wantFile >= 0 && m.from%8 != wantFile {
			continue
		}
		if wantRank >= 0 && m.from/8 != wantRank {
			continue
		}
		found = append(found, m)
	}
	switch len(found) {
	case 1:
		return uciOf(found[0]), nil
	case 0:
		return "", fmt.Errorf("corpus: aucun coup légal ne correspond à %q", san)
	default:
		return "", fmt.Errorf("corpus: %q est ambigu (%d coups)", san, len(found))
	}
}

func uciOf(m move) string {
	s := squareName(m.from) + squareName(m.to)
	if m.promo != 0 {
		s += string(m.promo)
	}
	return s
}

func squareName(sq int) string {
	return string(rune('a'+sq%8)) + string(rune('1'+sq/8))
}
