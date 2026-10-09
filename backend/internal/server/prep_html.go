package server

import "html/template"

// Page de préparation en ligne. Elle ne diffère de l'explorateur de parties que
// par sa barre de recherche : des pseudos Lichess / Chess.com au lieu des
// joueurs de la base, un bouton « Charger », et les favoris. Tout le reste —
// échiquier, arbre, partie rejouée — vient de explorerJS.
//
// Aucun accent grave dans ce fichier, et aucun attribut onclick= (CSP).
const prepHTML = `<!doctype html>
<html lang="fr"><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<meta name="robots" content="noindex,nofollow">
<title>Préparation en ligne</title>
<style>
` + explorerCSS + `
.search+.search{border-top:0;padding-top:0}
.chip .src{font-size:11px;letter-spacing:.04em;color:var(--dim);text-transform:uppercase}
.chip .star{color:var(--warm)}
.favs{padding:0 14px 12px}
.favs summary{cursor:pointer;color:var(--dim);font-size:13px;padding:4px 0}
.favs .chips{padding:8px 0 0}
.favs .chip .name{cursor:pointer}
.favs .chip .name:hover{text-decoration:underline}
.status{padding:0 14px 10px;color:var(--dim);font-size:13px;min-height:1.2em}
.status.err{color:var(--loss)}
</style></head><body>

<header>
  <h1>Préparation en ligne</h1>
  <span class="stat" id="meta">Parties Lichess et Chess.com, chargées à la demande</span>
  <span class="stat"><a href="/admin">&larr; tableau de bord</a></span>
</header>

<section class="search">
  <div class="field">
    <label for="src">Site</label>
    <select id="src"><option value="lichess">Lichess</option><option value="chesscom">Chess.com</option></select>
  </div>
  <div class="field grow">
    <label for="user">Pseudo</label>
    <input type="text" id="user" placeholder="pseudo, puis Entrée" autocomplete="off"
           autocapitalize="off" spellcheck="false" maxlength="30">
  </div>
  <div class="field"><label>&nbsp;</label><button id="add">Ajouter</button></div>
</section>
<div class="chips" id="chips"></div>
{{if .HasFavorites}}<details class="favs" id="favbox"><summary id="favsum">Favoris</summary>
  <div class="chips" id="favs"></div></details>{{end}}

<section class="search">
  <div class="field">
    <label>Cadences</label>
    <div class="bar-row" style="padding:0" id="speeds">
      <button data-sp="bullet">Bullet</button>
      <button data-sp="blitz" class="on">Blitz</button>
      <button data-sp="rapid" class="on">Rapide</button>
      <button data-sp="classical" class="on">Classique</button>
      <button data-sp="correspondence">Corresp.</button>
    </div>
  </div>
  <div class="field">
    <label for="months">Période</label>
    <select id="months">
      <option value="3">3 mois</option><option value="6">6 mois</option>
      <option value="12" selected>12 mois</option><option value="24">2 ans</option>
      <option value="0">Tout</option>
    </select>
  </div>
  <div class="field"><label for="max">Parties max</label><input type="number" id="max" value="300" min="10" max="1000"></div>
  <div class="field"><label>Classement</label><button id="rated" class="on">Classées seulement</button></div>
  <div class="field"><label>&nbsp;</label><button id="load" class="on">Charger</button></div>
</section>
<section class="search">
  <div class="field">
    <label>Couleur jouée</label>
    <div class="bar-row" style="padding:0">
      <button id="c-all" class="on">Les deux</button>
      <button id="c-w">Blancs</button>
      <button id="c-b">Noirs</button>
    </div>
  </div>
</section>
<div class="status" id="status"></div>

<main>
  <div class="left">
    <div id="board"></div>
    <div class="bar-row">
      <button id="b-back">&larr;</button>
      <button id="b-fwd">&rarr;</button>
      <button id="b-home">Début</button>
      <button id="b-flip">Retourner</button>
    </div>
    <div class="panel"><h2>Ligne</h2><div id="line"></div></div>
  </div>

  <div class="panel">
    <h2 id="moves-head">Coups</h2>
    <div class="body" id="moves"></div>
  </div>

  <div class="panel">
    <h2 id="games-head">Parties</h2>
    <div class="body" id="games"></div>
  </div>
</main>

<script nonce="{{ .Nonce }}">
var PIECES = {{ pieces .Pieces }};
var START = {{ .Start }};
var API = "/admin/prepa/api";
var HAS_FAVS = {{ .HasFavorites }};

var accounts = [];         // comptes à charger : {source, username}
var favs = [];             // favoris enregistrés
var colour = "";           // couleur jouée par le joueur préparé
var setId = null;          // lot chargé côté serveur
var loadState = null;      // son avancement (pas « status » : window.status est une chaîne)
var rated = true;
var pollTimer = null;
var SITE = {lichess:"Lichess", chesscom:"Chess.com"};
var SPEED = {bullet:"Bullet", blitz:"Blitz", rapid:"Rapide", classical:"Classique", correspondence:"Corresp."};
var USER_RE = /^[A-Za-z0-9_-]{2,30}$/;

// ---- ce que la navigation commune attend de la page ------------------------
function ready(){ return !!(setId && loadState && loadState.state === "ready" && loadState.games); }
function emptyText(){
  if(loadState && loadState.state === "loading") return "Chargement… " + loadState.games + " parties reçues.";
  if(loadState && loadState.state === "ready" && !loadState.games) return "Aucune partie ne correspond.";
  return "Ajoute un ou plusieurs pseudos, puis Charger.";
}
function query(extra){ return "set=" + encodeURIComponent(setId || "") + "&colour=" + colour + (extra || ""); }
function sameAccount(a, src, name){ return a.source === src && a.username.toLowerCase() === String(name).toLowerCase(); }
function flipFor(g){ return accounts.some(function(a){ return sameAccount(a, g.source, g.black); }); }
function gameHeadMeta(g){
  return (SPEED[g.speed] || esc(g.speed)) + (g.date ? " &middot; " + esc(g.date) : "") +
    (g.opening ? " &middot; " + esc(g.opening) : (g.eco ? " &middot; " + esc(g.eco) : "")) +
    (g.url ? " &middot; <a href='" + esc(g.url) + "' target='_blank' rel='noopener noreferrer'>voir sur " +
      (SITE[g.source] || "le site") + "</a>" : "");
}
function gameRowMeta(g){
  return (SPEED[g.speed] || esc(g.speed)) + (g.date ? " &middot; " + esc(g.date) : "") +
         (g.opening ? " &middot; " + esc(g.opening) : "");
}

// ---- comptes ----------------------------------------------------------------
function isFav(a){ return favs.some(function(f){ return sameAccount(f, a.source, a.username); }); }
function addAccount(src, name){
  name = String(name || "").trim();
  if(!USER_RE.test(name)){ say("Pseudo invalide : lettres, chiffres, - et _ seulement.", true); return false; }
  if(accounts.some(function(a){ return sameAccount(a, src, name); })) return true;
  if(accounts.length >= 4){ say("Quatre comptes au maximum.", true); return false; }
  accounts.push({source:src, username:name});
  renderChips(); say("");
  return true;
}
function addFromInput(){
  if(addAccount(el("src").value, el("user").value)) el("user").value = "";
}
el("add").addEventListener("click", addFromInput);
el("user").addEventListener("keydown", function(e){
  if(e.key === "Enter"){ e.preventDefault(); addFromInput(); }
});
el("chips").addEventListener("click", function(e){
  var b = e.target.closest("button[data-i]");
  if(!b) return;
  var a = accounts[+b.dataset.i];
  if(b.dataset.act === "rm"){ accounts.splice(+b.dataset.i, 1); renderChips(); return; }
  if(b.dataset.act === "star") toggleFav(a);
});
function renderChips(){
  el("chips").innerHTML = accounts.map(function(a, i){
    return "<span class='chip'><span class='src'>" + SITE[a.source] + "</span>" + esc(a.username) +
      (HAS_FAVS ? "<button data-i='" + i + "' data-act='star' class='star' title='Favori'>" +
        (isFav(a) ? "&#9733;" : "&#9734;") + "</button>" : "") +
      "<button data-i='" + i + "' data-act='rm' title='Retirer'>&times;</button></span>";
  }).join("");
}

// ---- favoris ------------------------------------------------------------------
function loadFavs(){
  if(!HAS_FAVS) return;
  api(API + "/favorites").then(setFavs).catch(function(e){ say(e.message || e, true); });
}
function setFavs(list){
  favs = list || [];
  el("favsum").textContent = "Favoris (" + favs.length + ")";
  el("favs").innerHTML = favs.length ? favs.map(function(f, i){
    return "<span class='chip'><span class='src'>" + SITE[f.source] + "</span>" +
      "<span class='name' data-i='" + i + "'>" + esc(f.username) + "</span>" +
      (f.note ? " <span class='src'>" + esc(f.note) + "</span>" : "") +
      "<button data-i='" + i + "' title='Retirer des favoris'>&times;</button></span>";
  }).join("") : "<span class='status'>Aucun favori : l'étoile d'un compte l'ajoute ici.</span>";
  renderChips();
}
function toggleFav(a){
  var opts = isFav(a)
    ? [API + "/favorites?source=" + a.source + "&username=" + encodeURIComponent(a.username), {method:"DELETE"}]
    : [API + "/favorites", {method:"POST", headers:{"Content-Type":"application/json"},
        body: JSON.stringify({source:a.source, username:a.username})}];
  api(opts[0], opts[1]).then(setFavs).catch(function(e){ say(e.message || e, true); });
}
if(HAS_FAVS){
  el("favs").addEventListener("click", function(e){
    var f, t = e.target.closest("[data-i]");
    if(!t) return;
    f = favs[+t.dataset.i];
    if(t.tagName === "BUTTON") toggleFav(f);
    else addAccount(f.source, f.username);
  });
}

// ---- options de chargement ------------------------------------------------------
el("speeds").addEventListener("click", function(e){
  var b = e.target.closest("button[data-sp]");
  if(b) b.className = b.className === "on" ? "" : "on";
});
el("rated").addEventListener("click", function(){
  rated = !rated;
  el("rated").className = rated ? "on" : "";
  el("rated").textContent = rated ? "Classées seulement" : "Classées et amicales";
});
function setColour(c){
  colour = c;
  flipped = c === "b"; // on prépare de son côté de l'échiquier
  el("c-all").className = c === "" ? "on" : "";
  el("c-w").className = c === "w" ? "on" : "";
  el("c-b").className = c === "b" ? "on" : "";
  reload();
}
el("c-all").addEventListener("click", function(){ setColour(""); });
el("c-w").addEventListener("click", function(){ setColour("w"); });
el("c-b").addEventListener("click", function(){ setColour("b"); });

// ---- chargement -------------------------------------------------------------
function say(msg, err){
  el("status").textContent = msg || "";
  el("status").className = "status" + (err ? " err" : "");
}
el("load").addEventListener("click", function(){
  if(el("user").value.trim()) addFromInput();
  if(!accounts.length){ say("Ajoute au moins un pseudo.", true); return; }
  var speeds = [].slice.call(el("speeds").querySelectorAll("button.on")).map(function(b){ return b.dataset.sp; });
  if(!speeds.length){ say("Choisis au moins une cadence.", true); return; }
  clearTimeout(pollTimer);
  el("load").disabled = true;
  setId = null; loadState = null; reload();
  say("Chargement…");
  api(API + "/load", {method:"POST", headers:{"Content-Type":"application/json"}, body: JSON.stringify({
    accounts: accounts, speeds: speeds, rated: rated,
    months: +el("months").value, max: +el("max").value || 300
  })}).then(function(st){ setId = st.id; track(st); })
    .catch(function(e){ el("load").disabled = false; say(e.message || e, true); });
});
// Lichess envoie une vingtaine de parties par seconde : on suit l'avancement
// plutôt que d'attendre une réponse unique que le serveur couperait à 30 s.
function track(st){
  loadState = st;
  if(st.state === "loading"){
    say("Chargement… " + st.games + " parties reçues");
    renderMoves();
    pollTimer = setTimeout(function(){
      api(API + "/status?set=" + encodeURIComponent(setId)).then(track)
        .catch(function(e){ el("load").disabled = false; say(e.message || e, true); });
    }, 800);
    return;
  }
  el("load").disabled = false;
  if(st.state === "failed"){ say(st.error || "Échec du chargement.", true); setId = null; reload(); return; }
  say(st.games + " parties chargées" + (st.error ? " — " + st.error : ""), !!st.error);
  reload();
}

` + explorerJS + `
loadFavs();
renderChips();
render();
</script>
</body></html>`

var prepTmpl = template.Must(template.New("prepa").Funcs(template.FuncMap{
	"pieces": pieceJSON,
}).Parse(prepHTML))
