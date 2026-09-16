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
type Filter struct {
	PlayerID             int64
	Colour               Colour
	FromYear, ToYear     int
	ExcludeTitledTuesday bool
	Limit                int
}

// where construit la clause commune à la recherche et à l'arbre. Les deux
// doivent filtrer EXACTEMENT pareil, sinon l'arbre décrirait un autre jeu de
// parties que la liste affichée à côté.
func (f Filter) where() (string, []any) {
	var cond []string
	var args []any

	switch f.Colour {
	case White:
		cond = append(cond, "g.white_id = ?")
		args = append(args, f.PlayerID)
	case Black:
		cond = append(cond, "g.black_id = ?")
		args = append(args, f.PlayerID)
	default:
		cond = append(cond, "(g.white_id = ? OR g.black_id = ?)")
		args = append(args, f.PlayerID, f.PlayerID)
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
	cond, args := f.where()
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	rows, err := s.db.Query(gameSelect+cond+` ORDER BY g.date DESC, g.id DESC LIMIT ?`,
		append(args, limit)...)
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
		out = append(out, g)
	}
	return out, rows.Err()
}

// Continuation est une branche de l'arbre : un coup, combien de fois il a été
// joué, et ce qu'il a rapporté AU JOUEUR CHERCHÉ — pas aux Blancs. C'est ce qui
// rend l'arbre lisible quand on prépare quelqu'un : « il joue ça et il gagne ».
type Continuation struct {
	SAN    string  `json:"san"`
	UCI    string  `json:"uci"`
	Games  int     `json:"games"`
	Wins   int     `json:"wins"`
	Draws  int     `json:"draws"`
	Losses int     `json:"losses"`
	Score  float64 `json:"score"` // en pourcentage, du point de vue du joueur
}

// Tree renvoie les coups jouables après `path` (en UCI), agrégés sur toutes les
// parties du filtre.
//
// L'agrégation se fait en Go et non en SQL : le filtre a déjà réduit à quelques
// centaines de parties, et un préfixe de coups ne s'indexe pas. Les positions
// sont identifiées par le CHEMIN et non par un hash, donc les transpositions ne
// fusionnent pas — pour préparer un adversaire on descend une ligne, ce qui est
// exactement ce qu'on veut voir.
func (s *Store) Tree(f Filter, path []string) ([]Continuation, error) {
	cond, args := f.where()
	rows, err := s.db.Query(`
		SELECT g.uci, g.san, g.result, g.white_id
		FROM game g LEFT JOIN event e ON e.id = g.event_id
		WHERE `+cond, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type agg struct {
		san            string
		games, w, d, l int
	}
	order := []string{}
	byUCI := map[string]*agg{}

	for rows.Next() {
		var uciStr, sanStr string
		var result int
		var whiteID int64
		if err := rows.Scan(&uciStr, &sanStr, &result, &whiteID); err != nil {
			return nil, err
		}
		uci := strings.Fields(uciStr)
		san := strings.Fields(sanStr)
		if len(uci) <= len(path) || len(san) < len(uci) {
			continue
		}
		match := true
		for i, p := range path {
			if uci[i] != p {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		next := uci[len(path)]
		a := byUCI[next]
		if a == nil {
			a = &agg{san: san[len(path)]}
			byUCI[next] = a
			order = append(order, next)
		}
		a.games++
		// Le résultat est rapporté au joueur cherché : il compte ses gains, pas
		// ceux des Blancs.
		pov := result
		if whiteID != f.PlayerID {
			pov = -result
		}
		switch {
		case pov > 0:
			a.w++
		case pov < 0:
			a.l++
		default:
			a.d++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]Continuation, 0, len(order))
	for _, u := range order {
		a := byUCI[u]
		c := Continuation{SAN: a.san, UCI: u, Games: a.games, Wins: a.w, Draws: a.d, Losses: a.l}
		if a.games > 0 {
			c.Score = (float64(a.w) + float64(a.d)/2) * 100 / float64(a.games)
		}
		out = append(out, c)
	}
	// Tri par fréquence : c'est la question posée, « qu'est-ce qu'il joue le
	// plus souvent ».
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j].Games > out[j-1].Games; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out, nil
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
