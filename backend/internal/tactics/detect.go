package tactics

import (
	"sort"
	"strings"

	"github.com/notnil/chess"
)

// Detection thresholds (centipawns, side-to-move POV). Tunable on real data.
const (
	gapCP      = 150 // best must beat 2nd-best by this → a real "only move"
	winCP      = 300 // best must reach this (or mate) → decisive
	balancedCP = 130 // 2nd-best must stay under this → the shot CREATES the win
	missCP     = 200 // best − played ≥ this → he missed it
	foundEps   = 40  // best − played ≤ this → he found it

	maxSolutionPlies = 7   // published solution length cap — keep puzzles SHORT
	minSolutionPlies = 3   // drop trivial one-movers (mate-in-1 / single winning check)
	sacThreshold     = 200 // end-of-line material deficit (cp) to count as a real sac
)

func truncatePV(pv []string, max int) []string {
	if len(pv) > max {
		return append([]string{}, pv[:max]...)
	}
	return pv
}

// Tactic is a detected combination in ORIGINAL orientation (from either side of
// one of Alexandre's games). It is never published as-is — NewPuzzle mirrors it.
type Tactic struct {
	FEN       string
	Solution  []string // UCI, the winning line (engine PV)
	Kind      string   // "played" | "missed"
	Mate      bool
	Sacrifice bool
	Beauty    int
	Source    string
}

// classify decides whether a position is a tactical point and, if so, whether
// Alexandre found or missed it. Pure function → unit-tested without an engine.
func classify(best, second, played Line) (isTactic bool, kind string, mate bool) {
	b, s, p := best.Score(), second.Score(), played.Score()

	decisive := best.Mate > 0 || b >= winCP
	onlyMove := (b - s) >= gapCP
	createsWin := s <= balancedCP
	if !(decisive && onlyMove && createsWin) {
		return false, "", false
	}
	mate = best.Mate > 0
	switch {
	case b-p <= foundEps:
		return true, "played", mate
	case b-p >= missCP:
		return true, "missed", mate
	default:
		return false, "", false // borderline → skip to keep the set clean
	}
}

// AnalyzeGame replays one game and returns the tactics found in it — for BOTH
// sides (Alexandre's shots and the opponent's), played and missed.
func AnalyzeGame(ev Evaluator, g Game) ([]Tactic, error) {
	pgnFn, err := chess.PGN(strings.NewReader(g.PGN))
	if err != nil {
		return nil, err
	}
	game := chess.NewGame(pgnFn)
	positions := game.Positions()
	moves := game.Moves()
	uci := chess.UCINotation{}

	var out []Tactic
	for i, mv := range moves {
		pos := positions[i]
		fen := pos.String()
		playedUCI := uci.Encode(pos, mv)

		lines, err := ev.Lines(fen, 2)
		if err != nil || len(lines) == 0 {
			continue
		}
		best := lines[0]
		second := best
		if len(lines) > 1 {
			second = lines[1]
		}
		played, err := ev.Move(fen, playedUCI)
		if err != nil {
			continue
		}

		isTactic, kind, mate := classify(best, second, played)
		if !isTactic {
			continue
		}
		if len(best.PV) == 0 {
			continue
		}
		// Publish only the FORCING line, ending on the solver's decisive blow.
		// Skip one-movers — too trivial for the target audience.
		sol := truncatePV(forcingLine(fen, best.PV), maxSolutionPlies)
		if len(sol) < minSolutionPlies {
			continue
		}
		sac := detectSacrifice(fen, sol)
		// The first move must be a genuine tactical blow: a CHECK (that doesn't
		// simply grab the queen), or the start of a real SACRIFICE. Plain
		// material-winning captures and quiet moves are rejected.
		if !firstMoveQualifies(fen, sol[0], sac) {
			continue
		}
		out = append(out, Tactic{
			FEN:       fen,
			Solution:  sol,
			Kind:      kind,
			Mate:      mate,
			Sacrifice: sac,
			Beauty:    beautyScore(best, sac),
			Source:    g.Source,
		})
	}

	// At most ONE position per game: keep the most beautiful.
	if len(out) == 0 {
		return nil, nil
	}
	top := out[0]
	for _, t := range out[1:] {
		if t.Beauty > top.Beauty {
			top = t
		}
	}
	return []Tactic{top}, nil
}

// firstMoveQualifies reports whether the solution may START with this move.
// Accepted: a CHECK, or the first move of a real SACRIFICE. Rejected: plain
// material-winning captures, quiet moves, and — always — grabbing the opponent's
// QUEEN on move one ("just take the queen" is never an interesting puzzle).
func firstMoveQualifies(fen, uciMove string, sac bool) bool {
	opt, err := chess.FEN(fen)
	if err != nil {
		return false
	}
	game := chess.NewGame(opt)
	uci := chess.UCINotation{}
	target, err := uci.Decode(game.Position(), uciMove)
	if err != nil {
		return false
	}
	// Reject capturing a queen on the first move, even with check.
	if p, ok := game.Position().Board().SquareMap()[target.S2()]; ok && p.Type() == chess.Queen {
		return false
	}
	for _, m := range game.ValidMoves() {
		if m.S1() == target.S1() && m.S2() == target.S2() {
			if m.HasTag(chess.Check) {
				return true // a check is always a valid forcing start
			}
			return sac // otherwise only a sacrifice qualifies (no plain captures)
		}
	}
	return false
}

// forcingLine returns the leading, forcing portion of a PV and — crucially —
// makes the puzzle END on the SOLVER's decisive blow: a capture, a check, a
// promotion, or mate (never a quiet move, never the opponent's passive reply).
//
// It walks solver moves + forced opponent replies, and returns the line cut at
// the LAST solver move that was a capture/check/promotion (or mate). If no such
// finishing move exists, it returns nil so the fizzling line is dropped.
func forcingLine(fen string, pv []string) []string {
	opt, err := chess.FEN(fen)
	if err != nil {
		return nil
	}
	game := chess.NewGame(opt)
	uci := chess.UCINotation{}

	lastFinishing := 0 // length ending on the last solver capture/check/promotion/mate
	for i, u := range pv {
		if i >= maxSolutionPlies {
			break
		}
		dec, err := uci.Decode(game.Position(), u)
		if err != nil {
			break
		}
		// uci.Decode does NOT set move tags; only ValidMoves() carry them.
		var m *chess.Move
		for _, vm := range game.ValidMoves() {
			if vm.S1() == dec.S1() && vm.S2() == dec.S2() && vm.Promo() == dec.Promo() {
				m = vm
				break
			}
		}
		if m == nil {
			break
		}
		solverMove := i%2 == 0
		forcing := m.HasTag(chess.Check) || m.HasTag(chess.Capture) ||
			m.HasTag(chess.EnPassant) || m.Promo() != chess.NoPieceType
		// After the first move, a quiet solver move ends the forcing sequence.
		if solverMove && i > 0 && !forcing {
			break
		}
		if err := game.Move(m); err != nil {
			break
		}
		mate := game.Method() == chess.Checkmate
		if solverMove && (forcing || mate) {
			lastFinishing = i + 1 // end the puzzle here — on the solver's blow
		}
		if mate {
			break
		}
	}
	if lastFinishing < 1 {
		return nil
	}
	return append([]string{}, pv[:lastFinishing]...)
}

// TopPuzzles analyses all games, ranks tactics by beauty and returns the top n
// as anonymised (mirrored) puzzles, de-duplicated by puzzle ID.
func TopPuzzles(ev Evaluator, games []Game, n int) []Puzzle {
	var all []Tactic
	for _, g := range games {
		ts, err := AnalyzeGame(ev, g)
		if err != nil {
			continue
		}
		all = append(all, ts...)
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].Beauty > all[j].Beauty })

	seen := map[string]bool{}
	var puzzles []Puzzle
	for _, t := range all {
		p := NewPuzzle(t.FEN, t.Solution, t.Kind, t.Source, t.Mate, t.Sacrifice, t.Beauty)
		if seen[p.ID] {
			continue
		}
		seen[p.ID] = true
		puzzles = append(puzzles, p)
		if len(puzzles) >= n {
			break
		}
	}
	return puzzles
}

func beautyScore(best Line, sacrifice bool) int {
	score := 0
	if best.Mate > 0 {
		score += 25 - min(best.Mate, 12) // faster mate = prettier
	}
	cp := best.CP
	if best.Mate > 0 {
		cp = 1200
	}
	score += min(cp, 1500) / 100 // up to +15 for material swing
	if sacrifice {
		score += 12
	}
	score += min(len(best.PV), 12) // longer forced line
	return score
}

var pieceValue = map[chess.PieceType]int{
	chess.Pawn: 100, chess.Knight: 320, chess.Bishop: 330,
	chess.Rook: 500, chess.Queen: 900, chess.King: 0,
}

// detectSacrifice reports whether the solver ends the forcing line DOWN material
// yet winning — the mark of a genuine sacrifice. It compares the material balance
// (solver minus opponent) at the start vs the end of the solution, so ordinary
// trades and material-winning combinations are NOT mislabelled as sacrifices
// (e.g. capturing a queen and being recaptured stays balanced → not a sac).
func detectSacrifice(fen string, solution []string) bool {
	opt, err := chess.FEN(fen)
	if err != nil {
		return false
	}
	game := chess.NewGame(opt)
	us := game.Position().Turn()
	them := us.Other()
	balance := func() int {
		b := game.Position().Board()
		return material(b, us) - material(b, them)
	}
	start := balance()
	uci := chess.UCINotation{}
	for _, mv := range solution {
		m, err := uci.Decode(game.Position(), mv)
		if err != nil {
			break
		}
		if err := game.Move(m); err != nil {
			break
		}
	}
	return balance() <= start-sacThreshold
}

func material(b *chess.Board, c chess.Color) int {
	total := 0
	for _, p := range b.SquareMap() {
		if p.Color() == c {
			total += pieceValue[p.Type()]
		}
	}
	return total
}
