"""Génère backend/internal/games/folds.go.

Le repli des lettres accentuées vers leur base ASCII est DÉRIVÉ de la
décomposition Unicode, pas recopié : une table écrite à la main oublie toujours
une langue — ici c'était le vietnamien, dont les voyelles portent deux
diacritiques et ne figurent dans aucune liste « européenne ».

    python scripts/gen-folds.py

Aucune dépendance : `unicodedata` est dans la bibliothèque standard, et le Go
produit n'importe rien non plus.
"""
import io
import unicodedata

OUT = "backend/internal/games/folds.go"

# Celles-là n'ont pas de décomposition : elles ne se déplient pas, elles se
# traduisent. Sans ça, « Straße » et « Strasse » restent deux joueurs.
SPECIAL = {
    "ß": "ss", "æ": "ae", "œ": "oe", "ø": "o", "đ": "d", "ð": "d",
    "þ": "th", "ł": "l", "ħ": "h", "ŋ": "n", "ı": "i", "ĸ": "k",
}

HEADER = '''// Code généré par scripts/gen-folds.py — NE PAS MODIFIER À LA MAIN.
//
// Repli des lettres accentuées vers leur base ASCII, dérivé de la décomposition
// Unicode plutôt que recopié : une table écrite à la main oublie toujours une
// langue, et en l'occurrence c'était le vietnamien.
//
// Les clés sont échappées pour que le fichier reste lisible quel que soit
// l'encodage de l'éditeur ; le caractère est rappelé en commentaire.

package games

var folds = map[rune]string{'''


def main():
    folds = dict(SPECIAL)
    for cp in range(0x00C0, 0x2000):
        ch = chr(cp)
        if not ch.isalpha():
            continue
        low = ch.lower()
        if len(low) != 1 or low in folds:
            continue
        base = "".join(c for c in unicodedata.normalize("NFD", low)
                       if unicodedata.category(c) != "Mn")
        if base and base != low and all("a" <= c <= "z" for c in base):
            folds[low] = base

    esc = chr(92) + "u"   # antislash-u, sans l'écrire (les heredocs le mangent)
    tab = chr(9)
    lines = [HEADER]
    for k in sorted(folds):
        lines.append("%s'%s%04x': \"%s\", // %s" % (tab, esc, ord(k), folds[k], k))
    lines.append("}")
    lines.append("")

    with io.open(OUT, "w", encoding="utf-8", newline="\n") as f:
        f.write("\n".join(lines))
    print(len(folds), "replis écrits dans", OUT)
    for probe in "ễşńžđåœß":
        print("   ", probe, "->", folds.get(probe, "(inchangé)"))


if __name__ == "__main__":
    main()
