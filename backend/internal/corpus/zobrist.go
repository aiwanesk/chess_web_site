// Package corpus interroge le graphe de positions construit à partir des cours
// (corpus.db) : à une position donnée, quels coups y sont joués et dans quelles
// proportions, et ce que les cours en disent.
//
// Les positions y sont identifiées par leur clé Polyglot, calculée ici. C'est
// la même convention que python-chess, qui a servi à indexer la base — une
// divergence d'un seul bit rendrait la base muette plutôt que fausse, ce qui
// est exactement le genre de panne qu'on met deux jours à diagnostiquer. D'où
// le jeu de référence de testdata/zobrist.json, recoupé avec la base elle-même.
package corpus

import (
	"errors"
	"strings"
)

// Ordre des pièces dans la table Polyglot : noir puis blanc, pion → roi.
// L'index d'une pièce sur une case vaut 64*kind + 8*rangée + colonne.
var polyglotKind = map[byte]int{
	'p': 0, 'P': 1,
	'n': 2, 'N': 3,
	'b': 4, 'B': 5,
	'r': 6, 'R': 7,
	'q': 8, 'Q': 9,
	'k': 10, 'K': 11,
}

const (
	castleOffset = 768 // KQkq, dans cet ordre
	epOffset     = 772 // + colonne
	turnKey      = 780 // XOR si le trait est aux Blancs
)

var errFEN = errors.New("corpus: FEN invalide")

// Hash renvoie la clé Polyglot de la position décrite par fen.
func Hash(fen string) (uint64, error) {
	f := strings.Fields(fen)
	if len(f) < 4 {
		return 0, errFEN
	}
	board, err := parsePlacement(f[0])
	if err != nil {
		return 0, err
	}

	var h uint64
	for sq, pc := range board {
		if pc == 0 {
			continue
		}
		kind, ok := polyglotKind[pc]
		if !ok {
			return 0, errFEN
		}
		h ^= polyglotKeys[64*kind+sq]
	}

	whiteToMove := f[1] == "w"
	if !whiteToMove && f[1] != "b" {
		return 0, errFEN
	}

	for i, c := range []byte{'K', 'Q', 'k', 'q'} {
		if strings.IndexByte(f[2], c) >= 0 {
			h ^= polyglotKeys[castleOffset+i]
		}
	}

	// La colonne de prise en passant n'entre dans le hash que si un pion du
	// trait est effectivement à côté de la case. La LÉGALITÉ de la prise est
	// hors sujet : un pion cloué compte quand même. C'est la règle de Polyglot,
	// et c'est le piège de tout ce fichier — une implémentation qui ajoute la
	// colonne dès que le FEN mentionne une case ep marche sur la quasi-totalité
	// des positions, puis renvoie zéro coup juste après chaque poussée de deux.
	if sq, ok := parseSquare(f[3]); ok && epCapturable(board, sq, whiteToMove) {
		h ^= polyglotKeys[epOffset+sq%8]
	}

	if whiteToMove {
		h ^= polyglotKeys[turnKey]
	}
	return h, nil
}

// SignedHash renvoie la clé telle qu'elle est STOCKÉE : SQLite n'a pas d'entier
// non signé, donc l'indexeur enregistre le complément à deux. Plus de la moitié
// des lignes de `move` sont ainsi négatives.
func SignedHash(fen string) (int64, error) {
	h, err := Hash(fen)
	return int64(h), err
}

// parsePlacement remplit un échiquier indexé 0 = a1 … 63 = h8.
func parsePlacement(s string) ([64]byte, error) {
	var b [64]byte
	rows := strings.Split(s, "/")
	if len(rows) != 8 {
		return b, errFEN
	}
	for i, row := range rows {
		rank := 7 - i // le FEN commence par la 8e rangée
		file := 0
		for k := 0; k < len(row); k++ {
			c := row[k]
			switch {
			case c >= '1' && c <= '8':
				file += int(c - '0')
			default:
				if file > 7 {
					return b, errFEN
				}
				b[rank*8+file] = c
				file++
			}
		}
		if file != 8 {
			return b, errFEN
		}
	}
	return b, nil
}

func parseSquare(s string) (int, bool) {
	if len(s) != 2 || s[0] < 'a' || s[0] > 'h' || s[1] < '1' || s[1] > '8' {
		return 0, false
	}
	return int(s[1]-'1')*8 + int(s[0]-'a'), true
}

// epCapturable : un pion du trait occupe-t-il une case voisine de celle d'où il
// pourrait prendre en passant ? Pour les Blancs la case ep est en 6e rangée et
// le pion preneur en 5e, donc juste en dessous ; pour les Noirs, l'inverse.
func epCapturable(board [64]byte, ep int, whiteToMove bool) bool {
	var from int
	var pawn byte
	if whiteToMove {
		from, pawn = ep-8, 'P'
	} else {
		from, pawn = ep+8, 'p'
	}
	if from < 0 || from > 63 {
		return false
	}
	file := from % 8
	if file > 0 && board[from-1] == pawn {
		return true
	}
	if file < 7 && board[from+1] == pawn {
		return true
	}
	return false
}
