package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
)

// Explorateur de corpus : à une position donnée, ce qui s'y joue et ce que les
// cours en disent. Espace strictement privé — le contenu des commentaires
// appartient aux auteurs des cours, il est acheté et n'a rien à faire dehors.
//
// Découpage repris du prototype Python : le serveur est autoritaire sur les
// positions, le client ne fait que dessiner. Aucun moteur d'échecs en JS, donc
// aucun second endroit où diverger — et on ne peut jouer que les coups que le
// corpus connaît, ce qui est exactement ce qu'on veut d'un explorateur.

// privateHeaders : rien de tout ça ne doit être indexé ni mis en cache.
func privateHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func (s *Server) corpusJSON(w http.ResponseWriter, v any) {
	privateHeaders(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) corpusError(w http.ResponseWriter, msg string, code int) {
	privateHeaders(w)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// handleCorpusPos : ce que le corpus sait de la position demandée.
func (s *Server) handleCorpusPos(w http.ResponseWriter, r *http.Request) {
	fen := strings.TrimSpace(r.URL.Query().Get("fen"))
	if fen == "" {
		fen = corpusStartFEN
	}
	pos, err := s.corpus.Lookup(fen)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.corpusJSON(w, pos)
}

// handleCorpusGo joue un coup et renvoie la position suivante, déjà consultée.
// Un aller-retour au lieu de deux : sur un téléphone en 4G, ça se sent.
func (s *Server) handleCorpusGo(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	fen := strings.TrimSpace(q.Get("fen"))
	if fen == "" {
		fen = corpusStartFEN
	}
	uci := strings.TrimSpace(q.Get("uci"))

	// Le coup doit exister dans la position courante. C'est la seule validation
	// de légalité du système, et elle suffit : la base ne contient que des
	// coups joués dans de vraies parties.
	cur, err := s.corpus.Lookup(fen)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusBadRequest)
		return
	}
	san := ""
	for _, m := range cur.Moves {
		if m.UCI == uci {
			san = m.SAN
			break
		}
	}
	if san == "" {
		s.corpusError(w, "coup absent du corpus pour cette position", http.StatusBadRequest)
		return
	}

	next, err := corpus.ApplyUCI(fen, uci)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusBadRequest)
		return
	}
	pos, err := s.corpus.Lookup(next)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, struct {
		*corpus.Position
		SAN string `json:"san"`
	}{pos, san})
}

func (s *Server) handleCorpusMeta(w http.ResponseWriter, r *http.Request) {
	m, err := s.corpus.Meta()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, m)
}

const corpusStartFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func (s *Server) handleCorpusPage(w http.ResponseWriter, r *http.Request) {
	privateHeaders(w)
	// Le script de la page est inline, et la politique par défaut est
	// `script-src 'self'` : sans nonce il est purement et simplement bloqué, en
	// silence côté serveur. On reprend le mécanisme déjà utilisé pour les
	// scripts d'hydratation plutôt que d'ouvrir 'unsafe-inline' au site entier.
	nonce := newNonce()
	w.Header().Set("Content-Security-Policy", cspHeader(nonce))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := corpusTmpl.Execute(w, map[string]any{
		"Pieces": corpusPieceSVG,
		"Start":  corpusStartFEN,
		"Nonce":  nonce,
	}); err != nil {
		http.Error(w, "template", http.StatusInternalServerError)
	}
}

// pieceJSON sérialise les SVG pour le script de la page : une seule table,
// injectée une fois, plutôt qu'une requête par pièce.
func pieceJSON(m map[byte]string) template.JS {
	out := map[string]string{}
	for k, v := range m {
		out[string(k)] = v
	}
	b, _ := json.Marshal(out)
	return template.JS(b)
}

var corpusTmpl = template.Must(template.New("corpus").Funcs(template.FuncMap{
	"pieces": pieceJSON,
}).Parse(corpusHTML))
