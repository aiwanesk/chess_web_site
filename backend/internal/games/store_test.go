package games

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Une base minuscule, écrite ici : mega.db pèse des gigaoctets et vit sur le PC
// d'Alexandre, mais rien de ce qui est testé ci-dessous n'a besoin de volume.
func testDB(t *testing.T) *Store {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mega.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(Schema); err != nil {
		t.Fatal(err)
	}
	exec := func(q string, args ...any) {
		if _, err := db.Exec(q, args...); err != nil {
			t.Fatal(err)
		}
	}
	// Trois orthographes du même homme, comme le font MegaBase et TWIC.
	for id, name := range map[int]string{1: "Pahud, Cedric", 2: "Iwanesko, Alexandre", 3: "Vianin, Pascal"} {
		exec(`INSERT INTO player (id, name, norm) VALUES (?,?,?)`, id, name, Normalize(name))
	}
	exec(`INSERT INTO event (id,name,norm,titled_tuesday) VALUES (1,'Open de Monthey',?,0)`, Normalize("Open de Monthey"))
	exec(`INSERT INTO event (id,name,norm,titled_tuesday) VALUES (2,'Titled Tuesday Blitz',?,1)`, Normalize("Titled Tuesday Blitz"))

	add := func(id, w, b, ev, year, result int, uci, san string) {
		exec(`INSERT INTO game (id,white_id,black_id,event_id,year,date,result,eco,uci,san)
		      VALUES (?,?,?,?,?,?,?,'B23',?,?)`,
			id, w, b, ev, year, "2026.01.01", result, uci, san)
	}
	// Pahud avec les Blancs : deux fois e4, une fois d4.
	add(1, 1, 2, 1, 2024, 1, "e2e4 c7c5 g1f3", "e4 c5 Nf3")
	add(2, 1, 3, 1, 2025, 0, "e2e4 e7e5", "e4 e5")
	add(3, 1, 2, 1, 2020, -1, "d2d4 g8f6", "d4 Nf6")
	// Pahud avec les Noirs, et une partie de Titled Tuesday à écarter.
	add(4, 2, 1, 1, 2025, 1, "d2d4 d7d5", "d4 d5")
	add(5, 1, 3, 2, 2025, 1, "e2e4 c7c6", "e4 c6")

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSearchPlayersPrefixFirst(t *testing.T) {
	s := testDB(t)

	got, err := s.SearchPlayers("pahud", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Name != "Pahud, Cedric" {
		t.Fatalf("recherche par préfixe = %+v", got)
	}
	if got[0].Games != 5 {
		t.Errorf("parties de Pahud = %d, attendu 5", got[0].Games)
	}
	// La recherche doit ignorer accents et ponctuation des deux côtés.
	for _, q := range []string{"PAHUD", "pahud, c", "Pahud,C."} {
		if r, _ := s.SearchPlayers(q, 10); len(r) == 0 {
			t.Errorf("%q ne trouve personne", q)
		}
	}
	// Un fragment au milieu doit fonctionner aussi, via le repli « contient ».
	if r, _ := s.SearchPlayers("wanesko", 10); len(r) != 1 {
		t.Errorf("recherche « contient » = %+v", r)
	}
}

func TestSearchByColourAndYears(t *testing.T) {
	s := testDB(t)

	all, _ := s.Search(Filter{PlayerID: 1})
	if len(all) != 5 {
		t.Fatalf("toutes couleurs = %d parties, attendu 5", len(all))
	}
	white, _ := s.Search(Filter{PlayerID: 1, Colour: White})
	if len(white) != 4 {
		t.Fatalf("avec les Blancs = %d, attendu 4", len(white))
	}
	black, _ := s.Search(Filter{PlayerID: 1, Colour: Black})
	if len(black) != 1 {
		t.Fatalf("avec les Noirs = %d, attendu 1", len(black))
	}
	since, _ := s.Search(Filter{PlayerID: 1, Colour: White, FromYear: 2024})
	if len(since) != 3 {
		t.Fatalf("Blancs depuis 2024 = %d, attendu 3", len(since))
	}
	window, _ := s.Search(Filter{PlayerID: 1, FromYear: 2024, ToYear: 2024})
	if len(window) != 1 {
		t.Fatalf("année 2024 seule = %d, attendu 1", len(window))
	}
}

func TestExcludeTitledTuesday(t *testing.T) {
	s := testDB(t)

	with, _ := s.Search(Filter{PlayerID: 1, Colour: White})
	without, _ := s.Search(Filter{PlayerID: 1, Colour: White, ExcludeTitledTuesday: true})
	if len(with) != 4 || len(without) != 3 {
		t.Fatalf("avec = %d (attendu 4), sans = %d (attendu 3)", len(with), len(without))
	}
	for _, g := range without {
		if g.Event == "Titled Tuesday Blitz" {
			t.Fatal("une partie de Titled Tuesday a survécu au filtre")
		}
	}
}

func TestTreeCountsAndScoreFromPlayerPointOfView(t *testing.T) {
	s := testDB(t)

	// Pahud avec les Blancs, Titled Tuesday écarté : e4 deux fois, d4 une fois.
	root, err := s.Tree(Filter{PlayerID: 1, Colour: White, ExcludeTitledTuesday: true}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(root) != 2 {
		t.Fatalf("premiers coups = %+v", root)
	}
	if root[0].SAN != "e4" || root[0].Games != 2 {
		t.Fatalf("coup le plus joué = %+v, attendu e4 ×2", root[0])
	}
	// e4 : une victoire, une nulle → 75 %.
	if root[0].Wins != 1 || root[0].Draws != 1 || root[0].Losses != 0 {
		t.Fatalf("bilan de e4 = %+v", root[0])
	}
	if root[0].Score < 74.9 || root[0].Score > 75.1 {
		t.Errorf("score de e4 = %.1f %%, attendu 75", root[0].Score)
	}
	// d4 : une défaite pour les Blancs, donc pour Pahud.
	if root[1].SAN != "d4" || root[1].Losses != 1 || root[1].Score != 0 {
		t.Fatalf("bilan de d4 = %+v", root[1])
	}

	// Une branche plus bas.
	after, _ := s.Tree(Filter{PlayerID: 1, Colour: White, ExcludeTitledTuesday: true}, []string{"e2e4"})
	if len(after) != 2 {
		t.Fatalf("réponses à e4 = %+v", after)
	}
}

// Avec les Noirs, une victoire des Blancs est une DÉFAITE pour le joueur
// cherché. C'est l'inversion qu'on oublie, et elle rend l'arbre trompeur.
func TestTreeScoreIsInvertedWhenPlayerIsBlack(t *testing.T) {
	s := testDB(t)

	tree, err := s.Tree(Filter{PlayerID: 1, Colour: Black}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tree) != 1 || tree[0].SAN != "d4" {
		t.Fatalf("arbre côté Noirs = %+v", tree)
	}
	if tree[0].Losses != 1 || tree[0].Wins != 0 {
		t.Fatalf("les Blancs gagnent, donc Pahud perd : %+v", tree[0])
	}
}

// La normalisation est écrite DEUX FOIS : ici, pour la recherche, et dans
// l'indexeur Python, pour remplir player.norm. Si les deux divergent, la
// recherche ne trouve plus rien — exactement le même mode de panne silencieuse
// que le hash Zobrist du corpus. D'où ce fichier de référence partagé : le test
// Go le rejoue, et l'indexeur doit le rejouer aussi avant d'écrire quoi que ce
// soit.
func TestNormalizeAgainstSharedFixture(t *testing.T) {
	raw, err := os.ReadFile("testdata/normalize.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct{ In, Out string }
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) < 20 {
		t.Fatalf("jeu de référence trop maigre : %d cas", len(cases))
	}
	for _, c := range cases {
		if got := Normalize(c.In); got != c.Out {
			t.Errorf("Normalize(%q) = %q, attendu %q", c.In, got, c.Out)
		}
	}
	t.Logf("%d cas de normalisation vérifiés", len(cases))
}

func TestTitledTuesdayDetection(t *testing.T) {
	yes := []string{"Titled Tuesday Blitz", "titled tuesday 25th feb", "TITLED TUESDAY, Late"}
	no := []string{"Open de Monthey", "Titled Players Cup", "Tuesday Evening League"}
	for _, e := range yes {
		if !TitledTuesday(e) {
			t.Errorf("%q devrait être reconnu", e)
		}
	}
	for _, e := range no {
		if TitledTuesday(e) {
			t.Errorf("%q ne devrait PAS être reconnu", e)
		}
	}
}
