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
` + explorerCSS + `</style></head><body>

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
var API = "/admin/parties/api";

var players = [];          // joueurs sélectionnés
var colour = "";           // "", "w" ou "b"
var noTT = true;

// ---- ce que la navigation commune attend de la page ------------------------
function ready(){ return players.length > 0; }
function emptyText(){ return "Choisis au moins un joueur."; }
function query(extra){
  var p = "players=" + players.map(function(x){ return x.id; }).join(",") +
          "&colour=" + colour + "&noTT=" + (noTT ? "1" : "0") +
          "&from=" + encodeURIComponent(el("from").value) +
          "&to=" + encodeURIComponent(el("to").value);
  return p + (extra || "");
}
function flipFor(g){ return players.some(function(pl){ return pl.name === g.black; }); }
function gameHeadMeta(g){
  return esc(g.event) + (g.date ? " &middot; " + esc(g.date) : "") +
         (g.eco ? " &middot; " + esc(g.eco) : "");
}
function gameRowMeta(g){ return esc(g.event) + (g.year ? " &middot; " + g.year : ""); }

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

` + explorerJS + `
api(API + "/meta").then(function(m){
  el("meta").textContent = (m.games || "?") + " parties" +
    (m.twic ? " · TWIC " + m.twic : "") + (m.built ? " · " + m.built : "");
}).catch(function(){});
renderChips();
render();
</script>
</body></html>`
