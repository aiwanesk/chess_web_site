package server

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// scannerPrefixes are path prefixes hit only by automated vulnerability scanners
// (there is no PHP/WordPress/admin panel here).
var scannerPrefixes = []string{
	"/wp-", "/wordpress", "/phpmyadmin", "/vendor/", "/cgi-bin/",
	"/administrator", "/phpunit", "/.aws", "/actuator",
}

// isScannerPath reports whether a path is obvious bot-scan noise: any *.php URL,
// any hidden dot-path (except the legitimate /.well-known/), or a known probe
// prefix. The real site never uses these.
func isScannerPath(p string) bool {
	if strings.HasSuffix(p, ".php") {
		return true
	}
	if strings.HasPrefix(p, "/.") && !strings.HasPrefix(p, "/.well-known/") {
		return true
	}
	for _, s := range scannerPrefixes {
		if strings.HasPrefix(p, s) {
			return true
		}
	}
	return false
}

// blockScanners quietly 404s bot-scan probes. Registered BEFORE requestLogger so
// this constant background noise (wp-login.php, xmlrpc.php, /.env, …) never spams
// the logs — real traffic and genuine anomalies stay visible.
func blockScanners(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isScannerPath(r.URL.Path) {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ipRateLimiter is a tiny in-memory per-IP token bucket protecting the API from
// abuse/spam. No external dependency; fine for a single-instance service.
type ipRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	rate    float64 // tokens refilled per second
	burst   float64 // max tokens
}

type tokenBucket struct {
	tokens float64
	last   time.Time
}

func newIPRateLimiter(rate, burst float64) *ipRateLimiter {
	return &ipRateLimiter{buckets: make(map[string]*tokenBucket), rate: rate, burst: burst}
}

func (l *ipRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.buckets) > 50000 { // crude unbounded-growth guard
		l.buckets = make(map[string]*tokenBucket)
	}
	now := time.Now()
	b := l.buckets[ip]
	if b == nil {
		b = &tokenBucket{tokens: l.burst, last: now}
		l.buckets[ip] = b
	}
	b.tokens += now.Sub(b.last).Seconds() * l.rate
	if b.tokens > l.burst {
		b.tokens = l.burst
	}
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// rateLimit returns middleware that 429s clients exceeding the limiter.
func rateLimit(l *ipRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !l.allow(clientIP(r)) {
				w.Header().Set("Retry-After", "2")
				http.Error(w, "trop de requêtes", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the client IP. middleware.RealIP has already populated
// RemoteAddr from X-Forwarded-For / X-Real-IP behind the proxy.
func clientIP(r *http.Request) string {
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

// requestLogger emits a structured line per request.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		start := time.Now()
		defer func() {
			slog.Info("http",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.Status(),
				"bytes", ww.BytesWritten(),
				"dur_ms", time.Since(start).Milliseconds(),
			)
		}()
		next.ServeHTTP(ww, r)
	})
}

// securityHeaders sets baseline hardening headers. The CSP is intentionally
// strict: no 'unsafe-inline' for scripts. HTML responses override this header
// with a per-request nonce (see serveHTMLWithNonce) so vite-react-ssg's inline
// hydration scripts run without weakening the policy. 'unsafe-inline' on
// style-src covers the critical inlined CSS.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		h.Set("Content-Security-Policy", cspHeader(""))
		next.ServeHTTP(w, r)
	})
}

// cspHeader builds the Content-Security-Policy. With a nonce, inline scripts
// carrying that nonce are allowed (used for the SSG hydration scripts);
// without, script-src is 'self' only.
func cspHeader(nonce string) string {
	scriptSrc := "script-src 'self'"
	if nonce != "" {
		scriptSrc = "script-src 'self' 'nonce-" + nonce + "'"
	}
	return "default-src 'self'; " +
		"img-src 'self' data:; " +
		"style-src 'self' 'unsafe-inline'; " +
		scriptSrc + "; " +
		"font-src 'self'; " +
		"connect-src 'self'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"frame-src https://www.openstreetmap.org; " + // contact page map
		"frame-ancestors 'none'"
}
