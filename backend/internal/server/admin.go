package server

import (
	"crypto/subtle"
	"html/template"
	"net/http"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/stats"
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
	TopPages                               []stats.PageRow
	DailyViews                             []stats.DayRow
	TotalPageviews                         int
	Bookings                               []bookingRow
}

type bookingRow struct {
	Date, Time, Name, Email string
	Price                   int
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

	view.TopPages, _ = s.store.TopPages(20)
	view.DailyViews, _ = s.store.DailyViews(14)
	for _, d := range view.DailyViews {
		view.TotalPageviews += d.Count
	}

	if s.bookings != nil {
		bs, _ := s.bookings.Upcoming(time.Now().Format("2006-01-02"), 40)
		for _, b := range bs {
			view.Bookings = append(view.Bookings, bookingRow{
				Date:  b.Date,
				Time:  minToHHMM(b.StartMin) + "–" + minToHHMM(b.EndMin),
				Name:  b.Name,
				Email: b.Email,
				Price: b.Price,
			})
		}
	}
	_ = adminTmpl.Execute(w, view)
}

var adminTmpl = template.Must(template.New("admin").Parse(`<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><meta name="robots" content="noindex">
<title>Tableau de bord — privé</title>
<style>
 body{font:15px/1.5 system-ui,sans-serif;color:#1e293b;max-width:960px;margin:2rem auto;padding:0 1rem}
 h1{font-size:1.4rem} h2{font-size:1.05rem;margin-top:2.25rem} .sub{color:#64748b}
 .grid{display:grid;gap:1.5rem;grid-template-columns:1fr 1fr}
 @media(max-width:640px){.grid{grid-template-columns:1fr}}
 table{border-collapse:collapse;width:100%;margin-top:.6rem;font-variant-numeric:tabular-nums}
 th,td{padding:.45rem .7rem;border-bottom:1px solid #e2e8f0;text-align:right}
 th:first-child,td:first-child,th:nth-child(2),td:nth-child(2){text-align:left}
 thead th{border-bottom:2px solid #cbd5e1;font-size:.78rem;text-transform:uppercase;letter-spacing:.05em;color:#64748b}
 tfoot td{font-weight:700;border-top:2px solid #cbd5e1}
 .empty{color:#94a3b8;margin-top:.6rem}
 td.p{max-width:16rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
</style></head><body>
<h1>Tableau de bord <span class="sub">— privé</span></h1>

<h2>Fréquentation <span class="sub">— {{.TotalPageviews}} vues sur 14 jours (sans cookie, sans tiers)</span></h2>
<div class="grid">
 <div>
  {{if .DailyViews}}
  <table>
   <thead><tr><th>Jour</th><th>Vues</th></tr></thead>
   <tbody>{{range .DailyViews}}<tr><td>{{.Day}}</td><td>{{.Count}}</td></tr>{{end}}</tbody>
  </table>
  {{else}}<p class="empty">Aucune visite enregistrée.</p>{{end}}
 </div>
 <div>
  {{if .TopPages}}
  <table>
   <thead><tr><th>Page (top)</th><th>Vues</th></tr></thead>
   <tbody>{{range .TopPages}}<tr><td class="p">{{.Path}}</td><td>{{.Count}}</td></tr>{{end}}</tbody>
  </table>
  {{end}}
 </div>
</div>

<h2>Réservations à venir <span class="sub">— {{len .Bookings}}</span></h2>
{{if .Bookings}}
<table>
 <thead><tr><th>Date</th><th>Horaire</th><th>Élève</th><th>E-mail</th><th>Montant</th></tr></thead>
 <tbody>{{range .Bookings}}<tr><td>{{.Date}}</td><td>{{.Time}}</td><td>{{.Name}}</td><td class="p">{{.Email}}</td><td>{{.Price}} CHF</td></tr>{{end}}</tbody>
</table>
{{else}}<p class="empty">Aucune réservation à venir.</p>{{end}}

<h2>Tactiques <span class="sub">— Vues {{.TotalViews}} · Tentatives {{.TotalAttempts}} · Résolus {{.TotalSolved}}</span></h2>
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
