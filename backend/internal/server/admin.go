package server

import (
	"crypto/subtle"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/iwanesko/chess-web-site/backend/internal/newsletter"
	"github.com/iwanesko/chess-web-site/backend/internal/stats"
)

// adminAuth garde le tableau de bord et l'espace privé. Compte unique en Basic
// Auth : l'identifiant ET le mot de passe doivent correspondre, comparés en
// temps constant. Si le mot de passe est vidé par configuration, la route
// répond 404 plutôt que 401 — elle ne révèle pas qu'elle existe.
func (s *Server) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.AdminToken == "" {
			http.NotFound(w, r)
			return
		}
		user, pass, ok := r.BasicAuth()
		// Les deux comparaisons sont évaluées AVANT le test : un `&&` court-
		// circuiterait, et le temps de réponse trahirait alors lequel des deux
		// champs est faux.
		okUser := subtle.ConstantTimeCompare([]byte(user), []byte(s.cfg.AdminUser)) == 1
		okPass := subtle.ConstantTimeCompare([]byte(pass), []byte(s.cfg.AdminToken)) == 1
		if !ok || !okUser || !okPass {
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
	Referrers                          []referrerRow
	TotalHuman, TotalBot, TotalUniques int
	BotPct                             int
	// Outils de l'espace privé
	HasCorpus bool
	HasGames  bool
	// Nonce : le panneau de téléversement porte du JavaScript, et la politique
	// par défaut est script-src 'self'. Sans nonce il serait bloqué en silence.
	Nonce string
	// Réservations
	Bookings []bookingRow
	// Newsletter
	Subscribers                           []subscriberRow
	TotalSubs, ConfirmedSubs, PendingSubs int
	ConfirmPct                            int
	SubsFR, SubsEN                        int
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

type subscriberRow struct {
	Email, Lang, Status, Created, Confirmed string
	Pending                                 bool
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

	view.Nonce = newNonce()
	view.HasCorpus = s.cfg.CorpusDB != ""
	view.HasGames = s.cfg.GamesDB != ""

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
	if refs, err := s.store.TopReferrers(12); err == nil {
		for _, r := range refs {
			view.Referrers = append(view.Referrers, referrerRow{refLabel(r.Host), r.Host, r.Count})
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

	// Newsletter
	if s.news != nil {
		subs, _ := s.news.Subscribers()
		for _, sub := range subs {
			row := subscriberRow{
				Email:   sub.Email,
				Lang:    strings.ToUpper(sub.Lang),
				Status:  "confirmé",
				Created: sub.CreatedAt.Format("2006-01-02"),
			}
			if sub.Status == newsletter.StatusConfirmed {
				view.ConfirmedSubs++
				row.Confirmed = sub.ConfirmedAt.Format("2006-01-02")
			} else {
				view.PendingSubs++
				row.Status, row.Pending, row.Confirmed = "en attente", true, "—"
			}
			if strings.EqualFold(sub.Lang, "en") {
				view.SubsEN++
			} else {
				view.SubsFR++
			}
			view.Subscribers = append(view.Subscribers, row)
		}
		view.TotalSubs = len(view.Subscribers)
		if view.TotalSubs > 0 {
			view.ConfirmPct = view.ConfirmedSubs * 100 / view.TotalSubs
		}
	}

	// Le CSP par défaut interdit le script inline du panneau de téléversement :
	// on émet la politique avec le nonce de cette page.
	w.Header().Set("Content-Security-Policy", cspHeader(view.Nonce))
	_ = adminTmpl.Execute(w, view)
}

// flagEmoji turns a 2-letter country code into its flag emoji.
// referrerRow : l'hôte stocké, plus un libellé lisible pour le tableau.
type referrerRow struct {
	Label, Host string
	Count       int
}

// refLabel nomme les hôtes qu'on a une raison de suivre — les moteurs
// génératifs en premier, puisque c'est la question qu'on se pose : est-ce que
// le site se fait citer ? Un hôte inconnu garde son nom brut.
var refLabels = map[string]string{
	"chatgpt.com":           "ChatGPT",
	"openai.com":            "ChatGPT",
	"perplexity.ai":         "Perplexity",
	"claude.ai":             "Claude",
	"copilot.microsoft.com": "Copilot",
	"gemini.google.com":     "Gemini",
	"google.com":            "Google",
	"google.ch":             "Google",
	"google.fr":             "Google",
	"bing.com":              "Bing",
	"duckduckgo.com":        "DuckDuckGo",
	"ecosia.org":            "Ecosia",
	"qwant.com":             "Qwant",
	"lichess.org":           "Lichess",
	"chess.com":             "Chess.com",
	"chess-results.com":     "chess-results",
	"swisschess.ch":         "Swiss Chess",
	"ratings.fide.com":      "FIDE",
	"linkedin.com":          "LinkedIn",
	"facebook.com":          "Facebook",
	"instagram.com":         "Instagram",
	"t.co":                  "X / Twitter",
	"x.com":                 "X / Twitter",
	"reddit.com":            "Reddit",
}

func refLabel(host string) string {
	if l, ok := refLabels[host]; ok {
		return l
	}
	return host
}

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
 /* Même palette que l'explorateur de corpus : l'espace privé doit se lire
    comme un seul outil, pas comme deux pages sans rapport. */
 :root{--bg:#14161a;--panel:#1c1f26;--line:#2a2f39;--ink:#e6e8ec;--dim:#98a0ae;
   --accent:#7aa2f7;--warm:#e0af68;--ok:#5fb87a;--soft:#232833}
 *{box-sizing:border-box}
 body{font:15px/1.55 system-ui,-apple-system,Segoe UI,Roboto,sans-serif;
   color:var(--ink);background:var(--bg);margin:0;padding:0 1.25rem 3rem;
   max-width:1100px;margin:0 auto;-webkit-text-size-adjust:100%}
 a{color:var(--accent)}

 .top{display:flex;align-items:center;justify-content:space-between;gap:1rem;
   flex-wrap:wrap;padding:1rem 0 .85rem;border-bottom:1px solid var(--line);
   margin-bottom:1.1rem}
 h1{font-size:1.15rem;margin:0;font-weight:700;letter-spacing:-.01em}
 h2{font-size:.95rem;margin:1.6rem 0 .5rem;font-weight:600}
 .sub{color:var(--dim);font-weight:400}

 /* Les outils de l'espace privé : c'est le point d'entrée, il doit sauter aux yeux. */
 .tools{display:flex;gap:.5rem;flex-wrap:wrap}
 .tool{display:inline-flex;align-items:center;gap:.45rem;text-decoration:none;
   background:var(--soft);border:1px solid var(--line);border-radius:.6rem;
   padding:.5rem .85rem;color:var(--ink);font-weight:600;font-size:.88rem;
   transition:border-color .15s,background .15s}
 .tool:hover{background:#2b323e;border-color:var(--accent)}
 .tool .ic{color:var(--warm);font-size:1rem;line-height:1}
 .tool.off{opacity:.45;pointer-events:none}
 .tool.off .ic{color:var(--dim)}
 .tool small{color:var(--dim);font-weight:400}

 .tabin{position:absolute;width:0;height:0;opacity:0}
 .tabs{display:flex;gap:.35rem;flex-wrap:wrap;margin-bottom:1.25rem}
 .tabs label{cursor:pointer;padding:.5rem .9rem;border-radius:.55rem;
   font-weight:600;font-size:.9rem;color:var(--dim);user-select:none;
   border:1px solid transparent}
 .tabs label:hover{color:var(--ink);background:var(--soft)}
 .panel{display:none}
 #t1:checked~#p1,#t2:checked~#p2,#t3:checked~#p3,#t4:checked~#p4,#t5:checked~#p5{display:block}
 #t1:checked~.tabs label[for=t1],#t2:checked~.tabs label[for=t2],
 #t3:checked~.tabs label[for=t3],#t4:checked~.tabs label[for=t4],
 #t5:checked~.tabs label[for=t5]{
   color:var(--ink);background:var(--panel);border-color:var(--line)}

 .cards{display:grid;gap:.7rem;grid-template-columns:repeat(3,1fr);margin:.5rem 0 1.1rem}
 @media(max-width:620px){.cards{grid-template-columns:1fr}}
 .kpi{background:var(--panel);border:1px solid var(--line);border-radius:.75rem;
   padding:.9rem 1rem}
 .kpi .n{font-size:1.7rem;font-weight:800;font-variant-numeric:tabular-nums;
   letter-spacing:-.02em}
 .kpi .l{color:var(--dim);font-size:.78rem;margin-top:.15rem}

 .grid2{display:grid;gap:1.25rem;grid-template-columns:1fr 1fr}
 @media(max-width:720px){.grid2{grid-template-columns:1fr}}

 /* Un tableau doit pouvoir défiler seul sur téléphone plutôt que de pousser
    la page : sept colonnes n'entrent pas dans 390 px. */
 table{border-collapse:collapse;width:100%;margin-top:.4rem;
   font-variant-numeric:tabular-nums;background:var(--panel);
   border:1px solid var(--line);border-radius:.7rem;overflow:hidden;font-size:.9rem}
 th,td{padding:.55rem .7rem;border-bottom:1px solid var(--line);text-align:right}
 th:first-child,td:first-child,th.l,td.l{text-align:left}
 thead th{font-size:.72rem;text-transform:uppercase;letter-spacing:.05em;
   color:var(--dim);background:#20242d;font-weight:600}
 tbody tr:last-child td{border-bottom:0}
 tbody tr:hover{background:#ffffff08}
 .p{max-width:15rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
 .bar{height:8px;border-radius:4px;background:#2a2f39;overflow:hidden;min-width:70px}
 .bar>span{display:block;height:100%;background:var(--ok)}
 .empty{color:var(--dim);margin:.5rem 0;font-style:italic}
 .flag{font-size:1.1rem}
 .tag{display:inline-block;padding:.12rem .5rem;border-radius:.4rem;font-size:.75rem;
   font-weight:600;background:#1e3a2a;color:#7ee2a8}
 .tag.wait{background:#3a2f16;color:#e8c37a}
 .up{background:var(--panel);border:1px solid var(--line);border-radius:.75rem;
   padding:1rem;margin:.9rem 0}
 .up h3{margin:0 0 .7rem;font-size:.95rem}
 .up-row{display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
 .up input[type=file]{flex:1 1 16rem;color:var(--dim);font-size:.85rem}
 .up button{background:#2b323e;color:var(--ink);border:1px solid var(--line);
   border-radius:.5rem;padding:.55rem 1rem;cursor:pointer;font-size:.9rem}
 .up button:hover{background:#353d4b}
 .up button:disabled{opacity:.45;cursor:default}
 .up .bar{margin:.8rem 0 .4rem;height:8px;width:100%;min-width:0}
 .up .bar>span{background:var(--accent);transition:width .2s}
 .up .state{margin:0;color:var(--dim);font-size:.85rem}
 code{background:#20242d;padding:.05rem .3rem;border-radius:.25rem;font-size:.85em}
</style></head><body>
<div class="top">
 <h1>Tableau de bord <span class="sub">— privé</span></h1>
 <nav class="tools">
  {{if .HasCorpus}}<a class="tool" href="/admin/corpus/"><span class="ic">&#9816;</span> Explorateur de corpus</a>
  {{else}}<span class="tool off"><span class="ic">&#9816;</span> Explorateur de corpus <small>— CORPUS_DB non configuré</small></span>{{end}}
  {{if .HasGames}}<a class="tool" href="/admin/parties/"><span class="ic">&#9820;</span> Base de parties</a>
  {{else}}<span class="tool off"><span class="ic">&#9820;</span> Base de parties <small>— GAMES_DB non configuré</small></span>{{end}}
 </nav>
</div>

<input type="radio" name="tab" id="t1" class="tabin" checked>
<input type="radio" name="tab" id="t2" class="tabin">
<input type="radio" name="tab" id="t3" class="tabin">
<input type="radio" name="tab" id="t4" class="tabin">
<input type="radio" name="tab" id="t5" class="tabin">
<nav class="tabs">
 <label for="t1">Fréquentation</label>
 <label for="t2">Réservations{{if .Bookings}} ({{len .Bookings}}){{end}}</label>
 <label for="t3">Tactiques</label>
 <label for="t4">Newsletter{{if .TotalSubs}} ({{.ConfirmedSubs}}){{end}}</label>
 <label for="t5">Bases</label>
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
   <h2>Provenance</h2>
   {{if .Referrers}}
   <table><thead><tr><th class="l">Source</th><th>Visites</th></tr></thead>
   <tbody>{{range .Referrers}}<tr><td class="l">{{.Label}}<span class="sub"> {{.Host}}</span></td><td>{{.Count}}</td></tr>{{end}}</tbody></table>
   {{else}}<p class="empty">Aucune provenance externe enregistrée — toutes les visites sont directes ou internes.</p>{{end}}
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
<section class="panel" id="p4">
 <div class="cards">
  <div class="kpi"><div class="n">{{.ConfirmedSubs}}</div><div class="l">Inscrits confirmés</div></div>
  <div class="kpi"><div class="n">{{.PendingSubs}}</div><div class="l">En attente de confirmation</div></div>
  <div class="kpi"><div class="n">{{.ConfirmPct}} %</div><div class="l">Taux de confirmation ({{.TotalSubs}} inscriptions)</div></div>
 </div>
 <p class="sub" style="font-size:.8rem">Double opt-in : seuls les confirmés reçoivent quoi que ce soit. Se désinscrire supprime la ligne — un ancien inscrit ne laisse aucune trace ici. Les jetons de confirmation et de désinscription ne sont volontairement pas affichés.</p>
 <h2>Liste <span class="sub">— {{.SubsFR}} FR · {{.SubsEN}} EN</span></h2>
 {{if .Subscribers}}
 <table><thead><tr><th class="l">E-mail</th><th class="l">Langue</th><th class="l">Statut</th><th class="l">Inscrit le</th><th class="l">Confirmé le</th></tr></thead>
 <tbody>{{range .Subscribers}}<tr><td class="l p">{{.Email}}</td><td class="l">{{.Lang}}</td>
   <td class="l"><span class="tag{{if .Pending}} wait{{end}}">{{.Status}}</span></td>
   <td class="l">{{.Created}}</td><td class="l">{{.Confirmed}}</td></tr>{{end}}</tbody></table>
 {{else}}<p class="empty">Aucune inscription.</p>{{end}}
</section>
<section class="panel" id="p5">
 <h2>Bases de l'espace privé <span class="sub">— téléversement</span></h2>
 <p class="empty" style="font-style:normal">
  L'envoi se fait par morceaux de 8 Mo et reprend là où il s'est arrêté. La base
  en ligne n'est remplacée qu'APRÈS vérification du fichier reçu : si elle
  échoue, rien ne bouge. L'ancienne est conservée en <code>.bak</code>.
 </p>
 {{if .HasCorpus}}
 <div class="up" data-target="corpus">
  <h3>corpus.db <span class="sub">— théorie et commentaires de cours</span></h3>
  <div class="up-row">
   <input type="file" accept=".db,.sqlite,.sqlite3">
   <button class="send">Envoyer</button>
   <button class="cancel" hidden>Annuler</button>
  </div>
  <div class="bar"><span style="width:0"></span></div>
  <p class="state">—</p>
 </div>
 {{end}}
 {{if .HasGames}}
 <div class="up" data-target="games">
  <h3>mega.db <span class="sub">— base de parties</span></h3>
  <div class="up-row">
   <input type="file" accept=".db,.sqlite,.sqlite3">
   <button class="send">Envoyer</button>
   <button class="cancel" hidden>Annuler</button>
  </div>
  <div class="bar"><span style="width:0"></span></div>
  <p class="state">—</p>
  </div>
 {{end}}
 {{if not (or .HasCorpus .HasGames)}}
 <p class="empty">Ni CORPUS_DB ni GAMES_DB ne sont configurés : il n'y a nulle part où écrire.</p>
 {{end}}
</section>

<script nonce="{{ .Nonce }}">
// Téléversement par morceaux. Deux raisons de ne pas faire un seul POST : le
// timeout du frontal HAProxy couperait une requête de plusieurs centaines de
// mégaoctets, et un échec à 90 % obligerait à tout recommencer.
//
// Le serveur vérifie l'offset de chaque morceau et répond 409 avec l'offset
// réel en cas de désaccord : c'est ce qui permet de reprendre, et ce qui
// empêche deux envois concurrents d'écrire une base mélangée.
(function(){
  var CHUNK = 8 * 1024 * 1024;

  function size(n){
    if(n >= 1073741824) return (n/1073741824).toFixed(1) + " Go";
    if(n >= 1048576) return Math.round(n/1048576) + " Mo";
    if(n > 0) return Math.round(n/1024) + " Ko";
    return "0";
  }

  document.querySelectorAll(".up").forEach(function(box){
    var target = box.dataset.target;
    var file = box.querySelector("input[type=file]");
    var send = box.querySelector(".send");
    var cancel = box.querySelector(".cancel");
    var bar = box.querySelector(".bar > span");
    var state = box.querySelector(".state");
    var stop = false;

    function say(msg){ state.textContent = msg; }
    function api(path, opt){
      return fetch("/admin/upload/" + path + "?target=" + target,
                   Object.assign({credentials:"same-origin"}, opt || {}));
    }

    api("status").then(function(r){ return r.json(); }).then(function(s){
      var parts = [];
      if(s.live) parts.push("en ligne : " + size(s.live) + " (" + s.liveModified + ")");
      if(s.uploaded) parts.push("envoi interrompu à " + size(s.uploaded) + " — relancer reprendra là");
      say(parts.join(" · ") || "aucune base en ligne");
    }).catch(function(){});

    cancel.addEventListener("click", function(){
      stop = true;
      api("abort", {method:"POST"}).then(function(){ say("annulé"); });
    });

    send.addEventListener("click", function(){
      var f = file.files[0];
      if(!f){ say("choisis un fichier"); return; }
      stop = false;
      send.disabled = true; file.disabled = true; cancel.hidden = false;

      api("status").then(function(r){ return r.json(); }).then(function(s){
        // On ne reprend que si le début correspond : un autre fichier repart
        // de zéro, sinon on collerait deux bases bout à bout.
        var start = (s.uploaded && s.uploaded < f.size) ? s.uploaded : 0;
        return push(f, start);
      }).then(function(){
        if(stop) return;
        say("vérification…");
        return api("commit", {method:"POST"}).then(function(r){
          return r.json().then(function(j){
            if(!r.ok) throw new Error(j.error || r.status);
            say("remplacée — " + size(j.size) + ". L'ancienne est gardée en .bak.");
            bar.style.width = "100%";
          });
        });
      }).catch(function(e){
        say("échec : " + (e.message || e));
      }).then(function(){
        send.disabled = false; file.disabled = false; cancel.hidden = true;
      });
    });

    function push(f, offset){
      if(stop || offset >= f.size) return Promise.resolve();
      var end = Math.min(offset + CHUNK, f.size);
      return api("chunk&offset=" + offset, {method:"POST", body:f.slice(offset, end)})
        .then(function(r){
          return r.json().then(function(j){
            if(r.status === 409){
              // Le serveur sait mieux que nous où il en est : on se recale.
              return push(f, j.expected);
            }
            if(!r.ok) throw new Error(j.error || r.status);
            bar.style.width = (100 * j.uploaded / f.size) + "%";
            say(size(j.uploaded) + " / " + size(f.size));
            return push(f, j.uploaded);
          });
        });
    }
  });
})();
</script>

</body></html>`))
