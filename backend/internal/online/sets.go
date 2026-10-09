package online

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Un « lot » est le résultat d'une recherche : les parties d'un ou plusieurs
// comptes, gardées en mémoire le temps de la préparation. Rien n'est écrit sur
// le disque ; au-delà de quelques lots ou de quelques heures, les plus anciens
// sont simplement oubliés.

const (
	maxSets    = 8
	setTTL     = 3 * time.Hour
	loadBudget = 4 * time.Minute
)

// State est l'avancement du chargement d'un lot.
type State string

const (
	Loading State = "loading"
	Ready   State = "ready"
	Failed  State = "failed"
)

// Set est un lot de parties.
type Set struct {
	ID       string
	Accounts []Account

	mu      sync.RWMutex
	state   State
	err     string
	games   []Game
	created time.Time
}

// Status est ce que le navigateur interroge pendant le chargement.
type Status struct {
	ID       string    `json:"id"`
	State    State     `json:"state"`
	Games    int       `json:"games"`
	Error    string    `json:"error,omitempty"`
	Accounts []Account `json:"accounts"`
}

// Status renvoie l'avancement.
func (s *Set) Status() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return Status{ID: s.ID, State: s.state, Games: len(s.games), Error: s.err, Accounts: s.Accounts}
}

// Sets garde les lots en mémoire.
type Sets struct {
	client *Client
	mu     sync.Mutex
	sets   map[string]*Set
}

// NewSets prépare la mémoire des lots.
func NewSets(c *Client) *Sets { return &Sets{client: c, sets: map[string]*Set{}} }

// Get renvoie un lot, ou nil s'il n'existe pas (ou plus).
func (ss *Sets) Get(id string) *Set {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	return ss.sets[id]
}

// Load démarre le chargement d'un lot et rend la main tout de suite : Lichess
// n'envoie qu'une vingtaine de parties par seconde, et le serveur coupe toute
// réponse à 30 s. Le navigateur suit l'avancement avec Status.
func (ss *Sets) Load(queries []Query) (*Set, error) {
	if len(queries) == 0 {
		return nil, errors.New("aucun compte")
	}
	set := &Set{ID: newID(), state: Loading, created: time.Now()}
	for _, q := range queries {
		set.Accounts = append(set.Accounts, q.Account)
	}

	ss.mu.Lock()
	ss.evict()
	ss.sets[set.ID] = set
	ss.mu.Unlock()

	go ss.run(set, queries)
	return set, nil
}

func (ss *Sets) run(set *Set, queries []Query) {
	// Pas le contexte de la requête : fermer l'onglet n'a pas à couper un
	// chargement qu'on retrouvera en revenant.
	ctx, cancel := context.WithTimeout(context.Background(), loadBudget)
	defer cancel()

	var firstErr error
	loaded := 0
	for _, q := range queries {
		err := ss.client.Fetch(ctx, q, func(g Game) {
			set.mu.Lock()
			g.ID = len(set.games)
			set.games = append(set.games, g)
			set.mu.Unlock()
		})
		if err != nil {
			if firstErr == nil {
				firstErr = describe(q.Account, err)
			}
			continue
		}
		loaded++
	}

	set.mu.Lock()
	defer set.mu.Unlock()
	// Un compte introuvable ne gâche pas les autres ; l'erreur n'est fatale que
	// si rien n'est arrivé.
	switch {
	case loaded == 0 && firstErr != nil:
		set.state, set.err = Failed, firstErr.Error()
	default:
		set.state = Ready
		if firstErr != nil {
			set.err = firstErr.Error()
		}
	}
}

func describe(a Account, err error) error {
	site := "Lichess"
	if a.Source == ChessCom {
		site = "Chess.com"
	}
	return errors.New(site + " · " + a.Username + " : " + err.Error())
}

// evict oublie les lots expirés, puis les plus anciens au-delà de maxSets.
// À appeler sous ss.mu.
func (ss *Sets) evict() {
	now := time.Now()
	for id, s := range ss.sets {
		if now.Sub(s.created) > setTTL {
			delete(ss.sets, id)
		}
	}
	for len(ss.sets) >= maxSets {
		var oldest *Set
		for _, s := range ss.sets {
			if oldest == nil || s.created.Before(oldest.created) {
				oldest = s
			}
		}
		delete(ss.sets, oldest.ID)
	}
}

func newID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// ---- lecture d'un lot ------------------------------------------------------

// Filter restreint les parties d'un lot.
type Filter struct {
	Colour string   // "w", "b" ou "" : couleur jouée par le joueur préparé
	Speeds []string // vide = toutes
	Path   []string // en UCI
}

// side : la partie a-t-elle été jouée par l'un des comptes du lot, et avec
// quelle couleur ? On compare sans casse, sur le même site seulement.
func (s *Set) side(g Game) (white, ok bool) {
	for _, a := range s.Accounts {
		if a.Source != g.Source {
			continue
		}
		if strings.EqualFold(a.Username, g.White) {
			return true, true
		}
		if strings.EqualFold(a.Username, g.Black) {
			return false, true
		}
	}
	return false, false
}

func (s *Set) match(g Game, f Filter) (white, ok bool) {
	white, ok = s.side(g)
	if !ok {
		return false, false
	}
	if (f.Colour == "w" && !white) || (f.Colour == "b" && white) {
		return false, false
	}
	if len(f.Speeds) > 0 {
		hit := false
		for _, sp := range f.Speeds {
			if sp == g.Speed {
				hit = true
				break
			}
		}
		if !hit {
			return false, false
		}
	}
	return white, games.HasPrefix(g.UCI, f.Path)
}

// Tree construit l'arbre d'ouverture du lot, résultats vus du joueur préparé.
// Le chemin de départ est pris dans f.Path ; opt.Path est ignoré.
func (s *Set) Tree(f Filter, opt games.TreeOptions) []games.Node {
	opt.Path = f.Path
	t := games.NewTree(opt)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, g := range s.games {
		white, ok := s.match(g, f)
		if !ok {
			continue
		}
		pov := g.Result
		if !white {
			pov = -pov
		}
		t.Add(g.UCI, g.SAN, pov)
	}
	return t.Nodes()
}

// Games renvoie les parties du lot qui passent par f, les plus récentes
// d'abord.
func (s *Set) Games(f Filter, limit int) []Game {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Game{}
	for _, g := range s.games {
		if _, ok := s.match(g, f); ok {
			out = append(out, g)
		}
	}
	// Chaque compte arrive du plus récent au plus ancien ; plusieurs comptes se
	// mélangent, d'où un tri par date (stable, pour garder l'ordre d'un jour).
	sort.SliceStable(out, func(i, j int) bool { return out[i].Date > out[j].Date })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out
}

// Game renvoie une partie par son numéro dans le lot.
func (s *Set) Game(id int) (Game, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if id < 0 || id >= len(s.games) {
		return Game{}, false
	}
	return s.games[id], true
}
