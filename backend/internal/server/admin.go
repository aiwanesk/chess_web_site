package server

import (
	"crypto/subtle"
	"html/template"
	"net/http"
	"strings"
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
	// Tactiques
	Rows                                   []adminRow
	TotalViews, TotalAttempts, TotalSolved int
	// Fréquentation
	TopPages                           []stats.PageRow
	Traffic                            []trafficRow
	Countries                          []countryRow
	TotalHuman, TotalBot, TotalUniques int
	BotPct                             int
	// Réservations
	Bookings []bookingRow
}

type adminRow struct {
	Week, PuzzleID          string
	Views, Attempts, Solved int
	SolveRate               int
}

type trafficRow struct {
	Day                           string
	Human, Bot, Uniques, HumanPct int
}

type countryRow struct {
	Flag, Code string
	Count      int
}

type bookingRow struct {
	Date, Time, Name, Email string
	Price                   int
}

func (s *Server) handleAdmin(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Robots-Tag", "noindex, nofollow")

	if s.store == nil {
		_, _ = w.Write([]byte("<h1>Tableau de bord</h1><p>Base de données non configurée (définir <code>DB_PATH</code>).</p>"))
		return
	}

	view := adminView{}

	// Tactiques
	if rows, err := s.store.Summary(); err == nil {
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
	}

	// Fréquentation
	view.TopPages, _ = s.store.TopPages(15)
	if traffic, err := s.store.Traffic(14); err == nil {
		for _, t := range traffic {
			pct := 100
			if total := t.Human + t.Bot; total > 0 {
				pct = t.Human * 100 / total
			}
			view.Traffic = append(view.Traffic, trafficRow{t.Day, t.Human, t.Bot, t.Uniques, pct})
			view.TotalHuman += t.Human
			view.TotalBot += t.Bot
			view.TotalUniques += t.Uniques
		}
	}
	if total := view.TotalHuman + view.TotalBot; total > 0 {
		view.BotPct = view.TotalBot * 100 / total
	}
	if countries, err := s.store.TopCountries(12); err == nil {
		for _, c := range countries {
			view.Countries = append(view.Countries, countryRow{flagEmoji(c.Country), c.Country, c.Count})
		}
	}

	// Réservations
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

// flagEmoji turns a 2-letter country code into its flag emoji.
func flagEmoji(code string) string {
	if len(code) != 2 {
		return "🏳️"
	}
	code = strings.ToUpper(code)
	return string([]rune{rune(0x1F1E6 + int(code[0]-'A')), rune(0x1F1E6 + int(code[1]-'A'))})
}

var adminTmpl = template.Must(template.New("admin").Parse(`<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="robots" content="noindex"><title>Tableau de bord — privé</title>
<style>
 :root{--ink:#1e293b;--sub:#64748b;--line:#e2e8f0;--card:#fff;--bg:#f8fafc;--accent:#1e293b}
 *{box-sizing:border-box}
 body{font:15px/1.5 system-ui,sans-serif;color:var(--ink);background:var(--bg);margin:0;padding:1.5rem;max-width:1000px;margin:0 auto}
 h1{font-size:1.4rem;margin:.2rem 0 1rem} h2{font-size:1rem;margin:1.75rem 0 .5rem} .sub{color:var(--sub);font-weight:400}
 .tabin{position:absolute;width:0;height:0;opacity:0}
 .tabs{display:flex;gap:.4rem;flex-wrap:wrap;border-bottom:2px solid var(--line);margin-bottom:1.25rem}
 .tabs label{cursor:pointer;padding:.55rem 1rem;border-radius:.5rem .5rem 0 0;font-weight:600;color:var(--sub);user-select:none}
 .tabs label:hover{color:var(--ink)}
 .panel{display:none}
 #t1:checked~#p1,#t2:checked~#p2,#t3:checked~#p3{display:block}
 #t1:checked~.tabs label[for=t1],#t2:checked~.tabs label[for=t2],#t3:checked~.tabs label[for=t3]{color:var(--ink);box-shadow:inset 0 -2px 0 var(--accent)}
 .cards{display:grid;gap:.75rem;grid-template-columns:repeat(3,1fr);margin:.5rem 0 1rem}
 @media(max-width:620px){.cards{grid-template-columns:1fr}}
 .kpi{background:var(--card);border:1px solid var(--line);border-radius:.75rem;padding:.9rem 1rem}
 .kpi .n{font-size:1.6rem;font-weight:800;font-variant-numeric:tabular-nums} .kpi .l{color:var(--sub);font-size:.8rem}
 .grid2{display:grid;gap:1.25rem;grid-template-columns:1fr 1fr}
 @media(max-width:720px){.grid2{grid-template-columns:1fr}}
 table{border-collapse:collapse;width:100%;margin-top:.4rem;font-variant-numeric:tabular-nums;background:var(--card);border:1px solid var(--line);border-radius:.6rem;overflow:hidden}
 th,td{padding:.5rem .7rem;border-bottom:1px solid var(--line);text-align:right}
 th:first-child,td:first-child,th.l,td.l{text-align:left}
 thead th{font-size:.75rem;text-transform:uppercase;letter-spacing:.04em;color:var(--sub);background:#f1f5f9}
 tbody tr:last-child td{border-bottom:0}
 .p{max-width:15rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
 .bar{height:8px;border-radius:4px;background:#e2e8f0;overflow:hidden;min-width:70px}
 .bar>span{display:block;height:100%;background:#16a34a}
 .empty{color:#94a3b8;margin:.5rem 0}
 .flag{font-size:1.1rem}
</style></head><body>
<h1>Tableau de bord <span class="sub">— privé</span></h1>

<input type="radio" name="tab" id="t1" class="tabin" checked>
<input type="radio" name="tab" id="t2" class="tabin">
<input type="radio" name="tab" id="t3" class="tabin">
<nav class="tabs">
 <label for="t1">Fréquentation</label>
 <label for="t2">Réservations{{if .Bookings}} ({{len .Bookings}}){{end}}</label>
 <label for="t3">Tactiques</label>
</nav>

<section class="panel" id="p1">
 <div class="cards">
  <div class="kpi"><div class="n">{{.TotalUniques}}</div><div class="l">Visiteurs uniques (14 j, cumul/jour)</div></div>
  <div class="kpi"><div class="n">{{.TotalHuman}}</div><div class="l">Visites humaines (14 j)</div></div>
  <div class="kpi"><div class="n">{{.TotalBot}}</div><div class="l">Bots ({{.BotPct}} % du trafic)</div></div>
 </div>
 <p class="sub" style="font-size:.8rem">Sans cookie, sans tiers. Pays et « unique » via lookup hors-ligne + empreinte anonyme (aucune IP stockée).</p>

 <div class="grid2">
  <div>
   <h2>Par jour</h2>
   {{if .Traffic}}
   <table><thead><tr><th class="l">Jour</th><th>Uniques</th><th>Visites</th><th>Bots</th><th class="l">Humain</th></tr></thead>
   <tbody>{{range .Traffic}}<tr><td class="l">{{.Day}}</td><td>{{.Uniques}}</td><td>{{.Human}}</td><td>{{.Bot}}</td>
     <td class="l"><span class="bar"><span style="width:{{.HumanPct}}%"></span></span></td></tr>{{end}}</tbody></table>
   {{else}}<p class="empty">Aucune visite enregistrée.</p>{{end}}
  </div>
  <div>
   <h2>Pays</h2>
   {{if .Countries}}
   <table><thead><tr><th class="l">Pays</th><th>Visites</th></tr></thead>
   <tbody>{{range .Countries}}<tr><td class="l"><span class="flag">{{.Flag}}</span> {{.Code}}</td><td>{{.Count}}</td></tr>{{end}}</tbody></table>
   {{else}}<p class="empty">—</p>{{end}}
   <h2>Pages populaires</h2>
   {{if .TopPages}}
   <table><thead><tr><th class="l">Page</th><th>Vues</th></tr></thead>
   <tbody>{{range .TopPages}}<tr><td class="l p">{{.Path}}</td><td>{{.Count}}</td></tr>{{end}}</tbody></table>
   {{else}}<p class="empty">—</p>{{end}}
  </div>
 </div>
</section>

<section class="panel" id="p2">
 <h2>Réservations à venir <span class="sub">— {{len .Bookings}}</span></h2>
 {{if .Bookings}}
 <table><thead><tr><th class="l">Date</th><th class="l">Horaire</th><th class="l">Élève</th><th class="l">E-mail</th><th>Montant</th></tr></thead>
 <tbody>{{range .Bookings}}<tr><td class="l">{{.Date}}</td><td class="l">{{.Time}}</td><td class="l">{{.Name}}</td><td class="l p">{{.Email}}</td><td>{{.Price}} CHF</td></tr>{{end}}</tbody></table>
 {{else}}<p class="empty">Aucune réservation à venir.</p>{{end}}
</section>

<section class="panel" id="p3">
 <h2>Tactiques <span class="sub">— Vues {{.TotalViews}} · Tentatives {{.TotalAttempts}} · Résolus {{.TotalSolved}}</span></h2>
 {{if .Rows}}
 <table><thead><tr><th class="l">Semaine</th><th class="l">Puzzle</th><th>Vues</th><th>Tentatives</th><th>Résolus</th><th>Taux</th></tr></thead>
 <tbody>{{range .Rows}}<tr><td class="l">{{.Week}}</td><td class="l">{{.PuzzleID}}</td><td>{{.Views}}</td><td>{{.Attempts}}</td><td>{{.Solved}}</td><td>{{.SolveRate}}%</td></tr>{{end}}</tbody></table>
 {{else}}<p class="empty">Aucune interaction enregistrée.</p>{{end}}
</section>
</body></html>`))
