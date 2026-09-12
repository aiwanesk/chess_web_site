package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/url"
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
		// Ne compter ICI que les robots qui s'annoncent. Les humains sont comptés
		// par la balise JS (/api/hit) : un compteur côté serveur enregistre toute
		// requête HTML, donc tous les crawlers qui se font passer pour un
		// navigateur — ce qui gonflait « visites humaines » d'un facteur inconnu.
		if !looksLikeBot(r.UserAgent()) {
			return
		}
		fp := s.visitorFingerprint(r, ip)
		go func() {
			if err := s.store.RecordPageview(path, country, true, fp, ""); err != nil {
				slog.Error("pageview record failed", "err", err)
			}
		}()
	})
}

// hitRequest est le corps envoyé par la balise JS.
type hitRequest struct {
	Path string `json:"path"`
	// Hôte du référent, lu dans document.referrer par la page. Il ne peut PAS
	// venir de l'en-tête Referer : /api/hit est un appel same-origin, donc cet
	// en-tête désigne toujours le site lui-même — c'est même ce que vérifie
	// sameSite(). L'hôte externe n'existe que côté navigateur.
	Ref string `json:"ref,omitempty"`
}

// refHost valide un hôte annoncé par la page. Le corps de la requête n'est pas
// digne de confiance : on n'accepte qu'un nom d'hôte plausible, jamais une URL,
// et on écarte le site lui-même (une navigation interne n'est pas une
// provenance). Renvoie "" si l'hôte doit être ignoré.
func refHost(raw, self string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" || len(raw) > 100 {
		return ""
	}
	if !hostRe.MatchString(raw) { // exclut d'office les URL, ports et IP nues
		return ""
	}
	raw = strings.TrimPrefix(raw, "www.")
	if raw == bareHost(self) {
		return "" // navigation interne : ce n'est pas une provenance
	}
	return raw
}

// bareHost réduit un Host HTTP à son nom d'hôte comparable : port retiré, en
// minuscules, sans « www. ».
func bareHost(h string) string {
	if host, _, err := net.SplitHostPort(h); err == nil {
		h = host
	}
	return strings.TrimPrefix(strings.ToLower(h), "www.")
}

// Un nom d'hôte et rien d'autre : des labels, un point, un TLD alphabétique.
var hostRe = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)*\.[a-z]{2,}$`)

// sameSite vérifie que la requête vient bien d'une page du site : Origin (ou à
// défaut Referer) doit désigner le même hôte que celui appelé. Ça n'arrête pas
// un attaquant déterminé, mais ça écarte les appels directs, qui sont la façon
// évidente de gonfler des statistiques publiques.
func sameSite(r *http.Request) bool {
	for _, raw := range []string{r.Header.Get("Origin"), r.Header.Get("Referer")} {
		if raw == "" {
			continue
		}
		u, err := url.Parse(raw)
		if err != nil || u.Host == "" {
			continue
		}
		return u.Host == r.Host
	}
	return false
}

// handleHit enregistre UNE vue humaine, déclenchée par la balise JS de la page.
// Les robots n'exécutent quasiment jamais de JavaScript : c'est ce qui rend ce
// compteur-là honnête, là où le comptage serveur ne l'était pas. Mêmes garanties
// qu'avant — aucune IP, aucun cookie, pays déduit hors ligne, empreinte du jour.
func (s *Server) handleHit(w http.ResponseWriter, r *http.Request) {
	if s.store == nil || !sameSite(r) || looksLikeBot(r.UserAgent()) {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	var req hitRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<10)).Decode(&req); err != nil {
		http.Error(w, "requête invalide", http.StatusBadRequest)
		return
	}
	path := cleanHitPath(req.Path)
	if path == "" {
		http.Error(w, "chemin invalide", http.StatusBadRequest)
		return
	}
	ip := clientIP(r)
	country := ""
	if pip := net.ParseIP(ip); pip != nil {
		if c := string(iploc.Country(pip)); c != "" && c != "ZZ" {
			country = c
		}
	}
	fp := s.visitorFingerprint(r, ip)
	if err := s.store.RecordPageview(path, country, false, fp, refHost(req.Ref, r.Host)); err != nil {
		slog.Error("pageview record failed", "err", err)
	}
	w.WriteHeader(http.StatusNoContent)
}

// cleanHitPath n'accepte qu'un chemin interne : le corps de la requête vient du
// navigateur, donc il n'est pas digne de confiance. Renvoie "" si le chemin doit
// être ignoré (page d'admin, URL absolue, chemin invraisemblable).
func cleanHitPath(raw string) string {
	if i := strings.IndexAny(raw, "?#"); i >= 0 {
		raw = raw[:i]
	}
	if raw == "" || len(raw) > 200 || !strings.HasPrefix(raw, "/") ||
		strings.HasPrefix(raw, "//") || strings.Contains(raw, "..") ||
		strings.HasPrefix(raw, "/admin") || strings.HasPrefix(raw, "/newsletter") {
		return ""
	}
	return raw
}
