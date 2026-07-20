package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

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
		"frame-ancestors 'none'"
}
