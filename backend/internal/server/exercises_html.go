package server

import "html/template"

// Page des exercices : un « Puzzle Storm » positionnel sur les imprécisions
// d'Alexandre. Trois modes :
//   - Storm : trois minutes, une position après l'autre, une erreur coûte dix
//     secondes ; la correction vient à la fin, position par position ;
//   - Libre : sans chrono, la correction s'affiche après chaque coup ;
//   - Théoriques : relire les positions qu'il a marquées « choix conscient »,
//     et les remettre dans le circuit au besoin.
//
// On joue en cliquant la pièce puis la case. Les coups légaux viennent du
// serveur : le navigateur ne connaît aucune règle.
//
// Aucun accent grave dans ce fichier, et aucun attribut onclick= (CSP).
const exercisesHTML = `<!doctype html>
<html lang="fr"><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1,viewport-fit=cover">
<meta name="robots" content="noindex,nofollow">
<title>Mes imprécisions</title>
<style>
` + explorerCSS + `
main{grid-template-columns:minmax(0,560px) minmax(300px,1fr)}
@media(max-width:1240px){main{grid-template-columns:minmax(0,560px) minmax(300px,1fr)}}
@media(max-width:860px){main{grid-template-columns:minmax(0,1fr)}}
.left{max-width:560px}
#board{transition:box-shadow .15s}
#board.ok{box-shadow:0 0 0 5px var(--win)}
#board.ko{box-shadow:0 0 0 5px var(--loss)}
#board.unk{box-shadow:0 0 0 5px #5a6274}
.sq.sel{box-shadow:inset 0 0 0 4px var(--accent)}
.sq.tgt::after{content:"";position:absolute;width:28%;height:28%;border-radius:50%;background:#0005}
.sq.tgt.cap::after{width:88%;height:88%;background:transparent;box-shadow:inset 0 0 0 5px #0005}
.sq{cursor:pointer}
.big{font-size:34px;font-weight:700;font-variant-numeric:tabular-nums}
.row{display:flex;gap:16px;align-items:baseline;flex-wrap:wrap;padding:12px}
.cnt{font-variant-numeric:tabular-nums;color:var(--dim)}
.cnt b{color:var(--ink)}
.cnt .w{color:var(--win)} .cnt .l{color:var(--loss)}
.turn{padding:0 12px 10px;font-weight:600}
.info{padding:0 12px 12px;color:var(--dim);font-size:13px}
.actions{display:flex;gap:8px;flex-wrap:wrap;padding:0 12px 12px}
.corr{padding:10px 12px;border-top:1px solid var(--line)}
.corr h3{margin:0 0 6px;font-size:15px}
.corr .v-ok{color:var(--win)} .corr .v-ko{color:var(--loss)} .corr .v-unk{color:var(--dim)}
.cand{display:grid;grid-template-columns:4.2rem 4.4rem 1fr;gap:8px;padding:5px 0;
  border-bottom:1px solid #0003;font-size:14px;align-items:baseline}
.cand .ev{font-variant-numeric:tabular-nums;text-align:right}
.cand .ln{color:var(--dim);font-size:12.5px;overflow-wrap:anywhere}
.cand.good .mv2{color:var(--win);font-weight:600}
.cand.mine .mv2{color:var(--loss);font-weight:600}
.review .g{grid-template-columns:auto 1fr auto}
.tag{font-size:11px;border:1px solid var(--line);border-radius:999px;padding:1px 7px;color:var(--dim)}
.off{opacity:.55}
</style></head><body>

<header>
  <h1>Mes imprécisions</h1>
  <span class="stat" id="meta"></span>
  <span class="stat"><a href="/admin">&larr; tableau de bord</a></span>
</header>

<section class="search">
  <div class="field">
    <label>Mode</label>
    <div class="bar-row" style="padding:0" id="modes">
      <button data-mode="storm" class="on">Storm 3 min</button>
      <button data-mode="free">Libre</button>
      <button data-mode="theory">Théoriques</button>
    </div>
  </div>
  <div class="field">
    <label for="phase">Phase</label>
    <select id="phase">
      <option value="">Toutes</option><option value="opening">Ouverture</option>
      <option value="middlegame">Milieu de jeu</option><option value="endgame">Finale</option>
    </select>
  </div>
  <div class="field">
    <label for="speed">Cadence</label>
    <select id="speed"><option value="">Toutes</option><option value="blitz">Blitz</option><option value="rapid">Rapide</option></select>
  </div>
  <div class="field"><label>Filtre</label><button id="unsolved">Toutes les positions</button></div>
  <div class="field"><label>&nbsp;</label><button id="start" class="on">Commencer</button></div>
</section>

<main>
  <div class="left">
    <div id="board"></div>
    <div class="bar-row">
      <button id="b-flip">Retourner</button>
    </div>
  </div>
  <div>
    <div class="panel" id="play">
      <h2 id="play-head">Prêt</h2>
      <div class="row">
        <span class="big" id="clock"></span>
        <span class="cnt" id="counts"></span>
      </div>
      <div class="turn" id="turn"></div>
      <div class="info" id="info">{{if .Ready}}Choisis un mode, puis Commencer.{{else}}DB_PATH n'est pas configuré : les exercices sont indisponibles.{{end}}</div>
      <div class="actions">
        <button id="theory" disabled title="Choix conscient : ne plus revoir cette position">Théorique</button>
        <button id="skip" disabled>Passer</button>
        <button id="next" disabled>Suivante</button>
        <button id="stop" disabled>Arrêter</button>
      </div>
      <div id="corr"></div>
    </div>
    <div class="panel review" id="review-panel" hidden style="margin-top:14px">
      <h2 id="review-head">Revue</h2>
      <div class="body" id="review"></div>
    </div>
  </div>
</main>

<script nonce="{{ .Nonce }}">
var PIECES = {{ pieces .Pieces }};
var API = "/admin/exercices/api";
var STORM_SECONDS = 180, PENALTY = 10;
var PHASE = {opening:"ouverture", middlegame:"milieu de jeu", endgame:"finale"};
var SPEED = {bullet:"bullet", blitz:"blitz", rapid:"rapide", classical:"classique"};

var mode = "storm", unsolvedOnly = false;
var queue = [], cur = null, sel = null, flipped = false, busy = false;
var running = false, deadline = 0, tick = null, started = 0;
var counts = {correct:0, wrong:0, unknown:0};
var done = [];             // {puzzle, uci, san, verdict, ex} de la manche (pas « history » : window.history)

function el(id){ return document.getElementById(id); }
// Notation française, comme partout sur le site : seules les majuscules d'un
// SAN sont des pièces (les colonnes sont en minuscules, le roque en O).
var FR = {N:"C", B:"F", R:"T", Q:"D", K:"R"};
function fr(san){ return String(san || "").replace(/[NBRQK]/g, function(c){ return FR[c]; }); }
function esc(s){
  return String(s).replace(/[&<>"\x27]/g, function(c){
    return {"&":"&amp;","<":"&lt;",">":"&gt;","\"":"&quot;","\x27":"&#39;"}[c];
  });
}
function api(p, body){
  var o = {credentials:"same-origin"};
  if(body !== undefined){ o.method = "POST"; o.headers = {"Content-Type":"application/json"}; o.body = JSON.stringify(body); }
  return fetch(p, o).then(function(r){
    if(!r.ok) return r.json().then(function(j){ throw new Error(j.error || r.status); },
                                   function(){ throw new Error("erreur " + r.status); });
    return r.json();
  });
}

// ---- échiquier ---------------------------------------------------------------
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
function whiteToMove(fen){ return fen.split(" ")[1] === "w"; }
// mark : {from, to} à surligner (le coup de la correction).
function drawBoard(fen, mark){
  var grid = squares(fen), html = "", targets = {};
  if(sel && cur){
    cur.legal.forEach(function(m){ if(m.slice(0,2) === sel) targets[m.slice(2,4)] = true; });
  }
  for(var i = 0; i < 8; i++){
    for(var j = 0; j < 8; j++){
      var r = flipped ? 7 - i : i, f = flipped ? 7 - j : j;
      var sq = "abcdefgh"[f] + (8 - r), pc = grid[r][f];
      var cls = "sq " + ((r + f) % 2 ? "d" : "l");
      if(sq === sel) cls += " sel";
      if(targets[sq]) cls += " tgt" + (pc ? " cap" : "");
      if(mark && (mark.from === sq || mark.to === sq)) cls += " from";
      html += "<div class='" + cls + "' data-sq='" + sq + "'>" + (pc ? PIECES[pc] : "") +
              (j === 0 ? "<span class='coord'>" + (8 - r) + "</span>" : "") + "</div>";
    }
  }
  el("board").innerHTML = html;
}
function pieceAt(fen, sq){
  var g = squares(fen), f = sq.charCodeAt(0) - 97, r = 8 - +sq[1];
  return g[r][f];
}
el("board").addEventListener("click", function(e){
  var d = e.target.closest(".sq");
  if(!d || !cur || busy || cur.answered) return;
  var sq = d.dataset.sq, pc = pieceAt(cur.fen, sq);
  var mine = pc && ((pc === pc.toUpperCase()) === whiteToMove(cur.fen));
  if(sel){
    var cands = cur.legal.filter(function(m){ return m.slice(0,2) === sel && m.slice(2,4) === sq; });
    if(cands.length){
      // Promotion : la dame par défaut, c'est elle qu'on veut dans 99 % des cas.
      var mv = cands.filter(function(m){ return m.length === 4 || m[4] === "q"; })[0] || cands[0];
      sel = null; answer(mv); return;
    }
  }
  sel = mine && sq !== sel ? sq : null;
  drawBoard(cur.fen);
});
el("b-flip").addEventListener("click", function(){ flipped = !flipped; if(cur) drawBoard(cur.fen, cur.mark); });

// ---- réglages ------------------------------------------------------------------
el("modes").addEventListener("click", function(e){
  var b = e.target.closest("button[data-mode]");
  if(!b || running) return;
  mode = b.dataset.mode;
  [].forEach.call(el("modes").children, function(x){ x.className = x === b ? "on" : ""; });
});
el("unsolved").addEventListener("click", function(){
  unsolvedOnly = !unsolvedOnly;
  el("unsolved").className = unsolvedOnly ? "on" : "";
  el("unsolved").textContent = unsolvedOnly ? "Jamais réussies" : "Toutes les positions";
});

// ---- manche --------------------------------------------------------------------
el("start").addEventListener("click", start);
el("stop").addEventListener("click", function(){ finish(); });
el("skip").addEventListener("click", function(){ if(cur && !busy) nextPuzzle(); });
el("next").addEventListener("click", function(){ if(cur && cur.answered) nextPuzzle(); });
el("theory").addEventListener("click", function(){
  if(!cur) return;
  var want = !cur.theory;
  api(API + "/tag", {id: cur.id, theory: want}).then(function(){
    cur.theory = want;
    var h = done.filter(function(x){ return x.puzzle.id === cur.id; })[0];
    if(h) h.puzzle.theory = want;
    // Marquer une position en cours de manche la retire sans la compter.
    if(want && !cur.answered && mode !== "theory"){ nextPuzzle(); return; }
    renderButtons();
  }).catch(function(e){ say(e.message || e); });
});

function start(){
  var q = "?n=" + (mode === "storm" ? 80 : 60) +
          "&phase=" + el("phase").value + "&speed=" + el("speed").value +
          (unsolvedOnly ? "&unsolved=1" : "") + (mode === "theory" ? "&theory=1" : "");
  el("review-panel").hidden = true;
  say("Chargement…");
  api(API + "/batch" + q).then(function(list){
    queue = list; done = []; counts = {correct:0, wrong:0, unknown:0};
    if(!queue.length){ say(mode === "theory" ? "Aucune position marquée théorique." : "Aucune position ne correspond."); return; }
    running = true; started = Date.now();
    el("start").disabled = true; el("stop").disabled = false;
    if(mode === "storm"){
      deadline = Date.now() + STORM_SECONDS * 1000;
      tick = setInterval(clock, 200);
    }
    clock(); nextPuzzle();
  }).catch(function(e){ say(e.message || e); });
}
function clock(){
  if(mode !== "storm"){ el("clock").textContent = ""; renderCounts(); return; }
  var left = Math.max(0, Math.round((deadline - Date.now()) / 1000));
  el("clock").textContent = Math.floor(left / 60) + ":" + ("0" + left % 60).slice(-2);
  renderCounts();
  if(left <= 0) finish();
}
function renderCounts(){
  el("counts").innerHTML = "<b class='w'>" + counts.correct + "</b> justes &middot; <b class='l'>" +
    counts.wrong + "</b> fausses" + (counts.unknown ? " &middot; " + counts.unknown + " à évaluer" : "");
}
function nextPuzzle(){
  el("board").className = "";
  el("corr").innerHTML = "";
  cur = queue.shift() || null;
  sel = null;
  if(!cur){ finish(); return; }
  cur.answered = false; cur.mark = null;
  flipped = !whiteToMove(cur.fen);
  el("play-head").textContent = mode === "storm" ? "Storm" : (mode === "free" ? "Libre" : "Théoriques");
  el("turn").textContent = "Trait aux " + (whiteToMove(cur.fen) ? "Blancs" : "Noirs") + " — trouve mieux que la partie";
  el("info").textContent = (PHASE[cur.phase] || "") + " · " + (SPEED[cur.speed] || cur.speed) + " · " + cur.date;
  drawBoard(cur.fen);
  renderButtons();
}
function renderButtons(){
  el("theory").disabled = !cur;
  el("theory").textContent = cur && cur.theory ? "Remettre en circuit" : "Théorique";
  el("skip").disabled = !cur || cur.answered || mode === "theory";
  el("next").disabled = !cur || !cur.answered;
}
function answer(uci){
  busy = true;
  var p = cur;
  api(API + "/answer", {id: p.id, uci: uci}).then(function(res){
    busy = false;
    counts[res.verdict]++;
    done.push({puzzle: p, uci: uci, san: res.san, verdict: res.verdict, ex: res.exercise});
    p.answered = true;
    el("board").className = res.verdict === "correct" ? "ok" : (res.verdict === "wrong" ? "ko" : "unk");
    if(mode === "storm"){
      if(res.verdict === "wrong") deadline -= PENALTY * 1000;
      clock();
      setTimeout(function(){ if(running && cur === p) nextPuzzle(); }, 350);
      return;
    }
    showCorrection(res.exercise, uci, res.verdict);
    renderButtons(); renderCounts();
  }).catch(function(e){ busy = false; say(e.message || e); });
}
function finish(){
  if(!running) return;
  running = false;
  clearInterval(tick);
  el("start").disabled = false; el("stop").disabled = true;
  var seconds = Math.round((Date.now() - started) / 1000);
  var total = counts.correct + counts.wrong + counts.unknown;
  if(mode !== "theory") api(API + "/run", {mode: mode, total: total, correct: counts.correct,
    wrong: counts.wrong, unknown: counts.unknown, seconds: seconds}).catch(function(){});
  cur = null; renderButtons();
  el("play-head").textContent = "Manche terminée";
  var judged = counts.correct + counts.wrong;
  el("turn").textContent = judged ? Math.round(100 * counts.correct / judged) + " % de précision" :
    (total ? "Aucun coup encore tranché par le moteur" : "");
  el("info").textContent = total + " positions en " + seconds + " s";
  el("corr").innerHTML = "";
  renderReview();
  loadStats();
}

// ---- correction --------------------------------------------------------------
function ev(c){
  if(c.mate) return "#" + c.mate;
  var v = c.cp / 100;
  return (v > 0 ? "+" : "") + v.toFixed(2);
}
function gameLink(ex){
  var u = ex.gameUrl;
  if(ex.source === "lichess") u += "#" + (ex.ply + 1);
  return "<a href='" + esc(u) + "' target='_blank' rel='noopener noreferrer'>voir la partie</a>";
}
function showCorrection(ex, uci, verdict){
  var best = ex.moves.filter(function(m){ return m.ok; })[0] || ex.moves[0];
  if(cur && cur.id === ex.id){
    cur.mark = best ? {from: best.uci.slice(0,2), to: best.uci.slice(2,4)} : null;
    drawBoard(ex.fen, cur.mark);
  }
  var label = {correct:"Juste", wrong:"Faux", unknown:"Coup jamais évalué — le moteur tranchera"}[verdict] || "";
  var cls = {correct:"v-ok", wrong:"v-ko", unknown:"v-unk"}[verdict] || "";
  var mine = done.filter(function(h){ return h.puzzle.id === ex.id; }).pop();
  var html = "<div class='corr'>" + (verdict ? "<h3 class='" + cls + "'>" + label +
    (mine ? " — tu as joué " + esc(fr(mine.san)) : "") + "</h3>" : "");
  html += ex.moves.map(function(m){
    var isGame = m.uci === ex.playedUci;
    return "<div class='cand" + (m.ok ? " good" : "") + (isGame ? " mine" : "") + "'>" +
      "<span class='mv2'>" + esc(m.san ? fr(m.san) : m.uci) + "</span><span class='ev'>" + ev(m) + "</span>" +
      "<span class='ln'>" + (isGame ? "joué en partie, −" + (ex.loss / 100).toFixed(2) + " · " : "") +
      esc((m.line || []).map(fr).join(" ")) + "</span></div>";
  }).join("");
  html += "<div class='info' style='padding:8px 0 0'>Profondeur " + ex.depth + " · " +
    esc(ex.white) + " – " + esc(ex.black) + " · " + gameLink(ex) +
    (ex.attempts ? " · réussie " + ex.solved + "/" + ex.attempts : "") + "</div></div>";
  el("corr").innerHTML = html;
}

// ---- revue de fin de manche -----------------------------------------------------
function renderReview(){
  if(!done.length){ el("review-panel").hidden = true; return; }
  el("review-panel").hidden = false;
  el("review-head").textContent = "Revue — " + done.length + " positions";
  el("review").innerHTML = done.map(function(h, i){
    var mark = h.verdict === "correct" ? "<span class='v-ok'>&#10003;</span>" :
               (h.verdict === "wrong" ? "<span class='v-ko'>&#10007;</span>" : "?");
    return "<div class='g" + (h.puzzle.theory ? " off" : "") + "' data-i='" + i + "'><span>" + mark + "</span>" +
      "<div class='who'>" + esc(fr(h.san)) + " <span class='meta'>" + (PHASE[h.puzzle.phase] || "") +
      " · " + esc(h.puzzle.date) + "</span></div>" +
      (h.puzzle.theory ? "<span class='tag'>théorique</span>" : "<span></span>") + "</div>";
  }).join("");
}
el("review").addEventListener("click", function(e){
  var d = e.target.closest(".g[data-i]");
  if(!d) return;
  var h = done[+d.dataset.i];
  cur = h.puzzle; cur.answered = true; sel = null;
  flipped = !whiteToMove(cur.fen);
  el("turn").textContent = "Trait aux " + (whiteToMove(cur.fen) ? "Blancs" : "Noirs");
  el("info").textContent = (PHASE[cur.phase] || "") + " · " + (SPEED[cur.speed] || cur.speed) + " · " + cur.date;
  showCorrection(h.ex, h.uci, h.verdict);
  renderButtons();
  el("next").disabled = true;
});

// ---- divers ----------------------------------------------------------------------
function say(msg){ el("info").textContent = msg; }
function loadStats(){
  api(API + "/stats").then(function(s){
    var last = s.lastRuns.filter(function(r){ return r.mode === "storm" && r.correct + r.wrong > 0; })[0];
    el("meta").textContent = (s.total - s.theory) + " positions à travailler · " + s.theory + " théoriques · " +
      s.games + " parties analysées" + (s.pending ? " · " + s.pending + " coups à évaluer" : "") +
      (last ? " · dernier Storm : " + last.correct + "/" + (last.correct + last.wrong) : "");
  }).catch(function(){});
}
{{if .Ready}}loadStats();{{end}}
drawBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1");
</script>
</body></html>`

var exercisesTmpl = template.Must(template.New("exercices").Funcs(template.FuncMap{
	"pieces": pieceJSON,
}).Parse(exercisesHTML))
