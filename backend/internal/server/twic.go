package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/twic"
)

// Branchement de la mise à jour hebdomadaire.
//
// Le planificateur vit dans le paquet twic ; ce fichier ne fait que lui donner
// de quoi cohabiter avec le reste du serveur — fermer la base pendant qu'il
// écrit, et ne pas écrire pendant un téléversement.

func (s *Server) newImporter() *twic.Importer {
	if s.cfg.GamesDB == "" {
		return nil
	}
	target := uploadTarget{"games", s.cfg.GamesDB}
	return &twic.Importer{
		Path: s.cfg.GamesDB,
		// uploadMu sérialise l'import avec les téléversements : renommer un
		// fichier que l'importeur tient ouvert échoue sous Windows, et remplacer
		// la base sous ses pieds serait pire ailleurs.
		Before: func() {
			s.uploadMu.Lock()
			s.mu.Lock()
			if s.games != nil {
				_ = s.games.Close()
				s.games = nil
			}
			s.mu.Unlock()
		},
		After: func() {
			// L'explorateur reste muet (503) le temps de l'import, puis retrouve
			// la base — avec les parties de la semaine dedans.
			_ = s.reopen(target)
			s.uploadMu.Unlock()
		},
	}
}

// StartTWIC lance la boucle hebdomadaire. Sans base de parties configurée,
// c'est un no-op — comme l'annonceur sans SMTP.
func (s *Server) StartTWIC(ctx context.Context) {
	if s.twic == nil {
		return
	}
	s.twic.Start(ctx)
}

// handleTWICStatus dit où en est la base sans rien déclencher.
func (s *Server) handleTWICStatus(w http.ResponseWriter, r *http.Request) {
	if s.twic == nil {
		s.corpusError(w, "base de parties non configurée", http.StatusServiceUnavailable)
		return
	}
	s.corpusJSON(w, s.twic.Status())
}

// handleTWICRun force une mise à jour depuis /admin. Utile la veille d'un
// tournoi : on ne va pas attendre mardi pour avoir les parties de l'adversaire.
func (s *Server) handleTWICRun(w http.ResponseWriter, r *http.Request) {
	if s.twic == nil {
		s.corpusError(w, "base de parties non configurée", http.StatusServiceUnavailable)
		return
	}
	// Volontairement PAS le contexte de la requête : fermer l'onglet ne doit pas
	// interrompre un import à mi-chemin. La borne de temps vit ici.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	rep, err := s.twic.Run(ctx)
	switch {
	case errors.Is(err, twic.ErrBusy):
		s.corpusError(w, "un import est déjà en cours", http.StatusConflict)
	case err != nil:
		// Un code d'erreur, mais AVEC le rapport : un échec au troisième numéro
		// ne doit pas faire croire que les deux premiers sont perdus.
		privateHeaders(w)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(map[string]any{"error": err.Error(), "report": rep})
	default:
		s.corpusJSON(w, map[string]any{"report": rep})
	}
}
