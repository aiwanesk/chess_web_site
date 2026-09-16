package corpus

import (
	"testing"
)

// La validation d'ApplyUCI ne vient pas de cas écrits à la main mais de la base
// elle-même : chaque arête porte le hash de la position d'arrivée (`child`).
// Rejouer le graphe en Go et vérifier que chaque FEN produit retombe sur ce
// hash teste d'un coup le roque, la prise en passant, la promotion, les droits
// de roque perdus et la case ep annoncée — sur des dizaines de milliers de
// positions réelles plutôt que sur cinq exemples choisis.
func TestApplyUCIReproducesStoredChildren(t *testing.T) {
	s := openTestStore(t)

	type node struct{ fen string }
	frontier := []node{{startFEN}}
	seen := map[string]bool{startFEN: true}

	var checked, castles, promos, enPassant int
	for depth := 0; depth < 10 && len(frontier) > 0; depth++ {
		var next []node
		for _, nd := range frontier {
			key, err := SignedHash(nd.fen)
			if err != nil {
				t.Fatalf("%s: %v", nd.fen, err)
			}
			rows, err := s.db.Query(
				`SELECT uci, san, child FROM move WHERE pos = ? ORDER BY n DESC LIMIT 6`, key)
			if err != nil {
				t.Fatal(err)
			}
			type edge struct {
				uci, san string
				child    int64
			}
			var edges []edge
			for rows.Next() {
				var e edge
				if err := rows.Scan(&e.uci, &e.san, &e.child); err != nil {
					rows.Close()
					t.Fatal(err)
				}
				edges = append(edges, e)
			}
			rows.Close()

			for _, e := range edges {
				got, err := ApplyUCI(nd.fen, e.uci)
				if err != nil {
					t.Fatalf("ApplyUCI(%q, %q) : %v", nd.fen, e.uci, err)
				}
				gotKey, err := SignedHash(got)
				if err != nil {
					t.Fatalf("%s: %v", got, err)
				}
				if gotKey != e.child {
					t.Fatalf("après %s (%s) depuis\n  %s\nobtenu\n  %s\nhash %d, attendu %d",
						e.san, e.uci, nd.fen, got, gotKey, e.child)
				}
				checked++
				switch {
				case len(e.uci) == 5:
					promos++
				case e.san == "O-O" || e.san == "O-O-O":
					castles++
				}
				if nEP(nd.fen, e.uci) {
					enPassant++
				}
				if !seen[got] && len(next) < 400 {
					seen[got] = true
					next = append(next, node{got})
				}
			}
		}
		frontier = next
	}

	if checked < 5000 {
		t.Fatalf("couverture trop faible : %d arêtes seulement", checked)
	}
	if castles == 0 {
		t.Error("aucun roque rencontré : le cas n'est pas couvert")
	}
	t.Logf("%d arêtes rejouées, dont %d roques, %d promotions, %d prises en passant",
		checked, castles, promos, enPassant)
}

// nEP : le coup est-il une prise en passant ? Un pion qui change de colonne
// vers une case vide.
func nEP(fen, uci string) bool {
	f := parseFields(fen)
	if len(f) < 4 {
		return false
	}
	board, err := parsePlacement(f[0])
	if err != nil {
		return false
	}
	from, to, _, err := parseUCI(uci)
	if err != nil {
		return false
	}
	return board[from]|0x20 == 'p' && from%8 != to%8 && board[to] == 0
}

func TestApplyUCIRejectsNonsense(t *testing.T) {
	for _, c := range []struct{ fen, uci, why string }{
		{startFEN, "e2e5x", "coup mal formé"},
		{startFEN, "e3e4", "case de départ vide"},
		{startFEN, "e7e5", "pièce du camp adverse"},
		{"pas un fen", "e2e4", "FEN invalide"},
		{startFEN, "e7e8k", "promotion en roi"},
	} {
		if _, err := ApplyUCI(c.fen, c.uci); err == nil {
			t.Errorf("accepté alors que %s : %q", c.why, c.uci)
		}
	}
}

// Les droits de roque se perdent aussi quand la tour se fait PRENDRE sur sa
// case d'origine — un cas qu'on oublie facilement et que le hash trahirait.
func TestCastlingRightsLostWhenRookIsCaptured(t *testing.T) {
	const fen = "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1"
	got, err := ApplyUCI(fen, "a1a8") // la tour blanche prend la tour noire en a8
	if err != nil {
		t.Fatal(err)
	}
	if f := parseFields(got); f[2] != "Kk" {
		t.Fatalf("droits de roque = %q, attendu \"Kk\" (Q perdu par la tour partie, q par la tour prise)", f[2])
	}
}

func parseFields(s string) []string {
	var out []string
	for _, p := range splitSpace(s) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func splitSpace(s string) []string {
	var out, cur = []string{}, ""
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			out, cur = append(out, cur), ""
			continue
		}
		cur += string(s[i])
	}
	return append(out, cur)
}

// Le parcours d'ouvertures ne rencontre ni promotion ni prise en passant : ces
// cas-là sont donc écrits explicitement, avec des FEN attendus produits par
// python-chess plutôt que par moi.
func TestApplyUCIRareCases(t *testing.T) {
	cases := []struct{ fen, uci, want, note string }{
		{"rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3", "e5d6",
			"rnbqkbnr/ppp1pppp/3P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3", "prise en passant, blancs"},
		{"rnbqkbnr/pppp1ppp/8/8/3pP3/8/PPP2PPP/RNBQKBNR b KQkq e3 0 3", "d4e3",
			"rnbqkbnr/pppp1ppp/8/8/8/4p3/PPP2PPP/RNBQKBNR w KQkq - 0 4", "prise en passant, noirs"},
		{"8/P6k/8/8/8/8/6K1/8 w - - 0 1", "a7a8q",
			"Q7/7k/8/8/8/8/6K1/8 b - - 0 1", "promotion en dame"},
		{"1r5k/P7/8/8/8/8/6K1/8 w - - 0 1", "a7b8n",
			"1N5k/8/8/8/8/8/6K1/8 b - - 0 1", "promotion avec prise, en cavalier"},
		{"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1", "e1c1",
			"r3k2r/8/8/8/8/8/8/2KR3R b kq - 1 1", "grand roque blanc"},
		{"r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1", "e8g8",
			"r4rk1/8/8/8/8/8/8/R3K2R w KQ - 1 2", "petit roque noir"},
	}
	for _, c := range cases {
		got, err := ApplyUCI(c.fen, c.uci)
		if err != nil {
			t.Errorf("%s : %v", c.note, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s | obtenu %s | attendu %s", c.note, got, c.want)
		}
	}
}
