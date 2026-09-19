package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Explorateur de parties : même interface que l'explorateur de corpus —
// échiquier à gauche, coups à droite, on clique et ça avance — mais alimenté
// par mega.db au lieu du graphe de positions.
//
// Différence de fond dans la circulation des données. Le corpus interroge le
// serveur à chaque position, parce que ses 2,8 millions d'arêtes ne se
// téléchargent pas. Ici l'arbre d'un joueur tient en quelques dizaines de
// kilo-octets : on l'envoie EN ENTIER, et toute la navigation devient locale.
// Dans une salle de tournoi mal couverte, c'est la différence entre fluide et
// pénible.

// parseFilter lit le filtre depuis la requête. Tout vient du navigateur, donc
// rien n'est cru sur parole : identifiants entiers, années bornées, couleur
// dans un jeu fermé.
func parseFilter(r *http.Request) games.Filter {
	q := r.URL.Query()
	f := games.Filter{ExcludeTitledTuesday: q.Get("noTT") == "1"}

	for _, raw := range strings.Split(q.Get("players"), ",") {
		if id, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64); err == nil && id > 0 {
			f.PlayerIDs = append(f.PlayerIDs, id)
			if len(f.PlayerIDs) >= 20 {
				break
			}
		}
	}
	switch q.Get("colour") {
	case "w":
		f.Colour = games.White
	case "b":
		f.Colour = games.Black
	}
	f.FromYear = year(q.Get("from"))
	f.ToYear = year(q.Get("to"))
	if p := strings.TrimSpace(q.Get("path")); p != "" {
		f.Path = strings.Split(p, ",")
	}
	return f
}

// year écarte les valeurs qui ne peuvent pas être une année de partie : une
// borne absurde ne doit pas se transformer en requête absurde.
func year(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n < 1475 || n > 2200 {
		return 0
	}
	return n
}

func (s *Server) handlePartiesPlayers(w http.ResponseWriter, r *http.Request) {
	g := s.gamesStore()
	if !s.requireStore(w, g != nil, "base de parties") {
		return
	}
	found, err := g.SearchPlayers(r.URL.Query().Get("q"), 25)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if found == nil {
		found = []games.Player{}
	}
	s.corpusJSON(w, found)
}

// handlePartiesTree renvoie l'arbre d'ouverture, chaque nœud portant le FEN de
// sa position. Le client n'a donc aucune règle du jeu à connaître : il dessine
// le FEN du nœud où il se trouve.
func (s *Server) handlePartiesTree(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	q := r.URL.Query()
	opt := games.TreeOptions{
		Path:     f.Path,
		MaxDepth: atoiDefault(q.Get("depth"), 14),
		MinGames: atoiDefault(q.Get("min"), 1),
	}
	if opt.MaxDepth > 30 {
		opt.MaxDepth = 30
	}
	g := s.gamesStore()
	if !s.requireStore(w, g != nil, "base de parties") {
		return
	}
	tree, err := g.OpeningTree(f, opt)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// La position de départ de l'arbre : on rejoue le chemin demandé.
	root := corpusStartFEN
	for _, uci := range opt.Path {
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

// fillFENs descend l'arbre en appliquant les coups. C'est ici que vit la règle
// du jeu, pas dans le paquet games — qui ne fait que du SQL — ni dans le
// navigateur, qui n'a pas à la réimplémenter.
func fillFENs(nodes []games.Node, fen string) error {
	for i := range nodes {
		next, err := corpus.ApplyUCI(fen, nodes[i].UCI)
		if err != nil {
			// Une partie au score douteux ne doit pas faire tomber tout l'arbre :
			// on coupe la branche et on continue.
			nodes[i].Children = nil
			continue
		}
		nodes[i].FEN = next
		if err := fillFENs(nodes[i].Children, next); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) handlePartiesGames(w http.ResponseWriter, r *http.Request) {
	f := parseFilter(r)
	f.Limit = atoiDefault(r.URL.Query().Get("limit"), 60)
	g := s.gamesStore()
	if !s.requireStore(w, g != nil, "base de parties") {
		return
	}
	found, err := g.Search(f)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, found)
}

func (s *Server) handlePartiesMeta(w http.ResponseWriter, r *http.Request) {
	g := s.gamesStore()
	if !s.requireStore(w, g != nil, "base de parties") {
		return
	}
	m, err := g.Meta()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, m)
}

func (s *Server) handlePartiesPage(w http.ResponseWriter, r *http.Request) {
	privateHeaders(w)
	nonce := newNonce()
	w.Header().Set("Content-Security-Policy", cspHeader(nonce))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := partiesTmpl.Execute(w, map[string]any{
		"Pieces": corpusPieceSVG,
		"Start":  corpusStartFEN,
		"Nonce":  nonce,
	}); err != nil {
		http.Error(w, "template", http.StatusInternalServerError)
	}
}

func atoiDefault(s string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || n <= 0 {
		return def
	}
	return n
}

// handlePartiesGame renvoie une partie entière, avec le FEN APRÈS chaque
// demi-coup.
//
// Les positions sont calculées ici et pas dans le navigateur : c'est le même
// principe que pour l'arbre — la règle du jeu vit d'un seul côté. Quatre-vingts
// positions à produire, ça ne vaut pas une seconde implémentation des règles en
// JavaScript, avec ses propres bugs de roque et de prise en passant.
func (s *Server) handlePartiesGame(w http.ResponseWriter, r *http.Request) {
	g := s.gamesStore()
	if !s.requireStore(w, g != nil, "base de parties") {
		return
	}
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		s.corpusError(w, "identifiant invalide", http.StatusBadRequest)
		return
	}
	game, err := g.ByID(id)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusNotFound)
		return
	}

	type ply struct {
		SAN string `json:"san"`
		UCI string `json:"uci"`
		FEN string `json:"fen"`
	}
	sans, ucis := strings.Fields(game.SAN), strings.Fields(game.UCI)
	plies := make([]ply, 0, len(ucis))
	fen := corpusStartFEN
	for i, uci := range ucis {
		next, err := corpus.ApplyUCI(fen, uci)
		if err != nil {
			// Une partie abîmée s'arrête là où elle cesse d'être jouable, elle ne
			// renvoie pas une erreur : la moitié lisible vaut mieux que rien.
			break
		}
		fen = next
		san := uci
		if i < len(sans) {
			san = sans[i]
		}
		plies = append(plies, ply{SAN: san, UCI: uci, FEN: fen})
	}
	s.corpusJSON(w, map[string]any{
		"id": game.ID, "white": game.White, "black": game.Black,
		"whiteElo": game.WhiteElo, "blackElo": game.BlackElo,
		"event": game.Event, "date": game.Date, "year": game.Year,
		"eco": game.ECO, "result": game.Result, "plies": plies,
	})
}
