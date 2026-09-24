package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/iwanesko/chess-web-site/backend/internal/content"
)

// formToken fetches a valid anti-spam token so tests can exercise the forms.
func formToken(t *testing.T, h http.Handler) string {
	t.Helper()
	rec := get(t, h, "/api/form-token")
	if rec.Code != http.StatusOK {
		t.Fatalf("form-token: got %d", rec.Code)
	}
	var b struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &b); err != nil || b.Token == "" {
		t.Fatalf("form-token decode: %v (%q)", err, rec.Body.String())
	}
	return b.Token
}

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
	// Les pages légales et les listes de catégorie descendent en « Optional » :
	// un moteur à court de contexte doit tomber sur les pages utiles en premier.
	opt := strings.Index(body, "## Optional")
	money := strings.Index(body, "/cours-echecs-adultes-geneve")
	if opt < 0 || money < 0 || money > opt {
		t.Fatalf("llms.txt: les pages principales doivent précéder ## Optional\n%s", body)
	}
	if i := strings.Index(body, "/confidentialite"); i >= 0 && i < opt {
		t.Fatalf("llms.txt: la politique de confidentialité doit être sous ## Optional")
	}
}

// Les carnets anglais avaient disparu du llms.txt : il ne lisait que le dossier
// français et préfixait les URLs en dur avec /blog/. Même bug que celui déjà
// corrigé sur le sitemap — ce test empêche la rechute.
func TestLLMsTxtListsBothLocales(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}
	post := "---\ntitle: \"Titre\"\ndescription: \"Résumé\"\ndate: \"2026-01-01\"\n---\n\nCorps.\n"
	if err := os.WriteFile(filepath.Join(dir, "carnet-fr.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "en", "diary-en.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: dir}, fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}
	body := get(t, srv.Handler(), "/llms.txt").Body.String()
	for _, want := range []string{"https://iwanesko.ch/blog/carnet-fr", "https://iwanesko.ch/en/blog/diary-en"} {
		if !strings.Contains(body, want) {
			t.Fatalf("llms.txt sans %s :\n%s", want, body)
		}
	}
}

func TestContactValidation(t *testing.T) {
	h := testServer(t)
	tok := formToken(t, h)

	if rec := postJSON(t, h, "/api/contact",
		`{"name":"","email":"bad","message":"","token":"`+tok+`"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want 422 for invalid contact, got %d", rec.Code)
	}
	if rec := postJSON(t, h, "/api/contact",
		`{"name":"Jean","email":"jean@example.com","message":"Bonjour","token":"`+tok+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("want 200 for valid contact, got %d", rec.Code)
	}
	// No token → rejected as spam.
	if rec := postJSON(t, h, "/api/contact",
		`{"name":"Jean","email":"jean@example.com","message":"Bonjour"}`); rec.Code != http.StatusForbidden {
		t.Fatalf("missing token: want 403, got %d", rec.Code)
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
		AdminUser:  "admin",
		AdminToken: adminToken,
		HourlyRate: 120,
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

	badUser := httptest.NewRequest(http.MethodGet, "/admin", nil)
	badUser.SetBasicAuth("root", "s3cret-token") // bon mot de passe, mauvais identifiant
	recBadUser := httptest.NewRecorder()
	h.ServeHTTP(recBadUser, badUser)
	if recBadUser.Code != http.StatusUnauthorized {
		t.Fatalf("wrong user: want 401, got %d", recBadUser.Code)
	}

	ok := httptest.NewRequest(http.MethodGet, "/admin", nil)
	ok.SetBasicAuth("admin", "s3cret-token")
	recOK := httptest.NewRecorder()
	h.ServeHTTP(recOK, ok)
	if recOK.Code != http.StatusOK {
		t.Fatalf("correct token: want 200, got %d", recOK.Code)
	}
}

func TestServerSurvivesUnusableDB(t *testing.T) {
	// Parent path is a regular file, so SQLite can't create the DB below it.
	f := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	badPath := filepath.Join(f, "nested", "stats.db")

	static := fstest.MapFS{"index.html": {Data: []byte("<!doctype html><h1>Accueil</h1>")}}
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", DBPath: badPath, AdminUser: "admin", AdminToken: "tok"}, static)
	if err != nil {
		t.Fatalf("New must degrade, not fail, on an unusable DB: %v", err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	h := srv.Handler()

	// Site stays up.
	if rec := get(t, h, "/"); rec.Code != http.StatusOK {
		t.Fatalf("site should stay up with a broken DB, got %d", rec.Code)
	}
	// Stats disabled → event endpoint is a silent no-op (not a 500).
	if rec := postJSON(t, h, "/api/tactics/event", `{"week":"20-07-26","puzzleId":"abc","kind":"view"}`); rec.Code != http.StatusNoContent {
		t.Fatalf("event with disabled stats: want 204, got %d", rec.Code)
	}
	// Newsletter disabled → subscribe reports unavailable, doesn't crash.
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"a@b.com","consent":true}`); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("subscribe with disabled newsletter: want 503, got %d", rec.Code)
	}
}

func TestNewsletterSubscribeValidation(t *testing.T) {
	h := statsServer(t, "tok") // DBPath set → newsletter enabled; no SMTP → confirm mail is a no-op
	tok := formToken(t, h)

	// Missing explicit consent → rejected.
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"a@b.com","consent":false,"token":"`+tok+`"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("no consent: want 422, got %d", rec.Code)
	}
	// Invalid e-mail → rejected.
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"nope","consent":true,"token":"`+tok+`"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("bad email: want 422, got %d", rec.Code)
	}
	// Honeypot filled → silently accepted (no leak, no token needed).
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"a@b.com","consent":true,"company":"spam"}`); rec.Code != http.StatusOK {
		t.Fatalf("honeypot: want 200, got %d", rec.Code)
	}
	// Valid → 200 pending.
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"a@b.com","consent":true,"lang":"fr","token":"`+tok+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("valid: want 200, got %d", rec.Code)
	}
	// No token → rejected as spam.
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"c@b.com","consent":true}`); rec.Code != http.StatusForbidden {
		t.Fatalf("missing token: want 403, got %d", rec.Code)
	}
}

func TestNewsletterLinksRenderOnUnknownToken(t *testing.T) {
	h := statsServer(t, "tok")
	if rec := get(t, h, "/newsletter/confirm?token=nope"); rec.Code != http.StatusOK {
		t.Fatalf("confirm unknown token: want 200 page, got %d", rec.Code)
	}
	if rec := get(t, h, "/newsletter/unsubscribe?token=nope"); rec.Code != http.StatusOK {
		t.Fatalf("unsubscribe unknown token: want 200 page, got %d", rec.Code)
	}
}

func TestNewsletterDisabledWithoutDB(t *testing.T) {
	h := testServer(t) // no DBPath → newsletter disabled
	if rec := get(t, h, "/newsletter/confirm?token=x"); rec.Code != http.StatusNotFound {
		t.Fatalf("confirm without DB: want 404, got %d", rec.Code)
	}
	if rec := postJSON(t, h, "/api/newsletter/subscribe", `{"email":"a@b.com","consent":true}`); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("subscribe without DB: want 503, got %d", rec.Code)
	}
}

func TestBookingHandler(t *testing.T) {
	h := statsServer(t, "tok") // DBPath set → bookings enabled; no SMTP → emails are no-ops
	tok := formToken(t, h)
	body := func(fields string) string { return "{" + fields + `,"token":"` + tok + `"}` }

	// Valid 17:30–19:30 = 2h → 240 CHF.
	rec := postJSON(t, h, "/api/booking", body(`"date":"2999-01-01","start":"17:30","end":"19:30","name":"Jean","email":"j@e.com"`))
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "240") {
		t.Fatalf("valid booking: want 200 + price 240, got %d body=%s", rec.Code, rec.Body.String())
	}
	// Overlapping the same slot → 409.
	if rec := postJSON(t, h, "/api/booking", body(`"date":"2999-01-01","start":"18:00","end":"19:00","name":"X","email":"x@e.com"`)); rec.Code != http.StatusConflict {
		t.Fatalf("overlap: want 409, got %d", rec.Code)
	}
	// end <= start → 422.
	if rec := postJSON(t, h, "/api/booking", body(`"date":"2999-01-01","start":"18:00","end":"18:00","name":"X","email":"x@e.com"`)); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("empty slot: want 422, got %d", rec.Code)
	}
	// Outside 17:30–20:00 → 422.
	if rec := postJSON(t, h, "/api/booking", body(`"date":"2999-01-01","start":"17:00","end":"18:00","name":"X","email":"x@e.com"`)); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("out of range: want 422, got %d", rec.Code)
	}
	// Past date → 422.
	if rec := postJSON(t, h, "/api/booking", body(`"date":"2000-01-01","start":"17:30","end":"18:00","name":"X","email":"x@e.com"`)); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("past date: want 422, got %d", rec.Code)
	}
	// No token → rejected as spam.
	if rec := postJSON(t, h, "/api/booking", `{"date":"2999-01-02","start":"17:30","end":"18:00","name":"X","email":"x@e.com"}`); rec.Code != http.StatusForbidden {
		t.Fatalf("missing token: want 403, got %d", rec.Code)
	}
}

func TestBookingAvailability(t *testing.T) {
	h := statsServer(t, "tok")
	tok := formToken(t, h)
	if rec := postJSON(t, h, "/api/booking", `{"date":"2999-07-01","start":"18:00","end":"19:00","name":"X","email":"x@e.com","token":"`+tok+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("seed booking: %d %s", rec.Code, rec.Body.String())
	}
	if rec := get(t, h, "/api/booking/availability?date=2999-07-01"); !strings.Contains(rec.Body.String(), "18:00") || !strings.Contains(rec.Body.String(), "19:00") {
		t.Fatalf("availability must list the booked slot: %s", rec.Body.String())
	}
	if rec := get(t, h, "/api/booking/availability?date=2999-07-02"); !strings.Contains(rec.Body.String(), `"taken":[]`) {
		t.Fatalf("empty day should be free: %s", rec.Body.String())
	}
}

func TestBookingRules(t *testing.T) {
	static := fstest.MapFS{"index.html": {Data: []byte("x")}}
	srv, err := New(Config{
		BaseURL: "https://iwanesko.ch", DBPath: filepath.Join(t.TempDir(), "b.db"),
		HourlyRate: 120, BookingMinDate: "2999-06-01",
	}, static)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	h := srv.Handler()
	tok := formToken(t, h)

	// booking-config advertises the opening date + minimum duration.
	if rec := get(t, h, "/api/booking-config"); !strings.Contains(rec.Body.String(), "2999-06-01") || !strings.Contains(rec.Body.String(), "60") {
		t.Fatalf("booking-config: %s", rec.Body.String())
	}
	// Before the opening date → rejected.
	if rec := postJSON(t, h, "/api/booking", `{"date":"2999-05-01","start":"17:30","end":"18:30","name":"X","email":"x@e.com","token":"`+tok+`"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("before opening date: want 422, got %d", rec.Code)
	}
	// After opening, a full hour → accepted.
	if rec := postJSON(t, h, "/api/booking", `{"date":"2999-06-02","start":"17:30","end":"18:30","name":"X","email":"x@e.com","token":"`+tok+`"}`); rec.Code != http.StatusOK {
		t.Fatalf("1h after opening: want 200, got %d", rec.Code)
	}
	// Under one hour → rejected.
	if rec := postJSON(t, h, "/api/booking", `{"date":"2999-06-03","start":"17:30","end":"18:00","name":"X","email":"x@e.com","token":"`+tok+`"}`); rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("under 1h: want 422, got %d", rec.Code)
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

func TestSecurityTxt(t *testing.T) {
	rec := get(t, testServer(t), "/.well-known/security.txt")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Contact: mailto:") {
		t.Fatalf("security.txt: code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestScannerPathsAreBlocked(t *testing.T) {
	for _, p := range []string{"/wp-login.php", "/xmlrpc.php", "/.env", "/.git/config", "/wp-admin/", "/phpmyadmin/"} {
		if !isScannerPath(p) {
			t.Errorf("expected %q to be treated as a scanner path", p)
		}
	}
	for _, p := range []string{"/", "/tarifs", "/api/contact", "/blog/x", "/.well-known/security.txt"} {
		if isScannerPath(p) {
			t.Errorf("legit path %q wrongly flagged as scanner", p)
		}
	}
	// End-to-end: a scanner probe still 404s (now without logging).
	if rec := get(t, testServer(t), "/wp-login.php"); rec.Code != http.StatusNotFound {
		t.Fatalf("scanner probe: want 404, got %d", rec.Code)
	}
}

// La balise remplace le comptage serveur : c'est elle qui décide désormais ce
// qu'est une « visite humaine ». Elle doit donc être difficile à polluer.
func TestHitBeacon(t *testing.T) {
	h := statsServer(t, "tok")

	hit := func(body, origin, ua string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/hit", strings.NewReader(body))
		req.Host = "iwanesko.ch"
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		if ua != "" {
			req.Header.Set("User-Agent", ua)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}

	const browser = "Mozilla/5.0 (Windows NT 10.0) Chrome/120"
	if code := hit(`{"path":"/blog/x"}`, "https://iwanesko.ch", browser); code != http.StatusNoContent {
		t.Fatalf("vue légitime : want 204, got %d", code)
	}
	// Appel direct, sans page d'origine : c'est la façon évidente de gonfler des
	// statistiques publiques.
	if code := hit(`{"path":"/blog/x"}`, "", browser); code != http.StatusNoContent {
		t.Fatalf("sans Origin : want 204 (ignoré), got %d", code)
	}
	if code := hit(`{"path":"/blog/x"}`, "https://ailleurs.example", browser); code != http.StatusNoContent {
		t.Fatalf("origine étrangère : want 204 (ignoré), got %d", code)
	}
	// Un robot qui s'annonce reste hors des humains, même via la balise.
	if code := hit(`{"path":"/blog/x"}`, "https://iwanesko.ch", "Googlebot/2.1"); code != http.StatusNoContent {
		t.Fatalf("bot déclaré : want 204, got %d", code)
	}
	// Le corps vient du navigateur : il n'est pas digne de confiance.
	for _, bad := range []string{`{"path":"https://evil.example/x"}`, `{"path":"//evil.example"}`, `{"path":"/admin/secret"}`, `{"path":""}`} {
		if code := hit(bad, "https://iwanesko.ch", browser); code != http.StatusBadRequest {
			t.Fatalf("chemin %s : want 400, got %d", bad, code)
		}
	}
}

func TestCleanHitPath(t *testing.T) {
	ok := map[string]string{
		"/":              "/",
		"/blog/x?utm=1":  "/blog/x",
		"/blog/x#partie": "/blog/x",
		"/en/blog/y":     "/en/blog/y",
	}
	for in, want := range ok {
		if got := cleanHitPath(in); got != want {
			t.Fatalf("cleanHitPath(%q) = %q, want %q", in, got, want)
		}
	}
	for _, bad := range []string{"", "blog/x", "//evil", "/../etc", "/admin", "/newsletter/confirm", "https://evil.example"} {
		if got := cleanHitPath(bad); got != "" {
			t.Fatalf("cleanHitPath(%q) = %q, want \"\"", bad, got)
		}
	}
}

// L'annonceur ne lisait que content/blog : les articles anglais n'étaient
// jamais annoncés, et comme les articles français partaient sans langue, un
// abonné anglophone recevait un lien vers un texte français.
func TestAnnouncerCollectsBothLocales(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}
	post := "---\ntitle: \"Titre\"\ndescription: \"Résumé\"\ndate: \"2026-01-01\"\n---\n\nCorps.\n"
	if err := os.WriteFile(filepath.Join(dir, "bilan.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "en", "review.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: dir}, fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}

	byID := map[string]announceItem{}
	for _, it := range srv.collectContent() {
		byID[it.ID] = it
	}

	fr, ok := byID["blog:bilan"]
	if !ok {
		t.Fatal("article FR absent de la collecte")
	}
	if fr.Lang != "fr" {
		t.Errorf("article FR: Lang = %q, want fr — sans langue il part aussi aux anglophones", fr.Lang)
	}
	if fr.URL != "https://iwanesko.ch/blog/bilan" {
		t.Errorf("article FR: URL = %q", fr.URL)
	}

	en, ok := byID["blogEN:review"]
	if !ok {
		t.Fatal("article EN absent de la collecte")
	}
	if en.Lang != "en" {
		t.Errorf("article EN: Lang = %q, want en", en.Lang)
	}
	if en.URL != "https://iwanesko.ch/en/blog/review" {
		t.Errorf("article EN: URL = %q — un abonné anglophone doit recevoir la page anglaise", en.URL)
	}
}

// Élargir la collecte à une nouvelle catégorie ne doit jamais poster l'arriéré
// aux abonnés : le premier passage marque l'existant comme vu, sans envoi.
func TestSeedKindMarksBacklogWithoutSending(t *testing.T) {
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", DBPath: filepath.Join(t.TempDir(), "nl.db")}, fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}
	// Sous Windows, un handle SQLite encore ouvert empêche t.TempDir de
	// nettoyer : le test échouerait sur le ménage, pas sur son sujet.
	t.Cleanup(func() { _ = srv.Close() })
	if srv.news == nil {
		t.Fatal("newsletter store non ouvert")
	}

	items := []announceItem{
		{ID: "blogEN:ancien"}, {ID: "blogEN:autre"}, {ID: "blog:francais"},
	}
	match := func(it announceItem) bool { return strings.HasPrefix(it.ID, "blogEN:") }

	if err := srv.seedKind(enSeedSentinel, items, match); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"blogEN:ancien", "blogEN:autre"} {
		if done, _ := srv.news.IsNotified(id); !done {
			t.Errorf("%s aurait dû être amorcé", id)
		}
	}
	// Un article français ne doit pas être emporté par l'amorçage anglais.
	if done, _ := srv.news.IsNotified("blog:francais"); done {
		t.Error("l'amorçage a débordé sur les articles français")
	}

	// Deuxième passage : un article EN publié après l'amorçage doit rester à envoyer.
	if err := srv.seedKind(enSeedSentinel, append(items, announceItem{ID: "blogEN:nouveau"}), match); err != nil {
		t.Fatal(err)
	}
	if done, _ := srv.news.IsNotified("blogEN:nouveau"); done {
		t.Error("l'amorçage a rejoué et avalé un article publié depuis")
	}
}

// refHost reçoit une valeur choisie par le navigateur : tout ce qui n'est pas
// un nom d'hôte plausible doit tomber, et le site lui-même ne doit jamais être
// compté comme une provenance.
func TestRefHost(t *testing.T) {
	const self = "iwanesko.ch"
	cases := []struct{ in, want string }{
		{"chatgpt.com", "chatgpt.com"},
		{"www.Perplexity.AI", "perplexity.ai"},
		{"  google.ch  ", "google.ch"},
		{"iwanesko.ch", ""},               // navigation interne
		{"www.iwanesko.ch", ""},           // idem, préfixe retiré des deux côtés
		{"", ""},                          // visite directe
		{"https://chatgpt.com/c/abc", ""}, // une URL, pas un hôte
		{"chatgpt.com/c/abc", ""},         // un chemin non plus
		{"localhost", ""},                 // pas de TLD
		{"192.168.1.1", ""},               // IP nue
		{"evil.com:8080", ""},             // port
		{"a..b.com", ""},                  // label vide
		{"nope", ""},
	}
	for _, c := range cases {
		if got := refHost(c.in, self); got != c.want {
			t.Errorf("refHost(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	if got := refHost("chatgpt.com", "iwanesko.ch:8080"); got != "chatgpt.com" {
		t.Errorf("un Host avec port doit rester comparable, got %q", got)
	}
	if got := refHost(strings.Repeat("a", 120)+".com", self); got != "" {
		t.Errorf("hôte trop long accepté : %q", got)
	}
}

// Une page, une URL. /tarifs, /tarifs/ et /tarifs/index.html répondaient tous
// les trois 200 : trois adresses pour un seul document, que seule la balise
// canonical rattachait entre elles.
func TestCanonicalRedirects(t *testing.T) {
	h := testServer(t)
	cases := []struct{ from, to string }{
		{"/cours-echecs-adultes-geneve/", "/cours-echecs-adultes-geneve"},
		{"/cours-echecs-adultes-geneve/index.html", "/cours-echecs-adultes-geneve"},
		{"/index.html", "/"},
		{"//cours-echecs-adultes-geneve", "/cours-echecs-adultes-geneve"},
		// La campagne survit à la redirection, sinon l'attribution est perdue.
		{"/cours-echecs-adultes-geneve/?utm_source=x", "/cours-echecs-adultes-geneve?utm_source=x"},
	}
	for _, c := range cases {
		rec := get(t, h, c.from)
		if rec.Code != http.StatusMovedPermanently {
			t.Fatalf("%s: code = %d, attendu 301", c.from, rec.Code)
		}
		if got := rec.Header().Get("Location"); got != c.to {
			t.Fatalf("%s: Location = %q, attendu %q", c.from, got, c.to)
		}
	}
	// Les URL déjà canoniques ne bougent pas.
	for _, p := range []string{"/", "/cours-echecs-adultes-geneve", "/assets/app.abc123.js"} {
		if rec := get(t, h, p); rec.Code != http.StatusOK {
			t.Fatalf("%s: code = %d, attendu 200", p, rec.Code)
		}
	}
}

// Deux 404 remontées par la Search Console, deux causes différentes : un
// article supprimé dont l'URL reste indexée, et une URL qu'un lien interne
// fautif avait fabriquée (slug anglais sous le préfixe français).
func TestRedirectsForDeadArticleURLs(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}
	post := func(title string) []byte {
		return []byte("---\ntitle: \"" + title + "\"\ndate: \"2026-01-01\"\n---\n\nCorps.\n")
	}
	// Un article propre à chaque langue, et un slug partagé par les deux.
	write := func(rel string, body []byte) {
		if err := os.WriteFile(filepath.Join(dir, rel), body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("reprendre-les-echecs.md", post("Reprendre"))
	write("en/returning-to-chess.md", post("Returning"))
	write("open-badalona-2026.md", post("Badalona"))
	write("en/open-badalona-2026.md", post("Badalona"))

	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: dir}, fstest.MapFS{
		"404.html": {Data: []byte("<!doctype html><h1>404</h1>")},
	})
	if err != nil {
		t.Fatal(err)
	}
	h := srv.Handler()

	cases := []struct{ from, to string }{
		// Article supprimé → l'archive de sa catégorie, dans sa langue.
		{"/blog/sortir-du-plateau-1500-elo", "/blog/categorie/progresser"},
		{"/en/blog/breaking-the-1500-elo-plateau", "/en/blog/category/improve"},
		// Mauvais préfixe de langue, dans les deux sens.
		{"/blog/returning-to-chess", "/en/blog/returning-to-chess"},
		{"/en/blog/reprendre-les-echecs", "/blog/reprendre-les-echecs"},
		// Un seul bond : le slash final ne coûte pas une redirection de plus.
		{"/blog/returning-to-chess/", "/en/blog/returning-to-chess"},
	}
	for _, c := range cases {
		rec := get(t, h, c.from)
		if rec.Code != http.StatusMovedPermanently {
			t.Fatalf("%s: code = %d, attendu 301", c.from, rec.Code)
		}
		if got := rec.Header().Get("Location"); got != c.to {
			t.Fatalf("%s: Location = %q, attendu %q", c.from, got, c.to)
		}
	}

	// Un slug que les DEUX langues publient reste valide des deux côtés : pas
	// de redirection, sinon on casse une URL qui marche.
	for _, p := range []string{"/blog/open-badalona-2026", "/en/blog/open-badalona-2026"} {
		if rec := get(t, h, p); rec.Code == http.StatusMovedPermanently {
			t.Fatalf("%s redirige alors que l'article existe dans les deux langues", p)
		}
	}
	// Et une URL qui n'a jamais existé reste un vrai 404.
	if rec := get(t, h, "/blog/jamais-ecrit"); rec.Code != http.StatusNotFound {
		t.Fatalf("/blog/jamais-ecrit: code = %d, attendu 404", rec.Code)
	}
}

// Le sitemap datait TOUTES les URL du jour de la requête : relu une heure plus
// tard, le fichier entier semblait réécrit. Google cesse de croire un lastmod
// pris en flagrant délit, et les articles — dont la date, elle, est vraie —
// perdent le signal avec le reste.
func TestSitemapLastModIsStable(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "en"), 0o755); err != nil {
		t.Fatal(err)
	}
	post := "---\ntitle: \"Titre\"\ndescription: \"Résumé\"\ndate: \"2026-01-01\"\ncategory: \"carnet-de-tournoi\"\n---\n\nCorps.\n"
	if err := os.WriteFile(filepath.Join(dir, "carnet.md"), []byte(post), 0o644); err != nil {
		t.Fatal(err)
	}
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: dir}, fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}
	body := get(t, srv.Handler(), "/sitemap.xml").Body.String()

	// Une page statique annonce la date qu'elle déclare elle-même, pas l'heure
	// du serveur : c'est ce qui rend le fichier stable d'une lecture à l'autre.
	var tarifs content.Page
	for _, p := range content.StaticPages {
		if p.Path == "/tarifs" {
			tarifs = p
		}
	}
	if !strings.Contains(body, "<loc>https://iwanesko.ch/tarifs</loc>\n    <lastmod>"+tarifs.Updated+"</lastmod>") {
		t.Fatalf("/tarifs ne porte pas sa date déclarée (%s) :\n%s", tarifs.Updated, body)
	}
	// L'article porte sa propre date, et l'archive de catégorie hérite de celle
	// de son article le plus récent : c'est le jour où cette page a changé.
	if !strings.Contains(body, "<loc>https://iwanesko.ch/blog/carnet</loc>\n    <lastmod>2026-01-01</lastmod>") {
		t.Fatalf("l'article ne porte pas sa date :\n%s", body)
	}
	if !strings.Contains(body, "<loc>https://iwanesko.ch/blog/categorie/carnet-de-tournoi</loc>\n    <lastmod>2026-01-01</lastmod>") {
		t.Fatalf("l'archive de catégorie n'hérite pas de la date de son article :\n%s", body)
	}
	// changefreq et priority ne sont lus par personne — ils ne sont plus émis.
	if strings.Contains(body, "changefreq") || strings.Contains(body, "priority") {
		t.Fatalf("changefreq/priority encore présents :\n%s", body)
	}
}

// /api/tactics triait les noms JJ-MM-AA comme du texte : « 31-08-26 » passait
// L'index des tactiques hérite de la date de sa semaine la plus récente :
// c'est le jour où la page a changé, et il tombe tout seul chaque lundi.
// Encore faut-il savoir laquelle est la plus récente : le nom de fichier est
// JJ-MM-AA, donc « 31-08-26 » passe après « 14-09-26 » dans un tri de texte.
func TestSitemapTacticsIndexUsesNewestWeek(t *testing.T) {
	dir := tacticsDir(t, "31-08-26.json", "14-09-26.json")
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: "nope", TacticsDir: dir}, fstest.MapFS{})
	if err != nil {
		t.Fatal(err)
	}
	body := get(t, srv.Handler(), "/sitemap.xml").Body.String()

	want := "<loc>https://iwanesko.ch/tactiques</loc>\n    <lastmod>2026-09-14</lastmod>"
	if !strings.Contains(body, want) {
		t.Fatalf("sitemap : %q absent\n%s", want, body)
	}
}

// tacticsDir écrit des séries hebdomadaires vides dans un dossier temporaire.
func tacticsDir(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, name := range names {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(`{"puzzles":[]}`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}
