# Chess Web Site — Claude Code Memory

## Project Structure
- **Framework**: React (Vite/SSR)
- **Blog articles FR**: `content/blog/*.md`
- **Blog articles EN**: `content/blog/en/*.md`
- **Diagrams script**: `scripts/gen-diagrams.mjs` — generates SVG chess diagrams from FEN using Lichess cburnett piece set
- **Diagram output**: `public/images/blog/<tournament>/`
- **i18n**: FR (default) + EN, routes defined in `frontend/src/lib/i18n.tsx`
- **Blog slugs**: `frontend/src/lib/blogSlugs.ts` — imports from `content/blog/*.md` (FR) and `content/blog/en/*.md` (EN)
- **Content processing**: `frontend/src/lib/content.ts`
- **Routes**: `frontend/src/routes.tsx` — FR at `/blog/:slug`, EN at `/en/blog/:slug`
- **altSlug**: frontmatter field to link FR ↔ EN versions (enables hreflang)

## Thème clair / sombre
- **Interrupteur** : bouton lune/soleil dans l'en-tête (`frontend/src/components/ThemeToggle.tsx`). Le choix est mémorisé dans `localStorage`; sans choix, on suit `prefers-color-scheme`.
- **L'état vit sur `<html data-theme>`**, posé par le script inline de `frontend/index.html` **avant le premier rendu** — indispensable sur un site pré-rendu, sinon la page clignote en blanc avant de basculer.
- **Comment ça bascule** : aucune classe `dark:` dans les composants. Les utilitaires Tailwind v4 compilent vers `var(--color-…)`, donc il suffit de **redéfinir les tokens** sous `html[data-theme='dark']` dans `styles.css`. L'échelle `ink` et l'échelle `cream` sont inversées ; les ors 200→500 et la famille `slab` ne bougent pas.
- **Trois tokens de rôle à connaître** — la palette ne peut plus servir deux rôles à la fois :
  - `paper` = ce qui était codé en dur `bg-white` (cartes, en-tête, sections). **Ne jamais réintroduire `bg-white`** pour une surface.
  - `slab-*` = les panneaux qui restent **sombres dans les deux thèmes** (pied de page, bandeaux CTA) et le texte posé dessus (`text-slab-300/400`). Un `text-ink-*` sur un slab devient illisible en mode nuit.
  - `on-gold` = le texte d'un bouton doré, qui doit rester quasi noir en permanence.
- **Le bouton n'a aucun état React** : les deux icônes sont dans le DOM et le CSS choisit laquelle sort. C'est ce qui évite le saut d'hydratation.
- **Attention en test** : Chrome headless annonce `prefers-color-scheme: dark` par défaut — une capture « claire » sort sombre si on ne force pas le thème.

## Tournament Calendar
- **Tournament list (edit this)**: `frontend/src/lib/tournaments.ts` — the single source of truth, edited by hand when publishing. Baked in at build, so every tournament ships inside the pre-rendered HTML (no runtime fetch, no backend involved).
- **Page**: `/calendrier` (FR) + `/en/calendar` (EN) — `frontend/src/pages/Calendrier.tsx`
- **Grid component**: `frontend/src/components/TournamentCalendar.tsx` — month grid where a multi-day tournament renders as one continuous bar, plus « previous / next » cards and a full-season list below
- **Grid/date logic**: `frontend/src/lib/calendar.ts` (`buildMonth`, `formatRange`, `statusOf`)
- **Linking a diary**: set `slugFr` / `slugEn` to the article file name (no `/blog/` prefix). A past tournament with a slug becomes clickable straight to its diary.
- **Adding a tournament**: one entry in `TOURNAMENTS`; order doesn't matter, the grid sorts itself. `start === end` for a one-day event.
- **`__BUILD_DATE__`**: injected by `define` in `vite.config.ts`. The calendar's first render uses it so SSG output and hydration match; the visitor's real date takes over in an effect.

## Student Results (`/resultats`)
- Presents **coached cases** (what was worked on → what it produced), described by Alexandre — **not** reviews written by students.
- **No star ratings, no `aggregateRatingSchema`** on this page. Declaring ratings nobody gave is fabricated structured data: manual-action risk with Google, and misleading advertising under Swiss LCD art. 3 / the EU Omnibus directive.
- The builder in `lib/schema.ts` is kept for the day a student sends a **real written review** — add the quote, then wire the schema back.
- Every figure on the page must be traceable to something real (Flavien's +100 FIDE Elo, the tournament count read from `tournaments.ts`).

## Blog Writing Style (Tournament Diaries)
- **VOICE — read first**: Before writing OR rewriting any tournament diary (`content/blog/**`, FR and EN), read [`docs/voix-carnet-tournoi.md`](docs/voix-carnet-tournoi.md) and apply it on the FIRST draft (don't write "clean" then fix later). It is the canonical tone reference: bans reporting connectors ("À partir de là", "Le problème c'est que", "Vient alors le moment-clé", "Nouveau carrefour", "Résultat :"…), requires showing emotions via concrete detail rather than declaring them, favours self-deprecation, reader address and lived imagery. Gold standard: `content/blog/open-pontevedra-2026.md`. Does NOT apply to SEO/informational articles (see `docs/plan-editorial-blog.md`).
- **Tone**: Self-deprecating humor, storytelling, honest about mistakes
- **Framework per round**:
  1. How you felt before the game (fatigue, motivation, mental state)
  2. How the opening went
  3. The moment you felt the advantage (or the pressure)
  4. What was going through your head at key moments
  5. Why this game mattered (or didn't)
  6. Self-deprecation — never take yourself too seriously
  7. Funny expressions when they fit naturally
- **Chess notation**: French style (Fou=F, Cavalier=C, Tour=T, Dame=D, Roi=R) for FR articles, English style (B, N, R, Q, K) for EN articles
- **Played moves vs variations**: a move actually played is written `` `20.Tc2` ``; a move only calculated is written `` *`23.g4`* `` — it renders as `<em><code>` and is greyed with no background (rule `.prose em code` in `frontend/src/styles.css`). Never let a reader mistake a variation for the game.
- **Engine verdicts**: when a diary separates the player's own calculations from the engine's, the engine goes in a grey callout labelled « L'ordinateur, après coup » (`<div class="engine-note">` + `<p class="engine-note-label">` ; see `championnat-suisse-equipes-2026.md`), never inline in the prose.
- **Diagrams**: Multiple per round at narrative turning points, with witty captions. Markup: `<div class="diagram-container">` + `<img …>` + `<p class="diagram-caption">`. **Never inline styles** — a style attribute cannot follow the light/dark theme, and these blocks stayed white rectangles in night mode until they were moved to classes. Same for the summary table: rows use `class="row-head"` / `row-a` / `row-b`, cells `cell-muted` / `cell-warn` / `cell-alt`.
- **Lessons**: Short, concrete, honest — not generic advice
- **Author**: Alexandre Iwanesko, FM (FIDE Master), 33 years old

## Fréquentation (/admin)
- **Les humains sont comptés par une balise JS**, pas par le serveur : `frontend/src/lib/analytics.ts` envoie `POST /api/hit` à chaque changement de route, `handleHit` (`backend/internal/server/analytics.go`) l'enregistre.
- **Pourquoi** : le comptage serveur enregistrait toute réponse HTML, donc tous les crawlers déguisés en navigateur. Le tableau de bord affichait des milliers de « visites humaines » venues de datacenters US/DE/NL/SG, pour 3 clics réels dans Search Console. Les robots n'exécutent presque jamais de JS.
- Le middleware serveur ne compte plus que les **robots qui s'annoncent** (regex user-agent), ce qui garde la colonne « bots » lisible.
- Garde-fous sur `/api/hit` : même origine exigée (`Origin`, sinon `Referer`), user-agent de robot rejeté, chemin validé par `cleanHitPath` (interne, pas `/admin` ni `/newsletter`). Testé dans `server_test.go`.
- Rien ne change côté vie privée : aucune IP stockée, aucun cookie, pays déduit hors ligne, empreinte de visiteur tournante à la journée.

## Tableaux dans les articles
- **Écrire les tableaux en Markdown**, pas en HTML. `enhanceTables()` (`frontend/src/lib/content.ts`) les enveloppe dans un `.table-wrap` et colore les cellules qui ne contiennent qu'un nombre signé (`+3,6` en vert, `−59,0` en rouge).
- Le style vit dans `styles.css`, **scopé sous `.table-wrap`** : en-tête presque noir, lignes alternées, chiffres tabulaires, dernière colonne alignée à droite. Rien de tout ça ne touche les tableaux HTML écrits à la main dans les carnets, qui gardent leur mise en forme.
- Au-delà de 768 px, un tableau **déborde volontairement de la colonne de texte** (68 caractères) et se recentre sur la page : six colonnes n'y entrent pas, et c'est toujours la dernière — la plus utile — qui se retrouvait hors champ.
- **Pas de drapeaux emoji** : Windows ne les rend pas et les affiche en paires de lettres (« CH », « FR »).

## PGN Game Viewer
- **Games live in**: `content/games/<id>.pgn` — one file per game, standard PGN with headers. Paste the score, done: no build step, no diagram generation.
- **Embed in an article**: put `[[pgn:<id>]]` alone on its own line in the Markdown (FR or EN). `splitChunks()` in `frontend/src/lib/content.ts` splits the compiled HTML on that marker and `BlogPost.tsx` renders a `<PgnViewer>` in its place.
- **Several games, one board**: `[[pgn:<id-a>,<id-b>]]` renders a single viewer with a picker (tab per game, labelled `R<ronde> · <adversaire>`). Preferred for a team-match diary: one board at the foot of the article rather than one per round.
- **Anchor**: the first `[[pgn:…]]` block of a post gets `id="parties"`, so the intro can link down to it with `[texte](#parties)`.
- **Component**: `frontend/src/components/PgnViewer.tsx` — board, first/prev/play/next/last, **flip button**, clickable move list, arrow-key navigation. Read-only.
- **Rules engine**: `frontend/src/lib/chess.ts` (legal-move generation, only used to resolve SAN into from/to squares) + `frontend/src/lib/pgn.ts` (PGN → plies, main line only).
- **Main line only**: comments `{...}`, sidelines `(...)` and NAGs are stripped. A sideline worth telling goes in the prose, not in the viewer.
- **Notation**: the viewer renders French notation (C, F, T, D, R) on FR pages and English on EN pages automatically — always paste the PGN in **English** SAN.
- **Board orientation**: defaults to Alexandre's side (detected from the `White`/`Black` headers), overridable with a `[Orientation "black"]` header. The reader can flip it anyway.
- **Piece set**: `frontend/src/components/chessPieces.tsx` — shared with `PuzzleBoard.tsx` (Lichess cburnett, inlined SVG).

## Cartes de partage (Open Graph) et JSON-LD
- **`scripts/gen-og.mjs`** produit les PNG 1200×630. Deux gabarits : `card()` pour un **carnet** (un vrai diagramme d'échiquier imbriqué depuis `public/images/blog/<dir>/`, donc la position montrée est celle de l'article) et `chartCard()` pour un **article de fond** (la courbe Elo, **lue directement dans le Markdown** — carte et graphique ne peuvent pas diverger).
- Le SVG reste à côté du PNG : c'est la source éditable. **C'est le PNG qui va dans le front-matter** — les réseaux sociaux n'affichent pas d'aperçu SVG. Carnets : `/images/blog/<dir>/og-*.png`. Articles de fond : `/public/og/<slug>.png`.
- Sans `image:` en front-matter, la page retombe sur `og/default.png` — silencieusement. `check-content-links.mjs` vérifie qu'une `image:` déclarée existe, **pour tout article** et plus seulement pour les carnets.
- **JSON-LD** (`frontend/src/lib/schema.ts`, `articleSchema`) : une page d'article ne porte aucun nœud `Person` ni `Organization` (ils ne sont émis que sur l'accueil et `/a-propos`). L'`@id` seul y serait une **référence pendante** — le nom voyage donc avec. `inLanguage` suit la locale de la page (un article EN se déclarait `fr`), et `image` a toujours une valeur.

## Diagram Generation
- Run `node scripts/gen-diagrams.mjs` after modifying FEN positions
- Each diagram entry: `{ file, fen, lastMove?, flip?, dir? }`
- `dir` picks the tournament folder under `public/images/blog/` (default `pontevedra-2026`); CSE 2026 uses `dir: 'cse-2026'`
- `flip: true` for games played as Black
- `lastMove` format: `'e2e4'` (from-to squares)

## Tournament Schedule (2026)
- Pontevedra: July 25–30, 2026 (completed, 4.5/9)
- Badalona: next tournament (directly after Pontevedra)
