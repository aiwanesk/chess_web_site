package corpus

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Le jeu de référence est produit par scripts/gen-zobrist.py depuis
// python-chess, et chacune de ses positions a été recoupée avec le `child`
// stocké dans corpus.db : si ce test passe, Go calcule les mêmes clés que la
// base. Une divergence ne casserait rien visiblement — elle renverrait
// simplement zéro coup sur les positions concernées.
type goldenCase struct {
	FEN  string `json:"fen"`
	Hash string `json:"hash"`
	Note string `json:"note"`
}

func loadGolden(t *testing.T) []goldenCase {
	t.Helper()
	raw, err := os.ReadFile("testdata/zobrist.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []goldenCase
	if err := json.Unmarshal(raw, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) < 500 {
		t.Fatalf("jeu de référence trop maigre : %d positions", len(cases))
	}
	return cases
}

func TestHashMatchesPythonChess(t *testing.T) {
	cases := loadGolden(t)
	var epHashed, epIgnored int
	for _, c := range cases {
		want, err := strconv.ParseUint(strings.TrimPrefix(c.Hash, "0x"), 16, 64)
		if err != nil {
			t.Fatalf("%s: %v", c.FEN, err)
		}
		got, err := Hash(c.FEN)
		if err != nil {
			t.Fatalf("%s (%s): %v", c.FEN, c.Note, err)
		}
		if got != want {
			t.Fatalf("%s (%s)\n  attendu 0x%016x\n  obtenu  0x%016x", c.FEN, c.Note, want, got)
		}
		switch {
		case strings.HasPrefix(c.Note, "ep-compte"):
			epHashed++
		case strings.HasPrefix(c.Note, "ep-ignore"):
			epIgnored++
		}
	}
	// Les deux branches de la règle en passant doivent être exercées, sinon ce
	// test passerait aussi avec une implémentation qui les confond.
	if epHashed == 0 || epIgnored == 0 {
		t.Fatalf("couverture ep insuffisante : %d prises comptées, %d ignorées", epHashed, epIgnored)
	}
	t.Logf("%d positions, dont %d avec ep dans le hash et %d avec ep ignoré", len(cases), epHashed, epIgnored)
}

// La position initiale a une clé publiée : c'est le seul contrôle qui ne dépend
// ni de la base ni du générateur.
func TestHashStartPosition(t *testing.T) {
	const want = 0x463b96181691fc9c
	got, err := Hash("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("position initiale : attendu 0x%016x, obtenu 0x%016x", uint64(want), got)
	}
}

// Une case ep annoncée par le FEN mais sans pion pour la prendre ne doit rien
// changer au hash — c'est précisément là que se casse une implémentation naïve.
func TestEnPassantOnlyWhenAPawnIsThere(t *testing.T) {
	const withPawn = "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3"
	const noPawn = "rnbqkbnr/pp1ppppp/8/2p5/P7/8/1PPPPPPP/RNBQKBNR b KQkq a3 0 2"

	hWith, _ := Hash(withPawn)
	hWithout, _ := Hash(strings.Replace(withPawn, " d6 ", " - ", 1))
	if hWith == hWithout {
		t.Fatal("la colonne ep devrait entrer dans le hash quand un pion peut prendre")
	}

	hNo, _ := Hash(noPawn)
	hNoStripped, _ := Hash(strings.Replace(noPawn, " a3 ", " - ", 1))
	if hNo != hNoStripped {
		t.Fatal("sans pion adjacent, la case ep ne doit rien changer")
	}
}

func TestHashRejectsBadFEN(t *testing.T) {
	for _, bad := range []string{
		"",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w",          // champs manquants
		"rnbqkbnr/pppppppp/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", // 7 rangées
		"rnbqkbnr/ppppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1", // trait inconnu
	} {
		if _, err := Hash(bad); err == nil {
			t.Errorf("FEN invalide acceptée : %q", bad)
		}
	}
}
