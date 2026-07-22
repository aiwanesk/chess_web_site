import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { BookingForm } from '../components/BookingForm'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const COPY: Record<Locale, { title: string; description: string; eyebrow: string; lead: string; asideTitle: string; aside: string[] }> = {
  fr: {
    title: 'Réserver un cours',
    description: 'Réserve un cours d’échecs particulier avec Alexandre Iwanesko, Maître FIDE, en soirée (17h30–20h00). Confirmation immédiate par e-mail.',
    eyebrow: 'Réservation',
    lead: 'Choisis ton créneau en soirée et réserve en une minute. Tu reçois une confirmation par e-mail avec le montant.',
    asideTitle: 'Comment ça marche',
    aside: [
      'Créneaux tous les soirs de 17h30 à 20h00, par tranches de 30 minutes.',
      'Sélectionne une plage (ex. 17h30–19h30) : le tarif se calcule automatiquement (120 CHF/h).',
      'Confirmation immédiate par e-mail, à toi et à Alexandre.',
      'Le paiement se règle avec Alexandre (présentiel à Genève ou en ligne).',
    ],
  },
  en: {
    title: 'Book a lesson',
    description: 'Book a private chess lesson with Alexandre Iwanesko, FIDE Master, in the evening (17:30–20:00). Instant e-mail confirmation.',
    eyebrow: 'Booking',
    lead: 'Pick your evening slot and book in a minute. You’ll get an e-mail confirmation with the amount due.',
    asideTitle: 'How it works',
    aside: [
      'Slots every evening from 17:30 to 20:00, in 30-minute steps.',
      'Select a range (e.g. 17:30–19:30): the price is computed automatically (120 CHF/h).',
      'Instant e-mail confirmation, to you and to Alexandre.',
      'Payment is arranged with Alexandre (in person in Geneva or online).',
    ],
  },
}

export function Component() {
  const locale = useLocale()
  const c = COPY[locale]
  const path = pathFor('reserver', locale)
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path },
  ]

  return (
    <>
      <Seo title={c.title} description={c.description} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={c.eyebrow} title={c.title} lead={c.lead} />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.2fr_1fr]">
          <div className="rounded-3xl border border-ink-200/80 bg-white p-7 shadow-card sm:p-8">
            <BookingForm />
          </div>
          <aside className="h-fit rounded-3xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <h2 className="font-display text-xl font-bold text-ink-900">{c.asideTitle}</h2>
            <ul className="mt-5 space-y-3 text-ink-700">
              {c.aside.map((li) => (
                <li key={li} className="flex gap-3">
                  <span aria-hidden className="mt-2 h-1.5 w-1.5 flex-none rounded-full bg-gold-500" />
                  <span className="leading-relaxed">{li}</span>
                </li>
              ))}
            </ul>
          </aside>
        </Container>
      </Section>
    </>
  )
}
