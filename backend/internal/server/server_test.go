package server

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
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

// statsServer builds a server backed by a real (temp-file) SQLite store so the
// event → admin round-trip is exercised end to end.
func statsServer(t *testing.T, adminToken string) http.Handler {
	t.Helper()
	static := fstest.MapFS{"index.html": {Data: []byte("x")}}
	srv, err := New(Config{
		BaseURL:    "https://iwanesko.ch",
		ContentDir: "does-not-exist",
		DBPath:     filepath.Join(t.TempDir(), "stats.db"),
		AdminToken: adminToken,
	}, static)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	return srv.Handler()
}

func postJSON(t *testing.T, h http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, strings.NewReader(body)))
	return rec
}

func TestAdminDisabledWithoutToken(t *testing.T) {
	// No ADMIN_TOKEN → route must not exist at all (404, not 401).
	if rec := get(t, statsServer(t, ""), "/admin"); rec.Code != http.StatusNotFound {
		t.Fatalf("admin without token: want 404, got %d", rec.Code)
	}
}

func TestAdminRequiresValidToken(t *testing.T) {
	h := statsServer(t, "s3cret-token")

	if rec := get(t, h, "/admin"); rec.Code != http.StatusUnauthorized {
		t.Fatalf("no auth: want 401, got %d", rec.Code)
	}

	wrong := httptest.NewRequest(http.MethodGet, "/admin", nil)
	wrong.SetBasicAuth("admin", "nope")
	recWrong := httptest.NewRecorder()
	h.ServeHTTP(recWrong, wrong)
	if recWrong.Code != http.StatusUnauthorized {
		t.Fatalf("wrong token: want 401, got %d", recWrong.Code)
	}

	ok := httptest.NewRequest(http.MethodGet, "/admin", nil)
	ok.SetBasicAuth("admin", "s3cret-token")
	recOK := httptest.NewRecorder()
	h.ServeHTTP(recOK, ok)
	if recOK.Code != http.StatusOK {
		t.Fatalf("correct token: want 200, got %d", recOK.Code)
	}
}

func TestAdminBruteForceIsRateLimited(t *testing.T) {
	h := statsServer(t, "s3cret-token")

	// Hammer with wrong credentials from the same IP. The limiter (burst 5) must
	// start returning 429 before we've made many guesses.
	throttled := false
	for i := 0; i < 30; i++ {
		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		req.RemoteAddr = "203.0.113.7:1234"
		req.SetBasicAuth("admin", "guess")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			throttled = true
			break
		}
	}
	if !throttled {
		t.Fatal("brute-force on /admin was never rate-limited (expected 429)")
	}
}

func TestTacticsEventRecordsAndSurfacesInAdmin(t *testing.T) {
	h := statsServer(t, "tok")

	for _, k := range []string{"view", "attempt", "solved"} {
		rec := postJSON(t, h, "/api/tactics/event",
			`{"week":"20-07-26","puzzleId":"abc123","kind":"`+k+`"}`)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("event %q: want 204, got %d", k, rec.Code)
		}
	}

	// Bad kind and bad token are rejected.
	if rec := postJSON(t, h, "/api/tactics/event",
		`{"week":"20-07-26","puzzleId":"abc123","kind":"hack"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad kind: want 400, got %d", rec.Code)
	}
	if rec := postJSON(t, h, "/api/tactics/event",
		`{"week":"../etc","puzzleId":"abc123","kind":"view"}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad week token: want 400, got %d", rec.Code)
	}

	// The recorded interaction shows up in the dashboard.
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.SetBasicAuth("admin", "tok")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	body := rec.Body.String()
	if rec.Code != http.StatusOK || !strings.Contains(body, "20-07-26") || !strings.Contains(body, "abc123") {
		t.Fatalf("admin dashboard missing recorded event: code=%d body=%q", rec.Code, body)
	}
}
