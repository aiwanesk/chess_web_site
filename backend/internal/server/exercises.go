package server

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/exercises"
)

// Exercices tirés des parties d'Alexandre (/admin/exercices/). L'analyse tourne
// sur son PC (cmd/mistakes) et pousse ici ; le serveur stocke, sert les
// positions et corrige les réponses à partir des coups que Stockfish a déjà
// évalués. Il n'a pas de moteur : un coup jamais vu part dans une file que la
// prochaine analyse tranche.

func (s *Server) exercisesReady(w http.ResponseWriter) bool {
	if s.exercises == nil {
		s.corpusError(w, "exercices indisponibles (DB_PATH non configuré)", http.StatusServiceUnavailable)
		return false
	}
	return true
}

// exJSON lit un corps JSON. Même exigence que la préparation : le Content-Type
// JSON impose un pré-vol CORS, qu'un formulaire d'un autre site ne passe pas.
func exJSON(w http.ResponseWriter, r *http.Request, limit int64, v any) bool {
	mt, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if mt != "application/json" {
		http.Error(w, "JSON attendu", http.StatusUnsupportedMediaType)
		return false
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, limit)).Decode(v); err != nil {
		http.Error(w, "corps illisible", http.StatusBadRequest)
		return false
	}
	return true
}

func (s *Server) handleExercisesPage(w http.ResponseWriter, r *http.Request) {
	privateHeaders(w)
	nonce := newNonce()
	w.Header().Set("Content-Security-Policy", cspHeader(nonce))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := exercisesTmpl.Execute(w, map[string]any{
		"Pieces": corpusPieceSVG,
		"Nonce":  nonce,
		"Ready":  s.exercises != nil,
	}); err != nil {
		http.Error(w, "template", http.StatusInternalServerError)
	}
}

// ---- côté analyseur ---------------------------------------------------------

func (s *Server) handleExGames(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	list, err := s.exercises.AnalyzedGames()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, list)
}

func (s *Server) handleExImport(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	var req struct {
		Game      string               `json:"game"`
		Exercises []exercises.Exercise `json:"exercises"`
	}
	// Une partie de blitz bourrée d'imprécisions reste sous le mégaoctet.
	if !exJSON(w, r, 4<<20, &req) {
		return
	}
	if !strings.HasPrefix(req.Game, "https://") {
		s.corpusError(w, "adresse de partie invalide", http.StatusBadRequest)
		return
	}
	for _, e := range req.Exercises {
		if e.ID == "" || e.FEN == "" {
			s.corpusError(w, "exercice incomplet", http.StatusBadRequest)
			return
		}
		if _, err := corpus.LegalUCIs(e.FEN); err != nil {
			s.corpusError(w, "position illisible : "+e.FEN, http.StatusBadRequest)
			return
		}
	}
	if err := s.exercises.Import(req.Game, req.Exercises); err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, map[string]int{"imported": len(req.Exercises)})
}

func (s *Server) handleExPending(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	list, err := s.exercises.Pending()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, list)
}

func (s *Server) handleExResolve(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	var req struct {
		ID   string              `json:"id"`
		Move exercises.Candidate `json:"move"`
	}
	if !exJSON(w, r, 64<<10, &req) {
		return
	}
	if err := s.exercises.Resolve(req.ID, req.Move); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, exercises.ErrNotFound) {
			code = http.StatusNotFound
		}
		s.corpusError(w, err.Error(), code)
		return
	}
	s.corpusJSON(w, map[string]bool{"ok": true})
}

// ---- côté entraînement --------------------------------------------------------

// exPuzzle : ce que la page reçoit AVANT de répondre. Pas de solution dedans :
// on ne s'entraîne pas en lisant le réseau. Les coups légaux y sont, pour que
// le navigateur accepte un clic sans connaître les règles.
type exPuzzle struct {
	ID     string   `json:"id"`
	FEN    string   `json:"fen"`
	Legal  []string `json:"legal"`
	Phase  string   `json:"phase"`
	Speed  string   `json:"speed"`
	Date   string   `json:"date"`
	Theory bool     `json:"theory"`
}

func (s *Server) handleExBatch(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	q := r.URL.Query()
	f := exercises.Filter{Unsolved: q.Get("unsolved") == "1", Theory: q.Get("theory") == "1",
		Limit: atoiDefault(q.Get("n"), 60)}
	switch q.Get("phase") {
	case "opening", "middlegame", "endgame":
		f.Phase = q.Get("phase")
	}
	switch q.Get("speed") {
	case "bullet", "blitz", "rapid", "classical":
		f.Speed = q.Get("speed")
	}
	list, err := s.exercises.List(f)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	out := make([]exPuzzle, 0, len(list))
	for _, e := range list {
		legal, err := corpus.LegalUCIs(e.FEN)
		if err != nil {
			continue
		}
		out = append(out, exPuzzle{ID: e.ID, FEN: e.FEN, Legal: legal,
			Phase: e.Phase, Speed: e.Speed, Date: e.Date, Theory: e.Theory})
	}
	s.corpusJSON(w, out)
}

func (s *Server) handleExAnswer(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	var req struct {
		ID  string `json:"id"`
		UCI string `json:"uci"`
	}
	if !exJSON(w, r, 4<<10, &req) {
		return
	}
	e, err := s.exercises.Get(req.ID)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusNotFound)
		return
	}
	// Un coup illégal n'est pas une réponse : il ne compte pas comme un essai.
	if san, err := corpus.UCIToSAN(e.FEN, req.UCI); err != nil {
		s.corpusError(w, "coup illégal", http.StatusBadRequest)
		return
	} else if e, verdict, err := s.exercises.Answer(req.ID, req.UCI); err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
	} else {
		s.corpusJSON(w, map[string]any{"verdict": verdict, "san": san, "exercise": e})
	}
}

func (s *Server) handleExTag(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	var req struct {
		ID     string `json:"id"`
		Theory bool   `json:"theory"`
	}
	if !exJSON(w, r, 4<<10, &req) {
		return
	}
	if err := s.exercises.SetTheory(req.ID, req.Theory); err != nil {
		s.corpusError(w, err.Error(), http.StatusNotFound)
		return
	}
	s.corpusJSON(w, map[string]bool{"theory": req.Theory})
}

func (s *Server) handleExReveal(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	e, err := s.exercises.Get(r.URL.Query().Get("id"))
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusNotFound)
		return
	}
	s.corpusJSON(w, e)
}

func (s *Server) handleExRun(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	var run exercises.Run
	if !exJSON(w, r, 4<<10, &run) {
		return
	}
	if run.Total <= 0 {
		s.corpusJSON(w, map[string]bool{"ok": true}) // manche vide : rien à garder
		return
	}
	if err := s.exercises.SaveRun(run); err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, map[string]bool{"ok": true})
}

func (s *Server) handleExStats(w http.ResponseWriter, r *http.Request) {
	if !s.exercisesReady(w) {
		return
	}
	st, err := s.exercises.Stats()
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, st)
}
