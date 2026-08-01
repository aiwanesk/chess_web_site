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

## Blog Writing Style (Tournament Diaries)
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
