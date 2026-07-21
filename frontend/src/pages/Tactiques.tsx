import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { IconArrowRight } from '../components/icons'
import { weeks, formatWeek } from '../lib/tactics'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, t, type Locale } from '../lib/i18n'

const STR: Record<Locale, {
  title: string; desc: string; eyebrow: string; heroTitle: string; heroLead: string; empty: string; puzzles: string; solve: string
}> = {
  fr: {
    title: 'Tactiques de la semaine', desc: 'Chaque lundi, une nouvelle série des plus belles tactiques d’échecs, à résoudre. Sélectionnées par un Maître FIDE.',
    eyebrow: 'Puzzles hebdo', heroTitle: 'Les tactiques de la semaine',
    heroLead: 'Chaque lundi, une nouvelle série des plus belles combinaisons — tirées de vraies parties. À toi de les trouver.',
    empty: 'Les premières tactiques arrivent lundi matin.', puzzles: 'positions', solve: 'Résoudre',
  },
  en: {
    title: 'Tactics of the week', desc: 'Every Monday, a fresh set of the best chess tactics to solve. Hand-picked by a FIDE Master.',
    eyebrow: 'Weekly puzzles', heroTitle: 'Tactics of the week',
    heroLead: 'Every Monday, a fresh set of the best combinations — from real games. Your turn to find them.',
    empty: 'The first tactics arrive Monday morning.', puzzles: 'positions', solve: 'Solve',
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const path = locale === 'en' ? '/en/tactics' : '/tactiques'
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: s.title, path },
  ]

  return (
    <>
      <Seo title={s.title} description={s.desc} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={s.eyebrow} title={s.heroTitle} lead={s.heroLead} primaryCta={{ to: locale === 'en' ? '/en/contact' : '/contact', label: locale === 'en' ? 'Book a lesson' : 'Réserver un cours' }} />

      <Section>
        <Container>
          {weeks.length === 0 ? (
            <p className="text-center text-ink-500">{s.empty}</p>
          ) : (
            <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {weeks.map((w) => (
                <li key={w.slug}>
                  <Link
                    to={`${path}/${w.slug}`}
                    className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-white p-6 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
                  >
                    <span className="text-xs font-semibold uppercase tracking-[0.1em] text-gold-700">{s.eyebrow}</span>
                    <h2 className="mt-2 font-display text-xl font-bold text-ink-900">{formatWeek(w.slug, locale)}</h2>
                    <p className="mt-1 text-sm text-ink-600">{w.puzzles.length} {s.puzzles}</p>
                    <span className="mt-4 inline-flex items-center gap-1.5 text-sm font-semibold text-gold-700">
                      {s.solve}
                      <IconArrowRight size={15} className="transition-transform group-hover:translate-x-0.5" />
                    </span>
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </Container>
      </Section>
    </>
  )
}
