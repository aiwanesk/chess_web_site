// Command tactics builds the weekly anonymised tactics puzzles from Alexandre's
// own games. Run locally:  go run ./cmd/tactics
//
// Reads usernames from the environment (or ../.env locally). It NEVER writes any
// username to its output — only mirrored puzzles (see internal/tactics).
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/tactics"
)

func main() {
	loadDotEnv("../.env")

	cfg := tactics.LoadConfig()
	if cfg.ChessComUser == "" && cfg.LichessUser == "" {
		log.Fatal("set CHESSCOM_USERNAME and/or LICHESS_USERNAME (in .env or env)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	games, err := tactics.FetchWeek(ctx, cfg, 7)
	if err != nil {
		log.Fatalf("fetch: %v", err)
	}

	byClass := map[string]int{}
	for _, g := range games {
		byClass[g.TimeClass]++
	}
	classes := make([]string, 0, len(byClass))
	for k := range byClass {
		classes = append(classes, k)
	}
	sort.Strings(classes)

	fmt.Printf("Parties récupérées (7 derniers jours) : %d\n", len(games))
	for _, c := range classes {
		fmt.Printf("  %-10s %d\n", c, byClass[c])
	}

	// Detection needs Stockfish. Skip cleanly if it isn't available.
	enginePath := os.Getenv("STOCKFISH_PATH")
	if enginePath == "" {
		enginePath = "stockfish"
	}
	moveTime := 200 // ms per position; raise for stronger analysis
	engine, err := tactics.NewStockfish(enginePath, moveTime)
	if err != nil {
		fmt.Printf("\nStockfish introuvable (%q) : détection sautée.\n", enginePath)
		fmt.Println("Installe Stockfish ou définis STOCKFISH_PATH, puis relance.")
		return
	}
	defer engine.Close()

	fmt.Println("\nAnalyse Stockfish en cours…")
	puzzles := tactics.TopPuzzles(engine, games, 10)

	y, w := time.Now().ISOWeek()
	week := fmt.Sprintf("%04d-S%02d", y, w)
	outDir := os.Getenv("TACTICS_DIR")
	if outDir == "" {
		outDir = "../content/tactiques"
	}
	// Defense-in-depth: never write a file containing a username.
	path, err := tactics.WritePuzzles(outDir, week, puzzles, []string{cfg.ChessComUser, cfg.LichessUser})
	if err != nil {
		log.Fatalf("write: %v", err)
	}
	fmt.Printf("%d tactiques écrites → %s\n", len(puzzles), path)
}

// loadDotEnv is a minimal .env loader for local runs (KEY=value lines). It never
// overrides an already-set env var. Silent if the file is absent.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k, v = strings.TrimSpace(k), strings.TrimSpace(v)
		if _, exists := os.LookupEnv(k); !exists {
			_ = os.Setenv(k, v)
		}
	}
}
