package server

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
)

// countPageviews records a privacy-first page view for successful HTML page
// responses: aggregate counts per (day, path) only — no IP, no user-agent, no
// cookie, no per-visitor record. No-op when stats are disabled. Recording runs
// off the request path so it never slows a response.
func (s *Server) countPageviews(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.store == nil || r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		path := r.URL.Path
		if ww.Status() != http.StatusOK ||
			!strings.HasPrefix(ww.Header().Get("Content-Type"), "text/html") ||
			strings.HasPrefix(path, "/admin") ||
			strings.HasPrefix(path, "/newsletter") {
			return
		}
		go func(p string) {
			if err := s.store.RecordPageview(p); err != nil {
				slog.Error("pageview record failed", "err", err)
			}
		}(path)
	})
}
