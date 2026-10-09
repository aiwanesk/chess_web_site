package online

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// fakeSites imite les deux API sur un seul serveur de test.
func fakeSites(t *testing.T) *Client {
	t.Helper()
	mux := http.NewServeMux()

	// Lichess : deux parties de « Prepared » + une Chess960 à écarter.
	mux.HandleFunc("/api/games/user/Prepared", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/x-ndjson" {
			t.Errorf("Accept = %q", r.Header.Get("Accept"))
		}
		if got := r.URL.Query().Get("perfType"); got != "blitz" {
			t.Errorf("perfType = %q, veut blitz", got)
		}
		lines := []string{
			`{"id":"aaa","rated":true,"variant":"standard","speed":"blitz","createdAt":1759000000000,"status":"mate","winner":"white","moves":"e4 e5 Nf3 Nc6 Bb5","players":{"white":{"user":{"name":"Prepared"},"rating":2100},"black":{"user":{"name":"Other"},"rating":2050}},"opening":{"eco":"C60","name":"Ruy Lopez"}}`,
			`{"id":"bbb","rated":true,"variant":"standard","speed":"blitz","createdAt":1758000000000,"status":"resign","winner":"white","moves":"d4 Nf6 c4 e6","players":{"white":{"user":{"name":"Other"},"rating":2000},"black":{"user":{"name":"prepared"},"rating":2090}}}`,
			`{"id":"ccc","rated":true,"variant":"chess960","speed":"blitz","createdAt":1757000000000,"status":"mate","winner":"black","moves":"e4 e5","players":{"white":{"user":{"name":"Prepared"},"rating":2100},"black":{"user":{"name":"X"},"rating":2000}}}`,
			``,
		}
		_, _ = fmt.Fprint(w, strings.Join(lines, "\n"))
	})
	mux.HandleFunc("/api/games/user/Nobody", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Chess.com : deux mois, le plus récent d'abord une fois inversé.
	mux.HandleFunc("/pub/player/hikaru/games/archives", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"archives": []string{
			"https://api.chess.com/pub/player/hikaru/games/2026/08",
			"https://api.chess.com/pub/player/hikaru/games/2026/09",
		}})
	})
	pgn := func(moves string) string {
		return "[Event \"Live Chess\"]\n[ECO \"B20\"]\n[ECOUrl \"https://www.chess.com/openings/Sicilian-Defense\"]\n\n" +
			moves + " 1-0\n"
	}
	mux.HandleFunc("/pub/player/hikaru/games/2026/09", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"games": []map[string]any{
			{"url": "https://www.chess.com/game/live/1", "pgn": pgn("1. e4 {[%clk 0:02:59]} 1... c5 {[%clk 0:02:58]} 2. Nf3"),
				"time_class": "blitz", "rated": true, "rules": "chess", "end_time": 1790467200,
				"white": map[string]any{"username": "Hikaru", "rating": 3300, "result": "win"},
				"black": map[string]any{"username": "Foe", "rating": 3000, "result": "resigned"}},
			{"url": "https://www.chess.com/game/live/2", "pgn": pgn("1. e4 e5"),
				"time_class": "blitz", "rated": true, "rules": "chess960", "end_time": 1790467300,
				"white": map[string]any{"username": "Hikaru", "result": "win"},
				"black": map[string]any{"username": "Foe", "result": "checkmated"}},
		}})
	})
	mux.HandleFunc("/pub/player/hikaru/games/2026/08", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"games": []map[string]any{
			{"url": "https://www.chess.com/game/daily/3", "pgn": pgn("1. d4 d5"),
				"time_class": "daily", "rated": true, "rules": "chess", "end_time": 1787000000,
				"white": map[string]any{"username": "Foe", "result": "agreed"},
				"black": map[string]any{"username": "Hikaru", "result": "agreed"}},
		}})
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return &Client{HTTP: srv.Client(), LichessURL: srv.URL, ChessComURL: srv.URL}
}

func collect(t *testing.T, c *Client, q Query) []Game {
	t.Helper()
	var out []Game
	if err := c.Fetch(context.Background(), q, func(g Game) { out = append(out, g) }); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestLichess(t *testing.T) {
	c := fakeSites(t)
	got := collect(t, c, Query{Account: Account{Lichess, "Prepared"}, Max: 50,
		Speeds: map[string]bool{"blitz": true}})
	if len(got) != 2 {
		t.Fatalf("%d parties, veut 2 (la Chess960 est écartée)", len(got))
	}
	g := got[0]
	if g.URL != "https://lichess.org/aaa" || g.Result != 1 || g.Date != "2025-09-27" || g.ECO != "C60" {
		t.Errorf("partie mal lue : %+v", g)
	}
	if strings.Join(g.UCI, " ") != "e2e4 e7e5 g1f3 b8c6 f1b5" {
		t.Errorf("UCI = %v", g.UCI)
	}
}

func TestLichessIntrouvable(t *testing.T) {
	c := fakeSites(t)
	err := c.Fetch(context.Background(), Query{Account: Account{Lichess, "Nobody"}, Max: 10}, func(Game) {})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, veut ErrNotFound", err)
	}
}

func TestChessCom(t *testing.T) {
	c := fakeSites(t)
	got := collect(t, c, Query{Account: Account{ChessCom, "Hikaru"}, Max: 50})
	if len(got) != 2 {
		t.Fatalf("%d parties, veut 2 (la Chess960 est écartée)", len(got))
	}
	// Le mois le plus récent d'abord.
	if got[0].URL != "https://www.chess.com/game/live/1" {
		t.Errorf("ordre : %s en premier", got[0].URL)
	}
	if strings.Join(got[0].SAN, " ") != "e4 c5 Nf3" || got[0].Result != 1 || got[0].Opening != "Sicilian Defense" {
		t.Errorf("partie mal lue : %+v", got[0])
	}
	if got[1].Speed != "correspondence" || got[1].Result != 0 {
		t.Errorf("daily = %q, résultat %d", got[1].Speed, got[1].Result)
	}

	// La borne de date coupe le mois d'août.
	since := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	if n := len(collect(t, c, Query{Account: Account{ChessCom, "Hikaru"}, Max: 50, Since: since})); n != 1 {
		t.Errorf("avec since : %d parties, veut 1", n)
	}
	// Le maximum aussi.
	if n := len(collect(t, c, Query{Account: Account{ChessCom, "Hikaru"}, Max: 1})); n != 1 {
		t.Errorf("max 1 : %d parties", n)
	}
}

func TestLotArbreEtCouleur(t *testing.T) {
	ss := NewSets(fakeSites(t))
	set, err := ss.Load([]Query{
		{Account: Account{Lichess, "Prepared"}, Max: 50, Speeds: map[string]bool{"blitz": true}},
		{Account: Account{Lichess, "Nobody"}, Max: 50},
	})
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for set.Status().State == Loading && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	st := set.Status()
	if st.State != Ready || st.Games != 2 {
		t.Fatalf("statut %+v : un compte introuvable ne doit pas tout faire échouer", st)
	}
	if !strings.Contains(st.Error, "Nobody") {
		t.Errorf("l'erreur du compte manquant doit rester visible : %q", st.Error)
	}

	// Avec les Noirs, la victoire des Blancs est une défaite du joueur préparé.
	tree := set.Tree(Filter{Colour: "b"}, games.TreeOptions{})
	if len(tree) != 1 || tree[0].SAN != "d4" || tree[0].Losses != 1 {
		t.Fatalf("arbre Noirs : %+v", tree)
	}
	tree = set.Tree(Filter{Colour: "w"}, games.TreeOptions{})
	if len(tree) != 1 || tree[0].SAN != "e4" || tree[0].Wins != 1 {
		t.Fatalf("arbre Blancs : %+v", tree)
	}
	if n := len(set.Games(Filter{Path: []string{"e2e4"}}, 0)); n != 1 {
		t.Errorf("parties par 1.e4 : %d", n)
	}
}

func TestFavoris(t *testing.T) {
	f, err := OpenFavorites(filepath.Join(t.TempDir(), "fav.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := f.Add(Favorite{Source: Lichess, Username: "Prepared", Note: "ronde 3"}); err != nil {
		t.Fatal(err)
	}
	// Même pseudo, autre casse : mise à jour, pas doublon.
	if err := f.Add(Favorite{Source: Lichess, Username: "prepared", Note: "ronde 4"}); err != nil {
		t.Fatal(err)
	}
	list, _ := f.List()
	if len(list) != 1 || list[0].Note != "ronde 4" {
		t.Fatalf("favoris : %+v", list)
	}
	for i := 1; i < MaxFavorites; i++ {
		if err := f.Add(Favorite{Source: ChessCom, Username: fmt.Sprintf("u%03d", i)}); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Add(Favorite{Source: ChessCom, Username: "onetoomany"}); !errors.Is(err, ErrTooManyFavorites) {
		t.Fatalf("101e favori : err = %v", err)
	}
	if err := f.Remove(Lichess, "PREPARED"); err != nil {
		t.Fatal(err)
	}
	if list, _ := f.List(); len(list) != MaxFavorites-1 {
		t.Errorf("après retrait : %d", len(list))
	}
}

func TestPseudo(t *testing.T) {
	for _, u := range []string{"DrNykterstein", "alireza2003", "Hikaru", "a-b_c"} {
		if !ValidUsername(u) {
			t.Errorf("%q refusé", u)
		}
	}
	for _, u := range []string{"", "x", "../admin", "a b", "é", strings.Repeat("a", 31), "u?x=1"} {
		if ValidUsername(u) {
			t.Errorf("%q accepté", u)
		}
	}
}

func TestNomOuvertureChessCom(t *testing.T) {
	for in, want := range map[string]string{
		"https://www.chess.com/openings/Sicilian-Defense-Kan-Modern-Variation...9.Nc3-Bg7-10.Be3-O-O": "Sicilian Defense Kan Modern Variation",
		"https://www.chess.com/openings/Sicilian-Defense-Kan-Variation-5...Nf6-6.O-O-Qc7":             "Sicilian Defense Kan Variation",
		"https://www.chess.com/openings/Caro-Kann-Defense-Two-Knights":                                "Caro Kann Defense Two Knights",
		"": "",
	} {
		if got := openingName(in); got != want {
			t.Errorf("%q → %q, veut %q", in, got, want)
		}
	}
}
