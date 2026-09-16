package server

// Page de l'explorateur. Reprise de la disposition du prototype Python —
// échiquier à gauche, coups et commentaires à droite, barre de fréquence sous
// chaque coup — avec trois différences qui viennent du téléphone :
//
//   - l'échiquier est fluide (repeat(8,1fr) dans un carré) et non figé à 514 px,
//     qui débordait d'un écran de 390 px ;
//   - les pièces sont les SVG cburnett du blog, parce que les glyphes Unicode
//     sont rendus de façon inégale selon les plateformes ;
//   - on peut jouer en touchant l'échiquier, pas seulement la liste des coups.
//
// Aucun gabarit JS à accents graves dans ce fichier : il est lui-même une chaîne
// Go entre accents graves.
const corpusHTML = `<!doctype html>
<html lang="fr"><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<meta name="robots" content="noindex,nofollow">
<title>Explorateur de corpus</title>
<style>
:root{
  --bg:#14161a; --panel:#1c1f26; --line:#2a2f39; --ink:#e6e8ec;
  --dim:#98a0ae; --accent:#7aa2f7; --warm:#e0af68; --light:#b9c0cc;
  --sq-l:#f0d9b5; --sq-d:#b58863; --hl:#c8a04a;
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--ink);
  font:15px/1.55 system-ui,-apple-system,Segoe UI,Roboto,sans-serif;
  -webkit-text-size-adjust:100%}
header{display:flex;gap:14px;align-items:baseline;padding:10px 14px;
  border-bottom:1px solid var(--line);flex-wrap:wrap;position:sticky;top:0;
  background:var(--bg);z-index:5}
header h1{font-size:14px;margin:0;font-weight:600;letter-spacing:.02em}
header .stat{color:var(--dim);font-size:12px}

main{display:grid;gap:14px;padding:14px;align-items:start;
  grid-template-columns:minmax(0,514px) minmax(240px,300px) minmax(280px,1fr)}
main>*{min-width:0}
@media(max-width:1240px){main{grid-template-columns:minmax(0,514px) minmax(260px,1fr)}}
@media(max-width:860px){main{grid-template-columns:minmax(0,1fr);padding:10px;gap:10px}}

.left{width:100%;max-width:514px;margin:0 auto}
/* Fluide : c'est le carré qui impose la taille, pas des cases de 64 px. */
#board{display:grid;grid-template-columns:repeat(8,1fr);aspect-ratio:1;
  border:1px solid var(--line);border-radius:4px;overflow:hidden;
  touch-action:manipulation;user-select:none}
.sq{position:relative;display:flex;align-items:center;justify-content:center}
.sq.l{background:var(--sq-l)} .sq.d{background:var(--sq-d)}
.sq.from,.sq.to{box-shadow:inset 0 0 0 4px var(--hl)}
.sq.sel{box-shadow:inset 0 0 0 4px var(--accent)}
.sq svg{width:100%;height:100%;display:block;pointer-events:none}
/* Destination jouable : une pastille, plus lisible qu'un cadre au doigt. */
.sq.dest::after{content:"";position:absolute;width:30%;height:30%;
  border-radius:50%;background:#7aa2f7aa;pointer-events:none}
.sq.dest.occupied::after{width:82%;height:82%;background:none;
  border:5px solid #7aa2f7aa;border-radius:50%}
.coord{position:absolute;bottom:1px;right:3px;font-size:9px;color:#0006;
  pointer-events:none}

.bar-row{display:flex;gap:6px;padding:10px 0;flex-wrap:wrap}
button{background:#262b34;color:var(--ink);border:1px solid var(--line);
  border-radius:6px;padding:9px 13px;cursor:pointer;font-size:14px;min-height:40px}
button:hover{background:#2f3540}
button:disabled{opacity:.4;cursor:default}

.panel{background:var(--panel);border:1px solid var(--line);border-radius:8px}
.panel h2{margin:0;padding:9px 12px;font-size:11px;letter-spacing:.08em;
  text-transform:uppercase;color:var(--dim);border-bottom:1px solid var(--line)}
.body{max-height:min(70vh,560px);overflow:auto;-webkit-overflow-scrolling:touch}

.mv{display:grid;grid-template-columns:1fr auto auto;gap:10px;align-items:center;
  padding:11px 12px;cursor:pointer;border-bottom:1px solid #0003;position:relative}
.mv:hover{background:#ffffff0d}
.mv .san{font-weight:600;font-variant-numeric:tabular-nums}
.mv .pct{color:var(--dim);font-size:13px;font-variant-numeric:tabular-nums}
.mv .dot{color:var(--warm);font-size:11px}
.mv .bar{position:absolute;left:0;bottom:0;height:2px;background:var(--accent);opacity:.55}

.note{padding:10px 12px;border-bottom:1px solid #0003}
.note .src{font-size:11px;color:var(--accent);margin-bottom:3px;word-break:break-word}
.note .txt{color:var(--light);white-space:pre-wrap;overflow-wrap:anywhere}
.empty{padding:14px 12px;color:var(--dim);font-style:italic}

#line{padding:8px 12px;line-height:2.1;max-height:170px;overflow:auto;
  overflow-wrap:anywhere}
#line span{cursor:pointer;padding:3px 6px;border-radius:4px}
#line span.now{background:var(--accent);color:#0b0d10;font-weight:600}
.num{color:var(--dim);cursor:default !important}
input[type=text]{background:#0f1115;color:var(--ink);border:1px solid var(--line);
  border-radius:6px;padding:9px 10px;font-size:14px;width:100%}
.filter-row{padding:9px 12px;border-bottom:1px solid var(--line)}
</style></head><body>

<header>
  <h1>Explorateur de corpus</h1>
  <span class="stat" id="meta"></span>
  <span class="stat" id="reach"></span>
</header>

<main>
  <div class="left">
    <div id="board"></div>
    <div class="bar-row">
      <button id="b-back" onclick="back()">&larr;</button>
      <button id="b-fwd" onclick="fwd()">&rarr;</button>
      <button onclick="home()">Début</button>
      <button onclick="flip()">Retourner</button>
    </div>
    <div class="panel"><h2>Ligne</h2><div id="line"></div></div>
  </div>

  <div class="panel">
    <h2>Coups</h2>
    <div class="body" id="moves"></div>
  </div>

  <div class="panel">
    <h2 id="notes-head">Commentaires</h2>
    <div class="filter-row">
      <input type="text" id="filter" placeholder="filtrer par chapitre ou par mot…" oninput="render()">
    </div>
    <div class="body" id="notes"></div>
  </div>
</main>

<script nonce="{{ .Nonce }}">
var PIECES = {{ pieces .Pieces }};
var START = {{ .Start }};
var API = "/admin/corpus/api";

var line = [{fen:START, san:null, uci:null}], cur = 0, flipped = false;
var data = null, sel = null;

function api(path){
  return fetch(path, {credentials:"same-origin"}).then(function(r){
    if(!r.ok) return r.json().then(function(j){ throw new Error(j.error || r.status); });
    return r.json();
  });
}

function load(){
  sel = null;
  return api(API + "/pos?fen=" + encodeURIComponent(line[cur].fen))
    .then(function(d){ data = d; render(); })
    .catch(fail);
}

function play(uci){
  return api(API + "/go?fen=" + encodeURIComponent(line[cur].fen) + "&uci=" + uci)
    .then(function(d){
      line = line.slice(0, cur + 1);
      line.push({fen:d.fen, san:d.san, uci:uci});
      cur = line.length - 1;
      data = d; sel = null; render();
    })
    .catch(fail);
}

function fail(e){
  document.getElementById("moves").innerHTML =
    "<p class='empty'>" + String(e.message || e) + "</p>";
}

function back(){ if(cur > 0){ cur--; load(); } }
function fwd(){ if(cur < line.length - 1){ cur++; load(); } }
function home(){ cur = 0; load(); }
function jump(i){ cur = i; load(); }
function flip(){ flipped = !flipped; render(); }

document.addEventListener("keydown", function(e){
  if(e.target.tagName === "INPUT") return;
  if(e.key === "ArrowLeft") back();
  if(e.key === "ArrowRight") fwd();
});

function squares(fen){
  var rows = fen.split(" ")[0].split("/"), out = [];
  for(var r = 0; r < 8; r++){
    var row = [];
    for(var i = 0; i < rows[r].length; i++){
      var c = rows[r][i];
      if(c >= "1" && c <= "8"){ for(var k = 0; k < +c; k++) row.push(null); }
      else row.push(c);
    }
    out.push(row);
  }
  return out; // out[0] = 8e rangée
}

function name(r, f){ return "abcdefgh"[f] + (8 - r); }

function renderBoard(){
  var grid = squares(line[cur].fen);
  var last = line[cur].uci;
  var dests = {};
  if(sel && data){
    data.moves.forEach(function(m){
      if(m.uci.slice(0,2) === sel) dests[m.uci.slice(2,4)] = m.uci;
    });
  }
  var html = "";
  for(var i = 0; i < 8; i++){
    for(var j = 0; j < 8; j++){
      var r = flipped ? 7 - i : i, f = flipped ? 7 - j : j;
      var sq = name(r, f), pc = grid[r][f];
      var cls = "sq " + ((r + f) % 2 ? "d" : "l");
      if(last && (last.slice(0,2) === sq)) cls += " from";
      if(last && (last.slice(2,4) === sq)) cls += " to";
      if(sel === sq) cls += " sel";
      if(dests[sq]) cls += " dest" + (pc ? " occupied" : "");
      html += "<div class='" + cls + "' data-sq='" + sq + "'>" +
              (pc ? PIECES[pc] : "") +
              (j === 0 ? "<span class='coord'>" + (8 - r) + "</span>" : "") +
              "</div>";
    }
  }
  var b = document.getElementById("board");
  b.innerHTML = html;
  b.onclick = function(ev){
    var cell = ev.target.closest(".sq");
    if(!cell || !data) return;
    var sq = cell.dataset.sq;
    if(sel && dests[sq]){ play(dests[sq]); return; }
    // On ne sélectionne que les cases d'où part au moins un coup du corpus.
    var any = data.moves.some(function(m){ return m.uci.slice(0,2) === sq; });
    sel = (any && sel !== sq) ? sq : null;
    renderBoard();
  };
}

function renderMoves(){
  var el = document.getElementById("moves");
  if(!data || !data.moves.length){
    el.innerHTML = "<p class='empty'>Position inconnue du corpus.</p>";
    return;
  }
  var top = data.moves[0].n || 1;
  el.innerHTML = data.moves.map(function(m){
    return "<div class='mv' onclick=\"play('" + m.uci + "')\">" +
      "<span class='san'>" + m.san + "</span>" +
      "<span class='dot'>" + (m.notes ? "● " + m.notes : "") + "</span>" +
      "<span class='pct'>" + m.pct.toFixed(1) + " %<br>" + m.n + "</span>" +
      "<span class='bar' style='width:" + (100 * m.n / top) + "%'></span>" +
      "</div>";
  }).join("");
}

function renderNotes(){
  var el = document.getElementById("notes");
  var head = document.getElementById("notes-head");
  if(!data){ el.innerHTML = ""; return; }
  var q = document.getElementById("filter").value.toLowerCase();
  var list = data.notes.filter(function(n){
    return !q || n.label.toLowerCase().indexOf(q) >= 0 || n.text.toLowerCase().indexOf(q) >= 0;
  });
  head.textContent = "Commentaires" +
    (data.notes.length ? " — " + list.length + "/" + data.notesTotal : "");
  if(!list.length){
    el.innerHTML = "<p class='empty'>" +
      (data.notes.length ? "Rien ne correspond au filtre." : "Aucun commentaire ici.") + "</p>";
    return;
  }
  el.innerHTML = list.map(function(n){
    return "<div class='note'><div class='src'>" + esc(n.label) + "</div>" +
           "<div class='txt'>" + esc(n.text) + "</div></div>";
  }).join("");
}

function renderLine(){
  var out = "";
  for(var i = 0; i < line.length; i++){
    if(i === 0){
      out += "<span class='" + (cur === 0 ? "now" : "") + "' onclick='jump(0)'>début</span> ";
      continue;
    }
    if(i % 2 === 1) out += "<span class='num'>" + ((i + 1) / 2 | 0) + ".</span> ";
    out += "<span class='" + (i === cur ? "now" : "") + "' onclick='jump(" + i + ")'>" +
           line[i].san + "</span> ";
  }
  var el = document.getElementById("line");
  el.innerHTML = out;
  var now = el.querySelector(".now");
  if(now) now.scrollIntoView({block:"nearest"});
}

function esc(s){
  return String(s).replace(/[&<>"]/g, function(c){
    return {"&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;"}[c];
  });
}

function render(){
  renderBoard(); renderMoves(); renderNotes(); renderLine();
  document.getElementById("reach").textContent =
    data ? (data.reach ? data.reach.toLocaleString("fr-CH") + " parties passent ici" : "hors corpus") : "";
  document.getElementById("b-back").disabled = cur === 0;
  document.getElementById("b-fwd").disabled = cur >= line.length - 1;
}

api(API + "/meta").then(function(m){
  document.getElementById("meta").textContent =
    (m.games || "?") + " parties · " + (m.moves || "?") + " coups · construit le " + (m.built || "?");
}).catch(function(){});
load();
</script>
</body></html>`
