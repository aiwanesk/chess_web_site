package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

// Téléversement des bases de l'espace privé.
//
// Deux contraintes dictent toute la forme.
//
// D'abord la taille : corpus.db pèse 253 Mo et mega.db plusieurs gigaoctets.
// Un seul POST se ferait couper par le timeout d'HAProxy, et un échec à 90 %
// obligerait à tout recommencer. D'où l'envoi PAR MORCEAUX, repris là où il
// s'est arrêté.
//
// Ensuite SQLite : écraser le fichier pendant que le serveur le lit le
// corrompt, et on ne s'en aperçoit qu'à la requête suivante — ou pire, la
// veille d'une ronde. D'où la séquence écrire à côté / VÉRIFIER / basculer par
// renommage, avec l'ancienne base gardée en .bak.

const chunkLimit = 16 << 20 // 16 Mo : le client en envoie 8, le reste est de la marge

// uploadTarget désigne une base. Jamais un chemin venu du client : il choisit
// une clé dans un jeu fermé, et c'est la configuration qui dit où ça écrit.
type uploadTarget struct {
	name string
	path string
}

func (s *Server) target(key string) (uploadTarget, bool) {
	switch key {
	case "corpus":
		if s.cfg.CorpusDB == "" {
			return uploadTarget{}, false
		}
		return uploadTarget{"corpus", s.cfg.CorpusDB}, true
	case "games":
		if s.cfg.GamesDB == "" {
			return uploadTarget{}, false
		}
		return uploadTarget{"games", s.cfg.GamesDB}, true
	}
	return uploadTarget{}, false
}

func (t uploadTarget) partial() string { return t.path + ".upload" }
func (t uploadTarget) backup() string  { return t.path + ".bak" }

func (s *Server) resolveTarget(w http.ResponseWriter, r *http.Request) (uploadTarget, bool) {
	t, ok := s.target(r.URL.Query().Get("target"))
	if !ok {
		s.corpusError(w, "base inconnue ou non configurée", http.StatusBadRequest)
		return uploadTarget{}, false
	}
	return t, true
}

// handleUploadStatus dit où en est un envoi en cours. C'est ce qui permet de
// reprendre : le client demande l'offset, puis continue à partir de là.
func (s *Server) handleUploadStatus(w http.ResponseWriter, r *http.Request) {
	t, ok := s.resolveTarget(w, r)
	if !ok {
		return
	}
	out := map[string]any{"target": t.name, "uploaded": int64(0)}
	if fi, err := os.Stat(t.partial()); err == nil {
		out["uploaded"] = fi.Size()
	}
	if fi, err := os.Stat(t.path); err == nil {
		out["live"] = fi.Size()
		out["liveModified"] = fi.ModTime().UTC().Format("2006-01-02 15:04")
	}
	s.corpusJSON(w, out)
}

// handleUploadChunk écrit un morceau à la suite.
//
// L'offset attendu est la taille actuelle du fichier partiel, et il est
// VÉRIFIÉ : deux envois concurrents, ou un morceau rejoué après un timeout,
// écriraient sinon une base silencieusement mélangée. En cas de désaccord on
// renvoie 409 avec l'offset réel, ce qui suffit au client pour se recaler.
func (s *Server) handleUploadChunk(w http.ResponseWriter, r *http.Request) {
	t, ok := s.resolveTarget(w, r)
	if !ok {
		return
	}
	offset, err := strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		s.corpusError(w, "offset invalide", http.StatusBadRequest)
		return
	}

	s.uploadMu.Lock()
	defer s.uploadMu.Unlock()

	var size int64
	if fi, err := os.Stat(t.partial()); err == nil {
		size = fi.Size()
	}
	if offset == 0 && size > 0 {
		// Reprendre à zéro veut dire « on recommence » : on jette le partiel.
		if err := os.Remove(t.partial()); err != nil {
			s.corpusError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		size = 0
	}
	if offset != size {
		privateHeaders(w)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": "offset désynchronisé", "expected": size,
		})
		return
	}

	f, err := os.OpenFile(t.partial(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	n, err := io.Copy(f, http.MaxBytesReader(w, r.Body, chunkLimit))
	cerr := f.Close()
	if err != nil {
		s.corpusError(w, "morceau refusé : "+err.Error(), http.StatusBadRequest)
		return
	}
	if cerr != nil {
		s.corpusError(w, cerr.Error(), http.StatusInternalServerError)
		return
	}
	s.corpusJSON(w, map[string]any{"uploaded": size + n})
}

// handleUploadCommit vérifie le fichier envoyé PUIS remplace la base en place.
// Si la vérification échoue, la base en ligne n'a pas bougé d'un octet.
func (s *Server) handleUploadCommit(w http.ResponseWriter, r *http.Request) {
	t, ok := s.resolveTarget(w, r)
	if !ok {
		return
	}
	s.uploadMu.Lock()
	defer s.uploadMu.Unlock()

	fi, err := os.Stat(t.partial())
	if err != nil {
		s.corpusError(w, "aucun envoi en cours", http.StatusBadRequest)
		return
	}
	if err := verifyDatabase(t, t.partial()); err != nil {
		// Un envoi refusé est connu comme mauvais : le garder le proposerait à
		// la reprise, et on repartirait d'un fichier dont on sait qu'il est
		// cassé. Une coupure réseau, elle, laisse le partiel intact — c'est
		// justement ce qui permet de reprendre.
		_ = os.Remove(t.partial())
		s.corpusError(w, "base refusée : "+err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err := s.swapDatabase(t); err != nil {
		s.corpusError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info("base remplacée", "target", t.name, "bytes", fi.Size())
	s.corpusJSON(w, map[string]any{"ok": true, "target": t.name, "size": fi.Size()})
}

func (s *Server) handleUploadAbort(w http.ResponseWriter, r *http.Request) {
	t, ok := s.resolveTarget(w, r)
	if !ok {
		return
	}
	s.uploadMu.Lock()
	defer s.uploadMu.Unlock()
	_ = os.Remove(t.partial())
	s.corpusJSON(w, map[string]any{"ok": true})
}

// verifyDatabase est le garde-fou qui rend l'opération sûre. Ouvrir la base et
// lui poser une question réelle teste d'un coup qu'elle n'est pas tronquée, que
// le schéma est le bon, et — pour le corpus — que les clés de hash sont dans la
// bonne convention. Un fichier à moitié transféré échoue ici, pas en ligne.
func verifyDatabase(t uploadTarget, path string) error {
	switch t.name {
	case "corpus":
		c, err := corpus.Open(path)
		if err != nil {
			return err
		}
		defer c.Close()
		m, err := c.Meta()
		if err != nil {
			return err
		}
		if m["moves"] == "" {
			return fmt.Errorf("table meta sans nombre de coups")
		}
		return nil
	case "games":
		g, err := games.Open(path)
		if err != nil {
			return err
		}
		defer g.Close()
		if _, err := g.SearchPlayers("a", 1); err != nil {
			return fmt.Errorf("table player illisible (%w)", err)
		}
		return nil
	}
	return fmt.Errorf("base inconnue")
}

// swapDatabase remplace la base par renommages, qui sont atomiques sur un même
// système de fichiers. L'ordre compte : on ferme AVANT de renommer, sinon
// Windows refuse de bouger un fichier ouvert et la production garderait un
// descripteur sur une base qui n'existe plus.
func (s *Server) swapDatabase(t uploadTarget) error {
	s.mu.Lock()
	switch t.name {
	case "corpus":
		if s.corpus != nil {
			_ = s.corpus.Close()
			s.corpus = nil
		}
	case "games":
		if s.games != nil {
			_ = s.games.Close()
			s.games = nil
		}
	}
	s.mu.Unlock()

	_ = os.Remove(t.backup())
	if _, err := os.Stat(t.path); err == nil {
		if err := os.Rename(t.path, t.backup()); err != nil {
			s.reopen(t) // on remet l'ancienne en service plutôt que de rester muet
			return fmt.Errorf("sauvegarde impossible : %w", err)
		}
	}
	if err := os.Rename(t.partial(), t.path); err != nil {
		_ = os.Rename(t.backup(), t.path) // retour en arrière
		s.reopen(t)
		return fmt.Errorf("bascule impossible : %w", err)
	}
	if err := s.reopen(t); err != nil {
		// La nouvelle base a été acceptée à la vérification mais refuse de
		// s'ouvrir : on remet l'ancienne, on ne laisse pas l'outil mort.
		_ = os.Rename(t.path, t.partial())
		_ = os.Rename(t.backup(), t.path)
		_ = s.reopen(t)
		return fmt.Errorf("nouvelle base illisible après bascule : %w", err)
	}
	return nil
}

func (s *Server) reopen(t uploadTarget) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch t.name {
	case "corpus":
		c, err := corpus.Open(t.path)
		if err != nil {
			return err
		}
		s.corpus = c
	case "games":
		g, err := games.Open(t.path)
		if err != nil {
			return err
		}
		s.games = g
	}
	return nil
}

// corpusStore et gamesStore : les handlers passent par là plutôt que de lire le
// champ. Une bascule remplace le pointeur pendant qu'une requête tourne.
func (s *Server) corpusStore() *corpus.Store {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.corpus
}

func (s *Server) gamesStore() *games.Store {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.games
}

// requireStore répond 503 quand la base n'est pas chargée — le cas normal tant
// qu'elle n'a pas été téléversée. La route existe, l'outil n'a pas encore ses
// données, et on le dit.
func (s *Server) requireStore(w http.ResponseWriter, loaded bool, what string) bool {
	if loaded {
		return true
	}
	s.corpusError(w, what+" : aucune base chargée, téléverse-la depuis /admin",
		http.StatusServiceUnavailable)
	return false
}

func humanSize(n int64) string {
	switch {
	case n >= 1<<30:
		return fmt.Sprintf("%.1f Go", float64(n)/(1<<30))
	case n >= 1<<20:
		return fmt.Sprintf("%.0f Mo", float64(n)/(1<<20))
	case n > 0:
		return fmt.Sprintf("%.0f Ko", float64(n)/(1<<10))
	}
	return "—"
}

var _ = strings.TrimSpace // conservé : les helpers de parsing vivent ici aussi
