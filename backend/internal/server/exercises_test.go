package server

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/iwanesko/chess-web-site/backend/internal/exercises"
)

const exFEN = "r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3"

func exServer(t *testing.T) *Server {
	t.Helper()
	srv, err := New(Config{BaseURL: "https://iwanesko.ch", ContentDir: "does-not-exist",
		AdminUser: "admin", AdminToken: "secret", DBPath: filepath.Join(t.TempDir(), "site.db")},
		fstest.MapFS{"404.html": {Data: []byte("404")}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Close() })
	return srv
}

func TestExercicesBoucleComplete(t *testing.T) {
	h := exServer(t).Handler()
	ex := exercises.Exercise{ID: "e1", GameURL: "https://lichess.org/abc", Source: "lichess",
		FEN: exFEN, PlayedUCI: "a7a6", PlayedSAN: "a6", Phase: "opening", Speed: "blitz",
		Best: 10, PlayedScore: -30, Loss: 40, Tolerance: 15, Depth: 30,
		Moves: []exercises.Candidate{
			{UCI: "g8f6", SAN: "Nf6", CP: 10, Score: 10, OK: true},
			{UCI: "a7a6", SAN: "a6", CP: -30, Score: -30},
		}}
	body, _ := json.Marshal(map[string]any{"game": ex.GameURL, "exercises": []exercises.Exercise{ex}})
	if rec := adminDo(t, h, "POST", "/admin/exercices/api/import", "application/json", string(body)); rec.Code != 200 {
		t.Fatalf("import : %d %s", rec.Code, rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/exercices/api/games", "", ""); !strings.Contains(rec.Body.String(), "lichess.org/abc") {
		t.Fatalf("parties analysées : %s", rec.Body.String())
	}

	// La position servie ne contient pas la solution, mais bien les coups légaux.
	rec := adminDo(t, h, "GET", "/admin/exercices/api/batch", "", "")
	if strings.Contains(rec.Body.String(), "g8f6\",\"san") || strings.Contains(rec.Body.String(), "Nf6") ||
		!strings.Contains(rec.Body.String(), `"legal"`) {
		t.Fatalf("lot : %s", rec.Body.String())
	}

	if rec := adminDo(t, h, "POST", "/admin/exercices/api/answer", "application/json", `{"id":"e1","uci":"e1e8"}`); rec.Code != 400 {
		t.Fatalf("coup illégal : %d", rec.Code)
	}
	rec = adminDo(t, h, "POST", "/admin/exercices/api/answer", "application/json", `{"id":"e1","uci":"g8f6"}`)
	if !strings.Contains(rec.Body.String(), `"verdict":"correct"`) || !strings.Contains(rec.Body.String(), `"san":"Nf6"`) {
		t.Fatalf("réponse : %s", rec.Body.String())
	}
	rec = adminDo(t, h, "POST", "/admin/exercices/api/answer", "application/json", `{"id":"e1","uci":"h7h6"}`)
	if !strings.Contains(rec.Body.String(), `"verdict":"unknown"`) {
		t.Fatalf("coup inconnu : %s", rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/exercices/api/pending", "", ""); !strings.Contains(rec.Body.String(), "h7h6") {
		t.Fatalf("file : %s", rec.Body.String())
	}

	// Théorique : la position ne sort plus dans les manches.
	if rec := adminDo(t, h, "POST", "/admin/exercices/api/tag", "application/json", `{"id":"e1","theory":true}`); rec.Code != 200 {
		t.Fatalf("tag : %d", rec.Code)
	}
	if rec := adminDo(t, h, "GET", "/admin/exercices/api/batch", "", ""); strings.TrimSpace(rec.Body.String()) != "[]" {
		t.Fatalf("une position théorique revient : %s", rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/exercices/api/batch?theory=1", "", ""); !strings.Contains(rec.Body.String(), `"e1"`) {
		t.Fatalf("relecture des théoriques : %s", rec.Body.String())
	}
	if rec := adminDo(t, h, "GET", "/admin/exercices/", "", ""); rec.Code != 200 || !strings.Contains(rec.Body.String(), "Mes imprécisions") {
		t.Fatalf("page : %d", rec.Code)
	}
}
