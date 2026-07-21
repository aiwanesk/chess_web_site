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
	sacWindowPlies   = 6   // only the forcing start counts for sacrifice detection
	sacThreshold     = 250 // material must drop ≥ a minor piece to count as a sac
)

func truncatePV(pv []string, max int) []string {
	if len(pv) > max {
		return append([]string{}, pv[:max]...)
	}
	return pv
}

// Tactic is a detected combination in ORIGINAL orientation (side to move is
// Alexandre's colour). It is never published as-is — NewPuzzle mirrors it.
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

// AnalyzeGame replays one game and returns the tactics found in Alexandre's
// moves (played and missed).
func AnalyzeGame(ev Evaluator, g Game) ([]Tactic, error) {
	pgnFn, err := chess.PGN(strings.NewReader(g.PGN))
	if err != nil {
		return nil, err
	}
	game := chess.NewGame(pgnFn)
	positions := game.Positions()
	moves := game.Moves()
	mine := colorOf(g.MyColor)
	uci := chess.UCINotation{}

	var out []Tactic
	for i, mv := range moves {
		pos := positions[i]
		if pos.Turn() != mine {
			continue
		}
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
		sac := detectSacrifice(fen, best.PV)
		// The first move must be a genuine tactical blow: a CHECK, or the start
		// of a SACRIFICE. A plain material-winning capture ("just take the free
		// piece") or a quiet move is rejected — those aren't interesting puzzles.
		if !firstMoveQualifies(fen, best.PV[0], sac) {
			continue
		}
		// Only publish the FORCING part of the line: cut it as soon as the winning
		// side would have to play a non-forcing (quiet) move. Keeps solutions short.
		sol := truncatePV(forcingLine(fen, best.PV), maxSolutionPlies)
		if len(sol) == 0 {
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
// Accepted: a CHECK, or the first move of a SACRIFICE. Rejected: a plain
// material-winning capture (avoid "just take the free piece") and quiet moves.
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

func colorOf(c string) chess.Color {
	if strings.EqualFold(c, "black") {
		return chess.Black
	}
	return chess.White
}

var pieceValue = map[chess.PieceType]int{
	chess.Pawn: 100, chess.Knight: 320, chess.Bishop: 330,
	chess.Rook: 500, chess.Queen: 900, chess.King: 0,
}

// detectSacrifice replays the PV and reports whether Alexandre's material dips
// at least a pawn below its starting value along the forced line (i.e. he
// invested material) — the marker of a genuine sacrifice/combination.
func detectSacrifice(fen string, pv []string) bool {
	opt, err := chess.FEN(fen)
	if err != nil {
		return false
	}
	game := chess.NewGame(opt)
	us := game.Position().Turn()
	start := material(game.Position().Board(), us)
	minMat := start
	uci := chess.UCINotation{}
	// Only the forcing start of the line: a real sacrifice invests material
	// early; deep exchanges later are just normal play, not a sac.
	for i, mv := range pv {
		if i >= sacWindowPlies {
			break
		}
		m, err := uci.Decode(game.Position(), mv)
		if err != nil {
			break
		}
		if err := game.Move(m); err != nil {
			break
		}
		if mat := material(game.Position().Board(), us); mat < minMat {
			minMat = mat
		}
	}
	return minMat <= start-sacThreshold
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
