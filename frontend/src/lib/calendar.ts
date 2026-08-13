/**
 * Tournament calendar helpers. The entries themselves live in tournaments.ts
 * and are baked into the build, so they ship inside the pre-rendered HTML.
 *
 * Every date is a plain YYYY-MM-DD day with no time zone. All helpers here work
 * on day numbers (days since the epoch) rather than millisecond arithmetic, so
 * a DST boundary can never shift an event by one square on the grid.
 */
import type { Locale } from './i18n'

export interface CalEvent {
  id: string
  start: string // YYYY-MM-DD
  end: string // YYYY-MM-DD, equal to start for a one-day event
  name: string
  nameEn?: string // falls back to `name` when the event has no English name
  location?: string
  locationEn?: string
  format?: string // "9 rondes · 90 min + 30 s"
  formatEn?: string
  noteFr?: string // what he's going there to do, shown while it's still ahead
  noteEn?: string
  slugFr?: string // blog article recounting it, FR
  slugEn?: string // blog article recounting it, EN
  result?: string // "4.5/9"
}

export type EventStatus = 'past' | 'ongoing' | 'upcoming'

/** Days elapsed since 1970-01-01 for a YYYY-MM-DD string. NaN if malformed. */
export function dayNumber(iso: string): number {
  const m = /^(\d{4})-(\d{2})-(\d{2})$/.exec(iso)
  if (!m) return NaN
  return Math.floor(Date.UTC(+m[1]!, +m[2]! - 1, +m[3]!) / 86_400_000)
}

/** The YYYY-MM-DD string for a day number. */
export function isoFromDayNumber(n: number): string {
  const d = new Date(n * 86_400_000)
  const p = (x: number) => String(x).padStart(2, '0')
  return `${d.getUTCFullYear()}-${p(d.getUTCMonth() + 1)}-${p(d.getUTCDate())}`
}

/** Today in the visitor's own time zone, as YYYY-MM-DD. */
export function todayISO(): string {
  const d = new Date()
  const p = (x: number) => String(x).padStart(2, '0')
  return `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())}`
}

/** ISO date strings sort chronologically, so plain comparison is enough. */
export function statusOf(e: CalEvent, today: string): EventStatus {
  if (e.end < today) return 'past'
  if (e.start > today) return 'upcoming'
  return 'ongoing'
}

/** Localised fields. English falls back to the French value when untranslated. */
export const nameFor = (e: CalEvent, locale: Locale): string =>
  (locale === 'en' ? e.nameEn : e.name) || e.name

export const locationFor = (e: CalEvent, locale: Locale): string | undefined =>
  (locale === 'en' ? e.locationEn ?? e.location : e.location) || undefined

export const formatFor = (e: CalEvent, locale: Locale): string | undefined =>
  (locale === 'en' ? e.formatEn ?? e.format : e.format) || undefined

export const noteFor = (e: CalEvent, locale: Locale): string | undefined =>
  (locale === 'en' ? e.noteEn : e.noteFr) || undefined

export const slugFor = (e: CalEvent, locale: Locale): string | undefined =>
  (locale === 'en' ? e.slugEn : e.slugFr) || undefined

/** Path of the blog article recounting this tournament, if it has one. */
export function articlePathFor(e: CalEvent, locale: Locale): string | undefined {
  const slug = slugFor(e, locale)
  if (!slug) return undefined
  return locale === 'en' ? `/en/blog/${slug}` : `/blog/${slug}`
}

/**
 * The most recent event that is fully over. Ties on the same end date are broken
 * by start date, so a one-day blitz never hides the nine-day open it sits inside.
 */
export function previousEvent(events: CalEvent[], today: string): CalEvent | undefined {
  return events
    .filter((e) => e.end < today)
    .sort((a, b) => (a.end !== b.end ? (a.end < b.end ? 1 : -1) : a.start < b.start ? 1 : -1))[0]
}

/** The next event to come — an ongoing one counts as "now", handled separately. */
export function nextEvent(events: CalEvent[], today: string): CalEvent | undefined {
  return events
    .filter((e) => e.start > today)
    .sort((a, b) => (a.start !== b.start ? (a.start < b.start ? -1 : 1) : a.end < b.end ? -1 : 1))[0]
}

/** The event being played right now, if any. */
export function ongoingEvent(events: CalEvent[], today: string): CalEvent | undefined {
  return events.find((e) => e.start <= today && e.end >= today)
}

// --- Formatting -------------------------------------------------------------

const MONTHS: Record<Locale, string[]> = {
  fr: ['janvier', 'février', 'mars', 'avril', 'mai', 'juin', 'juillet', 'août', 'septembre', 'octobre', 'novembre', 'décembre'],
  en: ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'],
}

const MONTHS_SHORT: Record<Locale, string[]> = {
  fr: ['janv.', 'févr.', 'mars', 'avril', 'mai', 'juin', 'juil.', 'août', 'sept.', 'oct.', 'nov.', 'déc.'],
  en: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'],
}

/** Weekday initials, Monday first (the European week the grid is built on). */
export const WEEKDAYS: Record<Locale, string[]> = {
  fr: ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'],
  en: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
}

export const monthLabel = (year: number, month: number, locale: Locale): string =>
  `${MONTHS[locale][month]} ${year}`

/** "2–10 août 2026" / "2–10 August 2026", collapsing what both dates share. */
export function formatRange(start: string, end: string, locale: Locale): string {
  const a = /^(\d{4})-(\d{2})-(\d{2})$/.exec(start)
  const b = /^(\d{4})-(\d{2})-(\d{2})$/.exec(end)
  if (!a || !b) return `${start} → ${end}`
  const [ya, ma, da] = [+a[1]!, +a[2]! - 1, +a[3]!]
  const [yb, mb, db] = [+b[1]!, +b[2]! - 1, +b[3]!]
  const mon = MONTHS_SHORT[locale]
  if (start === end) return `${da} ${mon[ma]} ${ya}`
  if (ya === yb && ma === mb) return `${da}–${db} ${mon[mb]} ${yb}`
  if (ya === yb) return `${da} ${mon[ma]} – ${db} ${mon[mb]} ${yb}`
  return `${da} ${mon[ma]} ${ya} – ${db} ${mon[mb]} ${yb}`
}

/** Whole days from today until the event starts (0 = starts today). */
export const daysUntil = (e: CalEvent, today: string): number =>
  dayNumber(e.start) - dayNumber(today)

// --- Month grid -------------------------------------------------------------

export interface DayCell {
  iso: string
  dayOfMonth: number
  inMonth: boolean
  isToday: boolean
}

/** One event as it appears inside a single week row. */
export interface Segment {
  event: CalEvent
  col: number // 0-based column where the bar starts in this week
  span: number // number of columns it covers
  lane: number // vertical slot, so overlapping tournaments never sit on top of each other
  continuesLeft: boolean // started in an earlier week
  continuesRight: boolean // runs into the next week
}

export interface Week {
  days: DayCell[]
  segments: Segment[]
  lanes: number
}

/**
 * Builds the Monday-first grid for one month, placing every event that touches
 * it. Weeks are only emitted while they still contain a day of the month, so a
 * short month never renders a trailing empty row.
 */
export function buildMonth(year: number, month: number, events: CalEvent[], today: string): Week[] {
  const firstDayNum = Math.floor(Date.UTC(year, month, 1) / 86_400_000)
  const daysInMonth = new Date(Date.UTC(year, month + 1, 0)).getUTCDate()
  // getUTCDay(): 0 = Sunday. Shift so Monday is 0.
  const firstWeekday = (new Date(firstDayNum * 86_400_000).getUTCDay() + 6) % 7
  const gridStart = firstDayNum - firstWeekday
  const weekCount = Math.ceil((firstWeekday + daysInMonth) / 7)

  const weeks: Week[] = []
  for (let w = 0; w < weekCount; w++) {
    const weekStart = gridStart + w * 7
    const weekEnd = weekStart + 6

    const days: DayCell[] = []
    for (let i = 0; i < 7; i++) {
      const n = weekStart + i
      const iso = isoFromDayNumber(n)
      days.push({
        iso,
        dayOfMonth: new Date(n * 86_400_000).getUTCDate(),
        inMonth: n >= firstDayNum && n < firstDayNum + daysInMonth,
        isToday: iso === today,
      })
    }

    // Longest tournaments first so the big bars take the top lanes and the
    // one-day entries tuck in underneath — much easier to read.
    const touching = events
      .filter((e) => dayNumber(e.start) <= weekEnd && dayNumber(e.end) >= weekStart)
      .sort((a, b) => {
        const byStart = dayNumber(a.start) - dayNumber(b.start)
        if (byStart !== 0) return byStart
        return dayNumber(b.end) - dayNumber(a.end)
      })

    const laneEnds: number[] = [] // last column occupied in each lane
    const segments: Segment[] = []
    for (const e of touching) {
      const s = dayNumber(e.start)
      const en = dayNumber(e.end)
      const col = Math.max(0, s - weekStart)
      const span = Math.min(6, en - weekStart) - col + 1
      let lane = laneEnds.findIndex((end) => end < col)
      if (lane === -1) lane = laneEnds.length
      laneEnds[lane] = col + span - 1
      segments.push({
        event: e,
        col,
        span,
        lane,
        continuesLeft: s < weekStart,
        continuesRight: en > weekEnd,
      })
    }

    weeks.push({ days, segments, lanes: laneEnds.length })
  }
  return weeks
}
