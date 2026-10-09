// Package exercises transforme les imprécisions d'Alexandre dans ses parties
// en exercices : la position avant le coup, et les coups que Stockfish juge
// bons à profondeur 30.
//
// L'analyse tourne sur le PC (cmd/mistakes) : la profondeur 30 se compte en
// secondes par position, pas question de la faire sur le serveur du site. Le
// serveur ne fait que stocker les exercices et corriger les réponses, à partir
// de coups déjà évalués — il n'a pas de moteur.
package exercises

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/online"
	"github.com/iwanesko/chess-web-site/backend/internal/tactics"
)

// Engine est ce dont l'analyse a besoin d'un moteur : une recherche à
// profondeur fixe, éventuellement restreinte à certains coups. Stockfish
// (tactics.Stockfish) l'implémente ; les tests utilisent un faux moteur.
type Engine interface {
	Search(fen string, multipv, depth int, moves ...string) ([]tactics.Line, error)
}

// Options règle l'analyse. Les zéros prennent les valeurs par défaut.
type Options struct {
	QuickDepth int // passe de repérage (défaut 18)
	Depth      int // passe de confirmation (défaut 30)
	MultiPV    int // coups gardés à la confirmation (défaut 3)
	Threshold  int // perte minimale pour un exercice, en centipions (défaut 30)
	Tolerance  int // un coup à moins de ça du meilleur est juste (défaut 15)
	SkipMoves  int // premiers coups ignorés, la théorie (défaut 4)
}

func (o Options) withDefaults() Options {
	if o.QuickDepth <= 0 {
		o.QuickDepth = 18
	}
	if o.Depth <= 0 {
		o.Depth = 30
	}
	if o.MultiPV <= 0 {
		o.MultiPV = 3
	}
	if o.Threshold <= 0 {
		o.Threshold = 30
	}
	if o.Tolerance <= 0 {
		o.Tolerance = 15
	}
	if o.SkipMoves < 0 {
		o.SkipMoves = 0
	} else if o.SkipMoves == 0 {
		o.SkipMoves = 4
	}
	return o
}

// decided : au-delà de cinq pions, la position est jouée. Perdre 0,4 en passant
// de +8 à +7,6 n'apprend rien — sauf si le résultat change de camp.
const decided = 500

// quickMargin : la passe rapide laisse passer un peu moins que le seuil, pour
// ne pas écarter un coup que la profondeur 30 jugerait finalement fautif.
const quickMargin = 10

// AnalyzeGame renvoie les exercices d'une partie, du point de vue du joueur
// dont le pseudo est me. progress, s'il n'est pas nil, est appelé à chaque
// position confirmée à profondeur 30 — c'est là que passe le temps.
func AnalyzeGame(eng Engine, g online.Game, me string, opt Options, progress func(ply int)) ([]Exercise, error) {
	opt = opt.withDefaults()
	white := strings.EqualFold(g.White, me)
	if !white && !strings.EqualFold(g.Black, me) {
		return nil, fmt.Errorf("exercises: %s n'a pas joué %s", me, g.URL)
	}
	colour := "b"
	if white {
		colour = "w"
	}

	var out []Exercise
	fen := startFEN
	for ply, uci := range g.UCI {
		mine := (ply%2 == 0) == white
		if mine && ply/2 >= opt.SkipMoves {
			ex, ok, err := analyzeMove(eng, fen, uci, opt)
			if err != nil {
				return nil, fmt.Errorf("exercises: %s, demi-coup %d : %w", g.URL, ply+1, err)
			}
			if ok {
				ex.ID = exerciseID(g.URL, ply)
				ex.GameURL, ex.Source = g.URL, string(g.Source)
				ex.White, ex.Black, ex.Colour = g.White, g.Black, colour
				ex.Speed, ex.Date, ex.Played = g.Speed, g.Date, g.Played
				ex.Ply = ply
				ex.Phase = phaseOf(fen, ply)
				if ply < len(g.SAN) {
					ex.PlayedSAN = g.SAN[ply]
				}
				out = append(out, ex)
				if progress != nil {
					progress(ply)
				}
			}
		}
		next, err := corpus.ApplyUCI(fen, uci)
		if err != nil {
			break // la suite ne se rejoue pas : on garde ce qui précède
		}
		fen = next
	}
	return out, nil
}

// analyzeMove décide si le coup joué en fen est une imprécision. Passe rapide
// d'abord ; la profondeur 30 n'est payée que pour les candidats.
func analyzeMove(eng Engine, fen, played string, opt Options) (Exercise, bool, error) {
	quick, err := eng.Search(fen, 1, opt.QuickDepth)
	if err != nil || len(quick) == 0 || len(quick[0].PV) == 0 {
		return Exercise{}, false, err
	}
	if quick[0].PV[0] == played {
		return Exercise{}, false, nil // il a joué le coup du moteur
	}
	mine, err := eng.Search(fen, 1, opt.QuickDepth, played)
	if err != nil || len(mine) == 0 {
		return Exercise{}, false, err
	}
	if !worthIt(quick[0].Score(), mine[0].Score(), opt.Threshold-quickMargin) {
		return Exercise{}, false, nil
	}

	// Confirmation à profondeur 30, en gardant les trois meilleurs coups : en
	// positionnel, plusieurs coups se valent souvent à 0,1 près.
	top, err := eng.Search(fen, opt.MultiPV, opt.Depth)
	if err != nil || len(top) == 0 {
		return Exercise{}, false, err
	}
	best := top[0].Score()
	playedScore, playedLine, inTop := 0, tactics.Line{}, false
	for _, l := range top {
		if len(l.PV) > 0 && l.PV[0] == played {
			playedScore, playedLine, inTop = l.Score(), l, true
		}
	}
	if !inTop {
		deep, err := eng.Search(fen, 1, opt.Depth, played)
		if err != nil || len(deep) == 0 {
			return Exercise{}, false, err
		}
		playedScore, playedLine = deep[0].Score(), deep[0]
	}
	if !worthIt(best, playedScore, opt.Threshold) {
		return Exercise{}, false, nil
	}

	reached := top[0].Depth
	if reached == 0 {
		reached = opt.Depth
	}
	ex := Exercise{FEN: fen, PlayedUCI: played, Depth: reached, Tolerance: opt.Tolerance,
		Best: best, PlayedScore: playedScore, Loss: best - playedScore}
	for _, l := range top {
		if len(l.PV) == 0 || l.PV[0] == played {
			continue
		}
		ex.Moves = append(ex.Moves, candidate(fen, l, best, opt.Tolerance))
	}
	pc := candidate(fen, playedLine, best, opt.Tolerance)
	pc.UCI, pc.OK = played, false // le coup de la partie n'est jamais la réponse
	if san, err := corpus.UCIToSAN(fen, played); err == nil {
		pc.SAN = san
	}
	ex.Moves = append(ex.Moves, pc)
	return ex, true, nil
}

// worthIt : la perte dépasse le seuil, et la position n'était pas déjà jouée
// d'avance dans le même sens avant comme après le coup.
func worthIt(best, played, threshold int) bool {
	if best-played < threshold {
		return false
	}
	if best >= decided && played >= decided {
		return false
	}
	if best <= -decided && played <= -decided {
		return false
	}
	return true
}

func candidate(fen string, l tactics.Line, best, tolerance int) Candidate {
	c := Candidate{CP: l.CP, Mate: l.Mate, Score: l.Score()}
	if len(l.PV) > 0 {
		c.UCI = l.PV[0]
		if san, err := corpus.UCIToSAN(fen, c.UCI); err == nil {
			c.SAN = san
		}
		c.Line = sanLine(fen, l.PV, 8)
	}
	c.OK = best-c.Score <= tolerance
	return c
}

// sanLine écrit le début de la variante en SAN, pour la correction.
func sanLine(fen string, pv []string, max int) []string {
	var out []string
	for i, uci := range pv {
		if i >= max {
			break
		}
		san, err := corpus.UCIToSAN(fen, uci)
		if err != nil {
			break
		}
		out = append(out, san)
		if fen, err = corpus.ApplyUCI(fen, uci); err != nil {
			break
		}
	}
	return out
}

// phaseOf : ouverture tant qu'on est dans les dix premiers coups, finale quand
// il reste peu de pièces, milieu de jeu sinon.
func phaseOf(fen string, ply int) string {
	if ply < 20 {
		return "opening"
	}
	value := map[rune]int{'q': 9, 'r': 5, 'b': 3, 'n': 3}
	total, queens := 0, 0
	for _, c := range strings.Fields(fen)[0] {
		lc := c | 0x20
		total += value[lc]
		if lc == 'q' {
			queens++
		}
	}
	if total <= 20 || (queens == 0 && total <= 26) {
		return "endgame"
	}
	return "middlegame"
}

// exerciseID est stable : réanalyser la même partie remplace l'exercice au lieu
// de le dupliquer, et garde son étiquette « théorique ».
func exerciseID(gameURL string, ply int) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%s#%d", gameURL, ply)))
	return hex.EncodeToString(h[:8])
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
