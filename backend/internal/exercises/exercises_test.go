package exercises

import (
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/online"
	"github.com/iwanesko/chess-web-site/backend/internal/tactics"
)

// fakeEngine note chaque coup d'une position (vu du camp au trait). Une
// position absente de la table vaut 0 partout : aucun coup n'y perd rien.
// Dans une position notée, les coups non listés valent −1.
// Il compte les recherches par profondeur, pour vérifier que la profondeur 30
// n'est payée que pour les candidats.
type fakeEngine struct {
	scores map[string]map[string]int
	calls  map[int]int
}

func (f *fakeEngine) Search(fen string, multipv, depth int, moves ...string) ([]tactics.Line, error) {
	f.calls[depth]++
	legal, err := corpus.LegalUCIs(fen)
	if err != nil {
		return nil, err
	}
	if len(moves) > 0 {
		legal = moves
	}
	var lines []tactics.Line
	table, known := f.scores[fen]
	for _, m := range legal {
		cp, ok := table[m]
		if known && !ok {
			cp = -100 // dans une position notée, un coup non listé est mauvais
		}
		lines = append(lines, tactics.Line{CP: cp, PV: []string{m}})
	}
	sort.SliceStable(lines, func(i, j int) bool { return lines[i].CP > lines[j].CP })
	if len(lines) > multipv {
		lines = lines[:multipv]
	}
	return lines, nil
}

// La partie : 1.e4 e5 2.Nf3 Nc6 3.Bb5 a6 4.Ba4 Nf6. Le joueur préparé a les
// Noirs ; on fait de 3…a6 une imprécision de 0,4, et de 4…Nf6 un coup juste.
func game(t *testing.T) (online.Game, string) {
	t.Helper()
	san := strings.Fields("e4 e5 Nf3 Nc6 Bb5 a6 Ba4 Nf6")
	fen, before := startFEN, ""
	var uci []string
	for i, m := range san {
		u, err := corpus.SANToUCI(fen, m)
		if err != nil {
			t.Fatal(err)
		}
		if i == 5 {
			before = fen
		}
		uci = append(uci, u)
		fen, _ = corpus.ApplyUCI(fen, u)
	}
	return online.Game{URL: "https://lichess.org/abc", Source: online.Lichess,
		White: "Other", Black: "Alex", Speed: "blitz", Date: "2026-10-01",
		SAN: san, UCI: uci}, before
}

func TestAnalyse(t *testing.T) {
	g, before := game(t)
	eng := &fakeEngine{calls: map[int]int{}, scores: map[string]map[string]int{
		// Après 3.Bb5 : Nf6 et d6 bons (+0,1 / 0), a6 joué à −0,3 → perte 0,4.
		before: {"g8f6": 10, "d7d6": 0, "f8c5": -20, "a7a6": -30},
	}}
	got, err := AnalyzeGame(eng, g, "alex", Options{SkipMoves: 1}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("%d exercices, veut 1 : %+v", len(got), got)
	}
	e := got[0]
	if e.Colour != "b" || e.Ply != 5 || e.PlayedSAN != "a6" || e.Loss != 40 || e.Phase != "opening" {
		t.Errorf("exercice mal rempli : %+v", e)
	}
	if v, _ := e.Check("g8f6"); v != Correct {
		t.Errorf("Nf6 (meilleur) : %s", v)
	}
	if v, _ := e.Check("d7d6"); v != Correct {
		t.Errorf("d6 (à 0,1 du meilleur) : %s", v)
	}
	if v, _ := e.Check("f8c5"); v != Wrong {
		t.Errorf("Bc5 (à 0,3) : %s", v)
	}
	if v, _ := e.Check("a7a6"); v != Wrong {
		t.Errorf("le coup de la partie : %s", v)
	}
	if v, _ := e.Check("h7h6"); v != Unknown {
		t.Errorf("coup jamais évalué : %s", v)
	}
	if e.Moves[0].SAN != "Nf6" {
		t.Errorf("SAN : %+v", e.Moves)
	}
	// Profondeur 30 : une seule position confirmée, en deux recherches au plus
	// (les trois meilleurs, puis le coup joué hors de ce trio).
	if eng.calls[30] != 2 {
		t.Errorf("%d recherches à profondeur 30, veut 2", eng.calls[30])
	}
	// Mêmes partie et demi-coup : même identifiant, d'une analyse à l'autre.
	again, _ := AnalyzeGame(eng, g, "Alex", Options{SkipMoves: 1}, nil)
	if again[0].ID != e.ID {
		t.Error("identifiant instable")
	}
}

func TestPositionDejaJouee(t *testing.T) {
	if worthIt(900, 700, 30) {
		t.Error("+9 → +7 n'est pas un exercice")
	}
	if !worthIt(600, 200, 30) {
		t.Error("+6 → +2 en est un")
	}
	if !worthIt(40, 0, 30) || worthIt(40, 20, 30) {
		t.Error("seuil de 0,3 mal appliqué")
	}
}

func TestStore(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "site.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	g, before := game(t)
	eng := &fakeEngine{calls: map[int]int{}, scores: map[string]map[string]int{
		before: {"g8f6": 10, "d7d6": 0, "f8c5": -20, "a7a6": -30},
	}}
	list, _ := AnalyzeGame(eng, g, "Alex", Options{SkipMoves: 1}, nil)
	if err := s.Import(g.URL, list); err != nil {
		t.Fatal(err)
	}
	id := list[0].ID

	if _, v, _ := s.Answer(id, "g8f6"); v != Correct {
		t.Fatalf("réponse juste : %s", v)
	}
	if _, v, _ := s.Answer(id, "h7h6"); v != Unknown {
		t.Fatalf("réponse inconnue : %s", v)
	}
	p, _ := s.Pending()
	if len(p) != 1 || p[0].UCI != "h7h6" || p[0].FEN != before {
		t.Fatalf("file du moteur : %+v", p)
	}
	// Le moteur tranche : h6 vaut −0,05, donc dans la tolérance.
	if err := s.Resolve(id, Candidate{UCI: "h7h6", SAN: "h6", CP: -5, Score: -5}); err != nil {
		t.Fatal(err)
	}
	if p, _ := s.Pending(); len(p) != 0 {
		t.Fatalf("file non vidée : %+v", p)
	}
	if _, v, _ := s.Answer(id, "h7h6"); v != Correct {
		t.Fatalf("après résolution : %s", v)
	}

	// L'étiquette théorique sort la position des manches, et survit à une
	// nouvelle analyse de la partie.
	if err := s.SetTheory(id, true); err != nil {
		t.Fatal(err)
	}
	if err := s.Import(g.URL, list); err != nil {
		t.Fatal(err)
	}
	if l, _ := s.List(Filter{}); len(l) != 0 {
		t.Fatalf("une position théorique revient : %+v", l)
	}
	if l, _ := s.List(Filter{Theory: true}); len(l) != 1 || l[0].Solved != 2 {
		t.Fatalf("liste des théoriques : %+v", l)
	}
	if games, _ := s.AnalyzedGames(); len(games) != 1 {
		t.Fatalf("parties analysées : %v", games)
	}
	st, err := s.Stats()
	if err != nil || st.Total != 1 || st.Theory != 1 || st.Games != 1 {
		t.Fatalf("stats : %+v %v", st, err)
	}
}
