package games

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Export réécrit toute la base en PGN, dans l'ordre d'insertion.
//
// En FLUX, partie par partie : onze millions de parties font plusieurs
// gigaoctets de texte, qu'il n'est pas question de poser en mémoire ni sur le
// disque du conteneur. Le contexte arrête la lecture — fermer l'onglet qui
// télécharge doit libérer la base, pas la laisser moudre dans le vide.
//
// Ce qui sort, c'est ce que la base garde, et rien de plus : ni Site, ni Round,
// ni commentaires, ni variantes. Ces balises obligatoires du PGN sortent en
// « ? », comme le veut la norme pour une valeur inconnue.
func (s *Store) Export(ctx context.Context, out io.Writer) (int, error) {
	rows, err := s.db.QueryContext(ctx, gameSelect+`1 ORDER BY g.id`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	w := bufio.NewWriterSize(out, 256<<10)
	n := 0
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.White, &g.Black, &g.WhiteElo, &g.BlackElo,
			&g.Event, &g.Date, &g.Year, &g.ECO, &g.Result, &g.SAN, &g.UCI); err != nil {
			return n, err
		}
		if err := WritePGN(w, g); err != nil {
			return n, err
		}
		n++
	}
	if err := rows.Err(); err != nil {
		return n, err
	}
	return n, w.Flush()
}

// WritePGN écrit une partie au format d'export PGN : les sept balises
// obligatoires dans l'ordre de la norme, puis le texte des coups coupé à 80
// colonnes, qui est ce que ChessBase et Scid relisent sans broncher.
func WritePGN(w io.Writer, g Game) error {
	res := resultString(g.Result)
	date := g.Date
	if date == "" {
		date = "????.??.??"
	}
	var b strings.Builder
	tag := func(k, v string) {
		if v == "" {
			v = "?"
		}
		b.WriteString("[" + k + " \"" + escapeTag(v) + "\"]\n")
	}
	tag("Event", g.Event)
	tag("Site", "")
	tag("Date", date)
	tag("Round", "")
	tag("White", g.White)
	tag("Black", g.Black)
	tag("Result", res)
	if g.WhiteElo > 0 {
		tag("WhiteElo", strconv.Itoa(g.WhiteElo))
	}
	if g.BlackElo > 0 {
		tag("BlackElo", strconv.Itoa(g.BlackElo))
	}
	if g.ECO != "" {
		tag("ECO", g.ECO)
	}
	b.WriteByte('\n')

	line := 0
	word := func(s string) {
		if line > 0 && line+1+len(s) > 80 {
			b.WriteByte('\n')
			line = 0
		} else if line > 0 {
			b.WriteByte(' ')
			line++
		}
		b.WriteString(s)
		line += len(s)
	}
	for i, san := range strings.Fields(g.SAN) {
		if i%2 == 0 {
			word(fmt.Sprintf("%d.%s", i/2+1, san))
		} else {
			word(san)
		}
	}
	word(res)
	b.WriteString("\n\n")
	_, err := io.WriteString(w, b.String())
	return err
}

func resultString(r int) string {
	switch r {
	case 1:
		return "1-0"
	case -1:
		return "0-1"
	default:
		return "1/2-1/2"
	}
}

// escapeTag protège les deux seuls caractères que la norme demande d'échapper
// dans une valeur de balise. Un guillemet dans un nom d'événement, sans ça,
// ferait dérailler le lecteur sur toute la suite du fichier.
func escapeTag(s string) string {
	if !strings.ContainsAny(s, `"\`) {
		return s
	}
	return strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(s)
}
