package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/online"
)

// prepServer : un serveur avec admin et DB_PATH, dont le client de préparation
// parle à un faux Lichess.
func prepServer(t *testing.T) http.Handler {
	t.Helper()
	fake := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/games/user/") {
			http.NotFound(w, r)
			return
		}
		_, _ = fmt.Fprintln(w, `{"id":"aaa","rated":true,"variant":"standard","speed":"blitz","createdAt":1759000000000,"status":"mate","winner":"black","moves":"e4 c5 Nf3 d6","players":{"white":{"user":{"name":"Other"},"rating":2000},"black":{"user":{"name":"Prepared"},"rating":2100}}}`)
	}))
	t.Cleanup(fake.Close)

	srv, err := New(Config{
		BaseURL: "https://iwanesko.ch", ContentDir: "does-not-exist",
		AdminUser: "admin", AdminToken: "secret",
		DBPath: filepath.Join(t.TempDir(), "site.db"),
	}, fstest.MapFS{"404.html": {Data: []byte("404")}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	srv.prepSets = online.NewSets(&online.Client{HTTP: fake.Client(), LichessURL: fake.URL, ChessComURL: fake.URL})
	return srv.Handler()
}

func adminDo(t *testing.T, h http.Handler, method, path, ctype, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.SetBasicAuth("admin", "secret")
	if ctype != "" {
		req.Header.Set("Content-Type", ctype)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestPrepPageEtAuth(t *testing.T) {
	h := prepServer(t)
	if rec := get(t, h, "/admin/prepa/"); rec.Code != http.StatusUnauthorized {
		t.Fatalf("sans mot de passe : %d", rec.Code)
	}
	rec := adminDo(t, h, "GET", "/admin/prepa/", "", "")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Préparation en ligne") ||
		!strings.Contains(rec.Body.String(), "function renderBoard") {
		t.Fatalf("page : %d", rec.Code)
	}
}

func TestPrepChargementEtArbre(t *testing.T) {
	h := prepServer(t)

	// Un formulaire d'un autre site ne peut pas déclencher de chargement.
	if rec := adminDo(t, h, "POST", "/admin/prepa/api/load", "text/plain",
		`{"accounts":[{"source":"lichess","username":"Prepared"}]}`); rec.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("POST non JSON : %d", rec.Code)
	}
	if rec := adminDo(t, h, "POST", "/admin/prepa/api/load", "application/json",
		`{"accounts":[{"source":"lichess","username":"../admin"}]}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("pseudo invalide : %d", rec.Code)
	}

	rec := adminDo(t, h, "POST", "/admin/prepa/api/load", "application/json",
		`{"accounts":[{"source":"lichess","username":"Prepared"}],"speeds":["blitz"]}`)
	if rec.Code != 200 {
		t.Fatalf("load : %d %s", rec.Code, rec.Body.String())
	}
	var st online.Status
	_ = json.Unmarshal(rec.Body.Bytes(), &st)
	for i := 0; st.State == online.Loading && i < 200; i++ {
		time.Sleep(10 * time.Millisecond)
		_ = json.Unmarshal(adminDo(t, h, "GET", "/admin/prepa/api/status?set="+st.ID, "", "").Body.Bytes(), &st)
	}
	if st.State != online.Ready || st.Games != 1 {
		t.Fatalf("statut : %+v", st)
	}

	// Arbre vu des Noirs : la partie est gagnée par le joueur préparé.
	rec = adminDo(t, h, "GET", "/admin/prepa/api/tree?set="+st.ID+"&colour=b", "", "")
	var tree struct {
		FEN   string `json:"fen"`
		Moves []struct {
			SAN  string `json:"san"`
			Wins int    `json:"wins"`
			FEN  string `json:"fen"`
		} `json:"moves"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &tree)
	if len(tree.Moves) != 1 || tree.Moves[0].SAN != "e4" || tree.Moves[0].Wins != 1 || tree.Moves[0].FEN == "" {
		t.Fatalf("arbre : %s", rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/prepa/api/tree?set="+st.ID+"&colour=w", "", ""); strings.Contains(rec.Body.String(), `"e4"`) {
		t.Fatalf("avec les Blancs il n'a rien joué : %s", rec.Body.String())
	}
	rec = adminDo(t, h, "GET", "/admin/prepa/api/game?set="+st.ID+"&id=0", "", "")
	if !strings.Contains(rec.Body.String(), `"plies"`) || !strings.Contains(rec.Body.String(), "lichess.org/aaa") {
		t.Fatalf("partie : %s", rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/prepa/api/tree?set=inconnu", "", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("lot inconnu : %d", rec.Code)
	}
}

func TestPrepFavoris(t *testing.T) {
	h := prepServer(t)
	rec := adminDo(t, h, "POST", "/admin/prepa/api/favorites", "application/json",
		`{"source":"chesscom","username":"Hikaru"}`)
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Hikaru") {
		t.Fatalf("ajout : %d %s", rec.Code, rec.Body.String())
	}
	rec = adminDo(t, h, "DELETE", "/admin/prepa/api/favorites?source=chesscom&username=hikaru", "", "")
	if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("retrait : %d %s", rec.Code, rec.Body.String())
	}
}
