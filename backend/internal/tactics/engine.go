package tactics

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

// Line is one engine principal variation: an evaluation from the side-to-move's
// perspective plus the move sequence (UCI).
type Line struct {
	CP   int      // centipawns (side to move); ignored if Mate != 0
	Mate int      // >0 = side to move mates in N; <0 = gets mated in N; 0 = none
	PV   []string // principal variation, UCI moves
}

// Score collapses CP/Mate into a single comparable integer (side-to-move POV).
func (l Line) Score() int {
	switch {
	case l.Mate > 0:
		return 100000 - l.Mate
	case l.Mate < 0:
		return -100000 - l.Mate
	default:
		return l.CP
	}
}

// Evaluator is the minimal surface the detector needs. Stockfish implements it;
// tests use a mock so detection logic is verified without a real engine.
type Evaluator interface {
	Lines(fen string, multipv int) ([]Line, error) // top-N moves
	Move(fen, uci string) (Line, error)            // eval of one specific move
}

// Stockfish drives a Stockfish process over UCI.
type Stockfish struct {
	mu       sync.Mutex
	cmd      *exec.Cmd
	stdin    io.Writer
	out      *bufio.Scanner
	moveTime int // ms per position
}

// NewStockfish starts the engine at path (e.g. "stockfish" on PATH).
func NewStockfish(path string, moveTimeMS int) (*Stockfish, error) {
	cmd := exec.Command(path)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	s := &Stockfish{cmd: cmd, stdin: stdin, out: bufio.NewScanner(stdout), moveTime: moveTimeMS}
	s.out.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	s.send("uci")
	s.waitFor("uciok")
	s.send("setoption name UCI_AnalyseMode value true")
	s.send("isready")
	s.waitFor("readyok")
	return s, nil
}

func (s *Stockfish) Close() error {
	s.send("quit")
	return s.cmd.Wait()
}

func (s *Stockfish) send(cmd string) {
	_, _ = io.WriteString(s.stdin, cmd+"\n")
}

func (s *Stockfish) waitFor(token string) {
	for s.out.Scan() {
		if strings.HasPrefix(s.out.Text(), token) {
			return
		}
	}
}

// Lines returns the top-`multipv` moves for the position.
func (s *Stockfish) Lines(fen string, multipv int) ([]Line, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.send(fmt.Sprintf("setoption name MultiPV value %d", multipv))
	s.send("position fen " + fen)
	s.send(fmt.Sprintf("go movetime %d", s.moveTime))
	return s.readLines(multipv)
}

// Move returns the evaluation of a single move (searchmoves).
func (s *Stockfish) Move(fen, uci string) (Line, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.send("setoption name MultiPV value 1")
	s.send("position fen " + fen)
	s.send(fmt.Sprintf("go movetime %d searchmoves %s", s.moveTime, uci))
	lines, err := s.readLines(1)
	if err != nil || len(lines) == 0 {
		return Line{}, err
	}
	return lines[0], nil
}

// readLines consumes engine output until "bestmove", keeping the deepest info
// line per multipv index.
func (s *Stockfish) readLines(multipv int) ([]Line, error) {
	best := map[int]Line{}
	for s.out.Scan() {
		text := s.out.Text()
		if strings.HasPrefix(text, "bestmove") {
			break
		}
		if !strings.HasPrefix(text, "info ") || !strings.Contains(text, " pv ") {
			continue
		}
		idx, line, ok := parseInfo(text)
		if ok {
			best[idx] = line // later (deeper) lines overwrite shallower ones
		}
	}
	out := make([]Line, 0, multipv)
	for i := 1; i <= multipv; i++ {
		if l, ok := best[i]; ok {
			out = append(out, l)
		}
	}
	if err := s.out.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// parseInfo extracts (multipv index, Line) from a UCI "info ... pv ..." line.
func parseInfo(text string) (int, Line, bool) {
	toks := strings.Fields(text)
	idx := 1
	var line Line
	hasScore := false
	for i := 0; i < len(toks); i++ {
		switch toks[i] {
		case "multipv":
			if i+1 < len(toks) {
				idx, _ = strconv.Atoi(toks[i+1])
			}
		case "score":
			if i+2 < len(toks) {
				kind, val := toks[i+1], toks[i+2]
				n, _ := strconv.Atoi(val)
				if kind == "mate" {
					line.Mate = n
				} else {
					line.CP = n
				}
				hasScore = true
			}
		case "pv":
			line.PV = append([]string{}, toks[i+1:]...)
			return idx, line, hasScore
		}
	}
	return idx, line, hasScore && line.PV != nil
}
