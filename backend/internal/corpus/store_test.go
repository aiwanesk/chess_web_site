package corpus

import (
	"os"
	"testing"
)

// corpus.db pèse des centaines de mégaoctets et contient du contenu de cours
// acheté : il n'est pas dans le dépôt, et ne doit pas y entrer. Ces tests-là ne
// tournent donc que si on leur désigne une base, et se sautent partout ailleurs
// — en CI notamment.
//
//	CORPUS_DB=C:/Users/Alex/Desktop/outputs/corpus.db go test ./internal/corpus/
func openTestStore(t *testing.T) *Store {
	t.Helper()
	path := os.Getenv("CORPUS_DB")
	if path == "" {
		t.Skip("CORPUS_DB non défini : test sauté")
	}
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func TestLookupStartPosition(t *testing.T) {
	s := openTestStore(t)

	p, err := s.Lookup(startFEN)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Moves) == 0 {
		t.Fatal("aucun coup depuis la position initiale")
	}
	if p.Moves[0].SAN != "e4" {
		t.Errorf("coup le plus joué = %q, attendu e4", p.Moves[0].SAN)
	}
	// Les coups arrivent triés par fréquence, et les parts se somment à 100.
	var sum float64
	for i, m := range p.Moves {
		if i > 0 && m.N > p.Moves[i-1].N {
			t.Fatalf("coups mal triés : %s (%d) après %s (%d)", m.SAN, m.N, p.Moves[i-1].SAN, p.Moves[i-1].N)
		}
		sum += m.Pct
	}
	if sum < 99.5 || sum > 100.5 {
		t.Errorf("somme des parts = %.2f %%, attendu ~100", sum)
	}
	if p.Reach == 0 {
		t.Error("reach nul sur la position initiale")
	}
	t.Logf("%d coups, reach=%d, meilleur=%s (%.1f %%)", len(p.Moves), p.Reach, p.Moves[0].SAN, p.Moves[0].Pct)
}

// Une position jamais atteinte doit renvoyer une réponse vide, pas une erreur :
// c'est le cas normal dès qu'on sort de la théorie couverte.
func TestLookupUnknownPosition(t *testing.T) {
	s := openTestStore(t)

	p, err := s.Lookup("7k/8/8/8/8/8/6PP/7K w - - 0 1")
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Moves) != 0 || len(p.Notes) != 0 || p.Reach != 0 {
		t.Fatalf("position inconnue : attendu une réponse vide, obtenu %+v", p)
	}
}

// Le FEN vient du navigateur : il ne doit jamais faire paniquer le serveur.
func TestLookupRejectsBadFEN(t *testing.T) {
	s := openTestStore(t)

	if _, err := s.Lookup("n'importe quoi"); err == nil {
		t.Fatal("FEN invalide acceptée")
	}
}

func TestMeta(t *testing.T) {
	s := openTestStore(t)

	m, err := s.Meta()
	if err != nil {
		t.Fatal(err)
	}
	if m["games"] == "" {
		t.Errorf("meta sans nombre de parties : %+v", m)
	}
	t.Logf("meta = %+v", m)
}

// Refuser une base qui n'est pas un corpus est ce qui protégera le
// téléversement : on vérifie avant de remplacer la base en ligne.
func TestOpenRejectsForeignDatabase(t *testing.T) {
	path := t.TempDir() + "/pas-un-corpus.db"
	if err := os.WriteFile(path, []byte("SQLite format 3\x00 mais pas vraiment"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(path); err == nil {
		t.Fatal("une base étrangère a été acceptée")
	}
}
