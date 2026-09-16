package games

import "strings"

// Normalize réduit un nom à sa forme cherchable.
//
// C'est la fonction la plus importante du paquet, et celle qui casse tout si on
// l'applique de travers. La MegaBase écrit « Pahud, Cedric », TWIC écrit parfois
// « Pahud,C. », un autre lot écrira « Pahud Cédric ». Sans forme commune, ce
// sont trois joueurs distincts et une recherche qui n'en montre qu'un tiers.
//
// Elle doit donc être appliquée AUX DEUX BOUTS — à l'ingestion pour remplir
// `player.norm`, et à la recherche pour construire le motif — et rester
// identique entre les deux. D'où sa place ici plutôt que dans l'indexeur.
//
// Le parti pris : tout en minuscules, accents repliés, et toute ponctuation
// devient une espace. « Pahud, Cedric », « Pahud,C. » et « Pahud Cédric »
// donnent alors « pahud cedric », « pahud c » et « pahud cedric » — trois
// chaînes qui partagent le préfixe « pahud c », ce qui suffit à l'autocomplétion.
func Normalize(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	lastSpace := true // évite l'espace de tête
	for _, r := range strings.ToLower(s) {
		if rep, ok := folds[r]; ok {
			b.WriteString(rep)
			lastSpace = false
			continue
		}
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			lastSpace = false
		default:
			if !lastSpace {
				b.WriteByte(' ')
				lastSpace = true
			}
		}
	}
	return strings.TrimSpace(b.String())
}

// TitledTuesday reconnaît les tournois en ligne hebdomadaires de chess.com, à
// écarter d'une préparation : ce sont des parties de blitz par centaines, qui
// noient le répertoire sérieux d'un joueur sous son jeu en ligne.
//
// À appliquer À L'INGESTION, pour que la requête teste un entier plutôt que de
// comparer une chaîne sur chaque ligne.
func TitledTuesday(eventName string) bool {
	return strings.Contains(Normalize(eventName), "titled tuesday")
}
