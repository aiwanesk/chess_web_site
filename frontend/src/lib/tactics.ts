/**
 * Weekly tactics content. Each week is a file content/tactiques/JJ-MM-AA.json
 * (produced by the Go batch, already anonymised — mirrored FEN + solution).
 * Loaded at build so each week's article (/tactiques/JJ-MM-AA) is pre-rendered.
 */
export interface Puzzle {
  id: string
  fen: string
  sideToMove: 'w' | 'b'
  solution: string[]
  mate?: boolean
  sacrifice?: boolean
}

export interface TacticsWeek {
  slug: string // "20-07-26"
  puzzles: Puzzle[]
}

const files = import.meta.glob('../../../content/tactiques/*.json', { eager: true }) as Record<
  string,
  { default: { week?: string; puzzles?: Puzzle[] } }
>

// Chronological sort key from a JJ-MM-AA slug → "20AA-MM-JJ".
function sortKey(slug: string): string {
  const [d, m, y] = slug.split('-')
  return `20${y}-${m}-${d}`
}

export const weeks: TacticsWeek[] = Object.entries(files)
  .map(([path, mod]) => ({
    slug: path.split('/').pop()!.replace(/\.json$/, ''),
    puzzles: mod.default.puzzles ?? [],
  }))
  .sort((a, b) => (sortKey(a.slug) < sortKey(b.slug) ? 1 : -1)) // newest first

export const getWeek = (slug: string): TacticsWeek | undefined => weeks.find((w) => w.slug === slug)

export const weekStaticPaths = (): string[] => weeks.map((w) => `/tactiques/${w.slug}`)

/** "20-07-26" → a Date, or null. */
export function weekDate(slug: string): Date | null {
  const m = /^(\d{2})-(\d{2})-(\d{2})$/.exec(slug)
  if (!m) return null
  return new Date(Number(`20${m[3]}`), Number(m[2]) - 1, Number(m[1]))
}

/** Localised long date, e.g. "20 juillet 2026". */
export function formatWeek(slug: string, locale: string): string {
  const d = weekDate(slug)
  if (!d) return slug
  return d.toLocaleDateString(locale === 'en' ? 'en-GB' : 'fr-CH', { day: 'numeric', month: 'long', year: 'numeric' })
}
