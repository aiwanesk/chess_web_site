package server

// Ce que les deux explorateurs — la base de parties et la préparation en
// ligne — ont en commun : l'apparence, l'échiquier, l'arbre qu'on descend coup
// par coup, la partie qu'on rejoue. Chaque page n'écrit que sa barre de
// recherche et quelques fonctions de présentation (voir l'en-tête de
// explorerJS). Une seule navigation, donc un seul endroit à corriger.
//
// Aucun accent grave ici : ce sont des chaînes Go entre accents graves.

const explorerCSS = `:root{
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
.g{cursor:pointer}
.g:hover{background:#ffffff0d}
.g.on{background:var(--accent);color:#0b0d10}
.g.on .meta,.g.on .res{color:#0b0d1099}
/* Mode partie : les coups se lisent en colonnes serrées, pas en lignes larges —
   une partie fait quatre-vingts demi-coups, l'arbre en montre six. */
.plies{display:flex;flex-wrap:wrap;gap:2px;padding:8px 10px;line-height:1.9}
.plies .numlbl{color:var(--dim);font-size:12.5px;padding:3px 2px 3px 6px}
.plies .p{cursor:pointer;padding:3px 6px;border-radius:4px;min-width:3.2rem}
.plies .p:hover{background:#ffffff0d}
.plies .p.now{background:var(--accent);color:#0b0d10;font-weight:600}
.ghead{padding:10px 12px;border-bottom:1px solid var(--line)}
.ghead .who{font-weight:600}
.ghead .meta{color:var(--dim);font-size:12.5px;margin-top:2px}
.ghead button{margin-top:9px}
.empty{padding:14px 12px;color:var(--dim);font-style:italic}
#line{padding:8px 12px;line-height:2.1;max-height:150px;overflow:auto;overflow-wrap:anywhere}
#line span{cursor:pointer;padding:3px 6px;border-radius:4px}
#line span.now{background:var(--accent);color:#0b0d10;font-weight:600}
.numlbl{color:var(--dim);cursor:default !important}
`

const explorerJS = `// ---- navigation commune aux deux explorateurs -------------------------------
// La page définit AVANT ce bloc : API, START, PIECES, et les fonctions
// ready(), emptyText(), query(extra), flipFor(g), gameHeadMeta(g), gameRowMeta(g).
var root = null;           // racine de l'arbre chargé
var path = [];             // suite de nœuds parcourus
var flipped = false;
var gamesTimer = null;
var game = null;           // partie ouverte, ou null en mode arbre
var plies = [];            // ses demi-coups, en chaîne linéaire
var openId = null;         // pour garder la ligne surlignée dans la liste

function api(p, opts){
  var o = {credentials:"same-origin"};
  for(var k in (opts || {})) o[k] = opts[k];
  return fetch(p, o).then(function(r){
    if(!r.ok) return r.json().then(function(j){ throw new Error(j.error || r.status); },
                                   function(){ throw new Error("erreur " + r.status); });
    return r.json();
  });
}
function esc(s){
  return String(s).replace(/[&<>"']/g, function(c){
    return {"&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;","'":"&#39;"}[c];
  });
}
function el(id){ return document.getElementById(id); }
function node(){ return path.length ? path[path.length-1] : root; }
function fen(){ var n = node(); return n && n.fen ? n.fen : START; }

// ---- chargement de l'arbre ------------------------------------------------
function reload(){
  path = []; game = null; plies = []; openId = null;
  if(!ready()){ root = null; render(); return; }
  el("moves").innerHTML = "<p class='empty'>Chargement…</p>";
  api(API + "/tree?" + query()).then(function(d){
    root = {fen:d.fen, children:d.moves};
    render();
  }).catch(function(e){
    root = null; render();
    el("moves").innerHTML = "<p class='empty'>" + esc(e.message || e) + "</p>";
  });
}

// ---- navigation -----------------------------------------------------------
// L'arbre arrive borné en profondeur : le charger entier serait une réponse
// énorme dont on ne lit que les premiers niveaux. On le PROLONGE donc quand on
// arrive au bout d'une branche, en rechargeant depuis le chemin courant. C'est
// ce qui lève la limite des sept coups sans jamais transférer un gros arbre.
function extend(){
  if(game) return;
  var n = node();
  if(!n || n.children || n.loading) return;
  n.loading = true;
  var p = path.map(function(x){ return x.uci; }).join(",");
  api(API + "/tree?" + query("&path=" + p)).then(function(d){
    n.loading = false;
    n.children = d.moves || [];
    if(node() === n) render();
  }).catch(function(){
    n.loading = false; n.children = [];
    if(node() === n) render();
  });
}
function play(i){
  var kids = (node() || {}).children || [];
  if(kids[i]){ path.push(kids[i]); render(); extend(); }
}
function fwd(){
  if(game){ if(path.length < plies.length) jump(path.length + 1); return; }
  play(0); // en mode arbre, avancer = suivre le coup le plus joué
}
function back(){ if(path.length){ path.pop(); render(); } }
function home(){ path = []; render(); }
function jump(n){ path = plies.length ? plies.slice(0, n) : path.slice(0, n); render(); }
el("b-back").addEventListener("click", back);
el("b-fwd").addEventListener("click", fwd);
el("b-home").addEventListener("click", home);
el("b-flip").addEventListener("click", function(){ flipped = !flipped; renderBoard(); });
document.addEventListener("keydown", function(e){
  var t = e.target.tagName;
  if(t === "INPUT" || t === "SELECT" || t === "TEXTAREA") return;
  if(e.key === "ArrowLeft"){ e.preventDefault(); back(); }
  if(e.key === "ArrowRight"){ e.preventDefault(); fwd(); }
});
el("moves").addEventListener("click", function(e){
  var d = e.target.closest(".mv");
  if(d){ play(+d.dataset.i); return; }
  var q = e.target.closest(".p");
  if(q) jump(+q.dataset.ply);
});

// ---- une partie entière ---------------------------------------------------
// Les coups sont déjà connus du serveur, qui y ajoute le FEN après chaque
// demi-coup : le navigateur n'a toujours aucune règle du jeu à connaître.
// La partie devient une chaîne de nœuds à un seul enfant : la navigation de
// l'arbre marche dessus sans rien changer.
function openGame(id){
  el("moves").innerHTML = "<p class='empty'>Chargement…</p>";
  api(API + "/game?" + query("&id=" + id)).then(function(g){
    game = g; openId = id;
    plies = g.plies.map(function(m){ return {san:m.san, uci:m.uci, fen:m.fen}; });
    for(var i = 0; i < plies.length - 1; i++) plies[i].children = [plies[i+1]];
    if(plies.length) plies[plies.length-1].children = [];
    root = {fen: START, children: plies.length ? [plies[0]] : []};
    path = [];
    // On oriente du côté du joueur cherché : préparer, c'est se mettre à sa place.
    flipped = flipFor(g);
    render();
  }).catch(function(e){
    el("moves").innerHTML = "<p class='empty'>" + esc(e.message || e) + "</p>";
  });
}
function closeGame(){ game = null; plies = []; openId = null; reload(); }
el("games").addEventListener("click", function(e){
  if(e.target.closest("#back-to-tree")){ closeGame(); return; }
  if(e.target.closest("a")) return; // un lien externe s'ouvre, sans ouvrir la partie
  var d = e.target.closest(".g[data-id]");
  if(d) openGame(+d.dataset.id);
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
  if(game){
    // Toute la partie d'un coup : sauter au vingtième demi-coup ne doit pas
    // demander vingt clics.
    el("moves-head").textContent = "Partie — " + plies.length + " demi-coups";
    var out = "";
    for(var i = 0; i < plies.length; i++){
      if(i % 2 === 0) out += "<span class='numlbl'>" + (i/2 + 1) + ".</span>";
      out += "<span class='p" + (i === path.length-1 ? " now" : "") +
             "' data-ply='" + (i+1) + "'>" + esc(plies[i].san) + "</span>";
    }
    el("moves").innerHTML = "<div class='plies'>" + (out || "—") + "</div>";
    return;
  }
  if(!ready()){
    el("moves").innerHTML = "<p class='empty'>" + esc(emptyText()) + "</p>";
    el("moves-head").textContent = "Coups";
    return;
  }
  var kids = (n || {}).children || [];
  if(!kids.length){
    el("moves").innerHTML = "<p class='empty'>" +
      ((n || {}).loading ? "Chargement…" : "Plus aucune partie ne va plus loin.") +
      "</p>";
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
function who(g){
  return esc(g.white) + (g.whiteElo ? " (" + g.whiteElo + ")" : "") + " &ndash; " +
         esc(g.black) + (g.blackElo ? " (" + g.blackElo + ")" : "");
}
function resultText(r){ return r > 0 ? "1-0" : (r < 0 ? "0-1" : "½-½"); }
function renderGames(){
  clearTimeout(gamesTimer);
  if(game){
    el("games-head").textContent = "Partie";
    el("games").innerHTML = "<div class='ghead'><div class='who'>" + who(game) +
      " &nbsp;" + resultText(game.result) + "</div><div class='meta'>" + gameHeadMeta(game) +
      "</div><button id='back-to-tree'>&larr; Retour à l'arbre</button></div>";
    return;
  }
  if(!ready()){ el("games").innerHTML = "<p class='empty'>—</p>"; return; }
  // La liste est le seul aller-retour restant : on la laisse respirer pendant
  // qu'on descend vite dans l'arbre.
  gamesTimer = setTimeout(function(){
    var p = path.map(function(n){ return n.uci; }).join(",");
    api(API + "/games?" + query("&path=" + p)).then(function(list){
      el("games-head").textContent = "Parties — " + list.length;
      if(!list.length){ el("games").innerHTML = "<p class='empty'>Aucune.</p>"; return; }
      el("games").innerHTML = list.map(function(g){
        return "<div class='g" + (g.id === openId ? " on" : "") +
          "' data-id='" + g.id + "'><div class='who'>" + who(g) +
          "<div class='meta'>" + gameRowMeta(g) + "</div>" +
          "</div><div class='res'>" + resultText(g.result) + "</div></div>";
      }).join("");
    }).catch(function(){ el("games").innerHTML = "<p class='empty'>—</p>"; });
  }, 220);
}
function render(){
  renderBoard(); renderMoves(); renderLine(); renderGames();
  el("b-back").disabled = path.length === 0;
  el("b-fwd").disabled = game
    ? path.length >= plies.length
    : !((node() || {}).children || []).length;
}
`
