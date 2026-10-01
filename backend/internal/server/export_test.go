package server

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/games"
)

func gamesServer(t *testing.T) *Server {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mega.db")
	w, err := games.OpenWriter(path)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := w.Begin()
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range []games.ImportGame{
		{White: "Iwanesko, Alexandre", Black: `O"Brien, Pat`, Event: "Open de Test",
			Date: "2026.09.15", Result: 1, ECO: "C42", WhiteElo: 2128,
			SAN: []string{"e4", "e5", "Nf3"}, UCI: []string{"e2e4", "e7e5", "g1f3"}},
		{White: "Carlsen, Magnus", Black: "Iwanesko, Alexandre",
			Date: "2026.09.16", Result: -1,
			SAN: []string{"d4", "Nf6"}, UCI: []string{"d2d4", "g8f6"}},
	} {
		if _, err := tx.Insert(g); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()

	srv, err := New(Config{
		BaseURL: "https://iwanesko.ch", ContentDir: "does-not-exist",
		GamesDB: path, AdminUser: "admin", AdminToken: "tok",
	}, fstest.MapFS{"index.html": {Data: []byte("x")}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	return srv
}

// Le test passe par un VRAI serveur HTTP dont la limite d'écriture est déjà
// écoulée quand le gestionnaire commence : si un middleware de la chaîne
// empêchait de la lever, la réponse serait coupée et le zip illisible.
func TestExportPGNSurUnVraiServeur(t *testing.T) {
	ts := httptest.NewUnstartedServer(gamesServer(t).Handler())
	ts.Config.WriteTimeout = time.Nanosecond
	ts.Start()
	t.Cleanup(ts.Close)

	req, _ := http.NewRequest(http.MethodGet, ts.URL+exportPath, nil)
	req.SetBasicAuth("admin", "tok")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatalf("réponse coupée : %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("corps coupé : %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("statut %d : %s", resp.StatusCode, body)
	}
	if ce := resp.Header.Get("Content-Encoding"); ce != "" {
		t.Fatalf("un zip recompressé en %q", ce)
	}

	zr, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil || len(zr.File) != 1 {
		t.Fatalf("archive illisible : %v", err)
	}
	f, _ := zr.File[0].Open()
	pgn, _ := io.ReadAll(f)
	got := string(pgn)

	for _, want := range []string{
		`[Event "Open de Test"]`,
		`[Site "?"]`,
		`[Black "O\"Brien, Pat"]`, // le guillemet échappé, sinon tout le fichier déraille
		`[WhiteElo "2128"]`,
		"1.e4 e5 2.Nf3 1-0",
		`[Event "?"]`, // événement inconnu
		"1.d4 Nf6 0-1",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("absent de l'export : %s\n---\n%s", want, got)
		}
	}
	if strings.Count(got, "[Event ") != 2 {
		t.Fatalf("%d parties exportées, 2 attendues", strings.Count(got, "[Event "))
	}
}

func TestExportExigeLAdmin(t *testing.T) {
	h := gamesServer(t).Handler()
	if rec := get(t, h, exportPath); rec.Code != http.StatusUnauthorized {
		t.Fatalf("export sans authentification : %d", rec.Code)
	}
}

func TestPGNCoupeA80Colonnes(t *testing.T) {
	san := strings.Repeat("Nf3 Nf6 Ng1 Ng8 ", 30)
	var b bytes.Buffer
	if err := games.WritePGN(&b, games.Game{White: "A", Black: "B", SAN: san}); err != nil {
		t.Fatal(err)
	}
	for _, l := range strings.Split(b.String(), "\n") {
		if len(l) > 80 {
			t.Fatalf("ligne de %d colonnes : %q", len(l), l)
		}
	}
}
