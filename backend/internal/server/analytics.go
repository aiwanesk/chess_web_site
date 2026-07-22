package server

import (
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/phuslu/iploc"
)

// botUA matches the user-agents of crawlers, scrapers and monitors.
var botUA = regexp.MustCompile(`(?i)bot|crawl|spider|slurp|bing|google|yandex|baidu|duckduck|facebook|embedly|python-requests|curl|wget|headless|semrush|ahrefs|mj12|dotbot|petalbot|monitor|uptime|feed`)

func looksLikeBot(ua string) bool {
	return strings.TrimSpace(ua) == "" || botUA.MatchString(ua)
}

// visitorFingerprint is a non-reversible, daily-rotating hash of IP + user-agent
// (salted with the per-process key). It lets us count UNIQUE visitors without
// storing any IP, cookie or per-visitor identifier — and can't be linked across
// days (the day is mixed into the hash).
func (s *Server) visitorFingerprint(r *http.Request, ip string) string {
	h := sha256.New()
	h.Write(s.formKey)
	h.Write([]byte(time.Now().UTC().Format("2006-01-02")))
	h.Write([]byte(ip))
	h.Write([]byte(r.UserAgent()))
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// countPageviews records a privacy-first page view for successful HTML page
// responses: aggregate counts only — no IP, no cookie, no per-visitor record.
// It also derives the country (offline lookup) and whether the request looks
// like a bot. No-op when stats are disabled; recording runs off the request path.
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
		ip := clientIP(r)
		country := ""
		if pip := net.ParseIP(ip); pip != nil {
			if c := string(iploc.Country(pip)); c != "" && c != "ZZ" {
				country = c
			}
		}
		isBot := looksLikeBot(r.UserAgent())
		fp := s.visitorFingerprint(r, ip)
		go func() {
			if err := s.store.RecordPageview(path, country, isBot, fp); err != nil {
				slog.Error("pageview record failed", "err", err)
			}
		}()
	})
}
