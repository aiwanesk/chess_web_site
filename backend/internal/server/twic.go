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
		//
		// En revanche on NE FERME PAS la base en lecture pendant l'import. C'est
		// précisément ce que le mode WAL permet — un écrivain, des lecteurs, en
		// même temps — et la fermer rendait l'explorateur muet pendant les deux
		// minutes du rattrapage initial, c'est-à-dire pile au moment où on vient
		// de téléverser et où on veut vérifier que ça marche.
		Before: func() { s.uploadMu.Lock() },
		After: func() {
			// Rouvrir après coup, en revanche : la connexion en lecture seule
			// avait été ouverte avant que l'importeur ne bascule le fichier en
			// WAL, et c'est le seul moyen sûr de la voir repartir propre.
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

const twicRunPath = "/admin/parties/twic"

// handleTWICRun force une mise à jour depuis /admin. Utile la veille d'un
// tournoi : on ne va pas attendre mardi pour avoir les parties de l'adversaire.
func (s *Server) handleTWICRun(w http.ResponseWriter, r *http.Request) {
	if s.twic == nil {
		s.corpusError(w, "base de parties non configurée", http.StatusServiceUnavailable)
		return
	}
	// Un rattrapage de plusieurs numéros dépasse les 30 s de WriteTimeout : sans
	// ça, l'import allait au bout mais le bouton affichait « échec ».
	noLimit(w)
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
