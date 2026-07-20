package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testServer(t *testing.T) http.Handler {
	t.Helper()
	static := fstest.MapFS{
		"index.html":                             {Data: []byte("<!doctype html><h1>Accueil</h1>")},
		"cours-echecs-adultes-geneve/index.html": {Data: []byte("<!doctype html><h1>Cours adultes</h1><script>window.x=1</script>")},
		"404.html":                               {Data: []byte("<!doctype html><h1>404</h1>")},
		"assets/app.abc123.js":                   {Data: []byte("console.log(1)")},
	}
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: "does-not-exist"}, static)
	if err != nil {
		t.Fatal(err)
	}
	return srv.Handler()
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestServesHomeAndRoute(t *testing.T) {
	h := testServer(t)

	if rec := get(t, h, "/"); rec.Code != 200 || !strings.Contains(rec.Body.String(), "Accueil") {
		t.Fatalf("home: code=%d body=%q", rec.Code, rec.Body.String())
	}
	rec := get(t, h, "/cours-echecs-adultes-geneve")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Cours adultes") {
		t.Fatalf("route: code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestUnknownRouteServes404(t *testing.T) {
	rec := get(t, testServer(t), "/nope")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", rec.Code)
	}
}

func TestHTMLNonceMatchesCSP(t *testing.T) {
	rec := get(t, testServer(t), "/cours-echecs-adultes-geneve")
	body := rec.Body.String()
	csp := rec.Header().Get("Content-Security-Policy")

	// The inline script must have gained a nonce...
	if !strings.Contains(body, `<script nonce="`) {
		t.Fatalf("inline script not nonced: %q", body)
	}
	// ...and the CSP must allow exactly that nonce (no 'unsafe-inline').
	i := strings.Index(body, `<script nonce="`) + len(`<script nonce="`)
	nonce := body[i : i+strings.Index(body[i:], `"`)]
	if !strings.Contains(csp, "'nonce-"+nonce+"'") {
		t.Fatalf("CSP %q does not carry the page nonce %q", csp, nonce)
	}
	// script-src must not fall back to 'unsafe-inline'.
	if strings.Contains(csp, "script-src 'self' 'unsafe-inline'") {
		t.Fatalf("script-src weakened with unsafe-inline: %q", csp)
	}
}

func TestImmutableCacheForAssets(t *testing.T) {
	rec := get(t, testServer(t), "/assets/app.abc123.js")
	if got := rec.Header().Get("Cache-Control"); !strings.Contains(got, "immutable") {
		t.Fatalf("asset cache header = %q", got)
	}
}

func TestSitemapContainsMoneyPage(t *testing.T) {
	rec := get(t, testServer(t), "/sitemap.xml")
	body := rec.Body.String()
	if rec.Code != 200 || !strings.Contains(body, "https://iwanesko.ch/cours-echecs-adultes-geneve") {
		t.Fatalf("sitemap missing money page: %s", body)
	}
	if rec.Header().Get("Content-Type") != "application/xml; charset=utf-8" {
		t.Fatalf("sitemap content-type = %q", rec.Header().Get("Content-Type"))
	}
}

func TestRobotsAllowsAICrawlers(t *testing.T) {
	body := get(t, testServer(t), "/robots.txt").Body.String()
	for _, ua := range []string{"GPTBot", "PerplexityBot", "ClaudeBot", "Google-Extended", "CCBot"} {
		if !strings.Contains(body, ua) {
			t.Fatalf("robots.txt missing %s crawler:\n%s", ua, body)
		}
	}
	if !strings.Contains(body, "Sitemap: https://iwanesko.ch/sitemap.xml") {
		t.Fatalf("robots.txt missing sitemap ref:\n%s", body)
	}
}

func TestLLMsTxt(t *testing.T) {
	body := get(t, testServer(t), "/llms.txt").Body.String()
	if !strings.Contains(body, "Maître FIDE") || !strings.Contains(body, "/cours-echecs-adultes-geneve") {
		t.Fatalf("llms.txt missing key facts:\n%s", body)
	}
}

func TestContactValidation(t *testing.T) {
	h := testServer(t)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/contact",
		strings.NewReader(`{"name":"","email":"bad","message":""}`)))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 for invalid contact, got %d", rec.Code)
	}

	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/contact",
		strings.NewReader(`{"name":"Jean","email":"jean@example.com","message":"Bonjour"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200 for valid contact, got %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestContactHoneypot(t *testing.T) {
	rec := httptest.NewRecorder()
	testServer(t).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/contact",
		strings.NewReader(`{"name":"Bot","email":"b@b.com","message":"x","company":"spam"}`)))
	if rec.Code != http.StatusOK {
		t.Fatalf("honeypot should silently succeed, got %d", rec.Code)
	}
}
