import { useEffect, useMemo, useState } from 'react'
import { Link } from 'vite-react-ssg'
import { IconArrowRight, IconTrophy } from './icons'
import { useLocale, type Locale } from '../lib/i18n'
import { TOURNAMENTS } from '../lib/tournaments'
import {
  articlePathFor,
  buildMonth,
  daysUntil,
  formatFor,
  formatRange,
  locationFor,
  monthLabel,
  nameFor,
  nextEvent,
  noteFor,
  ongoingEvent,
  previousEvent,
  statusOf,
  todayISO,
  WEEKDAYS,
  type CalEvent,
  type EventStatus,
} from '../lib/calendar'

const STR: Record<Locale, {
  empty: string
  prevMonth: string
  nextMonth: string
  today: string
  jumpPrev: string
  jumpNext: string
  lastPlayed: string
  playingNow: string
  upNext: string
  readDiary: string
  diaryComing: string
  noDiaryYet: string
  result: string
  inDays: (n: number) => string
  legendUpcoming: string
  legendOngoing: string
  legendPast: string
  legendPastNoArticle: string
  selectedTitle: string
  close: string
  goingThere: string
  nothingPlanned: string
  seasonTitle: string
  seasonUpcoming: string
  seasonPast: string
  showOnCalendar: string
}> = {
  fr: {
    empty: 'Aucun tournoi au calendrier pour le moment.',
    prevMonth: 'Mois précédent',
    nextMonth: 'Mois suivant',
    today: 'Aujourd’hui',
    jumpPrev: 'Tournoi précédent',
    jumpNext: 'Prochain tournoi',
    lastPlayed: 'Dernier tournoi joué',
    playingNow: 'En ce moment',
    upNext: 'Prochain tournoi',
    readDiary: 'Lire le carnet',
    diaryComing: 'Carnet en cours d’écriture',
    noDiaryYet: 'Pas encore de carnet pour ce tournoi.',
    result: 'Résultat',
    inDays: (n) => (n === 0 ? 'commence aujourd’hui' : n === 1 ? 'commence demain' : `dans ${n} jours`),
    legendUpcoming: 'À venir',
    legendOngoing: 'En cours',
    legendPast: 'Joué — carnet disponible',
    legendPastNoArticle: 'Joué — sans carnet',
    selectedTitle: 'Tournoi sélectionné',
    close: 'Fermer',
    goingThere: 'Ce que j’y joue',
    nothingPlanned: 'Rien de prévu pour l’instant — le prochain tournoi sera annoncé ici.',
    seasonTitle: 'Toute la saison',
    seasonUpcoming: 'À venir',
    seasonPast: 'Déjà joué',
    showOnCalendar: 'Voir dans le calendrier',
  },
  en: {
    empty: 'No tournaments on the calendar yet.',
    prevMonth: 'Previous month',
    nextMonth: 'Next month',
    today: 'Today',
    jumpPrev: 'Previous tournament',
    jumpNext: 'Next tournament',
    lastPlayed: 'Last tournament played',
    playingNow: 'Playing right now',
    upNext: 'Up next',
    readDiary: 'Read the diary',
    diaryComing: 'Diary being written',
    noDiaryYet: 'No diary for this tournament yet.',
    result: 'Result',
    inDays: (n) => (n === 0 ? 'starts today' : n === 1 ? 'starts tomorrow' : `in ${n} days`),
    legendUpcoming: 'Upcoming',
    legendOngoing: 'Ongoing',
    legendPast: 'Played — diary available',
    legendPastNoArticle: 'Played — no diary',
    selectedTitle: 'Selected tournament',
    close: 'Close',
    goingThere: 'What I’m playing there',
    nothingPlanned: 'Nothing scheduled yet — the next tournament will be announced here.',
    seasonTitle: 'The whole season',
    seasonUpcoming: 'Upcoming',
    seasonPast: 'Already played',
    showOnCalendar: 'Show on the calendar',
  },
}

/** Bar colours per status. Past tournaments with a diary read as clickable. */
function barClass(status: EventStatus, hasArticle: boolean): string {
  if (status === 'ongoing') return 'bg-gold-500 text-ink-950 ring-2 ring-ink-900 ring-offset-1'
  if (status === 'upcoming') return 'bg-gold-500 text-ink-950 hover:bg-gold-400'
  return hasArticle ? 'bg-ink-800 text-white hover:bg-ink-700' : 'bg-ink-200 text-ink-700 hover:bg-ink-300'
}

/** Compact card used for the previous / ongoing / next highlights. */
function HighlightCard({
  label,
  event,
  locale,
  s,
  tone,
  today,
}: {
  label: string
  event: CalEvent | undefined
  locale: Locale
  s: (typeof STR)[Locale]
  tone: 'past' | 'next'
  today: string
}) {
  const dark = tone === 'past'
  const article = event ? articlePathFor(event, locale) : undefined
  const note = event ? noteFor(event, locale) : undefined

  return (
    <div
      className={`flex flex-col rounded-2xl border p-6 shadow-soft ${
        dark ? 'border-ink-800 bg-ink-950 text-ink-300' : 'border-gold-300 bg-gold-50 text-ink-700'
      }`}
    >
      <p
        className={`flex items-center gap-2 text-[0.68rem] font-semibold uppercase tracking-[0.16em] ${
          dark ? 'text-gold-400' : 'text-gold-700'
        }`}
      >
        <span aria-hidden className="h-px w-4 bg-gold-500" />
        {label}
      </p>

      {!event ? (
        <p className={`mt-4 flex-1 text-sm leading-relaxed ${dark ? 'text-ink-400' : 'text-ink-600'}`}>
          {tone === 'next' ? s.nothingPlanned : s.empty}
        </p>
      ) : (
        <>
          <h3 className={`mt-4 font-display text-xl font-bold ${dark ? 'text-white' : 'text-ink-900'}`}>
            {nameFor(event, locale)}
          </h3>
          <p className="mt-1.5 text-sm">
            {formatRange(event.start, event.end, locale)}
            {locationFor(event, locale) ? ` · ${locationFor(event, locale)}` : ''}
          </p>
          {formatFor(event, locale) ? (
            <p className="mt-1 text-sm opacity-80">{formatFor(event, locale)}</p>
          ) : null}

          {tone === 'next' && note ? (
            <p className="mt-4 border-l-2 border-gold-500 pl-3 text-sm italic leading-relaxed text-ink-700">
              « {note} »
            </p>
          ) : null}

          {tone === 'next' && !note ? (
            <p className="mt-4 text-sm leading-relaxed text-ink-600">
              {daysUntil(event, today) >= 0 ? s.inDays(daysUntil(event, today)) : ''}
            </p>
          ) : null}

          {tone === 'past' && event.result ? (
            <p className="mt-4 font-display text-2xl font-extrabold text-gold-400">
              {event.result}
              <span className="ml-2 align-middle text-xs font-semibold uppercase tracking-[0.14em] text-ink-400">
                {s.result}
              </span>
            </p>
          ) : null}

          <div className="mt-auto pt-5">
            {article ? (
              <Link
                to={article}
                className={`group/cta inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-semibold transition-colors ${
                  dark ? 'bg-gold-500 text-ink-950 hover:bg-gold-400' : 'bg-ink-900 text-white hover:bg-ink-800'
                }`}
              >
                {s.readDiary}
                <span aria-hidden className="transition-transform group-hover/cta:translate-x-1">
                  <IconArrowRight size={16} />
                </span>
              </Link>
            ) : tone === 'past' ? (
              <p className="text-sm text-ink-400">{s.noDiaryYet}</p>
            ) : (
              <p className="text-sm font-semibold text-gold-700">
                {daysUntil(event, today) >= 0 ? s.inDays(daysUntil(event, today)) : ''}
              </p>
            )}
          </div>
        </>
      )}
    </div>
  )
}

export function TournamentCalendar() {
  const locale = useLocale()
  const s = STR[locale]

  const events = TOURNAMENTS
  const [selected, setSelected] = useState<CalEvent | null>(null)

  // First render — pre-rendered HTML and hydration alike — uses the build date,
  // so both produce identical markup. The visitor's real date takes over in the
  // effect below, which matters once the build is a few days old.
  const [today, setToday] = useState(__BUILD_DATE__)
  useEffect(() => {
    const real = todayISO()
    setToday((t) => (t === real ? t : real))
  }, [])

  const prev = useMemo(() => previousEvent(events, today), [events, today])
  const next = useMemo(() => nextEvent(events, today), [events, today])
  const now = useMemo(() => ongoingEvent(events, today), [events, today])

  // Open on the month that actually has something to show: what's being played,
  // else what's coming, else the last thing played.
  const defaultMonth = useMemo(() => {
    const iso = (now ?? next ?? prev)?.start ?? today
    return { y: +iso.slice(0, 4), m: +iso.slice(5, 7) - 1 }
  }, [now, next, prev, today])

  const [cursor, setCursor] = useState<{ y: number; m: number } | null>(null)
  const shown = cursor ?? defaultMonth

  const weeks = useMemo(
    () => buildMonth(shown.y, shown.m, events, today),
    [shown.y, shown.m, events, today],
  )

  const shiftMonth = (delta: number) => {
    const total = shown.y * 12 + shown.m + delta
    setCursor({ y: Math.floor(total / 12), m: ((total % 12) + 12) % 12 })
  }

  const jumpTo = (e: CalEvent | undefined) => {
    if (!e) return
    setCursor({ y: +e.start.slice(0, 4), m: +e.start.slice(5, 7) - 1 })
    setSelected(e)
  }

  if (events.length === 0) {
    return <p className="rounded-2xl border border-ink-200 bg-white p-6 text-ink-600">{s.empty}</p>
  }

  const navBtn =
    'inline-flex h-9 items-center gap-1.5 rounded-full border border-ink-200 bg-white px-3.5 text-sm font-semibold text-ink-700 transition-colors hover:border-gold-400 hover:text-ink-950 disabled:cursor-not-allowed disabled:opacity-40'

  return (
    <div>
      {/* Previous / next at a glance — the shortcut into the diaries. */}
      <div className="grid gap-6 md:grid-cols-2">
        <HighlightCard label={s.lastPlayed} event={prev} locale={locale} s={s} tone="past" today={today} />
        <HighlightCard
          label={now ? s.playingNow : s.upNext}
          event={now ?? next}
          locale={locale}
          s={s}
          tone="next"
          today={today}
        />
      </div>

      {/* Month grid */}
      <div className="mt-10 overflow-hidden rounded-3xl border border-ink-200/80 bg-white shadow-card">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-ink-100 bg-cream-50 px-4 py-4 sm:px-6">
          <div className="flex items-center gap-2">
            <button type="button" onClick={() => shiftMonth(-1)} aria-label={s.prevMonth} className={navBtn}>
              <span aria-hidden>←</span>
            </button>
            <h3 className="min-w-[9.5rem] text-center font-display text-lg font-bold capitalize text-ink-900">
              {monthLabel(shown.y, shown.m, locale)}
            </h3>
            <button type="button" onClick={() => shiftMonth(1)} aria-label={s.nextMonth} className={navBtn}>
              <span aria-hidden>→</span>
            </button>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <button type="button" onClick={() => jumpTo(prev)} disabled={!prev} className={navBtn}>
              <span aria-hidden>«</span> {s.jumpPrev}
            </button>
            <button
              type="button"
              onClick={() => setCursor({ y: +today.slice(0, 4), m: +today.slice(5, 7) - 1 })}
              className={navBtn}
            >
              {s.today}
            </button>
            <button type="button" onClick={() => jumpTo(now ?? next)} disabled={!now && !next} className={navBtn}>
              {s.jumpNext} <span aria-hidden>»</span>
            </button>
          </div>
        </div>

        {/* Weekday header */}
        <div className="grid grid-cols-7 border-b border-ink-100 bg-white">
          {WEEKDAYS[locale].map((d) => (
            <div
              key={d}
              className="px-1 py-2.5 text-center text-[0.68rem] font-semibold uppercase tracking-[0.1em] text-ink-500"
            >
              <span className="hidden sm:inline">{d}</span>
              <span aria-hidden className="sm:hidden">{d.charAt(0)}</span>
            </div>
          ))}
        </div>

        {/* Weeks */}
        <div>
          {weeks.map((week, wi) => {
            // Room for the date number plus one row per stacked tournament.
            const minHeight = `${2.25 + Math.max(1, week.lanes) * 1.9}rem`
            return (
              <div key={wi} className="relative border-b border-ink-100 last:border-b-0">
                <div className="grid grid-cols-7">
                  {week.days.map((d) => (
                    <div
                      key={d.iso}
                      style={{ minHeight }}
                      className={`border-r border-ink-100 last:border-r-0 px-1.5 pt-1.5 ${
                        d.inMonth ? 'bg-white' : 'bg-cream-50/60'
                      }`}
                    >
                      <span
                        className={`inline-flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold tabular-nums ${
                          d.isToday
                            ? 'bg-ink-900 text-white'
                            : d.inMonth
                              ? 'text-ink-700'
                              : 'text-ink-300'
                        }`}
                      >
                        {d.dayOfMonth}
                      </span>
                    </div>
                  ))}
                </div>

                {/* Event bars, laid over the day cells so one tournament is one
                    continuous bar across all the days it actually runs. */}
                <div className="pointer-events-none absolute inset-x-0 top-9 grid grid-cols-7 gap-y-1 px-1">
                  {week.segments.map((seg) => {
                    const status = statusOf(seg.event, today)
                    const article = articlePathFor(seg.event, locale)
                    const isLink = status === 'past' && !!article
                    const label = `${nameFor(seg.event, locale)} · ${formatRange(seg.event.start, seg.event.end, locale)}`
                    const shape = `${seg.continuesLeft ? 'rounded-l-none' : ''} ${
                      seg.continuesRight ? 'rounded-r-none' : ''
                    }`
                    const inner = (
                      <>
                        {seg.continuesLeft ? <span aria-hidden className="opacity-70">…</span> : null}
                        <span className="truncate">{nameFor(seg.event, locale)}</span>
                        {isLink ? (
                          <span aria-hidden className="ml-auto flex-none opacity-80">
                            <IconArrowRight size={12} />
                          </span>
                        ) : null}
                      </>
                    )
                    const cls = `pointer-events-auto flex h-7 w-full items-center gap-1 rounded-md px-2 text-left text-xs font-semibold shadow-sm transition-colors ${barClass(
                      status,
                      !!article,
                    )} ${shape} ${selected?.id === seg.event.id ? 'ring-2 ring-gold-600 ring-offset-1' : ''}`

                    return (
                      <div
                        key={seg.event.id + wi}
                        style={{ gridColumnStart: seg.col + 1, gridColumnEnd: `span ${seg.span}`, gridRow: seg.lane + 1 }}
                        className="min-w-0"
                      >
                        {isLink ? (
                          <Link to={article} title={label} aria-label={`${label} — ${s.readDiary}`} className={cls}>
                            {inner}
                          </Link>
                        ) : (
                          <button
                            type="button"
                            title={label}
                            aria-label={label}
                            onClick={() => setSelected(seg.event)}
                            className={cls}
                          >
                            {inner}
                          </button>
                        )}
                      </div>
                    )
                  })}
                </div>
              </div>
            )
          })}
        </div>

        {/* Legend */}
        <div className="flex flex-wrap items-center gap-x-5 gap-y-2 border-t border-ink-100 bg-cream-50 px-4 py-3 text-xs text-ink-600 sm:px-6">
          {[
            { c: 'bg-gold-500', l: s.legendUpcoming },
            { c: 'bg-gold-500 ring-2 ring-ink-900', l: s.legendOngoing },
            { c: 'bg-ink-800', l: s.legendPast },
            { c: 'bg-ink-200', l: s.legendPastNoArticle },
          ].map((item) => (
            <span key={item.l} className="flex items-center gap-2">
              <span aria-hidden className={`h-3 w-5 rounded ${item.c}`} />
              {item.l}
            </span>
          ))}
        </div>
      </div>

      {/* Full season list. The grid only ever shows one month, so this is what
          puts every tournament — and every diary link — into the page itself. */}
      <section className="mt-10">
        <h3 className="font-display text-2xl font-bold text-ink-900">{s.seasonTitle}</h3>
        {(
          [
            [s.seasonUpcoming, events.filter((e) => statusOf(e, today) !== 'past')],
            [s.seasonPast, [...events.filter((e) => statusOf(e, today) === 'past')].reverse()],
          ] as const
        ).map(([heading, list]) =>
          list.length === 0 ? null : (
            <div key={heading} className="mt-6">
              <p className="flex items-center gap-2 text-[0.68rem] font-semibold uppercase tracking-[0.16em] text-gold-700">
                <span aria-hidden className="h-px w-4 bg-gold-500" />
                {heading}
              </p>
              <ul className="mt-3 divide-y divide-ink-100 overflow-hidden rounded-2xl border border-ink-200/80 bg-white">
                {list.map((e) => {
                  const article = articlePathFor(e, locale)
                  const place = locationFor(e, locale)
                  return (
                    <li
                      key={e.id}
                      className="flex flex-wrap items-baseline gap-x-3 gap-y-1 px-5 py-4 transition-colors hover:bg-cream-50"
                    >
                      <span className="w-40 flex-none text-sm font-semibold tabular-nums text-ink-500">
                        {formatRange(e.start, e.end, locale)}
                      </span>
                      <span className="font-semibold text-ink-900">{nameFor(e, locale)}</span>
                      {place ? <span className="text-sm text-ink-500">{place}</span> : null}
                      {e.result ? (
                        <span className="rounded-full bg-ink-900 px-2.5 py-0.5 text-xs font-bold text-gold-400">
                          {e.result}
                        </span>
                      ) : null}
                      <span className="ml-auto flex items-center gap-4">
                        {article ? (
                          <Link
                            to={article}
                            className="group/l inline-flex items-center gap-1.5 text-sm font-semibold text-gold-700 underline decoration-gold-400 decoration-1 underline-offset-4 hover:text-ink-950"
                          >
                            {s.readDiary}
                            <span aria-hidden className="transition-transform group-hover/l:translate-x-0.5">
                              <IconArrowRight size={14} />
                            </span>
                          </Link>
                        ) : null}
                        <button
                          type="button"
                          onClick={() => jumpTo(e)}
                          className="text-sm text-ink-500 underline decoration-ink-300 decoration-1 underline-offset-4 hover:text-ink-900"
                        >
                          {s.showOnCalendar}
                        </button>
                      </span>
                    </li>
                  )
                })}
              </ul>
            </div>
          ),
        )}
      </section>

      {/* Detail for a tournament clicked in the grid that has no article to open. */}
      {selected ? (
        <div className="mt-6 rounded-2xl border border-gold-300 bg-gold-50 p-6 shadow-soft">
          <div className="flex items-start justify-between gap-4">
            <div>
              <p className="flex items-center gap-2 text-[0.68rem] font-semibold uppercase tracking-[0.16em] text-gold-700">
                <IconTrophy size={14} />
                {s.selectedTitle}
              </p>
              <h3 className="mt-2 font-display text-xl font-bold text-ink-900">{nameFor(selected, locale)}</h3>
              <p className="mt-1.5 text-sm text-ink-700">
                {formatRange(selected.start, selected.end, locale)}
                {locationFor(selected, locale) ? ` · ${locationFor(selected, locale)}` : ''}
                {formatFor(selected, locale) ? ` · ${formatFor(selected, locale)}` : ''}
              </p>
              {noteFor(selected, locale) ? (
                <p className="mt-4 border-l-2 border-gold-500 pl-3 text-sm italic leading-relaxed text-ink-700">
                  <span className="not-italic font-semibold text-ink-900">{s.goingThere} : </span>
                  « {noteFor(selected, locale)} »
                </p>
              ) : null}
              {selected.result ? (
                <p className="mt-3 text-sm font-semibold text-ink-900">
                  {s.result} : {selected.result}
                </p>
              ) : null}
              {statusOf(selected, today) === 'past' && !articlePathFor(selected, locale) ? (
                <p className="mt-3 text-sm text-ink-600">{s.diaryComing}</p>
              ) : null}
            </div>
            <button
              type="button"
              onClick={() => setSelected(null)}
              className="flex-none rounded-full border border-ink-300 px-3 py-1 text-xs font-semibold text-ink-600 transition-colors hover:border-ink-400 hover:text-ink-900"
            >
              {s.close}
            </button>
          </div>
        </div>
      ) : null}
    </div>
  )
}
