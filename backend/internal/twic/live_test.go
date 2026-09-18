package twic

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Le seul test qui touche le vrai theweekinchess.com. Désactivé par défaut — un
// test qui dépend d'un site tiers n'a rien à faire dans une suite qui doit
// passer hors ligne :
//
//	TWIC_LIVE=1662      go test ./internal/twic/ -run Live -v   (un numéro)
//	TWIC_LIVE=1640-1662 go test ./internal/twic/ -run Live -v   (une plage)
//
// Ce qu'il vérifie et que rien d'autre ne peut vérifier : que l'adresse des
// archives n'a pas changé, que le ZIP contient bien un .pgn, et surtout que le
// lecteur de SAN encaisse plusieurs milliers de parties réelles — promotions,
// sous-promotions, prises en passant, parties depuis une position de départ
// imposée. Le taux de rejet est le vrai résultat : il doit être nul.
func TestLiveIssue(t *testing.T) {
	v := os.Getenv("TWIC_LIVE")
	if v == "" {
		t.Skip("TWIC_LIVE non défini : test sauté")
	}
	from, to := v, v
	if i := strings.IndexByte(v, '-'); i > 0 {
		from, to = v[:i], v[i+1:]
	}
	n, err := strconv.Atoi(from)
	if err != nil {
		t.Fatalf("TWIC_LIVE doit être un numéro ou une plage : %v", err)
	}
	last, err := strconv.Atoi(to)
	if err != nil {
		t.Fatalf("TWIC_LIVE doit être un numéro ou une plage : %v", err)
	}

	path := filepath.Join(t.TempDir(), "mega.db")
	w, err := games.OpenWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	// On règle le curseur juste avant le premier numéro visé.
	if err := w.MetaSet(MetaLast, strconv.Itoa(n-1)); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()

	im := &Importer{Path: path}
	rep, err := im.Run(context.Background())
	if err != nil {
		t.Fatalf("import de TWIC %d : %v", n, err)
	}
	if len(rep.Issues) == 0 {
		t.Fatalf("TWIC %d introuvable", n)
	}
	var total, added, variants, rejected int
	for _, got := range rep.Issues {
		if got.Number > last {
			break
		}
		t.Logf("TWIC %d : %5d lues · %5d insérées · %3d doublons · %2d variantes · %d illisibles",
			got.Number, got.Games, got.Added, got.Skipped, got.Variants, got.Rejected)
		if got.Added+got.Skipped+got.Variants+got.Rejected != got.Games {
			t.Errorf("TWIC %d : les compteurs ne bouclent pas (%+v)", got.Number, got)
		}
		total += got.Games
		added += got.Added
		variants += got.Variants
		rejected += got.Rejected
	}
	t.Logf("TOTAL : %d parties lues, %d insérées, %d variantes écartées, %d illisibles",
		total, added, variants, rejected)
	if rejected > 0 {
		t.Errorf("%d parties sur %d n'ont pas pu être rejouées", rejected, total)
	}
}
