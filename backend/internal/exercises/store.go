package exercises

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

// Candidate est un coup évalué dans la position de l'exercice. Score est vu du
// joueur au trait (le mat compacté par tactics.Line.Score).
type Candidate struct {
	UCI   string   `json:"uci"`
	SAN   string   `json:"san"`
	CP    int      `json:"cp"`
	Mate  int      `json:"mate,omitempty"`
	Score int      `json:"score"`
	OK    bool     `json:"ok"`
	Line  []string `json:"line,omitempty"`
}

// Exercise est une position où Alexandre a perdu au moins le seuil.
type Exercise struct {
	ID          string      `json:"id"`
	GameURL     string      `json:"gameUrl"`
	Source      string      `json:"source"`
	White       string      `json:"white"`
	Black       string      `json:"black"`
	Colour      string      `json:"colour"` // sa couleur dans la partie
	Speed       string      `json:"speed"`
	Date        string      `json:"date"`
	Played      int64       `json:"played"`
	Ply         int         `json:"ply"`
	Phase       string      `json:"phase"`
	FEN         string      `json:"fen"`
	PlayedUCI   string      `json:"playedUci"`
	PlayedSAN   string      `json:"playedSan"`
	Depth       int         `json:"depth"`
	Best        int         `json:"best"`
	PlayedScore int         `json:"playedScore"`
	Loss        int         `json:"loss"`
	Tolerance   int         `json:"tolerance"`
	Moves       []Candidate `json:"moves"`

	// Tenu par le serveur, jamais par l'analyse.
	Theory   bool   `json:"theory"`
	Attempts int    `json:"attempts"`
	Solved   int    `json:"solved"`
	Last     string `json:"last,omitempty"`
}

// Verdict d'une réponse.
const (
	Correct = "correct"
	Wrong   = "wrong"
	Unknown = "unknown" // coup jamais évalué : ni juste ni faux, en attente du moteur
)

// Check corrige une réponse à partir des coups déjà évalués. Le coup de la
// partie est toujours faux : c'est lui qui a perdu le seuil.
func (e Exercise) Check(uci string) (string, *Candidate) {
	uci = strings.ToLower(uci)
	for i := range e.Moves {
		if e.Moves[i].UCI == uci {
			if e.Moves[i].OK && uci != e.PlayedUCI {
				return Correct, &e.Moves[i]
			}
			return Wrong, &e.Moves[i]
		}
	}
	if uci == e.PlayedUCI {
		return Wrong, nil
	}
	return Unknown, nil
}

const schema = `
CREATE TABLE IF NOT EXISTS ex_items (
	id       TEXT PRIMARY KEY,
	game_url TEXT NOT NULL,
	played   INTEGER NOT NULL,
	loss     INTEGER NOT NULL,
	phase    TEXT NOT NULL,
	speed    TEXT NOT NULL,
	data     TEXT NOT NULL,
	theory   INTEGER NOT NULL DEFAULT 0,
	attempts INTEGER NOT NULL DEFAULT 0,
	solved   INTEGER NOT NULL DEFAULT 0,
	last     TEXT NOT NULL DEFAULT '',
	created  INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_ex_items_game ON ex_items(game_url);
CREATE TABLE IF NOT EXISTS ex_games (
	url      TEXT PRIMARY KEY,
	analyzed INTEGER NOT NULL,
	count    INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS ex_pending (
	id      TEXT NOT NULL,
	uci     TEXT NOT NULL,
	created INTEGER NOT NULL,
	PRIMARY KEY (id, uci)
);
CREATE TABLE IF NOT EXISTS ex_runs (
	at      INTEGER NOT NULL,
	mode    TEXT NOT NULL,
	total   INTEGER NOT NULL,
	correct INTEGER NOT NULL,
	wrong   INTEGER NOT NULL,
	unknown INTEGER NOT NULL,
	seconds INTEGER NOT NULL
);`

// Store garde les exercices dans la base SQLite du site (DB_PATH).
type Store struct{ db *sql.DB }

// Open ouvre (et crée au besoin) les tables des exercices.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// ErrNotFound : exercice inconnu.
var ErrNotFound = errors.New("exercice introuvable")

// Import enregistre l'analyse d'une partie. Réanalyser une partie remplace ses
// exercices mais garde ce que l'entraînement a appris : étiquette théorique,
// essais, réussites.
func (s *Store) Import(gameURL string, list []Exercise) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	now := time.Now().Unix()
	for _, e := range list {
		e.Theory, e.Attempts, e.Solved, e.Last = false, 0, 0, ""
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(`INSERT INTO ex_items (id, game_url, played, loss, phase, speed, data, created)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (id) DO UPDATE SET loss = excluded.loss, phase = excluded.phase,
				speed = excluded.speed, data = excluded.data`,
			e.ID, gameURL, e.Played, e.Loss, e.Phase, e.Speed, string(data), now); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`INSERT INTO ex_games (url, analyzed, count) VALUES (?, ?, ?)
		ON CONFLICT (url) DO UPDATE SET analyzed = excluded.analyzed, count = excluded.count`,
		gameURL, now, len(list)); err != nil {
		return err
	}
	return tx.Commit()
}

// AnalyzedGames liste les parties déjà passées au moteur : l'analyseur les
// saute, ce qui rend une analyse interrompue reprenable.
func (s *Store) AnalyzedGames() ([]string, error) {
	rows, err := s.db.Query(`SELECT url FROM ex_games`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

const itemCols = `data, theory, attempts, solved, last`

func scanItem(sc interface{ Scan(...any) error }) (Exercise, error) {
	var data, last string
	var theory, attempts, solved int
	if err := sc.Scan(&data, &theory, &attempts, &solved, &last); err != nil {
		return Exercise{}, err
	}
	var e Exercise
	if err := json.Unmarshal([]byte(data), &e); err != nil {
		return Exercise{}, err
	}
	e.Theory, e.Attempts, e.Solved, e.Last = theory == 1, attempts, solved, last
	return e, nil
}

// Get renvoie un exercice.
func (s *Store) Get(id string) (Exercise, error) {
	e, err := scanItem(s.db.QueryRow(`SELECT `+itemCols+` FROM ex_items WHERE id = ?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Exercise{}, ErrNotFound
	}
	return e, err
}

// Filter choisit les exercices d'une manche.
type Filter struct {
	Phase    string // "", "opening", "middlegame", "endgame"
	Speed    string // "", "blitz", "rapid"…
	Unsolved bool   // seulement ceux jamais réussis
	Theory   bool   // seulement ceux marqués théoriques (pour les relire)
	Limit    int
}

// List tire des exercices au hasard. Sans Theory, les positions marquées
// théoriques ne reviennent jamais : c'est tout le sens de l'étiquette.
func (s *Store) List(f Filter) ([]Exercise, error) {
	cond := []string{"theory = ?"}
	args := []any{boolInt(f.Theory)}
	if f.Phase != "" {
		cond = append(cond, "phase = ?")
		args = append(args, f.Phase)
	}
	if f.Speed != "" {
		cond = append(cond, "speed = ?")
		args = append(args, f.Speed)
	}
	if f.Unsolved {
		cond = append(cond, "solved = 0")
	}
	limit := f.Limit
	if limit <= 0 || limit > 200 {
		limit = 60
	}
	rows, err := s.db.Query(`SELECT `+itemCols+` FROM ex_items WHERE `+strings.Join(cond, " AND ")+
		` ORDER BY random() LIMIT ?`, append(args, limit)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Exercise{}
	for rows.Next() {
		e, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Answer corrige une réponse et la compte. Un coup inconnu part dans la file du
// moteur : la prochaine analyse sur le PC tranchera.
func (s *Store) Answer(id, uci string) (Exercise, string, error) {
	e, err := s.Get(id)
	if err != nil {
		return e, "", err
	}
	verdict, _ := e.Check(uci)
	solved := 0
	if verdict == Correct {
		solved = 1
	}
	if _, err := s.db.Exec(`UPDATE ex_items SET attempts = attempts + 1, solved = solved + ?, last = ? WHERE id = ?`,
		solved, verdict, id); err != nil {
		return e, "", err
	}
	if verdict == Unknown {
		if _, err := s.db.Exec(`INSERT OR IGNORE INTO ex_pending (id, uci, created) VALUES (?, ?, ?)`,
			id, strings.ToLower(uci), time.Now().Unix()); err != nil {
			return e, "", err
		}
	}
	e.Attempts++
	e.Solved += solved
	e.Last = verdict
	return e, verdict, nil
}

// SetTheory pose ou retire l'étiquette « théorique / choix conscient ».
func (s *Store) SetTheory(id string, theory bool) error {
	res, err := s.db.Exec(`UPDATE ex_items SET theory = ? WHERE id = ?`, boolInt(theory), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrNotFound
	}
	return nil
}

// Pending est un coup proposé en réponse que le moteur n'a jamais vu.
type Pending struct {
	ID  string `json:"id"`
	FEN string `json:"fen"`
	UCI string `json:"uci"`
}

// Pending liste la file du moteur.
func (s *Store) Pending() ([]Pending, error) {
	rows, err := s.db.Query(`SELECT p.id, p.uci, i.data FROM ex_pending p JOIN ex_items i ON i.id = p.id ORDER BY p.created`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Pending{}
	for rows.Next() {
		var p Pending
		var data string
		if err := rows.Scan(&p.ID, &p.UCI, &data); err != nil {
			return nil, err
		}
		var e Exercise
		if err := json.Unmarshal([]byte(data), &e); err != nil {
			return nil, err
		}
		p.FEN = e.FEN
		out = append(out, p)
	}
	return out, rows.Err()
}

// Resolve range l'évaluation d'un coup de la file parmi les coups de
// l'exercice, juste ou faux selon la tolérance de l'analyse.
func (s *Store) Resolve(id string, c Candidate) error {
	e, err := s.Get(id)
	if err != nil {
		return err
	}
	tol := e.Tolerance
	if tol <= 0 {
		tol = 15
	}
	c.UCI = strings.ToLower(c.UCI)
	c.OK = c.UCI != e.PlayedUCI && e.Best-c.Score <= tol
	replaced := false
	for i := range e.Moves {
		if e.Moves[i].UCI == c.UCI {
			e.Moves[i], replaced = c, true
		}
	}
	if !replaced {
		e.Moves = append(e.Moves, c)
	}
	e.Theory, e.Attempts, e.Solved, e.Last = false, 0, 0, ""
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.Exec(`UPDATE ex_items SET data = ? WHERE id = ?`, string(data), id); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM ex_pending WHERE id = ? AND uci = ?`, id, c.UCI); err != nil {
		return err
	}
	return tx.Commit()
}

// Run est le bilan d'une manche.
type Run struct {
	At      int64  `json:"at"`
	Mode    string `json:"mode"`
	Total   int    `json:"total"`
	Correct int    `json:"correct"`
	Wrong   int    `json:"wrong"`
	Unknown int    `json:"unknown"`
	Seconds int    `json:"seconds"`
}

// SaveRun garde le bilan d'une manche.
func (s *Store) SaveRun(r Run) error {
	_, err := s.db.Exec(`INSERT INTO ex_runs (at, mode, total, correct, wrong, unknown, seconds) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		time.Now().Unix(), r.Mode, r.Total, r.Correct, r.Wrong, r.Unknown, r.Seconds)
	return err
}

// Stats résume la réserve d'exercices et les dernières manches.
type Stats struct {
	Total    int            `json:"total"`
	Theory   int            `json:"theory"`
	Solved   int            `json:"solved"`
	Pending  int            `json:"pending"`
	Games    int            `json:"games"`
	ByPhase  map[string]int `json:"byPhase"`
	LastRuns []Run          `json:"lastRuns"`
}

// Stats renvoie le résumé.
func (s *Store) Stats() (Stats, error) {
	st := Stats{ByPhase: map[string]int{}, LastRuns: []Run{}}
	if err := s.db.QueryRow(`SELECT COUNT(*), COALESCE(SUM(theory),0), COALESCE(SUM(solved > 0),0) FROM ex_items`).
		Scan(&st.Total, &st.Theory, &st.Solved); err != nil {
		return st, err
	}
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM ex_pending`).Scan(&st.Pending)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM ex_games`).Scan(&st.Games)
	rows, err := s.db.Query(`SELECT phase, COUNT(*) FROM ex_items WHERE theory = 0 GROUP BY phase`)
	if err != nil {
		return st, err
	}
	for rows.Next() {
		var p string
		var n int
		if err := rows.Scan(&p, &n); err != nil {
			rows.Close()
			return st, err
		}
		st.ByPhase[p] = n
	}
	rows.Close()
	runs, err := s.db.Query(`SELECT at, mode, total, correct, wrong, unknown, seconds FROM ex_runs ORDER BY at DESC LIMIT 20`)
	if err != nil {
		return st, err
	}
	defer runs.Close()
	for runs.Next() {
		var r Run
		if err := runs.Scan(&r.At, &r.Mode, &r.Total, &r.Correct, &r.Wrong, &r.Unknown, &r.Seconds); err != nil {
			return st, err
		}
		st.LastRuns = append(st.LastRuns, r)
	}
	return st, runs.Err()
}
