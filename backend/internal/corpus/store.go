package corpus

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Store lit le graphe de positions. La base est ouverte en LECTURE SEULE, à
// deux niveaux : `mode=ro` côté pilote et `query_only` côté moteur. Elle pèse
// plusieurs centaines de mégaoctets, se régénère uniquement hors ligne, et rien
// dans le serveur web n'a de raison d'y écrire — autant rendre l'accident
// impossible plutôt que improbable.
type Store struct {
	db *sql.DB
}

// Move est une arête du graphe : un coup jouable depuis la position demandée.
type Move struct {
	UCI   string  `json:"uci"`
	SAN   string  `json:"san"`
	N     int     `json:"n"`     // nombre de parties qui l'ont joué
	Pct   float64 `json:"pct"`   // part de ce coup dans la position
	Notes int     `json:"notes"` // commentaires attachés à la position d'arrivée
}

// Note est un commentaire de cours attaché à la position.
type Note struct {
	Label string `json:"label"` // le chapitre d'où il vient
	Text  string `json:"text"`
	N     int    `json:"n"`
}

// Position est la réponse complète pour une position.
type Position struct {
	FEN        string `json:"fen"`
	Moves      []Move `json:"moves"`
	Notes      []Note `json:"notes"`
	NotesTotal int    `json:"notesTotal"`
	Reach      int    `json:"reach"` // parties passées par cette position
}

// noteLimit borne la réponse : certaines positions d'ouverture portent des
// centaines de commentaires, et personne ne les lit sur un téléphone.
const noteLimit = 400

// Open ouvre corpus.db en lecture seule et vérifie que c'en est bien une.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?mode=ro&_pragma=query_only(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.check(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

// check est le contrôle de recevabilité d'une base : le schéma attendu, et
// surtout la position initiale qui doit répondre. Cette seule requête teste à
// la fois la présence de la table, celle de l'index et la convention de hash —
// c'est ce qui permettra de refuser un téléversement tronqué sans y toucher.
func (s *Store) check() error {
	const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	key, err := SignedHash(startFEN)
	if err != nil {
		return err
	}
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM move WHERE pos = ?`, key).Scan(&n); err != nil {
		return fmt.Errorf("corpus: base illisible (%w)", err)
	}
	if n == 0 {
		return fmt.Errorf("corpus: la position initiale ne renvoie aucun coup — base tronquée ou clés incompatibles")
	}
	return nil
}

func (s *Store) Close() error { return s.db.Close() }

// Meta renvoie la table `meta` telle quelle (nombre de parties, date de build…).
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

// Lookup renvoie ce que le corpus sait d'une position.
func (s *Store) Lookup(fen string) (*Position, error) {
	key, err := SignedHash(fen)
	if err != nil {
		return nil, err
	}
	out := &Position{FEN: fen, Moves: []Move{}, Notes: []Note{}}

	rows, err := s.db.Query(
		`SELECT uci, san, child, n FROM move WHERE pos = ? ORDER BY n DESC`, key)
	if err != nil {
		return nil, err
	}
	var children []int64
	for rows.Next() {
		var m Move
		var child int64
		if err := rows.Scan(&m.UCI, &m.SAN, &child, &m.N); err != nil {
			rows.Close()
			return nil, err
		}
		out.Reach += m.N
		children = append(children, child)
		out.Moves = append(out.Moves, m)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range out.Moves {
		if out.Reach > 0 {
			out.Moves[i].Pct = float64(out.Moves[i].N) * 100 / float64(out.Reach)
		}
	}

	// Combien de commentaires attendent derrière chaque coup : c'est ce qui
	// permet d'afficher une pastille et de savoir où il y a quelque chose à
	// lire avant de jouer le coup.
	if err := s.countChildNotes(out.Moves, children); err != nil {
		return nil, err
	}

	if err := s.loadNotes(out, key); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) countChildNotes(moves []Move, children []int64) error {
	if len(children) == 0 {
		return nil
	}
	// Une requête par lot plutôt qu'une par coup : jusqu'à quelques dizaines de
	// coups par position, et autant d'allers-retours SQLite évités.
	q := `SELECT pos, COUNT(*) FROM note WHERE pos IN (?` +
		repeatComma(len(children)-1) + `) GROUP BY pos`
	args := make([]any, len(children))
	for i, c := range children {
		args[i] = c
	}
	rows, err := s.db.Query(q, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	counts := map[int64]int{}
	for rows.Next() {
		var pos int64
		var n int
		if err := rows.Scan(&pos, &n); err != nil {
			return err
		}
		counts[pos] = n
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for i := range moves {
		moves[i].Notes = counts[children[i]]
	}
	return nil
}

func (s *Store) loadNotes(out *Position, key int64) error {
	rows, err := s.db.Query(
		`SELECT label.name, note.text, note.n FROM note
		 JOIN label ON label.id = note.label
		 WHERE note.pos = ? ORDER BY note.n DESC, label.name LIMIT ?`, key, noteLimit)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var n Note
		var cnt sql.NullInt64
		if err := rows.Scan(&n.Label, &n.Text, &cnt); err != nil {
			return err
		}
		n.N = 1
		if cnt.Valid && cnt.Int64 > 0 {
			n.N = int(cnt.Int64)
		}
		out.Notes = append(out.Notes, n)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(out.Notes) < noteLimit {
		out.NotesTotal = len(out.Notes)
		return nil
	}
	return s.db.QueryRow(`SELECT COUNT(*) FROM note WHERE pos = ?`, key).Scan(&out.NotesTotal)
}

func repeatComma(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, 0, n*2)
	for i := 0; i < n; i++ {
		b = append(b, ',', '?')
	}
	return string(b)
}
