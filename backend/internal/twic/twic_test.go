package twic

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// La logique de « mardi prochain » est la seule chose qu'un cron d'hébergeur
// n'aurait pas laissé tester. Autant en profiter.
func TestNextRunTombeToujoursUnMardi(t *testing.T) {
	// Le 2026-09-19 est un samedi ; on balaie deux semaines heure par heure.
	start := time.Date(2026, 9, 19, 0, 0, 0, 0, time.UTC)
	for h := 0; h < 24*14; h++ {
		now := start.Add(time.Duration(h) * time.Hour)
		next := nextRun(now)
		if next.Weekday() != time.Tuesday {
			t.Fatalf("depuis %s : %s n'est pas un mardi", now, next)
		}
		if next.Hour() != runHour || next.Minute() != 0 {
			t.Fatalf("depuis %s : %s n'est pas à %dh00 UTC", now, next, runHour)
		}
		if !next.After(now) {
			t.Fatalf("depuis %s : %s n'est pas dans le futur", now, next)
		}
		if next.Sub(now) > 7*24*time.Hour {
			t.Fatalf("depuis %s : %s est à plus d'une semaine", now, next)
		}
	}
}

// Le cas qui casse une implémentation naïve : on est mardi, mais l'heure est
// passée. Reprogrammer « aujourd'hui » ferait tourner la boucle à vide en
// rafale jusqu'à minuit.
func TestNextRunUnMardiApresLHeure(t *testing.T) {
	mardi := time.Date(2026, 9, 22, runHour, 30, 0, 0, time.UTC)
	next := nextRun(mardi)
	if got := next.Sub(mardi); got < 6*24*time.Hour {
		t.Fatalf("mardi %s → %s, soit %v plus tard", mardi, next, got)
	}
}

func TestParseMovetextNeGardeQueLaLignePrincipale(t *testing.T) {
	in := `1. e4 {Le coup le plus joué.} e5 2. Nf3 $1 (2. Bc4 Nf6 (2... Nc6 3. d3) 3. d3) 2... Nc6
3.Bb5 a6 1/2-1/2`
	want := []string{"e4", "e5", "Nf3", "Nc6", "Bb5", "a6"}
	got := parseMovetext(in)
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Fatalf("obtenu %v\nattendu %v", got, want)
	}
}

func TestReadGamesSepareSansLigneVide(t *testing.T) {
	in := `[Event "A"]
[White "Un, Joueur"]
[Black "Deux, Joueur"]
[Result "1-0"]
1. e4 e5 1-0
[Event "B"]
[White "Trois, Joueur"]
[Black "Quatre, Joueur"]
[Result "0-1"]
1. d4 d5 0-1
`
	got := readGames(strings.NewReader(in))
	if len(got) != 2 {
		t.Fatalf("%d parties lues, 2 attendues", len(got))
	}
	if got[1].tags["White"] != "Trois, Joueur" || len(got[1].moves) != 2 {
		t.Fatalf("deuxième partie mal lue : %+v", got[1])
	}
}

// Un nom en Latin-1 qui repartirait tel quel serait introuvable à la recherche :
// la partie existerait dans la base sans jamais ressortir.
func TestLecteurDecodeLeLatin1(t *testing.T) {
	raw := []byte("[White \"Gr\xfcnfeld, Ernst\"]\n[Black \"X, Y\"]\n[Result \"1-0\"]\n1. d4 d5 1-0\n")
	got := readGames(bytes.NewReader(raw))
	if len(got) != 1 {
		t.Fatalf("%d parties lues", len(got))
	}
	if got[0].tags["White"] != "Grünfeld, Ernst" {
		t.Fatalf("nom décodé en %q", got[0].tags["White"])
	}
}

const gameA = `[Event "Tournoi de Test"]
[Date "2026.09.15"]
[White "Iwanesko, Alexandre"]
[Black "Vianin, Pierre"]
[Result "1-0"]
[WhiteElo "2128"]
[ECO "C42"]
1. e4 e5 2. Nf3 Nf6 3. Nxe5 d6 4. Nf3 Nxe4 5. d4 d5 6. Bd3 Be7 7. O-O Nc6 1-0
`

const gameB = `[Event "Titled Tuesday Blitz September 15 2026"]
[Date "2026.09.15"]
[White "Carlsen, Magnus"]
[Black "Iwanesko, Alexandre"]
[Result "0-1"]
1. d4 Nf6 2. c4 e6 3. Nf3 d5 4. Nc3 Be7 5. Bg5 h6 6. Bh4 O-O 0-1
`

func zipOf(t *testing.T, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("twic.pgn")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write([]byte(body)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// serveIssues répond pour les numéros donnés et 404 pour tout le reste — ce qui
// est exactement le comportement de TWIC pour un numéro pas encore paru.
func serveIssues(t *testing.T, issues map[int]string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var n int
		if _, err := fmt.Sscanf(filepath.Base(r.URL.Path), "twic%dg.zip", &n); err != nil {
			http.NotFound(w, r)
			return
		}
		body, ok := issues[n]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(zipOf(t, body))
	}))
	t.Cleanup(srv.Close)
	return srv
}

// newBase fabrique une mega.db vide mais VALIDE. Elle doit exister avant
// l'import : l'importeur refuse de créer la base, sinon un GAMES_DB mal écrit
// produirait silencieusement une base vide.
func newBase(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mega.db")
	w, err := games.OpenWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestImportAvanceLeCurseurEtInsere(t *testing.T) {
	srv := serveIssues(t, map[int]string{
		DefaultLast + 1: gameA,
		DefaultLast + 2: gameB,
	})
	path := newBase(t)
	im := &Importer{Path: path, BaseURL: srv.URL}

	rep, err := im.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rep.Last != DefaultLast+2 {
		t.Fatalf("curseur à %d, attendu %d", rep.Last, DefaultLast+2)
	}
	if rep.Added != 2 {
		t.Fatalf("%d parties ajoutées, 2 attendues (%+v)", rep.Added, rep.Issues)
	}

	s, err := games.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	players, err := s.SearchPlayers("iwanesko", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 1 {
		t.Fatalf("%d joueurs trouvés pour « iwanesko »", len(players))
	}
	rows, err := s.Search(games.Filter{PlayerIDs: []int64{players[0].ID}, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("%d parties retrouvées, 2 attendues", len(rows))
	}
	// Le filtre Titled Tuesday ne marche que si l'événement a été marqué À
	// L'INGESTION : c'est le moment de le vérifier.
	rows, err = s.Search(games.Filter{
		PlayerIDs: []int64{players[0].ID}, ExcludeTitledTuesday: true, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("%d parties hors Titled Tuesday, 1 attendue", len(rows))
	}
}

// Le cœur de la demande : ne jamais retélécharger ni redoubler. On rejoue le
// MÊME import — le curseur doit tenir, et rien ne doit entrer deux fois.
func TestImportEstRejouable(t *testing.T) {
	srv := serveIssues(t, map[int]string{DefaultLast + 1: gameA + gameB})
	path := newBase(t)
	im := &Importer{Path: path, BaseURL: srv.URL}

	if _, err := im.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	second, err := im.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Issues) != 0 || second.Added != 0 {
		t.Fatalf("le deuxième passage a retravaillé : %+v", second)
	}
	if second.Last != DefaultLast+1 {
		t.Fatalf("curseur retombé à %d", second.Last)
	}
}

// Et le garde-fou indépendant du curseur : si le même numéro revient malgré
// tout — curseur remis à zéro, base restaurée d'une sauvegarde — l'index
// UNIQUE doit absorber les doublons au lieu de gonfler la base.
func TestEmpreinteAbsorbeLesDoublons(t *testing.T) {
	srv := serveIssues(t, map[int]string{DefaultLast + 1: gameA + gameA + gameB})
	path := newBase(t)
	im := &Importer{Path: path, BaseURL: srv.URL}

	rep, err := im.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rep.Added != 2 || rep.Skipped != 1 {
		t.Fatalf("ajoutées %d / doublons %d, attendu 2 et 1", rep.Added, rep.Skipped)
	}
}

// Une partie dont le SAN ne se rejoue pas doit être écartée sans emporter la
// livraison : un seul mauvais score ne fait pas perdre la semaine.
func TestPartieIllisibleNeFaitPasPerdreLaLivraison(t *testing.T) {
	casse := strings.Replace(gameA, "3. Nxe5", "3. Nxe7", 1)
	srv := serveIssues(t, map[int]string{DefaultLast + 1: casse + gameB})
	path := newBase(t)
	im := &Importer{Path: path, BaseURL: srv.URL}

	rep, err := im.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rep.Added != 1 {
		t.Fatalf("%d parties ajoutées, 1 attendue", rep.Added)
	}
	if rep.Last != DefaultLast+1 {
		t.Fatalf("curseur à %d : la livraison n'a pas été validée", rep.Last)
	}
}

// Sans base téléversée, l'importeur ne doit RIEN créer : une mega.db vide
// fabriquée par accident ferait répondre « aucune partie » à l'explorateur au
// lieu de « aucune base ».
func TestImportNeCreeJamaisLaBase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "absente.db")
	im := &Importer{Path: path, BaseURL: "http://127.0.0.1:1"}
	if _, err := im.Run(context.Background()); err == nil {
		t.Fatal("import accepté sans base")
	}
	if _, err := games.Open(path); err == nil {
		t.Fatal("une base a été créée")
	}
	if st := im.Status(); st.Present {
		t.Fatal("Status annonce une base présente")
	}
}

// Le Chess960 voyage dans les mêmes fichiers que le reste. Il doit être écarté
// EXPLICITEMENT — pas se retrouver dans le tas des parties illisibles, où il
// masquerait un vrai défaut du lecteur de SAN.
func TestChess960EstEcarteEtNonRejete(t *testing.T) {
	const c960 = `[Event "Champions Chess960"]
[Date "2026.09.15"]
[White "Kasparov, Garry"]
[Black "Topalov, Veselin"]
[Result "1-0"]
[FEN "rkbnnbqr/pppppppp/8/8/8/8/PPPPPPPP/RKBNNBQR w KQkq - 0 1"]
1. d4 f5 2. Nd3 Nf6 1-0
`
	srv := serveIssues(t, map[int]string{DefaultLast + 1: c960 + gameA})
	im := &Importer{Path: newBase(t), BaseURL: srv.URL}

	rep, err := im.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	got := rep.Issues[0]
	if got.Variants != 1 {
		t.Fatalf("%d variantes écartées, 1 attendue (%+v)", got.Variants, got)
	}
	if got.Rejected != 0 {
		t.Fatalf("%d parties comptées comme illisibles (%+v)", got.Rejected, got)
	}
	if got.Added != 1 {
		t.Fatalf("%d parties ajoutées, 1 attendue", got.Added)
	}
}

// Index est le chemin PARALLÈLE : la traduction SAN→UCI est répartie sur
// plusieurs goroutines et l'écriture est groupée en transactions. Il doit
// donner exactement le même résultat que l'import séquentiel — et surtout
// rester rejouable, puisque c'est ce qui sert de reprise après interruption
// sur un PGN de plusieurs gigaoctets.
func TestIndexDonneLeMemeResultatQueLImport(t *testing.T) {
	const c960 = `[Event "Champions Chess960"]
[White "Kasparov, Garry"]
[Black "Topalov, Veselin"]
[Result "1-0"]
[FEN "rkbnnbqr/pppppppp/8/8/8/8/PPPPPPPP/RKBNNBQR w KQkq - 0 1"]
1. d4 f5 1-0
`
	casse := strings.Replace(gameA, "3. Nxe5", "3. Nxe7", 1)
	pgn := gameA + gameB + c960 + casse

	path := newBase(t)
	w, err := games.OpenWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Close()

	// Un lot d'une seule partie force plusieurs transactions : c'est là qu'une
	// erreur de découpage se verrait.
	opt := IndexOptions{Workers: 4, Batch: 1}
	st, err := Index(context.Background(), strings.NewReader(pgn), w, opt)
	if err != nil {
		t.Fatal(err)
	}
	if st.Read != 4 || st.Added != 2 || st.Variants != 1 || st.Rejected != 1 {
		t.Fatalf("premier passage : %+v", st)
	}

	again, err := Index(context.Background(), strings.NewReader(pgn), w, opt)
	if err != nil {
		t.Fatal(err)
	}
	if again.Added != 0 || again.Skipped != 2 {
		t.Fatalf("second passage : %+v — il devrait n'y avoir que des doublons", again)
	}
}
