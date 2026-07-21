/**
 * Fire-and-forget interaction tracking for the weekly tactics puzzles.
 * Posts to the Go backend (/api/tactics/event), which records aggregate counts
 * in a private SQLite store — no cookies, no identifiers, purely anonymous.
 * Every call is best-effort: failures and SSG/prerender are silently ignored.
 */
export type TacticsEventKind = 'view' | 'attempt' | 'solved'

// De-dupe "view" per puzzle within a page load so scroll/re-render don't inflate.
const seenViews = new Set<string>()

export function recordTacticsEvent(week: string, puzzleId: string, kind: TacticsEventKind): void {
  if (typeof window === 'undefined' || typeof fetch === 'undefined') return
  if (!week || !puzzleId) return

  if (kind === 'view') {
    const key = week + '/' + puzzleId
    if (seenViews.has(key)) return
    seenViews.add(key)
  }

  const body = JSON.stringify({ week, puzzleId, kind })
  try {
    fetch('/api/tactics/event', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body,
      keepalive: true,
    }).catch(() => {})
  } catch {
    /* ignore */
  }
}
