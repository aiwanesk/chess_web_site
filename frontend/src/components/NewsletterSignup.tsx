import { useEffect, useState, type FormEvent } from 'react'
import { Link } from 'vite-react-ssg'
import { useLocale, pathFor, type Locale } from '../lib/i18n'
import { getFormToken, postWithToken } from '../lib/formToken'

type Status = 'idle' | 'submitting' | 'ok' | 'error'

const STR: Record<Locale, {
  title: string; blurb: string; email: string; emailPh: string
  consentPre: string; privacy: string; consentPost: string
  submit: string; sending: string; ok: string; errConsent: string; errEmail: string; net: string
}> = {
  fr: {
    title: 'Reste au courant',
    blurb: 'Nouveaux articles, tactiques de la semaine et dates de stages — directement par e-mail.',
    email: 'E-mail', emailPh: 'ton@email.ch',
    consentPre: 'J’accepte de recevoir les nouveautés par e-mail et j’ai lu la ',
    privacy: 'politique de confidentialité', consentPost: '.',
    submit: 'S’inscrire', sending: 'Envoi…',
    ok: 'Presque terminé ! Vérifie ta boîte mail et clique sur le lien de confirmation.',
    errConsent: 'Merci de cocher la case de consentement.',
    errEmail: 'Merci d’indiquer une adresse e-mail valide.',
    net: 'Impossible de t’inscrire pour le moment. Réessaie plus tard.',
  },
  en: {
    title: 'Stay in the loop',
    blurb: 'New articles, weekly tactics and upcoming events — straight to your inbox.',
    email: 'Email', emailPh: 'you@email.com',
    consentPre: 'I agree to receive updates by e-mail and have read the ',
    privacy: 'privacy policy', consentPost: '.',
    submit: 'Subscribe', sending: 'Sending…',
    ok: 'Almost there! Check your inbox and click the confirmation link.',
    errConsent: 'Please tick the consent box.',
    errEmail: 'Please enter a valid e-mail address.',
    net: 'Could not sign you up right now. Please try again later.',
  },
}

/**
 * GDPR / nLPD-compliant newsletter sign-up (double opt-in). Consent is an
 * explicit, unchecked-by-default checkbox linking to the privacy policy; the
 * backend only stores a pending record and e-mails a confirmation link.
 */
export function NewsletterSignup() {
  const locale = useLocale()
  const s = STR[locale]
  const [status, setStatus] = useState<Status>('idle')
  const [message, setMessage] = useState('')

  useEffect(() => {
    void getFormToken()
  }, [])

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const form = e.currentTarget
    const data = Object.fromEntries(new FormData(form).entries())
    if (!data.consent) {
      setStatus('error')
      setMessage(s.errConsent)
      return
    }
    setStatus('submitting')
    try {
      const res = await postWithToken('/api/newsletter/subscribe', {
        email: data.email,
        lang: locale,
        consent: true,
        company: data.company,
      })
      if (res.ok) {
        setStatus('ok')
        setMessage(s.ok)
        form.reset()
      } else {
        const body = (await res.json().catch(() => ({}))) as { error?: string }
        setStatus('error')
        setMessage(body.error || (res.status === 422 ? s.errEmail : s.net))
      }
    } catch {
      setStatus('error')
      setMessage(s.net)
    }
  }

  return (
    <div>
      <p className="flex items-center gap-2 text-xs font-semibold uppercase tracking-[0.16em] text-white">
        <span aria-hidden className="h-px w-4 bg-gold-500" />
        {s.title}
      </p>
      <p className="mt-5 text-sm leading-relaxed text-ink-400">{s.blurb}</p>

      <form onSubmit={onSubmit} className="mt-4 space-y-3" noValidate>
        {/* Honeypot */}
        <div className="hidden" aria-hidden="true">
          <label>
            Do not fill
            <input type="text" name="company" tabIndex={-1} autoComplete="off" />
          </label>
        </div>

        <div>
          <label htmlFor="nl-email" className="sr-only">
            {s.email}
          </label>
          <input
            id="nl-email"
            name="email"
            type="email"
            required
            autoComplete="email"
            placeholder={s.emailPh}
            className="w-full rounded-lg border border-white/15 bg-white/5 px-3.5 py-2.5 text-sm text-white placeholder:text-ink-500 focus-visible:border-gold-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-gold-500/25"
          />
        </div>

        <label className="flex items-start gap-2.5 text-xs leading-relaxed text-ink-400">
          <input
            type="checkbox"
            name="consent"
            value="1"
            className="mt-0.5 h-4 w-4 flex-none rounded border-white/20 bg-white/5 text-gold-500 focus-visible:ring-2 focus-visible:ring-gold-500/40"
          />
          <span>
            {s.consentPre}
            <Link to={pathFor('confidentialite', locale)} className="text-gold-400 underline underline-offset-2">
              {s.privacy}
            </Link>
            {s.consentPost}
          </span>
        </label>

        <button
          type="submit"
          disabled={status === 'submitting'}
          className="w-full rounded-lg bg-gold-500 px-4 py-2.5 text-sm font-semibold text-ink-950 transition-colors hover:bg-gold-400 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {status === 'submitting' ? s.sending : s.submit}
        </button>

        {status === 'ok' || status === 'error' ? (
          <p role="status" className={`text-xs ${status === 'ok' ? 'text-green-400' : 'text-red-400'}`}>
            {message}
          </p>
        ) : null}
      </form>
    </div>
  )
}
