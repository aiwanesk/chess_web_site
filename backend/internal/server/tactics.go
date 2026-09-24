package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
	"github.com/iwanesko/chess-web-site/backend/internal/stats"
)

// handleTactics serves the most recent weekly puzzle file (JJ-MM-AA.json) from
// TacticsDir. The files are already anonymised (mirrored FEN + solution, no
// identifiers) — safe to serve publicly. Returns an empty set if none exist.
func (s *Server) handleTactics(w http.ResponseWriter, _ *http.Request) {
	dir := s.cfg.TacticsDir
	latest := latestTacticsFile(dir)
	if latest == "" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=300")
		_, _ = w.Write([]byte(`{"week":"","puzzles":[]}`))
		return
	}
	data, err := os.ReadFile(filepath.Join(dir, latest))
	if err != nil {
		http.Error(w, `{"error":"unavailable"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=600")
	_, _ = w.Write(data)
}

// --- interaction events (private stats) ------------------------------------

type eventRequest struct {
	Week     string `json:"week"`
	PuzzleID string `json:"puzzleId"`
	Kind     string `json:"kind"`
}

// handleTacticsEvent records a puzzle interaction. Public, but validated and
// rate-limited (via the /api limiter). No-op if stats are disabled.
func (s *Server) handleTacticsEvent(w http.ResponseWriter, r *http.Request) {
	if s.store == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var req eventRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<10)).Decode(&req); err != nil {
		http.Error(w, "requête invalide", http.StatusBadRequest)
		return
	}
	if !validToken(req.Week) || !validToken(req.PuzzleID) {
		http.Error(w, "requête invalide", http.StatusBadRequest)
		return
	}
	switch req.Kind {
	case stats.View, stats.Attempt, stats.Solved:
	default:
		http.Error(w, "kind invalide", http.StatusBadRequest)
		return
	}
	if err := s.store.Record(req.Week, req.PuzzleID, req.Kind); err != nil {
		slog.Error("stats record failed", "err", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// validToken accepts short slugs/ids: [A-Za-z0-9-], 1..32 chars.
func validToken(v string) bool {
	if v == "" || len(v) > 32 {
		return false
	}
	for _, c := range v {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
			return false
		}
	}
	return true
}

// latestTacticsFile returns the newest weekly puzzle file, or "" if none.
//
// It used to sort the names as text on the assumption, written into its own
// comment, that they were AAAA-Sxx. They are JJ-MM-AA: "31-08-26" sorted after
// "14-09-26", so the endpoint served an August week all September. The date is
// rebuilt from the name before comparing — content.TacticsWeeks does exactly
// that for the sitemap, and one parser for one file-naming scheme means the
// two can no longer disagree about which week is the current one.
func latestTacticsFile(dir string) string {
	weeks := content.TacticsWeeks(dir)
	if len(weeks) == 0 {
		return ""
	}
	return weeks[0].Slug + ".json"
}
