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
- **Diagrams**: Multiple per round at narrative turning points, with witty captions
- **Lessons**: Short, concrete, honest — not generic advice
- **Author**: Alexandre Iwanesko, FM (FIDE Master), 33 years old

## Diagram Generation
- Run `node scripts/gen-diagrams.mjs` after modifying FEN positions
- Each diagram entry: `{ file, fen, lastMove?, flip? }`
- `flip: true` for games played as Black
- `lastMove` format: `'e2e4'` (from-to squares)

## Tournament Schedule (2026)
- Pontevedra: July 25–30, 2026 (completed, 4.5/9)
- Badalona: next tournament (directly after Pontevedra)
