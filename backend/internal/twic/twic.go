package twic

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// DefaultLast est le numéro à partir duquel on reprend : l'importeur commence
// au suivant, donc au 1640.
//
// Ce n'est QU'UNE VALEUR DE DÉFAUT, utilisée tant que `meta.twic_last` est
// absent de la base. Conséquence à ne pas manquer : si l'indexeur qui construit
// mega.db écrit cette clé, c'est elle qui gagne et le rattrapage n'a pas lieu.
//
// Reprendre avant ce que contient déjà la MegaBase est délibéré, mais ça ne
// tient que par l'empreinte : les lignes déjà présentes doivent porter le
// `hash` de games.Key, sinon ces numéros rentrent une deuxième fois.
const DefaultLast = 1639

// MetaLast est la clé qui porte le curseur dans mega.db.
const MetaLast = "twic_last"

// metaChecked retient la dernière tentative, réussie ou non. Elle ne sert qu'à
// éviter de harceler theweekinchess.com : le serveur redémarre à chaque
// publication d'article, et un rattrapage au démarrage sans garde-fou ferait
// une requête par redéploiement.
const metaChecked = "twic_checked"

// maxCatchUp borne un rattrapage. Assez large pour que la reprise depuis 1640
// se fasse en UN passage — la couper en deux laisserait la base à moitié à jour
// entre deux mardis, ce qui est le pire des deux mondes. La borne n'est là que
// pour qu'un curseur aberrant ne déclenche pas le téléchargement de dix ans
// d'archives : la sortie normale reste le 404 du numéro pas encore paru.
const maxCatchUp = 40

// Importer va chercher les livraisons manquantes et les fait entrer dans
// mega.db. Une seule exécution à la fois : le déclenchement manuel depuis
// /admin et le réveil du mardi peuvent tomber ensemble.
type Importer struct {
	Path   string // chemin de mega.db ; vide = importeur désactivé
	Client *http.Client
	// BaseURL n'existe que pour les tests : ils servent de vraies archives ZIP
	// depuis un httptest plutôt que de tirer sur theweekinchess.com.
	BaseURL string
	// Before et After encadrent l'écriture. Le serveur s'en sert pour fermer sa
	// connexion en lecture seule puis la rouvrir : sous Windows, un import qui
	// tourne pendant un téléversement empêcherait le renommage de la base.
	Before func()
	After  func()

	mu      sync.Mutex
	running bool
}

// Report dit ce qu'a fait un passage. C'est ce que /admin affiche.
type Report struct {
	Last    int      `json:"last"`
	Issues  []Issue  `json:"issues"`
	Added   int      `json:"added"`
	Skipped int      `json:"skipped"`
	Notes   []string `json:"notes,omitempty"`
}

// Issue est une livraison hebdomadaire traitée.
type Issue struct {
	Number  int `json:"number"`
	Games   int `json:"games"`
	Added   int `json:"added"`
	Skipped int `json:"skipped"`
	// Variants compte les parties écartées à dessein (Chess960, positions
	// imposées) ; Rejected celles qu'on n'a PAS su lire. Les distinguer est tout
	// l'intérêt du compteur : le premier est du tri, le second une alerte.
	Variants int `json:"variants"`
	Rejected int `json:"rejected"`
}

var errNotPublished = errors.New("numéro pas encore publié")

// ErrBusy est renvoyée quand un import tourne déjà. Ce n'est pas une panne :
// c'est le cas normal d'un double clic ou d'un réveil pendant un import manuel.
var ErrBusy = errors.New("twic: un import est déjà en cours")

func (im *Importer) Enabled() bool { return im.Path != "" }

func (im *Importer) client() *http.Client {
	if im.Client != nil {
		return im.Client
	}
	// Un numéro pèse quelques mégaoctets ; la minute laisse de la marge sur une
	// liaison lente sans jamais laisser une goroutine pendue pour la nuit.
	return &http.Client{Timeout: 2 * time.Minute}
}

// Run récupère tout ce qui manque. Idempotente : appelée deux fois de suite,
// la seconde ne trouve rien à faire.
func (im *Importer) Run(ctx context.Context) (Report, error) {
	if !im.Enabled() {
		return Report{}, errors.New("twic: GAMES_DB non configuré")
	}
	im.mu.Lock()
	if im.running {
		im.mu.Unlock()
		return Report{}, ErrBusy
	}
	im.running = true
	im.mu.Unlock()
	defer func() {
		im.mu.Lock()
		im.running = false
		im.mu.Unlock()
	}()

	if im.Before != nil {
		im.Before()
	}
	if im.After != nil {
		defer im.After()
	}

	w, err := openExisting(im.Path)
	if err != nil {
		return Report{}, err
	}
	defer w.Close()

	last := DefaultLast
	if v, err := w.MetaGet(MetaLast); err != nil {
		return Report{}, err
	} else if n, err := strconv.Atoi(v); err == nil && n > 0 {
		last = n
	}
	rep := Report{Last: last}
	_ = w.MetaSet(metaChecked, time.Now().UTC().Format(time.RFC3339))

	for i := 0; i < maxCatchUp; i++ {
		n := rep.Last + 1
		issue, err := im.importIssue(ctx, w, n)
		switch {
		case errors.Is(err, errNotPublished):
			return rep, nil // à jour : c'est la sortie normale
		case err != nil:
			// On rend ce qui est déjà entré. Le curseur a été avancé pour chaque
			// numéro réussi, donc le prochain passage reprend pile ici.
			rep.Notes = append(rep.Notes, fmt.Sprintf("TWIC %d : %v", n, err))
			return rep, err
		}
		rep.Issues = append(rep.Issues, issue)
		rep.Added += issue.Added
		rep.Skipped += issue.Skipped
		rep.Last = n
		if err := w.MetaSet(MetaLast, strconv.Itoa(n)); err != nil {
			return rep, err
		}
		slog.Info("TWIC importé", "numero", n, "parties", issue.Games,
			"ajoutees", issue.Added, "doublons", issue.Skipped)
	}
	rep.Notes = append(rep.Notes,
		fmt.Sprintf("arrêt après %d numéros — le reste au prochain passage", maxCatchUp))
	return rep, nil
}

func (im *Importer) importIssue(ctx context.Context, w *games.Writer, n int) (Issue, error) {
	raw, err := im.fetch(ctx, n)
	if err != nil {
		return Issue{}, err
	}
	pgn, err := extractPGN(raw)
	if err != nil {
		return Issue{}, err
	}
	list := readGames(bytes.NewReader(pgn))
	if len(list) == 0 {
		return Issue{}, fmt.Errorf("archive sans partie lisible")
	}

	tx, err := w.Begin()
	if err != nil {
		return Issue{}, err
	}
	out := Issue{Number: n, Games: len(list)}
	for _, g := range list {
		conv, err := convert(g)
		if errors.Is(err, errVariant) {
			out.Variants++
			continue
		}
		if err != nil {
			// Une partie qui ne se rejoue pas est perdue, pas fatale : elle vaut
			// mieux dehors que fausse dans l'arbre d'ouvertures.
			out.Rejected++
			continue
		}
		added, err := tx.Insert(conv)
		if err != nil {
			out.Rejected++
			continue
		}
		if added {
			out.Added++
		} else {
			out.Skipped++
		}
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return Issue{}, err
	}
	if out.Rejected > 0 {
		// Un warn, parce que c'est anormal : chaque ligne ici est une partie que
		// le lecteur de SAN n'a pas su rejouer, donc une piste de correction.
		slog.Warn("TWIC : parties illisibles", "numero", n, "nombre", out.Rejected,
			"sur", len(list))
	}
	if out.Variants > 0 {
		slog.Info("TWIC : variantes écartées", "numero", n, "nombre", out.Variants)
	}
	return out, nil
}

// fetch tire une archive. Un 404 veut dire « pas encore publié », ce qui est la
// réponse attendue tous les jours sauf le lendemain de la parution.
func (im *Importer) fetch(ctx context.Context, n int) ([]byte, error) {
	base := im.BaseURL
	if base == "" {
		base = "https://theweekinchess.com/zips"
	}
	url := fmt.Sprintf("%s/twic%dg.zip", base, n)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Se présenter : c'est un site tenu par une personne, pas une API.
	req.Header.Set("User-Agent",
		"iwanesko-chess/1.0 (+https://iwanesko-chess.ch ; mise à jour hebdomadaire)")
	resp, err := im.client().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errNotPublished
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	// archive/zip a besoin d'un io.ReaderAt : le flux doit de toute façon être
	// entièrement en mémoire. Le plafond évite qu'une redirection vers une page
	// d'erreur géante, ou un fichier inattendu, ne remplisse le conteneur.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 256<<20))
	if err != nil {
		return nil, err
	}
	return body, nil
}

// openExisting ouvre mega.db en écriture SANS la créer.
//
// games.OpenWriter pose le schéma, donc il fabriquerait une base vide si le
// chemin ne pointe sur rien — et l'explorateur se mettrait à répondre « aucune
// partie » au lieu de « aucune base », ce qui ne se diagnostique pas. Pire, le
// téléversement suivant sauvegarderait cette coquille en .bak par-dessus rien.
func openExisting(path string) (*games.Writer, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("twic: base de parties absente (%s)", path)
	}
	return games.OpenWriter(path)
}

func extractPGN(raw []byte) ([]byte, error) {
	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, fmt.Errorf("archive illisible (%w)", err)
	}
	for _, f := range zr.File {
		if !strings.HasSuffix(strings.ToLower(f.Name), ".pgn") {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(io.LimitReader(rc, 512<<20))
	}
	return nil, fmt.Errorf("aucun .pgn dans l'archive")
}
