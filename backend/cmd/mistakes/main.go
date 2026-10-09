// Command mistakes analyse les dernières parties d'Alexandre avec Stockfish et
// envoie ses imprécisions à l'espace privé du site, où elles deviennent des
// exercices (/admin/exercices/).
//
// À lancer sur le PC, pas sur le serveur : la profondeur 30 coûte plusieurs
// secondes par position.
//
//	go build -o mistakes.exe ./cmd/mistakes
//	mistakes.exe                       # les 100 dernières parties blitz et rapides
//	mistakes.exe -games 20 -depth 26   # plus court
//	mistakes.exe -pending              # seulement les coups proposés en exercice
//
// Lit ../.env : CHESSCOM_USERNAME, LICHESS_USERNAME, ADMIN_USER, ADMIN_TOKEN,
// et au besoin STOCKFISH_PATH et SITE_URL. Relancer ne refait jamais une partie
// déjà analysée : le site tient la liste. Un Ctrl+C perd au plus la partie en
// cours.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/exercises"
	"github.com/iwanesko/chess-web-site/backend/internal/online"
	"github.com/iwanesko/chess-web-site/backend/internal/tactics"
)

func main() {
	loadDotEnv("../.env")
	nGames := flag.Int("games", 100, "nombre de parties récentes à considérer")
	speeds := flag.String("speeds", "blitz,rapid", "cadences, séparées par des virgules")
	depth := flag.Int("depth", 30, "profondeur de confirmation")
	quick := flag.Int("quick", 18, "profondeur de la passe de repérage")
	threshold := flag.Int("threshold", 30, "perte minimale pour un exercice (centipions)")
	tolerance := flag.Int("tolerance", 15, "écart au meilleur coup encore accepté (centipions)")
	skip := flag.Int("skip", 4, "premiers coups ignorés (théorie)")
	threads := flag.Int("threads", max(1, runtime.NumCPU()-2), "fils de Stockfish")
	hash := flag.Int("hash", 8192, "table de hachage de Stockfish (Mo)")
	maxTime := flag.Int("maxtime", 30, "plafond par recherche en secondes (0 = aucun)")
	site := flag.String("site", env("SITE_URL", "https://iwanesko.ch"), "adresse du site")
	pendingOnly := flag.Bool("pending", false, "ne traiter que la file des coups proposés")
	flag.Parse()

	api := &siteAPI{base: strings.TrimRight(*site, "/") + "/admin/exercices/api",
		user: os.Getenv("ADMIN_USER"), pass: os.Getenv("ADMIN_TOKEN")}
	if api.pass == "" {
		log.Fatal("ADMIN_TOKEN manquant (dans ../.env) : impossible d'envoyer les exercices au site")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	sf, err := tactics.NewStockfish(env("STOCKFISH_PATH", "stockfish"), 0)
	if err != nil {
		log.Fatalf("Stockfish introuvable (STOCKFISH_PATH) : %v", err)
	}
	defer sf.Close()
	sf.Configure(*threads, *hash)
	sf.SetMaxTime(*maxTime * 1000)
	log.Printf("Stockfish : %d fils, %d Mo de hachage, profondeur %d (plafond %d s par recherche)",
		*threads, *hash, *depth, *maxTime)

	if err := resolvePending(ctx, api, sf, *depth); err != nil {
		log.Printf("file des coups proposés : %v", err)
	}
	if *pendingOnly {
		return
	}

	games, err := recentGames(ctx, *nGames, *speeds)
	if err != nil {
		log.Fatal(err)
	}
	done := map[string]bool{}
	var analyzed []string
	if err := api.call(ctx, "GET", "/games", nil, &analyzed); err != nil {
		log.Fatalf("site : %v", err)
	}
	for _, u := range analyzed {
		done[u] = true
	}
	var todo []online.Game
	for _, g := range games {
		if !done[g.URL] {
			todo = append(todo, g)
		}
	}
	log.Printf("%d parties récentes, %d déjà analysées, %d à faire", len(games), len(games)-len(todo), len(todo))

	opt := exercises.Options{QuickDepth: *quick, Depth: *depth, Threshold: *threshold,
		Tolerance: *tolerance, SkipMoves: *skip}
	total, start := 0, time.Now()
	for i, g := range todo {
		if ctx.Err() != nil {
			log.Print("interrompu : relancer reprendra à la partie suivante")
			return
		}
		me := os.Getenv("LICHESS_USERNAME")
		if g.Source == online.ChessCom {
			me = os.Getenv("CHESSCOM_USERNAME")
		}
		t0 := time.Now()
		fmt.Printf("[%d/%d] %s %s – %s (%s) ", i+1, len(todo), g.Date, g.White, g.Black, g.Speed)
		list, err := exercises.AnalyzeGame(sf, g, me, opt, func(int) { fmt.Print(".") })
		if err != nil {
			fmt.Println()
			log.Printf("partie sautée : %v", err)
			continue
		}
		if err := api.call(ctx, "POST", "/import", map[string]any{"game": g.URL, "exercises": list}, nil); err != nil {
			fmt.Println()
			log.Fatalf("envoi au site : %v", err)
		}
		total += len(list)
		fmt.Printf(" %d exercice(s) en %s\n", len(list), time.Since(t0).Round(time.Second))
	}
	log.Printf("terminé : %d exercices en %s", total, time.Since(start).Round(time.Second))
}

// recentGames : les dernières parties des deux comptes, mélangées par date.
func recentGames(ctx context.Context, n int, speeds string) ([]online.Game, error) {
	want := map[string]bool{}
	for _, s := range strings.Split(speeds, ",") {
		if s = strings.TrimSpace(s); s != "" {
			want[s] = true
		}
	}
	c := online.NewClient()
	var all []online.Game
	for _, a := range []online.Account{
		{Source: online.Lichess, Username: os.Getenv("LICHESS_USERNAME")},
		{Source: online.ChessCom, Username: os.Getenv("CHESSCOM_USERNAME")},
	} {
		if a.Username == "" {
			continue
		}
		err := c.Fetch(ctx, online.Query{Account: a, Max: n, Speeds: want, Rated: true},
			func(g online.Game) { all = append(all, g) })
		if err != nil {
			return nil, fmt.Errorf("%s %s : %w", a.Source, a.Username, err)
		}
	}
	sort.SliceStable(all, func(i, j int) bool { return all[i].Played > all[j].Played })
	if len(all) > n {
		all = all[:n]
	}
	return all, nil
}

// resolvePending évalue les coups qu'Alexandre a proposés en exercice et que
// l'analyse n'avait pas vus : ils deviennent justes ou faux pour de bon.
func resolvePending(ctx context.Context, api *siteAPI, sf *tactics.Stockfish, depth int) error {
	var list []exercises.Pending
	if err := api.call(ctx, "GET", "/pending", nil, &list); err != nil {
		return err
	}
	if len(list) == 0 {
		return nil
	}
	log.Printf("%d coup(s) proposé(s) à évaluer", len(list))
	for _, p := range list {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		lines, err := sf.Search(p.FEN, 1, depth, p.UCI)
		if err != nil || len(lines) == 0 {
			log.Printf("%s %s : %v", p.ID, p.UCI, err)
			continue
		}
		c := exercises.Candidate{UCI: p.UCI, CP: lines[0].CP, Mate: lines[0].Mate, Score: lines[0].Score()}
		if san, err := corpus.UCIToSAN(p.FEN, p.UCI); err == nil {
			c.SAN = san
		}
		if err := api.call(ctx, "POST", "/resolve", map[string]any{"id": p.ID, "move": c}, nil); err != nil {
			return err
		}
	}
	return nil
}

// siteAPI parle à /admin/exercices/api avec les identifiants de l'admin.
type siteAPI struct {
	base, user, pass string
}

func (a *siteAPI) call(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, a.base+path, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(a.user, a.pass)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := (&http.Client{Timeout: time.Minute}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 300))
		return fmt.Errorf("%s %s : %d %s", method, path, resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

// loadDotEnv : même lecture que cmd/tactics — une variable déjà posée gagne.
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
