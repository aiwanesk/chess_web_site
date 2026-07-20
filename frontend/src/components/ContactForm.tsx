import { useState, type FormEvent } from 'react'
import { useLocale, type Locale } from '../lib/i18n'

type Status = 'idle' | 'submitting' | 'ok' | 'error'

const STR: Record<Locale, {
  name: string; email: string; level: string; levelPh: string; goal: string; goalPh: string
  submit: string; sending: string; ok: string; err: string; net: string
}> = {
  fr: {
    name: 'Nom', email: 'E-mail', level: 'Votre niveau (Elo approximatif)', levelPh: 'ex. 1450 Elo',
    goal: 'Votre objectif', goalPh: 'Ce que vous aimeriez travailler, vos disponibilités…',
    submit: 'Envoyer ma demande', sending: 'Envoi…',
    ok: 'Merci, votre message a bien été envoyé. Je vous réponds rapidement.',
    err: 'Merci de renseigner votre nom, un e-mail valide et un message.',
    net: 'Impossible d’envoyer le message. Vérifiez votre connexion et réessayez.',
  },
  en: {
    name: 'Name', email: 'Email', level: 'Your level (approx. Elo)', levelPh: 'e.g. 1450 Elo',
    goal: 'Your goal', goalPh: 'What you’d like to work on, your availability…',
    submit: 'Send my request', sending: 'Sending…',
    ok: 'Thanks, your message has been sent. I’ll get back to you shortly.',
    err: 'Please provide your name, a valid email and a message.',
    net: 'Could not send the message. Check your connection and try again.',
  },
}

/**
 * Progressive-enhancement contact form. The page HTML (labels, fields) is
 * pre-rendered; this component adds the async submit to /api/contact after
 * hydration. Includes a honeypot field ("company") to deter bots.
 */
export function ContactForm() {
  const s = STR[useLocale()]
  const [status, setStatus] = useState<Status>('idle')
  const [message, setMessage] = useState('')

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setStatus('submitting')
    const form = e.currentTarget
    const data = Object.fromEntries(new FormData(form).entries())
    try {
      const res = await fetch('/api/contact', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      })
      if (res.ok) {
        setStatus('ok')
        setMessage(s.ok)
        form.reset()
      } else {
        setStatus('error')
        setMessage(s.err)
      }
    } catch {
      setStatus('error')
      setMessage(s.net)
    }
  }

  const field =
    'w-full rounded-xl border border-ink-300 bg-white px-4 py-3 text-ink-900 shadow-sm transition-colors placeholder:text-ink-400 hover:border-ink-400 focus-visible:border-gold-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gold-500/25'

  return (
    <form onSubmit={onSubmit} className="space-y-4" noValidate>
      {/* Honeypot: hidden from humans, tempting for bots. */}
      <div className="hidden" aria-hidden="true">
        <label>
          Do not fill
          <input type="text" name="company" tabIndex={-1} autoComplete="off" />
        </label>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="name" className="mb-1 block text-sm font-medium text-ink-700">
            {s.name} <span className="text-gold-700">*</span>
          </label>
          <input id="name" name="name" required autoComplete="name" className={field} />
        </div>
        <div>
          <label htmlFor="email" className="mb-1 block text-sm font-medium text-ink-700">
            {s.email} <span className="text-gold-700">*</span>
          </label>
          <input id="email" name="email" type="email" required autoComplete="email" className={field} />
        </div>
      </div>

      <div>
        <label htmlFor="level" className="mb-1 block text-sm font-medium text-ink-700">
          {s.level}
        </label>
        <input id="level" name="level" placeholder={s.levelPh} className={field} />
      </div>

      <div>
        <label htmlFor="message" className="mb-1 block text-sm font-medium text-ink-700">
          {s.goal} <span className="text-gold-700">*</span>
        </label>
        <textarea id="message" name="message" required rows={5} placeholder={s.goalPh} className={field} />
      </div>

      <button
        type="submit"
        disabled={status === 'submitting'}
        className="inline-flex items-center justify-center gap-2 rounded-full bg-gold-500 px-7 py-3 font-semibold text-ink-950 shadow-soft transition-[background-color,box-shadow,transform] duration-200 hover:bg-gold-400 hover:shadow-gold active:translate-y-px disabled:cursor-not-allowed disabled:opacity-60 disabled:shadow-none"
      >
        {status === 'submitting' ? s.sending : s.submit}
      </button>

      {status === 'ok' || status === 'error' ? (
        <p
          role="status"
          className={`rounded-xl border px-4 py-3 text-sm ${
            status === 'ok' ? 'border-green-200 bg-green-50 text-green-800' : 'border-red-200 bg-red-50 text-red-800'
          }`}
        >
          {message}
        </p>
      ) : null}
    </form>
  )
}
