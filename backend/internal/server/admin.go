package server

import (
	"crypto/subtle"
	"html/template"
	"net/http"
)

// adminAuth guards the dashboard. If ADMIN_TOKEN is unset the route is disabled
// (404 — we don't reveal it exists). Otherwise Basic Auth: any username, the
// password must equal ADMIN_TOKEN, compared in constant time.
func (s *Server) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.AdminToken == "" {
			http.NotFound(w, r)
			return
		}
		_, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(pass), []byte(s.cfg.AdminToken)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="iwanesko-admin", charset="UTF-8"`)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type adminView struct {
	Rows                                   []adminRow
	TotalViews, TotalAttempts, TotalSolved int
}

type adminRow struct {
	Week, PuzzleID          string
	Views, Attempts, Solved int
	SolveRate               int // % solved / views
}

func (s *Server) handleAdmin(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")

	if s.store == nil {
		_, _ = w.Write([]byte("<h1>Stats tactiques</h1><p>Base de données non configurée (définir <code>DB_PATH</code>).</p>"))
		return
	}
	rows, err := s.store.Summary()
	if err != nil {
		http.Error(w, "erreur base de données", http.StatusInternalServerError)
		return
	}

	view := adminView{}
	for _, r := range rows {
		rate := 0
		if r.Views > 0 {
			rate = r.Solved * 100 / r.Views
		}
		view.Rows = append(view.Rows, adminRow{r.Week, r.PuzzleID, r.Views, r.Attempts, r.Solved, rate})
		view.TotalViews += r.Views
		view.TotalAttempts += r.Attempts
		view.TotalSolved += r.Solved
	}
	_ = adminTmpl.Execute(w, view)
}

var adminTmpl = template.Must(template.New("admin").Parse(`<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><meta name="robots" content="noindex">
<title>Stats tactiques — privé</title>
<style>
 body{font:15px/1.5 system-ui,sans-serif;color:#1e293b;max-width:960px;margin:2rem auto;padding:0 1rem}
 h1{font-size:1.4rem} .sub{color:#64748b}
 table{border-collapse:collapse;width:100%;margin-top:1rem;font-variant-numeric:tabular-nums}
 th,td{padding:.5rem .75rem;border-bottom:1px solid #e2e8f0;text-align:right}
 th:first-child,td:first-child,th:nth-child(2),td:nth-child(2){text-align:left}
 thead th{border-bottom:2px solid #cbd5e1;font-size:.8rem;text-transform:uppercase;letter-spacing:.05em;color:#64748b}
 tfoot td{font-weight:700;border-top:2px solid #cbd5e1}
 .empty{color:#94a3b8;margin-top:1rem}
</style></head><body>
<h1>Stats tactiques <span class="sub">— privé</span></h1>
<p class="sub">Vues : {{.TotalViews}} · Tentatives : {{.TotalAttempts}} · Résolus : {{.TotalSolved}}</p>
{{if .Rows}}
<table>
 <thead><tr><th>Semaine</th><th>Puzzle</th><th>Vues</th><th>Tentatives</th><th>Résolus</th><th>Taux</th></tr></thead>
 <tbody>
 {{range .Rows}}<tr><td>{{.Week}}</td><td>{{.PuzzleID}}</td><td>{{.Views}}</td><td>{{.Attempts}}</td><td>{{.Solved}}</td><td>{{.SolveRate}}%</td></tr>{{end}}
 </tbody>
 <tfoot><tr><td colspan="2">Total</td><td>{{.TotalViews}}</td><td>{{.TotalAttempts}}</td><td>{{.TotalSolved}}</td><td></td></tr></tfoot>
</table>
{{else}}<p class="empty">Aucune interaction enregistrée pour l'instant.</p>{{end}}
</body></html>`))
