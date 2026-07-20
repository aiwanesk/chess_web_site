import { Seo } from '../lib/seo'
import { SITE } from '../lib/site'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { ContactForm } from '../components/ContactForm'
import { IconMail, IconPhone, IconPin } from '../components/icons'
import { localBusinessSchema, breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, t, type Locale } from '../lib/i18n'

const C: Record<Locale, {
  path: string; title: string; description: string; eyebrow: string; h1: string; lead: string
  coords: string; emailL: string; phoneL: string; addressL: string; areaL: string; areaJoin: string; mapTodo: string
}> = {
  fr: {
    path: '/contact', title: 'Contact',
    description: 'Contactez Alexandre Iwanesko, Maître FIDE, pour un cours d’échecs à Genève ou en ligne. Premier échange pour définir vos objectifs.',
    eyebrow: 'Contact', h1: 'Parlons de vos objectifs',
    lead: 'Décrivez votre niveau et ce que vous souhaitez travailler. Je vous réponds rapidement pour définir le format le plus adapté — présentiel à Genève ou en ligne.',
    coords: 'Coordonnées', emailL: 'E-mail', phoneL: 'Téléphone', addressL: 'Adresse', areaL: 'Zone',
    areaJoin: 'et', mapTodo: 'Carte de Genève (chargement différé) — à intégrer',
  },
  en: {
    path: '/en/contact', title: 'Contact',
    description: 'Contact Alexandre Iwanesko, FIDE Master, for chess lessons in Geneva or online. A first conversation to define your goals.',
    eyebrow: 'Contact', h1: 'Let’s talk about your goals',
    lead: 'Tell me your level and what you’d like to work on. I’ll reply quickly to define the best format — in person in Geneva or online.',
    coords: 'Contact details', emailL: 'Email', phoneL: 'Phone', addressL: 'Address', areaL: 'Area',
    areaJoin: 'and', mapTodo: 'Map of Geneva (lazy-loaded) — to be added',
  },
}

export function Component() {
  const locale = useLocale()
  const c = C[locale]
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: c.title, path: c.path },
  ]

  const iconWrap = 'mt-0.5 flex h-9 w-9 flex-none items-center justify-center rounded-full bg-white text-gold-700 shadow-sm'
  const labelC = 'block text-xs font-semibold uppercase tracking-wide text-ink-500'

  return (
    <>
      <Seo title={c.title} description={c.description} path={c.path} jsonLd={[breadcrumbSchema(crumbs), localBusinessSchema()]} />
      <Breadcrumbs crumbs={crumbs} />

      <Section>
        <Container className="grid gap-12 lg:grid-cols-[1.2fr_1fr]">
          <div>
            <Eyebrow>{c.eyebrow}</Eyebrow>
            <h1 className="font-display text-[2rem] font-bold tracking-tight text-ink-900 sm:text-[2.5rem]">{c.h1}</h1>
            <p className="mt-5 max-w-xl text-lg leading-relaxed text-ink-600">{c.lead}</p>
            <div className="mt-9 rounded-3xl border border-ink-200/80 bg-white p-7 shadow-card sm:p-8">
              <ContactForm />
            </div>
          </div>

          <aside className="h-fit rounded-3xl border border-ink-200/80 bg-cream-100 p-7 shadow-soft">
            <h2 className="font-display text-xl font-bold text-ink-900">{c.coords}</h2>
            <address className="mt-6 space-y-5 not-italic text-ink-700">
              <p className="flex items-start gap-3">
                <span aria-hidden className={iconWrap}><IconMail size={18} /></span>
                <span>
                  <span className={labelC}>{c.emailL}</span>
                  <a href={`mailto:${SITE.contact.email}`} className="font-medium text-gold-700 underline underline-offset-2">{SITE.contact.email}</a>
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className={iconWrap}><IconPhone size={18} /></span>
                <span>
                  <span className={labelC}>{c.phoneL}</span>
                  <a href={SITE.contact.phoneHref} className="font-medium text-gold-700 underline underline-offset-2">{SITE.contact.phone}</a>
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className={iconWrap}><IconPin size={18} /></span>
                <span>
                  <span className={labelC}>{c.addressL}</span>
                  {SITE.address.street}
                  <br />
                  {SITE.address.postalCode} {SITE.address.locality}
                </span>
              </p>
              <p className="flex items-start gap-3">
                <span aria-hidden className={iconWrap}><IconPin size={18} /></span>
                <span>
                  <span className={labelC}>{c.areaL}</span>
                  {SITE.address.locality} {c.areaJoin} {SITE.areaServed.join(', ')}.
                </span>
              </p>
            </address>
            <div className="mt-7 flex aspect-video items-center justify-center rounded-2xl border border-dashed border-ink-300 bg-white text-center text-sm text-ink-500">
              {c.mapTodo}
            </div>
          </aside>
        </Container>
      </Section>
    </>
  )
}
