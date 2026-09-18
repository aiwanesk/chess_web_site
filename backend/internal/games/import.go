package games

import (
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

// Écriture dans mega.db. Le gros de la base est construit hors ligne, sur le PC
// d'Alexandre ; ce qui passe ici, ce sont les mises à jour hebdomadaires TWIC —
// quelques milliers de parties par semaine contre plusieurs millions au départ.
//
// Un Writer ouvre la MÊME base que Store, mais en écriture. Les deux coexistent
// pendant un import : le mode WAL laisse l'explorateur lire pendant qu'on écrit.

// ImportGame est une partie prête à insérer, telle que la produit un lecteur de
// PGN : les noms sont bruts (c'est le Writer qui normalise), uci et san sont des
// listes de coups séparés par des espaces et DE MÊME LONGUEUR.
type ImportGame struct {
	White, Black       string
	Event              string
	Date               string // « 2026.09.15 » tel qu'écrit dans le PGN
	Result             int    // 1 Blancs, 0 nulle, -1 Noirs
	ECO                string
	WhiteElo, BlackElo int
	UCI, SAN           []string
}

// Key est l'empreinte de dédoublonnage.
//
// Elle ne porte QUE ce qui identifie une partie indépendamment de sa source :
// les deux joueurs, la date, le résultat et les coups. Volontairement pas
// l'événement ni les Elo — le même tournoi est écrit « Titled Tuesday Blitz »
// ici et « chess.com Titled Tuesday » là, et une partie reprise d'une semaine
// sur l'autre par TWIC ne doit pas rentrer deux fois pour si peu.
//
// L'indexeur qui construit mega.db doit calculer la même chose, sinon les
// parties déjà présentes seront réimportées. En Python :
//
//	raw = "\x1f".join([white, black, date, str(result), " ".join(san)])
//	h = hashlib.sha256(raw.encode()).digest()[:8]
//	key = int.from_bytes(h, "big", signed=True)
func Key(white, black, date string, result int, san []string) int64 {
	raw := strings.Join([]string{
		white, black, date, strconv.Itoa(result), strings.Join(san, " "),
	}, "\x1f")
	sum := sha256.Sum256([]byte(raw))
	// SQLite ne connaît que des entiers signés : on lit les 8 octets comme tels
	// plutôt que de tronquer un uint64 qui déborderait à l'insertion.
	return int64(binary.BigEndian.Uint64(sum[:8]))
}

func (g ImportGame) Key() int64 { return Key(g.White, g.Black, g.Date, g.Result, g.SAN) }

// Year extrait l'année d'une date PGN, 0 si elle est inconnue (« ???? »).
func (g ImportGame) Year() int {
	if len(g.Date) < 4 {
		return 0
	}
	n, err := strconv.Atoi(g.Date[:4])
	if err != nil {
		return 0
	}
	return n
}

// Writer insère des parties. Les identifiants de joueurs et d'événements sont
// mis en cache : un tournoi TWIC, c'est quelques centaines de noms répétés sur
// des milliers de parties, et faire un SELECT par ligne triplerait le temps.
type Writer struct {
	db      *sql.DB
	players map[string]int64
	events  map[string]int64
}

// OpenWriter ouvre mega.db en écriture et pose le schéma si besoin.
func OpenWriter(path string) (*Writer, error) {
	db, err := sql.Open("sqlite", "file:"+path+
		"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(15000)&_pragma=synchronous(NORMAL)")
	if err != nil {
		return nil, err
	}
	// Une base construite avant l'ajout de l'empreinte n'a pas la colonne, et
	// CREATE TABLE IF NOT EXISTS ne la rattrape pas. L'erreur « duplicate
	// column » est le cas NORMAL : elle dit juste que la colonne est déjà là.
	if _, err := db.Exec(`ALTER TABLE game ADD COLUMN hash INTEGER`); err != nil &&
		!strings.Contains(err.Error(), "duplicate column") &&
		!strings.Contains(err.Error(), "no such table") {
		_ = db.Close()
		return nil, fmt.Errorf("games: migration de l'empreinte impossible (%w)", err)
	}
	if _, err := db.Exec(Schema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("games: schéma refusé (%w)", err)
	}
	return &Writer{db: db, players: map[string]int64{}, events: map[string]int64{}}, nil
}

func (w *Writer) Close() error { return w.db.Close() }

// MetaGet et MetaSet portent l'état de l'importeur — le dernier numéro TWIC
// avalé, notamment. Le ranger DANS la base plutôt qu'à côté a une conséquence
// utile : un téléversement de mega.db depuis /admin remplace aussi ce curseur,
// donc une base reconstruite sur le PC repart du bon numéro sans rien régler.
func (w *Writer) MetaGet(key string) (string, error) {
	var v string
	err := w.db.QueryRow(`SELECT value FROM meta WHERE key = ?`, key).Scan(&v)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return v, err
}

func (w *Writer) MetaSet(key, value string) error {
	_, err := w.db.Exec(
		`INSERT INTO meta(key, value) VALUES(?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`, key, value)
	return err
}

// Tx regroupe l'import d'un numéro TWIC en une seule transaction : soit le
// numéro entre en entier, soit rien — un import coupé à mi-chemin laisserait
// sinon `twic_last` en désaccord avec le contenu réel.
type Tx struct {
	w  *Writer
	tx *sql.Tx
}

func (w *Writer) Begin() (*Tx, error) {
	tx, err := w.db.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{w: w, tx: tx}, nil
}

func (t *Tx) Rollback() {
	_ = t.tx.Rollback()
	// Les caches renvoient des identifiants qui viennent d'être annulés : les
	// garder ferait pointer les imports suivants sur des lignes inexistantes.
	t.w.players = map[string]int64{}
	t.w.events = map[string]int64{}
}

func (t *Tx) Commit() error { return t.tx.Commit() }

func (t *Tx) playerID(name string) (int64, error) {
	if id, ok := t.w.players[name]; ok {
		return id, nil
	}
	norm := Normalize(name)
	var id int64
	err := t.tx.QueryRow(`SELECT id FROM player WHERE name = ?`, name).Scan(&id)
	if err == sql.ErrNoRows {
		res, err := t.tx.Exec(`INSERT INTO player(name, norm) VALUES(?, ?)`, name, norm)
		if err != nil {
			return 0, err
		}
		if id, err = res.LastInsertId(); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	t.w.players[name] = id
	return id, nil
}

func (t *Tx) eventID(name string) (int64, error) {
	if name == "" {
		return 0, nil
	}
	if id, ok := t.w.events[name]; ok {
		return id, nil
	}
	var id int64
	err := t.tx.QueryRow(`SELECT id FROM event WHERE name = ?`, name).Scan(&id)
	if err == sql.ErrNoRows {
		tt := 0
		if TitledTuesday(name) {
			tt = 1
		}
		res, err := t.tx.Exec(
			`INSERT INTO event(name, norm, titled_tuesday) VALUES(?, ?, ?)`,
			name, Normalize(name), tt)
		if err != nil {
			return 0, err
		}
		if id, err = res.LastInsertId(); err != nil {
			return 0, err
		}
	} else if err != nil {
		return 0, err
	}
	t.w.events[name] = id
	return id, nil
}

// Insert ajoute une partie et dit si elle était nouvelle. Le doublon n'est pas
// une erreur : TWIC republie régulièrement des parties déjà parues, et l'index
// UNIQUE sur l'empreinte est là précisément pour absorber ça sans bruit.
func (t *Tx) Insert(g ImportGame) (bool, error) {
	if len(g.UCI) != len(g.SAN) {
		return false, fmt.Errorf("games: %d coups UCI pour %d SAN", len(g.UCI), len(g.SAN))
	}
	if g.White == "" || g.Black == "" || len(g.SAN) == 0 {
		return false, fmt.Errorf("games: partie sans joueurs ou sans coups")
	}
	wID, err := t.playerID(g.White)
	if err != nil {
		return false, err
	}
	bID, err := t.playerID(g.Black)
	if err != nil {
		return false, err
	}
	eID, err := t.eventID(g.Event)
	if err != nil {
		return false, err
	}
	var event any
	if eID != 0 {
		event = eID
	}
	res, err := t.tx.Exec(`
		INSERT OR IGNORE INTO game
			(hash, white_id, black_id, event_id, year, date, result, eco,
			 white_elo, black_elo, uci, san)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`,
		g.Key(), wID, bID, event, g.Year(), g.Date, g.Result, g.ECO,
		nullInt(g.WhiteElo), nullInt(g.BlackElo),
		strings.Join(g.UCI, " "), strings.Join(g.SAN, " "))
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

func nullInt(n int) any {
	if n <= 0 {
		return nil
	}
	return n
}
