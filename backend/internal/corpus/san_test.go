package corpus

import (
	"strings"
	"testing"
)

// Même méthode que pour ApplyUCI : la base sert d'oracle. Chaque arête porte le
// SAN ET l'UCI, écrits par python-chess. Rejouer le graphe et vérifier que
// SANToUCI retombe sur l'UCI stocké compare mon générateur de coups à une
// implémentation de référence, sur des dizaines de milliers de positions
// réelles — clouages, désambiguïsations, promotions compris.
//
// Le test compte aussi les cas particuliers rencontrés : un générateur de coups
// qui n'a jamais vu de roque ni de promotion n'a pas été testé, il a été
// parcouru.
func TestSANToUCIReproducesStoredMoves(t *testing.T) {
	s := openTestStore(t)

	frontier := []string{startFEN}
	seen := map[string]bool{startFEN: true}

	var checked, castles, promos, enPassant, disamb int
	for depth := 0; depth < 12 && len(frontier) > 0; depth++ {
		var next []string
		for _, fen := range frontier {
			key, err := SignedHash(fen)
			if err != nil {
				t.Fatalf("%s: %v", fen, err)
			}
			rows, err := s.db.Query(
				`SELECT uci, san FROM move WHERE pos = ? ORDER BY n DESC LIMIT 12`, key)
			if err != nil {
				t.Fatal(err)
			}
			type edge struct{ uci, san string }
			var edges []edge
			for rows.Next() {
				var e edge
				if err := rows.Scan(&e.uci, &e.san); err != nil {
					rows.Close()
					t.Fatal(err)
				}
				edges = append(edges, e)
			}
			rows.Close()

			for _, e := range edges {
				got, err := SANToUCI(fen, e.san)
				if err != nil {
					t.Fatalf("SANToUCI(%q, %q) : %v", fen, e.san, err)
				}
				if got != e.uci {
					t.Fatalf("depuis\n  %s\n%s traduit en %s, attendu %s",
						fen, e.san, got, e.uci)
				}
				checked++
				switch {
				case strings.HasPrefix(e.san, "O-O"):
					castles++
				case strings.Contains(e.san, "="):
					promos++
				}
				// Prise en passant : un pion prend sur une case vide.
				if e.san[0] >= 'a' && e.san[0] <= 'h' && strings.Contains(e.san, "x") {
					p, _ := parsePosition(fen)
					if to, ok := parseSquare(e.uci[2:4]); ok && p.board[to] == 0 {
						enPassant++
					}
				}
				// Désambiguïsation : « Cbd2 », « T1e2 », « Dh4e1 ».
				if len(e.san) > 0 && e.san[0] >= 'A' && e.san[0] <= 'Z' &&
					len(strings.TrimRight(strings.ReplaceAll(e.san, "x", ""), "+#")) > 3 {
					disamb++
				}

				child, err := ApplyUCI(fen, got)
				if err != nil {
					t.Fatalf("ApplyUCI(%q, %q) : %v", fen, got, err)
				}
				if !seen[child] && len(next) < 4000 {
					seen[child] = true
					next = append(next, child)
				}
			}
		}
		frontier = next
	}

	t.Logf("%d coups traduits · %d roques · %d promotions · %d prises en passant · %d désambiguïsations",
		checked, castles, promos, enPassant, disamb)
	if checked < 10000 {
		t.Fatalf("échantillon trop maigre : %d coups", checked)
	}
	if castles == 0 || disamb == 0 {
		t.Fatalf("cas particuliers jamais rencontrés (roques %d, désambiguïsations %d)",
			castles, disamb)
	}
}

// Les cas que le graphe d'ouverture ne fournit pas : une base de théorie ne
// contient presque aucune promotion ni prise en passant en profondeur 12, et
// aucun coup rendu non ambigu par un clouage. Écrits à la main, vérifiés à la
// main.
func TestSANToUCIEdgeCases(t *testing.T) {
	cases := []struct {
		name, fen, san, want string
	}{
		{"promotion dame", "8/4P3/8/8/8/8/8/K6k w - - 0 1", "e8=Q", "e7e8q"},
		{"sous-promotion cavalier", "8/4P3/8/8/8/8/8/K6k w - - 0 1", "e8=N", "e7e8n"},
		{"promotion en prenant", "3r4/4P3/8/8/8/8/8/K6k w - - 0 1", "exd8=Q+", "e7d8q"},
		{"prise en passant", "8/8/8/3pP3/8/8/8/K6k w - d6 0 1", "exd6", "e5d6"},
		{"petit roque", "8/8/8/8/8/8/8/R3K2R w KQ - 0 1", "O-O", "e1g1"},
		{"grand roque", "8/8/8/8/8/8/8/R3K2R w KQ - 0 1", "O-O-O", "e1c1"},
		{"roque écrit en zéros", "8/8/8/8/8/8/8/R3K2R w KQ - 0 1", "0-0", "e1g1"},
		{"roque noir", "r3k2r/8/8/8/8/8/8/K7 b kq - 0 1", "O-O-O", "e8c8"},
		// Les cavaliers de b1 et f1 vont tous deux en d2 : la colonne tranche.
		{"désambiguïsation colonne", "4k3/8/8/8/8/8/8/1N2KN2 w - - 0 1", "Nbd2", "b1d2"},
		{"désambiguïsation colonne (bis)", "4k3/8/8/8/8/8/8/1N2KN2 w - - 0 1", "Nfd2", "f1d2"},
		// Les dames de a1 et a5 touchent e1, et partagent la colonne : il faut
		// donc la forme complète.
		{"désambiguïsation complète", "4k3/8/8/Q6Q/8/8/8/Q5K1 w - - 0 1", "Qa5e1", "a5e1"},
		// La tour de e2 est clouée par la tour noire de e8 : « Rd2 » ne peut donc
		// désigner que celle de a2, sans désambiguïsation. C'est le cas qui
		// distingue un vrai générateur légal d'un filtre pseudo-légal.
		{"clouage lève l'ambiguïté", "4r2k/8/8/8/8/8/R3R3/4K3 w - - 0 1", "Rd2", "a2d2"},
		{"échec annoncé ignoré", "4k3/8/8/8/8/8/8/R3K3 w Q - 0 1", "Ra8#", "a1a8"},
	}
	for _, c := range cases {
		got, err := SANToUCI(c.fen, c.san)
		if err != nil {
			t.Errorf("%s : %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s : %s → %s, attendu %s", c.name, c.san, got, c.want)
		}
	}
}

// Un SAN qui ne correspond à aucun coup légal doit être refusé, pas deviné :
// un fichier PGN abîmé doit sauter la partie, pas empoisonner la base.
func TestSANToUCIRejectsIllegal(t *testing.T) {
	cases := []struct{ name, fen, san string }{
		{"coup impossible", startFEN, "e5"},
		{"roque sans droit", "8/8/8/8/8/8/8/R3K2R w - - 0 1", "O-O"},
		{"roque à travers un échec", "3r4/8/8/8/8/8/8/R3K2R w KQ - 0 1", "O-O-O"},
		{"laisse le roi en échec", "4r3/8/8/8/8/8/4R3/4K3 w - - 0 1", "Rd2"},
		{"charabia", startFEN, "zz9"},
	}
	for _, c := range cases {
		if got, err := SANToUCI(c.fen, c.san); err == nil {
			t.Errorf("%s : %q accepté et traduit en %s", c.name, c.san, got)
		}
	}
}
