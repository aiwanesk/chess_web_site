// Package online prépare un adversaire à partir de ses parties en ligne.
//
// On tape un pseudo Lichess ou Chess.com, le serveur va chercher ses parties
// récentes et en construit l'arbre d'ouverture, exactement comme l'explorateur
// de la base. Rien n'est écrit sur le disque : les parties vivent en mémoire le
// temps de la préparation, puis disparaissent. Seuls les pseudos favoris sont
// gardés (voir favorites.go).
package online

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
)

// Source est le site d'où viennent les parties.
type Source string

const (
	Lichess  Source = "lichess"
	ChessCom Source = "chesscom"
)

// ParseSource n'accepte que les deux sites connus.
func ParseSource(s string) (Source, bool) {
	switch Source(strings.ToLower(strings.TrimSpace(s))) {
	case Lichess:
		return Lichess, true
	case ChessCom:
		return ChessCom, true
	}
	return "", false
}

// Les deux sites limitent les pseudos aux lettres, chiffres, tiret et souligné.
// Tout le reste est refusé AVANT d'être collé dans une URL.
var usernameRe = regexp.MustCompile(`^[A-Za-z0-9_-]{2,30}$`)

// ValidUsername dit si un pseudo peut être envoyé tel quel aux deux API.
func ValidUsername(u string) bool { return usernameRe.MatchString(u) }

// Account est un compte à charger.
type Account struct {
	Source   Source `json:"source"`
	Username string `json:"username"`
}

// Speeds regroupe les cadences des deux sites sous les mêmes noms.
var Speeds = []string{"bullet", "blitz", "rapid", "classical", "correspondence"}

// Query décrit ce qu'on va chercher pour un compte.
type Query struct {
	Account
	Max    int             // nombre maximal de parties
	Speeds map[string]bool // cadences retenues ; vide = toutes
	Since  time.Time       // zéro = pas de borne
	Rated  bool            // parties classées seulement
}

func (q Query) wants(speed string) bool { return len(q.Speeds) == 0 || q.Speeds[speed] }

// Game est une partie chargée. Result est vu des Blancs (+1, 0, -1).
type Game struct {
	ID       int      `json:"id"`
	Source   Source   `json:"source"`
	URL      string   `json:"url"`
	White    string   `json:"white"`
	Black    string   `json:"black"`
	WhiteElo int      `json:"whiteElo"`
	BlackElo int      `json:"blackElo"`
	Date     string   `json:"date"`
	Played   int64    `json:"played"` // instant de la partie, en secondes Unix
	Speed    string   `json:"speed"`
	Rated    bool     `json:"rated"`
	Result   int      `json:"result"`
	ECO      string   `json:"eco"`
	Opening  string   `json:"opening"`
	SAN      []string `json:"-"`
	UCI      []string `json:"-"`
}

// ErrNotFound : le pseudo n'existe pas sur ce site.
var ErrNotFound = errors.New("pseudo introuvable")

// ErrRateLimited : le site demande de patienter.
var ErrRateLimited = errors.New("le site limite les requêtes, réessaie dans une minute")

// Client va chercher les parties. Les adresses de base sont des champs pour que
// les tests puissent les remplacer par un faux serveur.
type Client struct {
	HTTP        *http.Client
	LichessURL  string
	ChessComURL string
}

// NewClient renvoie un client sur les vraies API.
func NewClient() *Client {
	return &Client{
		HTTP:        &http.Client{Timeout: 3 * time.Minute},
		LichessURL:  "https://lichess.org",
		ChessComURL: "https://api.chess.com",
	}
}

// userAgent : les deux sites demandent qu'un outil s'identifie.
const userAgent = "iwanesko.ch-prepa/1.0 (+https://iwanesko.ch)"

// Fetch charge les parties d'un compte, les plus récentes d'abord, et appelle
// onGame pour chacune. Une partie qui ne se rejoue pas est sautée.
func (c *Client) Fetch(ctx context.Context, q Query, onGame func(Game)) error {
	switch q.Source {
	case Lichess:
		return c.fetchLichess(ctx, q, onGame)
	case ChessCom:
		return c.fetchChessCom(ctx, q, onGame)
	}
	return fmt.Errorf("source inconnue : %q", q.Source)
}

func (c *Client) get(ctx context.Context, url, accept string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return resp, nil
	case http.StatusNotFound, http.StatusGone:
		resp.Body.Close()
		return nil, ErrNotFound
	case http.StatusTooManyRequests:
		resp.Body.Close()
		return nil, ErrRateLimited
	}
	resp.Body.Close()
	return nil, fmt.Errorf("%s a répondu %d", req.URL.Host, resp.StatusCode)
}

const startFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

// toUCI rejoue le SAN depuis la position initiale. Même générateur de coups
// légaux que l'import TWIC : l'arbre ne fait confiance qu'à ce qui se rejoue.
func toUCI(san []string) ([]string, error) {
	out := make([]string, 0, len(san))
	fen := startFEN
	for _, m := range san {
		uci, err := corpus.SANToUCI(fen, m)
		if err != nil {
			return nil, err
		}
		next, err := corpus.ApplyUCI(fen, uci)
		if err != nil {
			return nil, err
		}
		out = append(out, uci)
		fen = next
	}
	return out, nil
}
