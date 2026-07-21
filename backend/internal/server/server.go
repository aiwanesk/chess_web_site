package server

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	stdgzip "github.com/CAFxX/httpcompression/contrib/compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
	"github.com/iwanesko/chess-web-site/backend/internal/stats"
)

// Server wires the HTTP handler.
type Server struct {
	cfg    Config
	static fs.FS
	store  *stats.Store      // nil if stats are disabled (no DB_PATH)
	news   *newsletter.Store // nil if the newsletter is disabled (no DB_PATH)
}

// New builds a Server. static is the resolved frontend file source (embedded
// build or on-disk dev directory), provided by the caller.
func New(cfg Config, static fs.FS) (*Server, error) {
	s := &Server{cfg: cfg, static: static}
	if cfg.DBPath != "" {
		st, err := stats.Open(cfg.DBPath)
		if err != nil {
			return nil, fmt.Errorf("open stats db: %w", err)
		}
		s.store = st

		nl, err := newsletter.Open(cfg.DBPath)
		if err != nil {
			return nil, fmt.Errorf("open newsletter db: %w", err)
		}
		s.news = nl
	}
	return s, nil
}

// Close releases resources (the SQLite handles).
func (s *Server) Close() error {
	var err error
	if s.store != nil {
		err = s.store.Close()
	}
	if s.news != nil {
		if e := s.news.Close(); e != nil && err == nil {
			err = e
		}
	}
	return err
}

// Handler returns the fully-configured HTTP handler.
func (s *Server) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(requestLogger)
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders)
	r.Use(s.compression())

	// SEO / GEO endpoints — generated from the shared content package.
	r.Get("/healthz", s.handleHealth)
	r.Get("/robots.txt", s.handleRobots)
	r.Get("/sitemap.xml", s.handleSitemap)
	r.Get("/llms.txt", s.handleLLMs)

	// JSON API for dynamic parts (contact form, tactics, ...).
	// Per-IP rate limit: ~2 req/s sustained, burst 20 — generous for normal use,
	// blocks spam/abuse of the mutating endpoints.
	apiLimiter := newIPRateLimiter(2, 20)
	r.Route("/api", func(api chi.Router) {
		api.Use(rateLimit(apiLimiter))
		api.Get("/health", s.handleHealth)
		api.Post("/contact", s.handleContact)
		api.Get("/tactics", s.handleTactics)
		api.Post("/tactics/event", s.handleTacticsEvent)
		api.Post("/newsletter/subscribe", s.handleSubscribe)
	})

	// Newsletter double opt-in links (from e-mails) — server-rendered pages,
	// no rate limiter needed (idempotent, token-guarded, no enumeration value).
	r.Get("/newsletter/confirm", s.handleNewsletterConfirm)
	r.Get("/newsletter/unsubscribe", s.handleUnsubscribe)
	r.Post("/newsletter/unsubscribe", s.handleUnsubscribe) // RFC 8058 one-click

	// Private stats dashboard — Basic Auth with ADMIN_TOKEN (disabled if unset).
	// Rate-limited BEFORE auth so failed attempts count: ~1 try / 5s per IP,
	// burst 5. Brute-force guard (defence in depth — the real defence is a long
	// random ADMIN_TOKEN). The limiter runs first, so wrong passwords are throttled.
	adminLimiter := newIPRateLimiter(0.2, 5)
	r.With(rateLimit(adminLimiter), s.adminAuth).Get("/admin", s.handleAdmin)

	// Everything else is the pre-rendered SSG site.
	r.NotFound(s.handleStatic)
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	return r
}

// compression negotiates Brotli then gzip based on Accept-Encoding.
func (s *Server) compression() func(http.Handler) http.Handler {
	brEnc, _ := brotli.New(brotli.Options{Quality: 5})
	gzEnc, _ := stdgzip.New(stdgzip.Options{Level: 6})
	compress, _ := httpcompression.Adapter(
		httpcompression.Compressor("br", 1, brEnc),
		httpcompression.Compressor("gzip", 0, gzEnc),
		httpcompression.MinSize(512),
	)
	return compress
}
