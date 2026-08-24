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
- **Engine verdicts**: when a diary separates the player's own calculations from the engine's, the engine goes in a grey callout labelled « L'ordinateur, après coup » (inline-styled div; see `championnat-suisse-equipes-2026.md`), never inline in the prose.
- **Diagrams**: Multiple per round at narrative turning points, with witty captions
- **Lessons**: Short, concrete, honest — not generic advice
- **Author**: Alexandre Iwanesko, FM (FIDE Master), 33 years old

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

## Diagram Generation
- Run `node scripts/gen-diagrams.mjs` after modifying FEN positions
- Each diagram entry: `{ file, fen, lastMove?, flip?, dir? }`
- `dir` picks the tournament folder under `public/images/blog/` (default `pontevedra-2026`); CSE 2026 uses `dir: 'cse-2026'`
- `flip: true` for games played as Black
- `lastMove` format: `'e2e4'` (from-to squares)

## Tournament Schedule (2026)
- Pontevedra: July 25–30, 2026 (completed, 4.5/9)
- Badalona: next tournament (directly after Pontevedra)
