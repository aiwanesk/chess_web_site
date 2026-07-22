package server

import (
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	stdgzip "github.com/CAFxX/httpcompression/contrib/compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iwanesko/chess-web-site/backend/internal/booking"
	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
	"github.com/iwanesko/chess-web-site/backend/internal/stats"
)

// Server wires the HTTP handler.
type Server struct {
	cfg      Config
	static   fs.FS
	store    *stats.Store      // nil if stats are disabled (no DB_PATH)
	news     *newsletter.Store // nil if the newsletter is disabled (no DB_PATH)
	bookings *booking.Store    // nil if bookings are disabled (no DB_PATH)
	formKey  []byte            // HMAC key for anti-spam form tokens (per-process)
}

// New builds a Server. static is the resolved frontend file source (embedded
// build or on-disk dev directory), provided by the caller.
func New(cfg Config, static fs.FS) (*Server, error) {
	s := &Server{cfg: cfg, static: static, formKey: newFormKey()}
	if cfg.DBPath != "" {
		// A DB failure (e.g. an unwritable /data volume) must NOT take the site
		// down: log loudly and degrade — stats + newsletter simply stay off.
		// They come back automatically once the path is fixed and the app boots.
		if st, err := stats.Open(cfg.DBPath); err != nil {
			slog.Error("stats DB unavailable — stats & newsletter disabled, site stays up", "path", cfg.DBPath, "err", err)
		} else {
			s.store = st
			if nl, err := newsletter.Open(cfg.DBPath); err != nil {
				slog.Error("newsletter DB unavailable — newsletter disabled", "path", cfg.DBPath, "err", err)
			} else {
				s.news = nl
			}
			if bk, err := booking.Open(cfg.DBPath); err != nil {
				slog.Error("booking DB unavailable — booking disabled", "path", cfg.DBPath, "err", err)
			} else {
				s.bookings = bk
			}
		}
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
	if s.bookings != nil {
		if e := s.bookings.Close(); e != nil && err == nil {
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
	r.Use(blockScanners) // quietly drop bot-scan noise before it hits the logger
	r.Use(requestLogger)
	r.Use(s.countPageviews) // privacy-first, aggregate page-view analytics
	r.Use(middleware.Recoverer)
	r.Use(securityHeaders)
	r.Use(s.compression())

	// SEO / GEO endpoints — generated from the shared content package.
	r.Get("/healthz", s.handleHealth)
	r.Get("/robots.txt", s.handleRobots)
	r.Get("/sitemap.xml", s.handleSitemap)
	r.Get("/llms.txt", s.handleLLMs)
	r.Get("/.well-known/security.txt", s.handleSecurityTxt)

	// JSON API for dynamic parts (contact form, tactics, ...).
	// Per-IP rate limit: ~2 req/s sustained, burst 20 — generous for normal use,
	// blocks spam/abuse of the mutating endpoints.
	apiLimiter := newIPRateLimiter(2, 20)
	// Stricter per-IP budget for form submissions (shared across the 3 forms):
	// ~6 up front, then 1 every 20 s. Plenty for a human, painful for a spammer.
	submitLimiter := newIPRateLimiter(0.05, 6)
	r.Route("/api", func(api chi.Router) {
		api.Use(rateLimit(apiLimiter))
		api.Get("/health", s.handleHealth)
		api.Get("/form-token", s.handleFormToken)
		api.Get("/tactics", s.handleTactics)
		api.Get("/booking-config", s.handleBookingConfig)
		api.Get("/booking/availability", s.handleAvailability)
		api.Post("/tactics/event", s.handleTacticsEvent)
		api.With(rateLimit(submitLimiter)).Post("/contact", s.handleContact)
		api.With(rateLimit(submitLimiter)).Post("/newsletter/subscribe", s.handleSubscribe)
		api.With(rateLimit(submitLimiter)).Post("/booking", s.handleBooking)
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
