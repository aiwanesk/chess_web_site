package server

import (
	"io/fs"
	"net/http"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	stdgzip "github.com/CAFxX/httpcompression/contrib/compress/gzip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server wires the HTTP handler.
type Server struct {
	cfg    Config
	static fs.FS
}

// New builds a Server. static is the resolved frontend file source (embedded
// build or on-disk dev directory), provided by the caller.
func New(cfg Config, static fs.FS) (*Server, error) {
	return &Server{cfg: cfg, static: static}, nil
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

	// JSON API for dynamic parts (contact form, booking, ...).
	r.Route("/api", func(api chi.Router) {
		api.Get("/health", s.handleHealth)
		api.Post("/contact", s.handleContact)
	})

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
