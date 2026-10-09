package corpus

import (
	"strings"
	"testing"
)

const start = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// Une partie rejouée SAN → UCI → SAN doit retomber sur le SAN d'origine, échecs,
// roques, prises en passant et désambiguïsations compris.
func TestUCIToSANAllerRetour(t *testing.T) {
	games := []string{
		// Partie de l'Opéra (Morphy) : échecs, roque, mat.
		"e4 e5 Nf3 d6 d4 Bg4 dxe5 Bxf3 Qxf3 dxe5 Bc4 Nf6 Qb3 Qe7 Nc3 c6 Bg5 b5 Nxb5 cxb5 Bxb5+ Nbd7 O-O-O Rd8 Rxd7 Rxd7 Rd1 Qe6 Bxd7+ Nxd7 Qb8+ Nxb8 Rd8#",
		// Prise en passant, désambiguïsation de cavaliers, petit roque.
		"e4 Nf6 e5 d5 exd6 Nc6 Nf3 Nd7 Nc3 Nde5 Be2 g6 O-O Bg7",
		// Promotion avec échec.
		"h4 g5 hxg5 h6 gxh6 Nf6 h7 Ng8 hxg8=Q",
	}
	for _, g := range games {
		fen := start
		for i, san := range strings.Fields(g) {
			uci, err := SANToUCI(fen, san)
			if err != nil {
				t.Fatalf("%s, coup %d : %v", san, i, err)
			}
			back, err := UCIToSAN(fen, uci)
			if err != nil {
				t.Fatalf("%s (%s) : %v", san, uci, err)
			}
			if back != san {
				t.Errorf("coup %d : %s → %s → %s", i, san, uci, back)
			}
			if fen, err = ApplyUCI(fen, uci); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestLegalUCIs(t *testing.T) {
	moves, err := LegalUCIs(start)
	if err != nil {
		t.Fatal(err)
	}
	if len(moves) != 20 {
		t.Fatalf("%d coups au départ, veut 20", len(moves))
	}
	if _, err := UCIToSAN(start, "e2e5"); err == nil {
		t.Error("e2e5 accepté")
	}
}
