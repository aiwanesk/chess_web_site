// Package twic tient mega.db à jour avec les livraisons hebdomadaires de
// The Week in Chess.
//
// Le découpage : ce fichier lit du PGN et n'en sort que des données, twic.go
// va chercher les archives et les fait entrer dans la base, schedule.go décide
// QUAND. Chacun se teste sans les deux autres.
package twic

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// rawGame est une partie telle qu'elle sort du fichier : des en-têtes et une
// suite de coups en SAN. La traduction en UCI vient après, parce qu'elle peut
// échouer et qu'il vaut mieux savoir de quelle partie on parle à ce moment-là.
type rawGame struct {
	tags  map[string]string
	moves []string
}

// readGames lit un flux PGN entier et rend tout d'un coup. Pratique pour les
// tests et pour une livraison TWIC ; à réserver aux flux dont on connaît la
// taille. Pour un fichier de plusieurs gigaoctets, voir scanGames.
func readGames(r io.Reader) []rawGame {
	var out []rawGame
	scanGames(r, func(g rawGame) bool { out = append(out, g); return true })
	return out
}

// scanGames lit un flux PGN À LA VOLÉE et appelle fn pour chaque partie.
// Renvoyer false arrête la lecture.
//
// C'est la forme qui compte pour l'indexation de la MegaBase : le PGN source
// pèse plusieurs gigaoctets, et tout charger en mémoire pour en ressortir une
// tranche ne passerait sur aucune machine.
//
// Volontairement tolérant : une partie mal formée ne doit pas faire perdre le
// reste du fichier. Ce qui ne se lit pas est sauté, ce qui se lit entre.
func scanGames(r io.Reader, fn func(rawGame) bool) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 8<<20) // une partie commentée peut être longue

	cur := rawGame{tags: map[string]string{}}
	var movetext strings.Builder
	inMoves := true
	stopped := false

	flush := func() {
		if !stopped && (len(cur.tags) > 0 || movetext.Len() > 0) {
			cur.moves = parseMovetext(movetext.String())
			if len(cur.moves) > 0 && !fn(cur) {
				stopped = true
			}
		}
		cur = rawGame{tags: map[string]string{}}
		movetext.Reset()
		inMoves = false
	}

	for !stopped && sc.Scan() {
		line := strings.TrimSpace(decodeLine(sc.Bytes()))
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") &&
			strings.Contains(line, `"`) {
			// Un en-tête qui suit des coups ouvre la partie suivante : c'est le
			// seul séparateur sur lequel on puisse compter, les lignes vides
			// étant parfois absentes entre deux parties.
			if inMoves {
				flush()
			}
			if k, v, ok := parseTag(line); ok {
				cur.tags[k] = v
			}
			continue
		}
		if line == "" {
			continue
		}
		inMoves = true
		movetext.WriteString(line)
		movetext.WriteByte(' ')
	}
	flush()
}

func parseTag(line string) (key, value string, ok bool) {
	inner := strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
	i := strings.IndexByte(inner, '"')
	if i < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(inner[:i])
	value = strings.TrimSuffix(inner[i+1:], `"`)
	return key, value, key != ""
}

// parseMovetext ne garde que la ligne principale.
//
// Tout le reste — commentaires entre accolades, variantes entre parenthèses,
// NAG, numéros de coups, résultat final — est du décor. Les variantes se jettent
// en comptant la profondeur : TWIC en imbrique, et un simple « jusqu'à la
// prochaine parenthèse » avalerait la fin de la partie.
func parseMovetext(s string) []string {
	var moves []string
	depth := 0
	for i := 0; i < len(s); {
		switch c := s[i]; c {
		case '{':
			if j := strings.IndexByte(s[i:], '}'); j >= 0 {
				i += j + 1
			} else {
				i = len(s)
			}
			continue
		case ';': // commentaire jusqu'à la fin de la ligne
			if j := strings.IndexByte(s[i:], '\n'); j >= 0 {
				i += j + 1
			} else {
				i = len(s)
			}
			continue
		case '(':
			depth++
			i++
			continue
		case ')':
			if depth > 0 {
				depth--
			}
			i++
			continue
		case ' ', '\t', '\r', '\n':
			i++
			continue
		}
		j := i
		for j < len(s) && !strings.ContainsRune(" \t\r\n{}();", rune(s[j])) {
			j++
		}
		tok := s[i:j]
		i = j
		if depth > 0 {
			continue
		}
		if m, ok := cleanMove(tok); ok {
			moves = append(moves, m)
		}
	}
	return moves
}

// cleanMove trie un jeton de la ligne principale : est-ce un coup, ou du
// remplissage ? Le cas tordu est « 12.Nf3 » collé, que TWIC écrit parfois sans
// espace, et « 12... » qui est purement du numérotage.
func cleanMove(tok string) (string, bool) {
	switch tok {
	case "", "1-0", "0-1", "1/2-1/2", "*", "--", "Z0":
		return "", false
	}
	if tok[0] == '$' {
		return "", false
	}
	if tok[0] >= '0' && tok[0] <= '9' {
		k := 0
		for k < len(tok) && tok[k] >= '0' && tok[k] <= '9' {
			k++
		}
		for k < len(tok) && tok[k] == '.' {
			k++
		}
		tok = tok[k:]
		if tok == "" {
			return "", false
		}
	}
	tok = strings.TrimRight(tok, "!?")
	if tok == "" {
		return "", false
	}
	return tok, true
}

// convert traduit une partie brute en partie prête à insérer. Une erreur ici
// veut dire que le SAN ne se rejoue pas : la partie est sautée, pas devinée.
func convert(g rawGame) (games.ImportGame, error) {
	out := games.ImportGame{
		White:    strings.TrimSpace(g.tags["White"]),
		Black:    strings.TrimSpace(g.tags["Black"]),
		Event:    strings.TrimSpace(g.tags["Event"]),
		Date:     strings.TrimSpace(g.tags["Date"]),
		ECO:      strings.TrimSpace(g.tags["ECO"]),
		Result:   resultOf(g.tags["Result"]),
		WhiteElo: atoi(g.tags["WhiteElo"]),
		BlackElo: atoi(g.tags["BlackElo"]),
	}
	if variant(g) {
		return games.ImportGame{}, errVariant
	}
	fen := startFEN
	for _, san := range g.moves {
		uci, err := corpus.SANToUCI(fen, san)
		if err != nil {
			return games.ImportGame{}, err
		}
		next, err := corpus.ApplyUCI(fen, uci)
		if err != nil {
			return games.ImportGame{}, err
		}
		out.UCI = append(out.UCI, uci)
		out.SAN = append(out.SAN, san)
		fen = next
	}
	return out, nil
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// errVariant écarte ce qui n'est pas une partie d'échecs standard.
//
// TWIC publie du Chess960 dans les mêmes fichiers : une bonne partie des
// Kasparov-Topalov de 2026 commencent sur « rkbnnbqr/… ». Ces parties ne sont
// pas mal écrites, elles sont d'un AUTRE jeu — le roque n'y obéit pas aux mêmes
// règles, et surtout elles n'ont rien à faire dans un arbre d'ouvertures
// construit depuis la position initiale : mélangées aux vraies, elles feraient
// croire à de la théorie là où il n'y en a pas.
//
// Même raisonnement pour toute partie commencée depuis une position imposée :
// elle ne se raccroche à aucune branche.
var errVariant = errors.New("partie hors échecs standard")

func variant(g rawGame) bool {
	if f := strings.TrimSpace(g.tags["FEN"]); f != "" && f != startFEN {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(g.tags["Variant"])) {
	case "", "standard", "normal", "chess":
		return false
	}
	return true
}

func resultOf(s string) int {
	switch strings.TrimSpace(s) {
	case "1-0":
		return 1
	case "0-1":
		return -1
	}
	return 0
}

func atoi(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return n
}

// decodeLine rattrape l'encodage des noms.
//
// TWIC et les exports de MegaBase publient en Latin-1 : « Vachier-Lagrave »
// passe, mais « Grünfeld » arrive en octet isolé, que Go recopierait tel quel.
// Un nom mal encodé n'est pas cosmétique ici — il ne se retrouve plus à la
// recherche, donc la partie devient invisible.
//
// Le choix se fait LIGNE PAR LIGNE, et pas sur le fichier entier. Un fichier
// presque entièrement en UTF-8 avec trois noms en Latin-1 serait déclaré
// invalide dans son ensemble, et la conversion transformerait alors tous les
// vrais caractères UTF-8 en charabia. Ligne par ligne, seules les lignes
// fautives sont converties — et ça se lit en flux.
func decodeLine(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var out strings.Builder
	out.Grow(len(b) * 2)
	for _, c := range b {
		out.WriteRune(rune(c))
	}
	return out.String()
}
