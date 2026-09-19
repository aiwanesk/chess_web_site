// megaindex construit mega.db depuis un ou plusieurs PGN.
//
// Il tourne sur le PC, pas sur le serveur : les gigaoctets de PGN ne montent
// jamais en ligne, seule la base dérivée le fait — via l'onglet « Bases » de
// /admin.
//
//	go build -o megaindex.exe ./cmd/megaindex
//	megaindex.exe -out mega.db "D:\bases\mega2026.pgn"
//	megaindex.exe -out mega.db "D:\bases\*.zip"
//
// Rejouable : relancer sur le même fichier ne fait que compter des doublons,
// grâce à l'index UNIQUE sur l'empreinte. C'est aussi la reprise après une
// interruption — on relance, et seul ce qui manque entre.
package main

import (
	"archive/zip"
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
	"github.com/iwanesko/chess-web-site/backend/internal/twic"
)

func main() {
	out := flag.String("out", "mega.db", "base SQLite à écrire (créée si absente)")
	workers := flag.Int("workers", 0, "traducteurs SAN→UCI en parallèle (défaut : nombre de cœurs)")
	batch := flag.Int("batch", 20000, "parties par transaction")
	twicLast := flag.Int("twic-last", 0,
		"numéro TWIC déjà contenu dans le PGN. À NE POSER QUE SI on en est sûr :\n"+
			"la clé écrite ici empêche le rattrapage automatique du serveur.")
	flag.Parse()

	files := expand(flag.Args())
	// Poser le curseur seul est une opération légitime — on rattrape une base
	// déjà construite — donc l'absence de fichier n'est une erreur que s'il n'y
	// a rien d'autre à faire.
	if len(files) == 0 && *twicLast <= 0 {
		fmt.Fprintln(os.Stderr, "usage : megaindex -out mega.db <fichier.pgn|fichier.zip>…")
		fmt.Fprintln(os.Stderr, "        megaindex -out mega.db -twic-last 1662   (curseur seul)")
		os.Exit(2)
	}

	w, err := games.OpenWriter(*out)
	if err != nil {
		fatal(err)
	}
	defer w.Close()

	// Ctrl+C ferme proprement : la transaction en cours est validée avant de
	// rendre la main, donc on ne perd pas les vingt mille dernières parties.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	start := time.Now()
	var total twic.IndexStats
	for _, f := range files {
		fmt.Printf("→ %s\n", f)
		st, err := indexFile(ctx, f, w, *workers, *batch, total)
		total = add(total, st)
		if err != nil {
			fmt.Fprintf(os.Stderr, "   %s : %v\n", filepath.Base(f), err)
			break
		}
		if ctx.Err() != nil {
			fmt.Fprintln(os.Stderr, "   interrompu")
			break
		}
	}

	if *twicLast > 0 {
		if err := w.MetaSet(twic.MetaLast, fmt.Sprint(*twicLast)); err != nil {
			fatal(err)
		}
		fmt.Printf("curseur TWIC posé à %d\n", *twicLast)
	}

	if len(files) == 0 {
		return // curseur posé, rien à indexer
	}
	d := time.Since(start)
	fmt.Printf("\n%s lues · %s insérées · %s doublons · %s variantes écartées · %s illisibles\n",
		num(total.Read), num(total.Added), num(total.Skipped),
		num(total.Variants), num(total.Rejected))
	fmt.Printf("en %s (%s parties/s)\n", d.Round(time.Second),
		num(int64(float64(total.Read)/d.Seconds())))
	if total.Rejected > 0 {
		// Ce compteur doit rester à zéro : chaque unité est une partie que le
		// lecteur de SAN n'a pas su rejouer, donc une piste de correction — pas
		// une fatalité du format.
		fmt.Fprintf(os.Stderr, "⚠ %s parties n'ont pas pu être rejouées\n", num(total.Rejected))
	}
}

// indexFile ouvre un .pgn ou un .zip et l'enfourne. Le .zip est le format dans
// lequel arrivent la plupart des exports, et le décompresser à la main pour
// occuper trois fois la place sur le disque n'a pas d'intérêt.
func indexFile(ctx context.Context, path string, w *games.Writer,
	workers, batch int, sofar twic.IndexStats) (twic.IndexStats, error) {

	opt := twic.IndexOptions{
		Workers: workers, Batch: batch,
		Progress: func(st twic.IndexStats) {
			c := add(sofar, st)
			fmt.Printf("\r   %s lues · %s insérées · %s doublons   ",
				num(c.Read), num(c.Added), num(c.Skipped))
		},
	}
	defer fmt.Println()

	if strings.EqualFold(filepath.Ext(path), ".zip") {
		zr, err := zip.OpenReader(path)
		if err != nil {
			return twic.IndexStats{}, err
		}
		defer zr.Close()
		var total twic.IndexStats
		for _, f := range zr.File {
			if !strings.EqualFold(filepath.Ext(f.Name), ".pgn") {
				continue
			}
			rc, err := f.Open()
			if err != nil {
				return total, err
			}
			st, err := twic.Index(ctx, bufio.NewReaderSize(rc, 4<<20), w, opt)
			rc.Close()
			total = add(total, st)
			if err != nil {
				return total, err
			}
		}
		return total, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return twic.IndexStats{}, err
	}
	defer f.Close()
	return twic.Index(ctx, bufio.NewReaderSize(f, 4<<20), w, opt)
}

// expand développe les jokers : cmd.exe ne le fait pas pour nous, contrairement
// à un shell Unix, et « *.pgn » arriverait tel quel.
func expand(args []string) []string {
	var out []string
	for _, a := range args {
		if !strings.ContainsAny(a, "*?") {
			out = append(out, a)
			continue
		}
		m, err := filepath.Glob(a)
		if err != nil || len(m) == 0 {
			fmt.Fprintf(os.Stderr, "aucun fichier pour %q\n", a)
			continue
		}
		out = append(out, m...)
	}
	return out
}

func add(a, b twic.IndexStats) twic.IndexStats {
	return twic.IndexStats{
		Read: a.Read + b.Read, Added: a.Added + b.Added,
		Skipped: a.Skipped + b.Skipped, Variants: a.Variants + b.Variants,
		Rejected: a.Rejected + b.Rejected,
	}
}

// num écrit les grands nombres avec des espaces fines : « 11 240 318 ». À onze
// millions, « 11240318 » ne se lit plus.
func num(n int64) string {
	s := fmt.Sprint(n)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			b.WriteRune(' ')
		}
		b.WriteRune(c)
	}
	return b.String()
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
