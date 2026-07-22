import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link } from 'vite-react-ssg'
import { useLocale, pathFor, type Locale } from '../lib/i18n'
import { getFormToken, postWithToken } from '../lib/formToken'

type Status = 'idle' | 'submitting' | 'ok' | 'error'

// 30-min boundaries from 17:30 to 20:00. A lesson is a contiguous range.
const TIMES = ['17:30', '18:00', '18:30', '19:00', '19:30', '20:00']
const HOURLY_RATE = 120 // CHF/h (indicative — the server recomputes authoritatively)
const DAY_END = 20 * 60

const toMin = (t: string) => Number(t.slice(0, 2)) * 60 + Number(t.slice(3))
const fmt = (m: number) => `${String(Math.floor(m / 60)).padStart(2, '0')}:${String(m % 60).padStart(2, '0')}`

const STR: Record<Locale, {
  intro: string; date: string; from: string; to: string; name: string; email: string
  duration: string; total: string; submit: string; sending: string
  ok: (p: number) => string; errSlot: string; errFields: string; conflict: string; net: string
  hoursLabel: (h: number) => string; consent: string; privacy: string
  taken: string; bookedLabel: string
}> = {
  fr: {
    intro: 'Choisis un jour et un créneau entre 17h30 et 20h00. Séance individuelle à ' + HOURLY_RATE + ' CHF/h.',
    date: 'Date', from: 'De', to: 'À', name: 'Nom', email: 'E-mail',
    duration: 'Durée', total: 'Total', submit: 'Réserver le cours', sending: 'Réservation…',
    ok: (p) => `Cours réservé ! Un e-mail de confirmation t’a été envoyé. Montant : ${p} CHF.`,
    errSlot: 'Choisis un créneau valide (l’heure de fin doit être après le début).',
    errFields: 'Merci d’indiquer ton nom et un e-mail valide.',
    conflict: 'Ce créneau est déjà réservé — choisis-en un autre.',
    net: 'Impossible de réserver pour le moment. Réessaie plus tard.',
    hoursLabel: (h) => (Number.isInteger(h) ? `${h} h` : `${Math.floor(h)} h 30`),
    consent: 'En réservant, j’accepte le traitement de mes données pour gérer ce cours (voir la ',
    privacy: 'politique de confidentialité',
    taken: 'pris', bookedLabel: 'Déjà réservé ce jour :',
  },
  en: {
    intro: 'Pick a day and a slot between 17:30 and 20:00. One-to-one lesson at ' + HOURLY_RATE + ' CHF/h.',
    date: 'Date', from: 'From', to: 'To', name: 'Name', email: 'Email',
    duration: 'Duration', total: 'Total', submit: 'Book the lesson', sending: 'Booking…',
    ok: (p) => `Lesson booked! A confirmation e-mail is on its way. Amount: ${p} CHF.`,
    errSlot: 'Pick a valid slot (the end time must be after the start).',
    errFields: 'Please provide your name and a valid e-mail.',
    conflict: 'That slot is already booked — please pick another.',
    net: 'Could not book right now. Please try again later.',
    hoursLabel: (h) => (Number.isInteger(h) ? `${h} h` : `${Math.floor(h)} h 30`),
    consent: 'By booking, I agree to my data being processed to manage this lesson (see the ',
    privacy: 'privacy policy',
    taken: 'booked', bookedLabel: 'Already booked that day:',
  },
}

export function BookingForm() {
  const locale = useLocale()
  const s = STR[locale]
  const today = new Date().toISOString().slice(0, 10)

  // Server-driven constraints (opening date, minimum duration, hourly rate).
  const [cfg, setCfg] = useState({ minDate: today, minMinutes: 60, hourlyRate: HOURLY_RATE })
  const [date, setDate] = useState('')
  const [start, setStart] = useState('17:30')
  const [end, setEnd] = useState('18:30')
  const [taken, setTaken] = useState<{ start: number; end: number }[]>([])
  const [status, setStatus] = useState<Status>('idle')
  const [message, setMessage] = useState('')

  useEffect(() => {
    void getFormToken()
    fetch('/api/booking-config')
      .then((r) => (r.ok ? r.json() : null))
      .then((c: { minDate?: string; minMinutes?: number; hourlyRate?: number } | null) => {
        if (!c) return
        setCfg({
          minDate: c.minDate || today,
          minMinutes: c.minMinutes || 60,
          hourlyRate: c.hourlyRate || HOURLY_RATE,
        })
        setDate((d) => d || c.minDate || '') // pre-fill the first open day
      })
      .catch(() => {})
  }, [today])

  // Load already-booked slots whenever the chosen day changes.
  useEffect(() => {
    if (!date) {
      setTaken([])
      return
    }
    let cancelled = false
    fetch('/api/booking/availability?date=' + encodeURIComponent(date))
      .then((r) => (r.ok ? r.json() : { taken: [] }))
      .then((b: { taken?: { start: string; end: string }[] }) => {
        if (!cancelled) setTaken((b.taken || []).map((t) => ({ start: toMin(t.start), end: toMin(t.end) })))
      })
      .catch(() => {})
    return () => {
      cancelled = true
    }
  }, [date])

  const isStartTaken = (m: number) => taken.some((t) => m >= t.start && m < t.end)
  // A booked block after `m` caps how far a lesson starting at `m` can run.
  const boundaryAfter = (m: number) => taken.filter((t) => t.start >= m).reduce((b, t) => Math.min(b, t.start), DAY_END)
  const startTaken = (t: string) => {
    const m = toMin(t)
    return isStartTaken(m) || boundaryAfter(m) - m < cfg.minMinutes
  }

  const endOptions = useMemo(() => {
    const bound = boundaryAfter(toMin(start))
    return TIMES.filter((t) => toMin(t) >= toMin(start) + cfg.minMinutes && toMin(t) <= bound)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [start, cfg.minMinutes, taken])

  // If the current selection is no longer valid (day changed), jump to a free slot.
  useEffect(() => {
    if (!startTaken(start) && toMin(end) <= boundaryAfter(toMin(start)) && toMin(end) - toMin(start) >= cfg.minMinutes) return
    const freeStart = TIMES.slice(0, -1).find((t) => !startTaken(t))
    if (!freeStart) return
    setStart(freeStart)
    const e = TIMES.find((t) => toMin(t) >= toMin(freeStart) + cfg.minMinutes && toMin(t) <= boundaryAfter(toMin(freeStart)))
    if (e) setEnd(e)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [taken, cfg.minMinutes])

  const minutes = Math.max(0, toMin(end) - toMin(start))
  const hours = minutes / 60
  const price = Math.round((cfg.hourlyRate * minutes) / 60)
  const overlapsTaken = taken.some((t) => toMin(start) < t.end && toMin(end) > t.start)
  const valid = minutes >= cfg.minMinutes && !isStartTaken(toMin(start)) && !overlapsTaken

  function onStartChange(v: string) {
    setStart(v)
    const bound = boundaryAfter(toMin(v))
    if (toMin(end) < toMin(v) + cfg.minMinutes || toMin(end) > bound) {
      const next = TIMES.find((t) => toMin(t) >= toMin(v) + cfg.minMinutes && toMin(t) <= bound)
      if (next) setEnd(next)
    }
  }

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const form = e.currentTarget
    const data = Object.fromEntries(new FormData(form).entries()) as Record<string, string>
    if (!valid) {
      setStatus('error')
      setMessage(s.errSlot)
      return
    }
    setStatus('submitting')
    try {
      const res = await postWithToken('/api/booking', {
        date,
        start,
        end,
        name: data.name,
        email: data.email,
        lang: locale,
        company: data.company,
      })
      if (res.ok) {
        const body = (await res.json().catch(() => ({}))) as { price?: number }
        setStatus('ok')
        setMessage(s.ok(body.price ?? price))
        form.reset()
        setDate(cfg.minDate)
      } else {
        const body = (await res.json().catch(() => ({}))) as { error?: string }
        setStatus('error')
        setMessage(res.status === 409 ? s.conflict : body.error || (res.status === 422 ? s.errFields : s.net))
      }
    } catch {
      setStatus('error')
      setMessage(s.net)
    }
  }

  const field =
    'w-full rounded-xl border border-ink-300 bg-white px-4 py-3 text-ink-900 shadow-sm transition-colors hover:border-ink-400 focus-visible:border-gold-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gold-500/25'
  const label = 'mb-1 block text-sm font-medium text-ink-700'

  return (
    <form onSubmit={onSubmit} className="space-y-4" noValidate>
      <div className="hidden" aria-hidden="true">
        <label>Do not fill<input type="text" name="company" tabIndex={-1} autoComplete="off" /></label>
      </div>

      <p className="text-ink-600">{s.intro}</p>

      <div>
        <label htmlFor="bk-date" className={label}>{s.date} <span className="text-gold-700">*</span></label>
        <input id="bk-date" name="date" type="date" required min={cfg.minDate} value={date} onChange={(e) => setDate(e.target.value)} className={field} />
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="bk-start" className={label}>{s.from}</label>
          <select id="bk-start" value={start} onChange={(e) => onStartChange(e.target.value)} className={field}>
            {TIMES.slice(0, -1).map((t) => (
              <option key={t} value={t} disabled={startTaken(t)}>
                {t}
                {startTaken(t) ? ` — ${s.taken}` : ''}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label htmlFor="bk-end" className={label}>{s.to}</label>
          <select id="bk-end" value={end} onChange={(e) => setEnd(e.target.value)} className={field}>
            {endOptions.map((t) => <option key={t} value={t}>{t}</option>)}
          </select>
        </div>
      </div>

      {taken.length > 0 ? (
        <p className="text-xs text-ink-500" role="status">
          {s.bookedLabel} {taken.map((t) => `${fmt(t.start)}–${fmt(t.end)}`).join(', ')}
        </p>
      ) : null}

      {/* Live duration + price */}
      <div className="flex items-center justify-between rounded-xl bg-cream-100 px-4 py-3 text-sm">
        <span className="text-ink-600">{s.duration} : <strong className="text-ink-900">{s.hoursLabel(hours)}</strong></span>
        <span className="text-ink-600">{s.total} : <strong className="text-ink-900">{price} CHF</strong></span>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="bk-name" className={label}>{s.name} <span className="text-gold-700">*</span></label>
          <input id="bk-name" name="name" required autoComplete="name" className={field} />
        </div>
        <div>
          <label htmlFor="bk-email" className={label}>{s.email} <span className="text-gold-700">*</span></label>
          <input id="bk-email" name="email" type="email" required autoComplete="email" className={field} />
        </div>
      </div>

      <p className="text-xs leading-relaxed text-ink-500">
        {s.consent}
        <Link to={pathFor('confidentialite', locale)} className="text-gold-700 underline underline-offset-2">{s.privacy}</Link>).
      </p>

      <button
        type="submit"
        disabled={status === 'submitting' || !valid}
        className="inline-flex items-center justify-center gap-2 rounded-full bg-gold-500 px-7 py-3 font-semibold text-ink-950 shadow-soft transition-[background-color,box-shadow,transform] duration-200 hover:bg-gold-400 hover:shadow-gold active:translate-y-px disabled:cursor-not-allowed disabled:opacity-60 disabled:shadow-none"
      >
        {status === 'submitting' ? s.sending : `${s.submit} — ${price} CHF`}
      </button>

      {status === 'ok' || status === 'error' ? (
        <p role="status" className={`rounded-xl border px-4 py-3 text-sm ${status === 'ok' ? 'border-green-200 bg-green-50 text-green-800' : 'border-red-200 bg-red-50 text-red-800'}`}>
          {message}
        </p>
      ) : null}
    </form>
  )
}
