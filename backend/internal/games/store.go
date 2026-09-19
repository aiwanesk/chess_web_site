// Package games interroge la base de parties (mega.db) : recherche par joueur,
// par couleur et par période, et arbre des ouvertures construit sur le résultat.
//
// C'est un métier différent de celui du paquet corpus, et la base l'est aussi.
// corpus.db est un graphe de positions interrogé par hash Zobrist, qui répond
// « qu'est-ce qui se joue ici et qu'en disent les cours ». mega.db est une table
// de parties interrogée en SQL sur des index, qui répond « qu'a joué ce type
// avec les Noirs depuis 2022 ». Les deux ne partagent que l'échiquier qui
// affiche le résultat.
package games

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// Schema est la forme attendue. L'indexeur tourne sur le PC d'Alexandre — les
// gigaoctets de PGN ne montent jamais sur le serveur, seule la base dérivée le
// fait — mais le schéma est défini ici pour qu'il n'y ait qu'une source de
// vérité, et pour que les tests puissent fabriquer une base minuscule.
const Schema = `
CREATE TABLE IF NOT EXISTS player (
	id   INTEGER PRIMARY KEY,
	name TEXT NOT NULL,          -- tel qu'écrit dans le PGN
	norm TEXT NOT NULL           -- sans accent, en minuscules : c'est là qu'on cherche
);
CREATE INDEX IF NOT EXISTS player_norm ON player(norm);

CREATE TABLE IF NOT EXISTS event (
	id   INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	norm TEXT NOT NULL,
	-- Marqué À L'INGESTION plutôt que testé par LIKE à chaque requête : un
	-- entier se lit, une chaîne se compare.
	titled_tuesday INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS game (
	id        INTEGER PRIMARY KEY,
	-- Empreinte de dédoublonnage : joueurs + date + résultat + coups (voir
	-- games.Key). NULLable à dessein — SQLite autorise plusieurs NULL sous un
	-- index UNIQUE, donc une base construite sans empreintes reste valide, elle
	-- perd seulement la protection contre les doublons.
	hash      INTEGER,
	white_id  INTEGER NOT NULL,
	black_id  INTEGER NOT NULL,
	event_id  INTEGER,
	year      INTEGER,
	date      TEXT,
	result    INTEGER NOT NULL,   -- 1 = les Blancs gagnent, 0 = nulle, -1 = les Noirs
	eco       TEXT,
	white_elo INTEGER,
	black_elo INTEGER,
	-- Deux colonnes et pas une : l'UCI se rejoue (donc donne les positions, via
	-- corpus.ApplyUCI), le SAN se lit. Les dériver l'un de l'autre demanderait
	-- un générateur de coups légaux, c'est-à-dire un second endroit où diverger.
	uci TEXT NOT NULL,            -- « e2e4 e7e5 g1f3 … »
	san TEXT NOT NULL             -- « e4 e5 Nf3 … »
);
-- Les deux index qui portent toute la recherche : « ce joueur, avec cette
-- couleur, sur cette période » devient un parcours d'index même sur onze
-- millions de lignes.
CREATE INDEX IF NOT EXISTS game_white ON game(white_id, year);
CREATE INDEX IF NOT EXISTS game_black ON game(black_id, year);
-- C'est cet index qui rend l'import TWIC rejouable : un INSERT OR IGNORE sur
-- une partie déjà connue ne fait rien au lieu de la doubler.
CREATE UNIQUE INDEX IF NOT EXISTS game_hash ON game(hash);

CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT);
`

// Store lit la base de parties. Ouverte en lecture seule ici : seul l'importeur
// TWIC y écrit, et il le fera avec sa propre connexion.
type Store struct{ db *sql.DB }

// Open ouvre mega.db en lecture seule et vérifie que le schéma est là.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro&_pragma=query_only(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM game`).Scan(&n); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("games: base illisible (%w)", err)
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Player est un joueur trouvé par la recherche.
type Player struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Games int    `json:"games"`
}

// SearchPlayers cherche par préfixe d'abord — c'est ce que sert l'index — puis
// élargit à « contient » si la moisson est maigre. L'ordre compte : taper
// « pahud » doit d'abord proposer Pahud, pas les quinze joueurs dont le nom
// contient ces lettres au milieu.
func (s *Store) SearchPlayers(q string, limit int) ([]Player, error) {
	q = Normalize(q)
	if q == "" {
		return nil, nil
	}
	if limit <= 0 {
		limit = 20
	}
	out, err := s.searchPlayers(q+"%", limit)
	if err != nil {
		return nil, err
	}
	if len(out) >= limit {
		return out, nil
	}
	// Le « contient » ne peut pas utiliser l'index et parcourt la table des
	// joueurs. C'est acceptable parce qu'elle est petite devant celle des
	// parties, et parce qu'on n'y tombe que si le préfixe n'a rien donné.
	seen := map[int64]bool{}
	for _, p := range out {
		seen[p.ID] = true
	}
	more, err := s.searchPlayers("%"+q+"%", limit)
	if err != nil {
		return out, nil // le préfixe a déjà donné quelque chose : on s'en contente
	}
	for _, p := range more {
		if !seen[p.ID] && len(out) < limit {
			out = append(out, p)
		}
	}
	return out, nil
}

func (s *Store) searchPlayers(pattern string, limit int) ([]Player, error) {
	rows, err := s.db.Query(`
		SELECT p.id, p.name,
		       (SELECT COUNT(*) FROM game g WHERE g.white_id = p.id OR g.black_id = p.id)
		FROM player p WHERE p.norm LIKE ? ORDER BY p.norm LIMIT ?`, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.Games); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Colour restreint la recherche à une couleur.
type Colour string

const (
	White Colour = "w"
	Black Colour = "b"
	Both  Colour = ""
)

// Filter décrit une recherche. Zéro sur les années = pas de borne.
//
// PlayerIDs accepte PLUSIEURS joueurs : on prépare parfois une équipe, ou on
// veut l'arbre de deux joueurs qui partagent le même répertoire.
type Filter struct {
	PlayerIDs            []int64
	Colour               Colour
	FromYear, ToYear     int
	ExcludeTitledTuesday bool
	Limit                int
	// Path restreint aux parties qui passent par ce début de partie (en UCI).
	// Le tri se fait en Go, après la requête : un préfixe de coups ne s'indexe
	// pas, et le filtre a déjà ramené l'ensemble à quelques centaines de lignes.
	Path []string
}

func placeholders(n int) string {
	if n <= 0 {
		return "NULL"
	}
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

// where construit la clause commune à la recherche et à l'arbre. Les deux
// doivent filtrer EXACTEMENT pareil, sinon l'arbre décrirait un autre jeu de
// parties que la liste affichée à côté.
func (f Filter) where() (string, []any) {
	ids := make([]any, len(f.PlayerIDs))
	for i, id := range f.PlayerIDs {
		ids[i] = id
	}
	ph := placeholders(len(ids))

	var cond []string
	var args []any
	switch f.Colour {
	case White:
		cond = append(cond, "g.white_id IN ("+ph+")")
		args = append(args, ids...)
	case Black:
		cond = append(cond, "g.black_id IN ("+ph+")")
		args = append(args, ids...)
	default:
		cond = append(cond, "(g.white_id IN ("+ph+") OR g.black_id IN ("+ph+"))")
		args = append(args, ids...)
		args = append(args, ids...)
	}
	if f.FromYear > 0 {
		cond = append(cond, "g.year >= ?")
		args = append(args, f.FromYear)
	}
	if f.ToYear > 0 {
		cond = append(cond, "g.year <= ?")
		args = append(args, f.ToYear)
	}
	if f.ExcludeTitledTuesday {
		// LEFT JOIN : une partie sans événement renseigné ne doit pas disparaître
		// parce qu'on exclut les Titled Tuesday.
		cond = append(cond, "COALESCE(e.titled_tuesday, 0) = 0")
	}
	return strings.Join(cond, " AND "), args
}

// Game est une partie trouvée.
type Game struct {
	ID       int64  `json:"id"`
	White    string `json:"white"`
	Black    string `json:"black"`
	WhiteElo int    `json:"whiteElo"`
	BlackElo int    `json:"blackElo"`
	Event    string `json:"event"`
	Date     string `json:"date"`
	Year     int    `json:"year"`
	ECO      string `json:"eco"`
	Result   int    `json:"result"`
	SAN      string `json:"san,omitempty"`
	UCI      string `json:"uci,omitempty"`
}

const gameSelect = `
	SELECT g.id, w.name, b.name, COALESCE(g.white_elo,0), COALESCE(g.black_elo,0),
	       COALESCE(e.name,''), COALESCE(g.date,''), COALESCE(g.year,0),
	       COALESCE(g.eco,''), g.result, g.san, g.uci
	FROM game g
	JOIN player w ON w.id = g.white_id
	JOIN player b ON b.id = g.black_id
	LEFT JOIN event e ON e.id = g.event_id
	WHERE `

// Search renvoie les parties correspondant au filtre, les plus récentes d'abord.
func (s *Store) Search(f Filter) ([]Game, error) {
	if len(f.PlayerIDs) == 0 {
		return []Game{}, nil
	}
	cond, args := f.where()
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	// Avec un chemin, on ramène plus large et on trie ensuite : sinon le LIMIT
	// couperait avant le filtrage et on perdrait des parties qui y passent.
	sqlLimit := limit
	if len(f.Path) > 0 {
		sqlLimit = 3000
	}
	rows, err := s.db.Query(gameSelect+cond+` ORDER BY g.date DESC, g.id DESC LIMIT ?`,
		append(args, sqlLimit)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Game{}
	for rows.Next() {
		var g Game
		if err := rows.Scan(&g.ID, &g.White, &g.Black, &g.WhiteElo, &g.BlackElo,
			&g.Event, &g.Date, &g.Year, &g.ECO, &g.Result, &g.SAN, &g.UCI); err != nil {
			return nil, err
		}
		if len(f.Path) > 0 && !hasPrefix(strings.Fields(g.UCI), f.Path) {
			continue
		}
		out = append(out, g)
		if len(out) >= limit {
			break
		}
	}
	return out, rows.Err()
}

// Node est un nœud de l'arbre d'ouverture : un coup, combien de parties l'ont
// joué, et ce qu'il a rapporté AUX JOUEURS CHERCHÉS — pas aux Blancs. C'est ce
// qui rend l'arbre lisible en préparation : « il joue ça, et il gagne ».
type Node struct {
	SAN    string  `json:"san"`
	UCI    string  `json:"uci"`
	Games  int     `json:"games"`
	Wins   int     `json:"wins"`
	Draws  int     `json:"draws"`
	Losses int     `json:"losses"`
	Score  float64 `json:"score"`
	// FEN est laissé vide ici : ce paquet ne fait que du SQL et de
	// l'agrégation. C'est la couche serveur qui le remplit, parce qu'elle a
	// déjà le moteur d'application des coups — et qu'un paquet de requêtes n'a
	// pas à connaître les règles du jeu.
	FEN      string `json:"fen,omitempty"`
	Children []Node `json:"children,omitempty"`
}

// TreeOptions borne l'arbre. Sans bornes il descendrait jusqu'au dernier coup
// de la plus longue partie : illisible, et une réponse énorme.
type TreeOptions struct {
	Path     []string // en UCI : d'où part l'arbre
	MaxDepth int      // en demi-coups sous Path (défaut 12, soit six coups)
	MinGames int      // une branche sous ce seuil est coupée (défaut 1)
}

// OpeningTree construit l'arbre d'ouverture des parties filtrées, imbriqué, en
// UNE seule lecture de la base — plutôt qu'un aller-retour par niveau déplié.
//
// L'agrégation se fait en Go et non en SQL : le filtre a déjà réduit à quelques
// centaines de parties, et un préfixe de coups ne s'indexe pas. Les nœuds sont
// identifiés par le CHEMIN parcouru et non par un hash de position, donc les
// transpositions ne fusionnent pas — en préparation on descend une ligne, et
// c'est exactement ce qu'on veut lire.
func (s *Store) OpeningTree(f Filter, opt TreeOptions) ([]Node, error) {
	if len(f.PlayerIDs) == 0 {
		return []Node{}, nil
	}
	if opt.MaxDepth <= 0 {
		opt.MaxDepth = 12
	}
	if opt.MinGames <= 0 {
		opt.MinGames = 1
	}
	selected := map[int64]bool{}
	for _, id := range f.PlayerIDs {
		selected[id] = true
	}

	cond, args := f.where()
	rows, err := s.db.Query(`
		SELECT g.uci, g.san, g.result, g.white_id, g.black_id
		FROM game g LEFT JOIN event e ON e.id = g.event_id
		WHERE `+cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	root := &builder{kids: map[string]*builder{}}
	for rows.Next() {
		var uciStr, sanStr string
		var result int
		var whiteID, blackID int64
		if err := rows.Scan(&uciStr, &sanStr, &result, &whiteID, &blackID); err != nil {
			return nil, err
		}
		uci, san := strings.Fields(uciStr), strings.Fields(sanStr)
		if len(san) < len(uci) {
			continue // scores dépareillés : on ne devine pas
		}
		if !hasPrefix(uci, opt.Path) {
			continue
		}
		// Le résultat est rapporté aux joueurs cherchés : ils comptent LEURS
		// gains. Avec les Noirs, une victoire des Blancs est une défaite.
		pov := result
		if !povIsWhite(whiteID, blackID, selected, f.Colour) {
			pov = -result
		}
		node := root
		for d := len(opt.Path); d < len(uci) && d-len(opt.Path) < opt.MaxDepth; d++ {
			node = node.child(uci[d], san[d])
			node.count(pov)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return root.harvest(opt.MinGames), nil
}

// povIsWhite : lequel des deux camps représente les joueurs cherchés ? La
// couleur demandée tranche ; sans elle on prend celui des deux qui est dans la
// sélection — et les Blancs si les deux y sont, ce qui arrive dès qu'on met
// dans le même arbre deux joueurs qui se sont affrontés.
func povIsWhite(whiteID, blackID int64, selected map[int64]bool, c Colour) bool {
	switch c {
	case White:
		return true
	case Black:
		return false
	default:
		if selected[whiteID] {
			return true
		}
		return !selected[blackID]
	}
}

func hasPrefix(uci, path []string) bool {
	if len(uci) < len(path) {
		return false
	}
	for i, p := range path {
		if uci[i] != p {
			return false
		}
	}
	return true
}

type builder struct {
	san            string
	games, w, d, l int
	order          []string
	kids           map[string]*builder
}

func (b *builder) child(uci, san string) *builder {
	k := b.kids[uci]
	if k == nil {
		k = &builder{san: san, kids: map[string]*builder{}}
		b.kids[uci] = k
		b.order = append(b.order, uci)
	}
	return k
}

func (b *builder) count(pov int) {
	b.games++
	switch {
	case pov > 0:
		b.w++
	case pov < 0:
		b.l++
	default:
		b.d++
	}
}

func (b *builder) harvest(min int) []Node {
	out := make([]Node, 0, len(b.order))
	for _, uci := range b.order {
		k := b.kids[uci]
		if k.games < min {
			continue
		}
		n := Node{SAN: k.san, UCI: uci, Games: k.games, Wins: k.w, Draws: k.d, Losses: k.l}
		n.Score = (float64(k.w) + float64(k.d)/2) * 100 / float64(k.games)
		n.Children = k.harvest(min)
		out = append(out, n)
	}
	// Tri par fréquence : c'est la question posée, « qu'est-ce qu'ils jouent le
	// plus souvent ».
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Games > out[j-1].Games; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// Meta renvoie la table meta (nombre de parties, dernier TWIC importé, …).
func (s *Store) Meta() (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, value FROM meta`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]string{}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

// ByID renvoie une partie entière, coups compris. C'est ce qui permet de
// cliquer une ligne de la liste et de la rejouer : le SAN et l'UCI complets
// sont déjà stockés, il n'y a rien à recalculer.
func (s *Store) ByID(id int64) (Game, error) {
	var g Game
	err := s.db.QueryRow(gameSelect+`g.id = ?`, id).Scan(
		&g.ID, &g.White, &g.Black, &g.WhiteElo, &g.BlackElo,
		&g.Event, &g.Date, &g.Year, &g.ECO, &g.Result, &g.SAN, &g.UCI)
	if err == sql.ErrNoRows {
		return g, fmt.Errorf("games: partie %d introuvable", id)
	}
	return g, err
}
