package server

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/games"
	"github.com/iwanesko/chess-web-site/backend/internal/online"
)

// Préparation en ligne : les parties d'un pseudo Lichess ou Chess.com, chargées
// à la demande et explorées comme la base. Mêmes réponses que l'explorateur de
// parties (arbre, liste, partie), pour que la page partage tout son code de
// navigation avec lui. Seuls les favoris touchent le disque.

const (
	prepMaxAccounts = 4
	prepMaxGames    = 1000
)

// prepJSONBody refuse tout POST qui n'est pas du JSON. L'admin est protégé par
// Basic Auth, que le navigateur renvoie tout seul : sans cette exigence, un
// formulaire sur un autre site pourrait déclencher des chargements. Un
// Content-Type JSON impose un pré-vol CORS, qu'un autre site ne passe pas.
func prepJSONBody(w http.ResponseWriter, r *http.Request, v any) bool {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		http.Error(w, "JSON attendu", http.StatusUnsupportedMediaType)
		return false
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(v); err != nil {
		http.Error(w, "corps illisible", http.StatusBadRequest)
		return false
	}
	return true
}

func (s *Server) handlePrepPage(w http.ResponseWriter, r *http.Request) {
	privateHeaders(w)
	nonce := newNonce()
	w.Header().Set("Content-Security-Policy", cspHeader(nonce))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := prepTmpl.Execute(w, map[string]any{
		"Pieces":       corpusPieceSVG,
		"Start":        corpusStartFEN,
		"Nonce":        nonce,
		"HasFavorites": s.favorites != nil,
	}); err != nil {
		http.Error(w, "template", http.StatusInternalServerError)
	}
}

// ---- chargement -------------------------------------------------------------

type prepLoadRequest struct {
	Accounts []online.Account `json:"accounts"`
	Max      int              `json:"max"`
	Speeds   []string         `json:"speeds"`
	Months   int              `json:"months"` // 0 = pas de borne
	Rated    bool             `json:"rated"`
}

func (s *Server) handlePrepLoad(w http.ResponseWriter, r *http.Request) {
	var req prepLoadRequest
	if !prepJSONBody(w, r, &req) {
		return
	}
	if len(req.Accounts) == 0 || len(req.Accounts) > prepMaxAccounts {
		s.corpusError(w, "entre 1 et 4 comptes", http.StatusBadRequest)
		return
	}
	if req.Max <= 0 {
		req.Max = 300
	}
	if req.Max > prepMaxGames {
		req.Max = prepMaxGames
	}
	speeds := map[string]bool{}
	for _, sp := range req.Speeds {
		for _, known := range online.Speeds {
			if sp == known {
				speeds[sp] = true
			}
		}
	}
	var since time.Time
	if req.Months > 0 && req.Months <= 120 {
		since = time.Now().AddDate(0, -req.Months, 0)
	}

	var queries []online.Query
	for _, a := range req.Accounts {
		src, ok := online.ParseSource(string(a.Source))
		a.Username = strings.TrimSpace(a.Username)
		if !ok || !online.ValidUsername(a.Username) {
			s.corpusError(w, "pseudo invalide : "+a.Username, http.StatusBadRequest)
			return
		}
		queries = append(queries, online.Query{
			Account: online.Account{Source: src, Username: a.Username},
			Max:     req.Max, Speeds: speeds, Since: since, Rated: req.Rated,
		})
	}
	set, err := s.prepSets.Load(queries)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.corpusJSON(w, set.Status())
}

// prepSet retrouve le lot demandé, ou répond 404 : un lot expiré se recharge
// d'un clic, la page le dit.
func (s *Server) prepSet(w http.ResponseWriter, r *http.Request) *online.Set {
	set := s.prepSets.Get(r.URL.Query().Get("set"))
	if set == nil {
		s.corpusError(w, "recherche expirée : relance le chargement", http.StatusNotFound)
	}
	return set
}

func prepFilter(r *http.Request) online.Filter {
	q := r.URL.Query()
	f := online.Filter{}
	switch q.Get("colour") {
	case "w", "b":
		f.Colour = q.Get("colour")
	}
	if p := strings.TrimSpace(q.Get("path")); p != "" {
		f.Path = strings.Split(p, ",")
	}
	return f
}

func (s *Server) handlePrepStatus(w http.ResponseWriter, r *http.Request) {
	if set := s.prepSet(w, r); set != nil {
		s.corpusJSON(w, set.Status())
	}
}

// ---- lecture : mêmes formes que /admin/parties/api --------------------------

func (s *Server) handlePrepTree(w http.ResponseWriter, r *http.Request) {
	set := s.prepSet(w, r)
	if set == nil {
		return
	}
	f := prepFilter(r)
	depth := atoiDefault(r.URL.Query().Get("depth"), 14)
	if depth > 30 {
		depth = 30
	}
	tree := set.Tree(f, games.TreeOptions{MaxDepth: depth})

	root := corpusStartFEN
	for _, uci := range f.Path {
		next, err := corpus.ApplyUCI(root, uci)
		if err != nil {
			s.corpusError(w, "chemin invalide : "+err.Error(), http.StatusBadRequest)
			return
		}
		root = next
	}
	if err := fillFENs(tree, root); err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, map[string]any{"fen": root, "moves": tree})
}

func (s *Server) handlePrepGames(w http.ResponseWriter, r *http.Request) {
	set := s.prepSet(w, r)
	if set == nil {
		return
	}
	limit := atoiDefault(r.URL.Query().Get("limit"), 80)
	if limit > 500 {
		limit = 500
	}
	s.corpusJSON(w, set.Games(prepFilter(r), limit))
}

func (s *Server) handlePrepGame(w http.ResponseWriter, r *http.Request) {
	set := s.prepSet(w, r)
	if set == nil {
		return
	}
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		s.corpusError(w, "identifiant invalide", http.StatusBadRequest)
		return
	}
	g, ok := set.Game(id)
	if !ok {
		s.corpusError(w, "partie introuvable", http.StatusNotFound)
		return
	}
	type ply struct {
		SAN string `json:"san"`
		UCI string `json:"uci"`
		FEN string `json:"fen"`
	}
	plies := make([]ply, 0, len(g.UCI))
	fen := corpusStartFEN
	for i, uci := range g.UCI {
		next, err := corpus.ApplyUCI(fen, uci)
		if err != nil {
			break
		}
		fen = next
		plies = append(plies, ply{SAN: g.SAN[i], UCI: uci, FEN: fen})
	}
	s.corpusJSON(w, map[string]any{
		"id": g.ID, "white": g.White, "black": g.Black,
		"whiteElo": g.WhiteElo, "blackElo": g.BlackElo,
		"date": g.Date, "speed": g.Speed, "url": g.URL, "source": g.Source,
		"eco": g.ECO, "opening": g.Opening, "result": g.Result, "plies": plies,
	})
}

// ---- favoris -----------------------------------------------------------------

func (s *Server) handlePrepFavorites(w http.ResponseWriter, r *http.Request) {
	if s.favorites == nil {
		s.corpusError(w, "favoris indisponibles (DB_PATH non configuré)", http.StatusServiceUnavailable)
		return
	}
	switch r.Method {
	case http.MethodGet:
		list, err := s.favorites.List()
		if err != nil {
			s.corpusError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.corpusJSON(w, list)
		return
	case http.MethodPost:
		var fav online.Favorite
		if !prepJSONBody(w, r, &fav) {
			return
		}
		src, ok := online.ParseSource(string(fav.Source))
		fav.Username = strings.TrimSpace(fav.Username)
		if !ok || !online.ValidUsername(fav.Username) {
			s.corpusError(w, "pseudo invalide", http.StatusBadRequest)
			return
		}
		fav.Source = src
		if err := s.favorites.Add(fav); err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, online.ErrTooManyFavorites) {
				code = http.StatusConflict
			}
			s.corpusError(w, err.Error(), code)
			return
		}
	case http.MethodDelete:
		q := r.URL.Query()
		src, ok := online.ParseSource(q.Get("source"))
		if !ok || !online.ValidUsername(q.Get("username")) {
			s.corpusError(w, "pseudo invalide", http.StatusBadRequest)
			return
		}
		if err := s.favorites.Remove(src, q.Get("username")); err != nil {
			s.corpusError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "méthode", http.StatusMethodNotAllowed)
		return
	}
	list, err := s.favorites.List()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, list)
}
