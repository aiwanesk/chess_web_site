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

	// TODO(next): analyse each game with Stockfish (UCI), detect played/missed
	// combinations via eval-swing + only-move-gap thresholds, score "beauty",
	// take top 10, build anonymised puzzles via tactics.NewPuzzle(...), and write
	// content/tactiques/<year>-S<week>.json for the frontend to render.
	fmt.Println("\nDétection Stockfish : à implémenter (prochaine étape).")
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
