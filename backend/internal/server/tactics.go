package server

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// handleTactics serves the most recent weekly puzzle file (AAAA-Sxx.json) from
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

// latestTacticsFile returns the lexically-greatest AAAA-Sxx.json name (weeks
// sort chronologically), or "" if none.
func latestTacticsFile(dir string) string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var names []string
	for _, e := range entries {
		n := e.Name()
		if !e.IsDir() && strings.HasSuffix(n, ".json") {
			names = append(names, n)
		}
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	return names[len(names)-1]
}
