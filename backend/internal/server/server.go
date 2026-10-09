package server

import (
	"io/fs"
	"log/slog"
	"net/http"
	"sync"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	stdgzip "github.com/CAFxX/httpcompression/contrib/compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iwanesko/chess-web-site/backend/internal/booking"
	"github.com/iwanesko/chess-web-site/backend/internal/corpus"
	"github.com/iwanesko/chess-web-site/backend/internal/exercises"
	"github.com/iwanesko/chess-web-site/backend/internal/games"
	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
	"github.com/iwanesko/chess-web-site/backend/internal/online"
	"github.com/iwanesko/chess-web-site/backend/internal/stats"
	"github.com/iwanesko/chess-web-site/backend/internal/twic"
)

// Server wires the HTTP handler.
type Server struct {
	cfg    Config
	static fs.FS
	// redirects : 301 des URLs indexées qui n'existent plus, calculées une fois
	// au démarrage (voir redirects.go). Le contenu vit dans l'image.
	redirects map[string]string
	store     *stats.Store      // nil if stats are disabled (no DB_PATH)
	news      *newsletter.Store // nil if the newsletter is disabled (no DB_PATH)
	bookings  *booking.Store    // nil if bookings are disabled (no DB_PATH)
	formKey   []byte            // HMAC key for anti-spam form tokens (per-process)
	corpus    *corpus.Store     // nil tant qu'aucune base n'est chargée
	games     *games.Store      // idem
	twic      *twic.Importer    // mise à jour hebdomadaire de games ; nil si GAMES_DB absent
	// mu protège corpus et games : un téléversement remplace le pointeur
	// pendant que des requêtes tournent. uploadMu sérialise les envois entre
	// eux, deux morceaux concurrents écrivant sinon une base mélangée.
	mu       sync.RWMutex
	uploadMu sync.Mutex
	exportMu sync.Mutex // un seul export PGN à la fois
	// Préparation en ligne : les parties chargées vivent en mémoire (prepSets),
	// seuls les pseudos favoris sont écrits (favorites, nil sans DB_PATH).
	prepSets  *online.Sets
	favorites *online.Favorites
	// Exercices tirés des parties d'Alexandre (nil sans DB_PATH).
	exercises *exercises.Store
}

// New builds a Server. static is the resolved frontend file source (embedded
// build or on-disk dev directory), provided by the caller.
func New(cfg Config, static fs.FS) (*Server, error) {
	s := &Server{cfg: cfg, static: static, formKey: newFormKey(), prepSets: online.NewSets(online.NewClient())}
	s.redirects = buildRedirects(cfg.ContentDir)
	if cfg.GamesDB != "" {
		if g, err := games.Open(cfg.GamesDB); err != nil {
			slog.Error("base de parties indisponible — explorateur désactivé", "path", cfg.GamesDB, "err", err)
		} else {
			s.games = g
		}
	}
	if cfg.CorpusDB != "" {
		// Même principe que pour les autres bases : une base absente ou illisible
		// désactive l'explorateur, elle ne fait pas tomber le site.
		if c, err := corpus.Open(cfg.CorpusDB); err != nil {
			slog.Error("corpus indisponible — explorateur désactivé", "path", cfg.CorpusDB, "err", err)
		} else {
			s.corpus = c
		}
	}
	s.twic = s.newImporter()
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
			if ex, err := exercises.Open(cfg.DBPath); err != nil {
				slog.Error("exercices indisponibles", "path", cfg.DBPath, "err", err)
			} else {
				s.exercises = ex
			}
			if fav, err := online.OpenFavorites(cfg.DBPath); err != nil {
				slog.Error("favoris de préparation indisponibles", "path", cfg.DBPath, "err", err)
			} else {
				s.favorites = fav
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
	if s.favorites != nil {
		if e := s.favorites.Close(); e != nil && err == nil {
			err = e
		}
	}
	if s.exercises != nil {
		if e := s.exercises.Close(); e != nil && err == nil {
			err = e
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.games != nil {
		if e := s.games.Close(); e != nil && err == nil {
			err = e
		}
		s.games = nil
	}
	if s.corpus != nil {
		if e := s.corpus.Close(); e != nil && err == nil {
			err = e
		}
		s.corpus = nil
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
		api.Post("/hit", s.handleHit) // balise de fréquentation (voir analytics.go)
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

	// L'explorateur interroge le serveur À CHAQUE COUP JOUÉ : le limiteur du
	// tableau de bord (un appel toutes les 5 s) le rendrait inutilisable. Celui-ci
	// est dimensionné pour une navigation au doigt, la force brute restant
	// couverte par le limiteur d'authentification ci-dessus.
	if s.cfg.CorpusDB != "" {
		corpusLimiter := newIPRateLimiter(15, 90)
		r.Route("/admin/corpus", func(c chi.Router) {
			c.Use(rateLimit(corpusLimiter), s.adminAuth)
			c.Get("/", s.handleCorpusPage)
			c.Get("/api/pos", s.handleCorpusPos)
			c.Get("/api/go", s.handleCorpusGo)
			c.Get("/api/meta", s.handleCorpusMeta)
		})
	}
	if s.cfg.GamesDB != "" {
		partiesLimiter := newIPRateLimiter(15, 90)
		r.Route("/admin/parties", func(c chi.Router) {
			c.Use(rateLimit(partiesLimiter), s.adminAuth)
			c.Get("/", s.handlePartiesPage)
			c.Get("/api/players", s.handlePartiesPlayers)
			c.Get("/api/tree", s.handlePartiesTree)
			c.Get("/api/games", s.handlePartiesGames)
			c.Get("/api/game", s.handlePartiesGame)
			c.Get("/api/meta", s.handlePartiesMeta)
			c.Get("/api/twic", s.handleTWICStatus)
			c.Post("/twic", s.handleTWICRun)
			c.Get("/export", s.handlePartiesExport)
		})
	}
	// Préparation en ligne : ne dépend d'aucune base, elle va chercher les
	// parties sur Lichess et Chess.com. Même limiteur que l'explorateur, la
	// navigation y est aussi coup par coup.
	prepLimiter := newIPRateLimiter(15, 90)
	r.Route("/admin/prepa", func(c chi.Router) {
		c.Use(rateLimit(prepLimiter), s.adminAuth)
		c.Get("/", s.handlePrepPage)
		c.Post("/api/load", s.handlePrepLoad)
		c.Get("/api/status", s.handlePrepStatus)
		c.Get("/api/tree", s.handlePrepTree)
		c.Get("/api/games", s.handlePrepGames)
		c.Get("/api/game", s.handlePrepGame)
		c.Get("/api/favorites", s.handlePrepFavorites)
		c.Post("/api/favorites", s.handlePrepFavorites)
		c.Delete("/api/favorites", s.handlePrepFavorites)
	})
	// Exercices : l'analyseur du PC pousse ses parties d'un coup (une requête
	// par partie), la page tire une position à la fois. Même limiteur que la
	// préparation.
	exLimiter := newIPRateLimiter(15, 90)
	r.Route("/admin/exercices", func(c chi.Router) {
		c.Use(rateLimit(exLimiter), s.adminAuth)
		c.Get("/", s.handleExercisesPage)
		c.Get("/api/games", s.handleExGames)
		c.Post("/api/import", s.handleExImport)
		c.Get("/api/pending", s.handleExPending)
		c.Post("/api/resolve", s.handleExResolve)
		c.Get("/api/batch", s.handleExBatch)
		c.Post("/api/answer", s.handleExAnswer)
		c.Get("/api/reveal", s.handleExReveal)
		c.Post("/api/tag", s.handleExTag)
		c.Post("/api/run", s.handleExRun)
		c.Get("/api/stats", s.handleExStats)
	})
	if s.cfg.CorpusDB != "" || s.cfg.GamesDB != "" {
		uploadLimiter := newIPRateLimiter(30, 120) // un morceau toutes les 2 s en régime
		r.Route("/admin/upload", func(c chi.Router) {
			c.Use(rateLimit(uploadLimiter), s.adminAuth)
			c.Get("/status", s.handleUploadStatus)
			c.Post("/chunk", s.handleUploadChunk)
			c.Post("/commit", s.handleUploadCommit)
			c.Post("/abort", s.handleUploadAbort)
		})
	}

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
	return func(next http.Handler) http.Handler {
		compressed := compress(next)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Les réponses longues passent à côté : voir noLimit.
			if r.URL.Path == exportPath || r.URL.Path == twicRunPath {
				next.ServeHTTP(w, r)
				return
			}
			compressed.ServeHTTP(w, r)
		})
	}
}
