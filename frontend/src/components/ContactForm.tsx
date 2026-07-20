import { useState, type FormEvent } from 'react'

type Status = 'idle' | 'submitting' | 'ok' | 'error'

/**
 * Progressive-enhancement contact form. The page HTML (labels, fields) is
 * pre-rendered; this component adds the async submit to /api/contact after
 * hydration. Includes a honeypot field ("company") to deter bots.
 */
export function ContactForm() {
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
      const body = (await res.json().catch(() => ({}))) as { message?: string; error?: string }
      if (res.ok) {
        setStatus('ok')
        setMessage(body.message ?? 'Merci, votre message a bien été envoyé.')
        form.reset()
      } else {
        setStatus('error')
        setMessage(body.error ?? 'Une erreur est survenue. Merci de réessayer.')
      }
    } catch {
      setStatus('error')
      setMessage('Impossible d’envoyer le message. Vérifiez votre connexion et réessayez.')
    }
  }

  const field =
    'w-full rounded-xl border border-ink-300 bg-white px-4 py-3 text-ink-900 shadow-sm transition-colors placeholder:text-ink-400 hover:border-ink-400 focus-visible:border-gold-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gold-500/25'

  return (
    <form onSubmit={onSubmit} className="space-y-4" noValidate>
      {/* Honeypot: hidden from humans, tempting for bots. */}
      <div className="hidden" aria-hidden="true">
        <label>
          Ne pas remplir
          <input type="text" name="company" tabIndex={-1} autoComplete="off" />
        </label>
      </div>

      <div className="grid gap-4 sm:grid-cols-2">
        <div>
          <label htmlFor="name" className="mb-1 block text-sm font-medium text-ink-700">
            Nom <span className="text-gold-700">*</span>
          </label>
          <input id="name" name="name" required autoComplete="name" className={field} />
        </div>
        <div>
          <label htmlFor="email" className="mb-1 block text-sm font-medium text-ink-700">
            E-mail <span className="text-gold-700">*</span>
          </label>
          <input id="email" name="email" type="email" required autoComplete="email" className={field} />
        </div>
      </div>

      <div>
        <label htmlFor="level" className="mb-1 block text-sm font-medium text-ink-700">
          Votre niveau (Elo approximatif)
        </label>
        <input id="level" name="level" placeholder="ex. 1450 Elo" className={field} />
      </div>

      <div>
        <label htmlFor="message" className="mb-1 block text-sm font-medium text-ink-700">
          Votre objectif <span className="text-gold-700">*</span>
        </label>
        <textarea
          id="message"
          name="message"
          required
          rows={5}
          placeholder="Ce que vous aimeriez travailler, vos disponibilités…"
          className={field}
        />
      </div>

      <button
        type="submit"
        disabled={status === 'submitting'}
        className="inline-flex items-center justify-center gap-2 rounded-full bg-gold-500 px-7 py-3 font-semibold text-ink-950 shadow-soft transition-[background-color,box-shadow,transform] duration-200 hover:bg-gold-400 hover:shadow-gold active:translate-y-px disabled:cursor-not-allowed disabled:opacity-60 disabled:shadow-none"
      >
        {status === 'submitting' ? 'Envoi…' : 'Envoyer ma demande'}
      </button>

      {status === 'ok' || status === 'error' ? (
        <p
          role="status"
          className={`rounded-xl border px-4 py-3 text-sm ${
            status === 'ok'
              ? 'border-green-200 bg-green-50 text-green-800'
              : 'border-red-200 bg-red-50 text-red-800'
          }`}
        >
          {message}
        </p>
      ) : null}
    </form>
  )
}
