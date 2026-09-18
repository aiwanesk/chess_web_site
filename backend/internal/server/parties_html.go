package server

// Page de l'explorateur de parties. Même apparence que celle du corpus — et
// c'est voulu : ce sont deux vues du même geste, avancer dans une ouverture.
//
// La différence est dans la circulation des données. Le corpus demande au
// serveur à chaque position ; ici l'arbre du joueur arrive en entier et toute
// la navigation se fait en mémoire. Seule la liste des parties, en bas, va
// rechercher quelque chose quand on change de ligne.
//
// Aucun accent grave dans ce fichier : il est lui-même une chaîne Go entre
// accents graves. Et aucun attribut onclick= : le CSP les bloque, un nonce ne
// couvrant que les balises script.
const partiesHTML = `<!doctype html>
<html lang="fr"><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<meta name="robots" content="noindex,nofollow">
<title>Explorateur de parties</title>
<style>
:root{
  --bg:#14161a; --panel:#1c1f26; --line:#2a2f39; --ink:#e6e8ec;
  --dim:#98a0ae; --accent:#7aa2f7; --warm:#e0af68; --light:#b9c0cc;
  --win:#5fb87a; --loss:#d97a7a;
  --sq-l:#f0d9b5; --sq-d:#b58863; --hl:#c8a04a;
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--ink);
  font:15px/1.55 system-ui,-apple-system,Segoe UI,Roboto,sans-serif;
  -webkit-text-size-adjust:100%;touch-action:manipulation}
header{display:flex;gap:14px;align-items:baseline;padding:10px 14px;
  border-bottom:1px solid var(--line);flex-wrap:wrap}
header h1{font-size:14px;margin:0;font-weight:600}
header .stat{color:var(--dim);font-size:12px}
a{color:var(--accent)}

.search{padding:12px 14px;border-bottom:1px solid var(--line);
  display:flex;gap:10px;flex-wrap:wrap;align-items:flex-end}
.field{display:flex;flex-direction:column;gap:4px;min-width:0}
.field label{font-size:11px;letter-spacing:.06em;text-transform:uppercase;color:var(--dim)}
.grow{flex:1 1 240px}
input,select{background:#0f1115;color:var(--ink);border:1px solid var(--line);
  border-radius:6px;padding:10px;font-size:16px;width:100%}
input[type=number]{width:6.5rem}
button{background:#262b34;color:var(--ink);border:1px solid var(--line);
  border-radius:6px;padding:10px 14px;cursor:pointer;font-size:15px;min-height:44px}
button:hover{background:#2f3540}
button:disabled{opacity:.4;cursor:default}
button.on{background:var(--accent);color:#0b0d10;border-color:var(--accent);font-weight:600}

.chips{display:flex;gap:6px;flex-wrap:wrap;padding:0 14px 12px}
.chip{display:inline-flex;align-items:center;gap:8px;background:var(--panel);
  border:1px solid var(--line);border-radius:999px;padding:6px 8px 6px 12px;font-size:14px}
.chip button{min-height:0;padding:2px 8px;border-radius:999px;font-size:13px;line-height:1.2}
.sugg{position:relative}
.sugg-list{position:absolute;z-index:20;left:0;right:0;top:100%;margin-top:4px;
  background:var(--panel);border:1px solid var(--line);border-radius:8px;
  max-height:260px;overflow:auto}
.sugg-list div{padding:10px 12px;cursor:pointer;border-bottom:1px solid #0003}
.sugg-list div:hover{background:#ffffff0d}
.sugg-list .n{color:var(--dim);font-size:12px}

main{display:grid;gap:14px;padding:14px;align-items:start;
  grid-template-columns:minmax(0,514px) minmax(260px,340px) minmax(280px,1fr)}
main>*{min-width:0}
@media(max-width:1240px){main{grid-template-columns:minmax(0,514px) minmax(280px,1fr)}}
@media(max-width:860px){main{grid-template-columns:minmax(0,1fr);padding:10px;gap:10px}}

.left{width:100%;max-width:514px;margin:0 auto}
#board{display:grid;grid-template-columns:repeat(8,minmax(0,1fr));
  grid-template-rows:repeat(8,minmax(0,1fr));aspect-ratio:1;
  border:1px solid var(--line);border-radius:4px;overflow:hidden;user-select:none}
.sq{position:relative;display:flex;align-items:center;justify-content:center}
.sq.l{background:var(--sq-l)} .sq.d{background:var(--sq-d)}
.sq.from,.sq.to{box-shadow:inset 0 0 0 4px var(--hl)}
.sq svg{width:100%;height:100%;display:block;pointer-events:none}
.coord{position:absolute;top:1px;left:3px;font-size:9px;color:#0006;pointer-events:none}
.bar-row{display:flex;gap:6px;padding:10px 0;flex-wrap:wrap}

.panel{background:var(--panel);border:1px solid var(--line);border-radius:8px}
.panel h2{margin:0;padding:9px 12px;font-size:11px;letter-spacing:.08em;
  text-transform:uppercase;color:var(--dim);border-bottom:1px solid var(--line)}
.body{max-height:min(70vh,560px);max-height:min(70dvh,560px);overflow:auto}

.mv{display:grid;grid-template-columns:auto 1fr auto;gap:10px;align-items:center;
  padding:11px 12px;cursor:pointer;border-bottom:1px solid #0003;position:relative}
.mv:hover{background:#ffffff0d}
.mv .san{font-weight:600;min-width:3.4rem}
.mv .wdl{display:flex;height:9px;border-radius:3px;overflow:hidden;min-width:70px}
.mv .wdl i{display:block}
.mv .wdl .w{background:var(--win)} .mv .wdl .d{background:#5a6274} .mv .wdl .l{background:var(--loss)}
.mv .num{color:var(--dim);font-size:12.5px;text-align:right;font-variant-numeric:tabular-nums}
.mv .bar{position:absolute;left:0;bottom:0;height:2px;background:var(--accent);opacity:.5}

.g{display:grid;grid-template-columns:1fr auto;gap:8px;padding:9px 12px;
  border-bottom:1px solid #0003;font-size:13.5px}
.g .who{overflow:hidden;text-overflow:ellipsis}
.g .meta{color:var(--dim);font-size:12px}
.g .res{font-variant-numeric:tabular-nums;color:var(--dim)}
.empty{padding:14px 12px;color:var(--dim);font-style:italic}
#line{padding:8px 12px;line-height:2.1;max-height:150px;overflow:auto;overflow-wrap:anywhere}
#line span{cursor:pointer;padding:3px 6px;border-radius:4px}
#line span.now{background:var(--accent);color:#0b0d10;font-weight:600}
.numlbl{color:var(--dim);cursor:default !important}
</style></head><body>

<header>
  <h1>Explorateur de parties</h1>
  <span class="stat" id="meta"></span>
  <span class="stat"><a href="/admin">&larr; tableau de bord</a></span>
</header>

<section class="search">
  <div class="field grow sugg">
    <label for="q">Joueurs</label>
    <input type="text" id="q" placeholder="taper un nom, puis choisir" autocomplete="off">
    <div class="sugg-list" id="sugg" hidden></div>
  </div>
  <div class="field">
    <label>Couleur</label>
    <div class="bar-row" style="padding:0">
      <button id="c-all" class="on">Les deux</button>
      <button id="c-w">Blancs</button>
      <button id="c-b">Noirs</button>
    </div>
  </div>
  <div class="field"><label for="from">De</label><input type="number" id="from" placeholder="1990"></div>
  <div class="field"><label for="to">À</label><input type="number" id="to" placeholder="2026"></div>
  <div class="field">
    <label>Parties en ligne</label>
    <button id="tt" class="on">Titled Tuesday exclus</button>
  </div>
</section>
<div class="chips" id="chips"></div>

<main>
  <div class="left">
    <div id="board"></div>
    <div class="bar-row">
      <button id="b-back">&larr;</button>
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
var API = "/admin/parties/api";

var players = [];          // joueurs sélectionnés
var colour = "";           // "", "w" ou "b"
var noTT = true;
var root = null;           // racine de l'arbre chargé
var path = [];             // suite de nœuds parcourus
var flipped = false;
var gamesTimer = null;

function api(p){
  return fetch(p, {credentials:"same-origin"}).then(function(r){
    if(!r.ok) return r.json().then(function(j){ throw new Error(j.error || r.status); });
    return r.json();
  });
}
function esc(s){
  return String(s).replace(/[&<>"]/g, function(c){
    return {"&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;"}[c];
  });
}
function el(id){ return document.getElementById(id); }
function node(){ return path.length ? path[path.length-1] : root; }
function fen(){ var n = node(); return n && n.fen ? n.fen : START; }
function query(extra){
  var p = "players=" + players.map(function(x){ return x.id; }).join(",") +
          "&colour=" + colour + "&noTT=" + (noTT ? "1" : "0") +
          "&from=" + encodeURIComponent(el("from").value) +
          "&to=" + encodeURIComponent(el("to").value);
  return p + (extra || "");
}

// ---- recherche de joueurs -------------------------------------------------
var suggTimer = null;
el("q").addEventListener("input", function(){
  clearTimeout(suggTimer);
  var q = el("q").value.trim();
  if(q.length < 2){ el("sugg").hidden = true; return; }
  suggTimer = setTimeout(function(){
    api(API + "/players?q=" + encodeURIComponent(q)).then(function(list){
      if(!list.length){ el("sugg").hidden = true; return; }
      el("sugg").innerHTML = list.map(function(p){
        return "<div data-id='" + p.id + "' data-name='" + esc(p.name) + "'>" +
               esc(p.name) + " <span class='n'>" + p.games + " parties</span></div>";
      }).join("");
      el("sugg").hidden = false;
    }).catch(function(){ el("sugg").hidden = true; });
  }, 180);
});
el("sugg").addEventListener("click", function(e){
  var d = e.target.closest("div[data-id]");
  if(!d) return;
  var id = +d.dataset.id;
  if(!players.some(function(p){ return p.id === id; })){
    players.push({id:id, name:d.dataset.name});
  }
  el("q").value = ""; el("sugg").hidden = true;
  renderChips(); reload();
});
el("chips").addEventListener("click", function(e){
  var b = e.target.closest("button[data-id]");
  if(!b) return;
  players = players.filter(function(p){ return p.id !== +b.dataset.id; });
  renderChips(); reload();
});
function renderChips(){
  el("chips").innerHTML = players.map(function(p){
    return "<span class='chip'>" + esc(p.name) +
           "<button data-id='" + p.id + "'>&times;</button></span>";
  }).join("");
}

// ---- filtres --------------------------------------------------------------
function setColour(c){
  colour = c;
  el("c-all").className = c === "" ? "on" : "";
  el("c-w").className = c === "w" ? "on" : "";
  el("c-b").className = c === "b" ? "on" : "";
  reload();
}
el("c-all").addEventListener("click", function(){ setColour(""); });
el("c-w").addEventListener("click", function(){ setColour("w"); });
el("c-b").addEventListener("click", function(){ setColour("b"); });
el("tt").addEventListener("click", function(){
  noTT = !noTT;
  el("tt").className = noTT ? "on" : "";
  el("tt").textContent = noTT ? "Titled Tuesday exclus" : "Titled Tuesday inclus";
  reload();
});
el("from").addEventListener("change", reload);
el("to").addEventListener("change", reload);

// ---- chargement de l'arbre ------------------------------------------------
function reload(){
  path = [];
  if(!players.length){ root = null; render(); return; }
  el("moves").innerHTML = "<p class='empty'>Chargement…</p>";
  api(API + "/tree?" + query()).then(function(d){
    root = {fen:d.fen, children:d.moves};
    render();
  }).catch(function(e){
    root = null; render();
    el("moves").innerHTML = "<p class='empty'>" + esc(e.message || e) + "</p>";
  });
}

// ---- navigation, entièrement locale ---------------------------------------
function play(i){
  var kids = (node() || {}).children || [];
  if(kids[i]){ path.push(kids[i]); render(); }
}
function back(){ if(path.length){ path.pop(); render(); } }
function home(){ path = []; render(); }
function jump(n){ path = path.slice(0, n); render(); }
el("b-back").addEventListener("click", back);
el("b-home").addEventListener("click", home);
el("b-flip").addEventListener("click", function(){ flipped = !flipped; renderBoard(); });
document.addEventListener("keydown", function(e){
  if(e.target.tagName === "INPUT") return;
  if(e.key === "ArrowLeft") back();
});
el("moves").addEventListener("click", function(e){
  var d = e.target.closest(".mv");
  if(d) play(+d.dataset.i);
});
el("line").addEventListener("click", function(e){
  var d = e.target.closest("span[data-n]");
  if(d) jump(+d.dataset.n);
});

// ---- rendu ----------------------------------------------------------------
function squares(f){
  var rows = f.split(" ")[0].split("/"), out = [];
  for(var r = 0; r < 8; r++){
    var row = [];
    for(var i = 0; i < rows[r].length; i++){
      var c = rows[r][i];
      if(c >= "1" && c <= "8"){ for(var k = 0; k < +c; k++) row.push(null); }
      else row.push(c);
    }
    out.push(row);
  }
  return out;
}
function renderBoard(){
  var grid = squares(fen());
  var last = path.length ? path[path.length-1].uci : null;
  var html = "";
  for(var i = 0; i < 8; i++){
    for(var j = 0; j < 8; j++){
      var r = flipped ? 7 - i : i, f = flipped ? 7 - j : j;
      var sq = "abcdefgh"[f] + (8 - r), pc = grid[r][f];
      var cls = "sq " + ((r + f) % 2 ? "d" : "l");
      if(last && last.slice(0,2) === sq) cls += " from";
      if(last && last.slice(2,4) === sq) cls += " to";
      html += "<div class='" + cls + "'>" + (pc ? PIECES[pc] : "") +
              (j === 0 ? "<span class='coord'>" + (8 - r) + "</span>" : "") + "</div>";
    }
  }
  el("board").innerHTML = html;
}
function renderMoves(){
  var n = node();
  if(!players.length){
    el("moves").innerHTML = "<p class='empty'>Choisis au moins un joueur.</p>";
    el("moves-head").textContent = "Coups";
    return;
  }
  var kids = (n || {}).children || [];
  if(!kids.length){
    el("moves").innerHTML = "<p class='empty'>Plus rien à cette profondeur.</p>";
    return;
  }
  var top = kids[0].games || 1;
  el("moves-head").textContent = "Coups — " + kids.reduce(function(a, k){ return a + k.games; }, 0) + " parties";
  el("moves").innerHTML = kids.map(function(k, i){
    var t = k.games || 1;
    return "<div class='mv' data-i='" + i + "'>" +
      "<span class='san'>" + esc(k.san) + "</span>" +
      "<span class='wdl'>" +
        "<i class='w' style='width:" + (100*k.wins/t) + "%'></i>" +
        "<i class='d' style='width:" + (100*k.draws/t) + "%'></i>" +
        "<i class='l' style='width:" + (100*k.losses/t) + "%'></i></span>" +
      "<span class='num'>" + k.games + "<br>" + k.score.toFixed(0) + " %</span>" +
      "<span class='bar' style='width:" + (100*k.games/top) + "%'></span>" +
      "</div>";
  }).join("");
}
function renderLine(){
  var out = "<span class='" + (path.length ? "" : "now") + "' data-n='0'>début</span> ";
  for(var i = 0; i < path.length; i++){
    if(i % 2 === 0) out += "<span class='numlbl'>" + (i/2 + 1) + ".</span> ";
    out += "<span class='" + (i === path.length-1 ? "now" : "") + "' data-n='" + (i+1) + "'>" +
           esc(path[i].san) + "</span> ";
  }
  el("line").innerHTML = out;
}
function renderGames(){
  clearTimeout(gamesTimer);
  if(!players.length){ el("games").innerHTML = "<p class='empty'>—</p>"; return; }
  // La liste est le seul aller-retour restant : on la laisse respirer pendant
  // qu'on descend vite dans l'arbre.
  gamesTimer = setTimeout(function(){
    var p = path.map(function(n){ return n.uci; }).join(",");
    api(API + "/games?" + query("&path=" + p)).then(function(list){
      el("games-head").textContent = "Parties — " + list.length;
      if(!list.length){ el("games").innerHTML = "<p class='empty'>Aucune.</p>"; return; }
      el("games").innerHTML = list.map(function(g){
        var res = g.result > 0 ? "1-0" : (g.result < 0 ? "0-1" : "½-½");
        return "<div class='g'><div class='who'>" +
          esc(g.white) + (g.whiteElo ? " (" + g.whiteElo + ")" : "") + " &ndash; " +
          esc(g.black) + (g.blackElo ? " (" + g.blackElo + ")" : "") +
          "<div class='meta'>" + esc(g.event) + (g.year ? " &middot; " + g.year : "") + "</div>" +
          "</div><div class='res'>" + res + "</div></div>";
      }).join("");
    }).catch(function(){ el("games").innerHTML = "<p class='empty'>—</p>"; });
  }, 220);
}
function render(){
  renderBoard(); renderMoves(); renderLine(); renderGames();
  el("b-back").disabled = path.length === 0;
}

api(API + "/meta").then(function(m){
  el("meta").textContent = (m.games || "?") + " parties" +
    (m.twic ? " · TWIC " + m.twic : "") + (m.built ? " · " + m.built : "");
}).catch(function(){});
renderChips();
render();
</script>
</body></html>`
