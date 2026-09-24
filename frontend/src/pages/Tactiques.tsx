import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { IconArrowRight } from '../components/icons'
import { weeks, formatWeek } from '../lib/tactics'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const STR: Record<Locale, {
  title: string; desc: string; eyebrow: string; heroTitle: string; heroLead: string; empty: string; puzzles: string; solve: string
  introTitle: string
  intro: (series: number, positions: number) => string
  howTitle: string; how: string[]
  archiveTitle: string
  ctaLabel: string
}> = {
  fr: {
    title: 'Tactiques de la semaine', desc: 'Chaque lundi, une nouvelle série des plus belles tactiques d’échecs, à résoudre. Sélectionnées par un Maître FIDE.',
    eyebrow: 'Puzzles hebdo', heroTitle: 'Les tactiques de la semaine',
    heroLead: 'Chaque lundi, une nouvelle série des plus belles combinaisons — tirées de vraies parties. À toi de les trouver.',
    empty: 'Les premières tactiques arrivent lundi matin.', puzzles: 'positions', solve: 'Résoudre',
    introTitle: 'D’où viennent ces positions',
    intro: (series, positions) =>
      `Aucune de ces positions n’a été composée. Elles sont extraites de parties réellement jouées, ` +
      `relues par un moteur pour ne garder que celles où un coup précis change le résultat, puis triées à la main : ` +
      `une combinaison qui ne se voit pas du premier coup d’œil est plus utile qu’un mat en deux forcé. ` +
      `${series} séries sont en ligne, soit ${positions} positions.`,
    howTitle: 'Comment s’en servir',
    how: [
      'Cherchez sur l’échiquier, pas dans votre tête : la position est jouable directement, cliquez la pièce puis sa case d’arrivée.',
      'Donnez-vous trois minutes par position avant de regarder la solution — c’est à peu près le temps que vous aurez en partie.',
      'Une position ratée vaut mieux qu’une position résolue : notez ce que vous n’avez pas vu, c’est là que se trouve le travail.',
    ],
    archiveTitle: 'Toutes les séries',
    ctaLabel: 'Réserver un cours',
  },
  en: {
    title: 'Tactics of the week', desc: 'Every Monday, a fresh set of the best chess tactics to solve. Hand-picked by a FIDE Master.',
    eyebrow: 'Weekly puzzles', heroTitle: 'Tactics of the week',
    heroLead: 'Every Monday, a fresh set of the best combinations — from real games. Your turn to find them.',
    empty: 'The first tactics arrive Monday morning.', puzzles: 'positions', solve: 'Solve',
    introTitle: 'Where these positions come from',
    intro: (series, positions) =>
      `None of these positions was composed. They are taken from games actually played, ` +
      `re-read by an engine to keep only those where one precise move changes the result, then sorted by hand: ` +
      `a combination you don't see at first glance teaches more than a forced mate in two. ` +
      `${series} sets are online, ${positions} positions in total.`,
    howTitle: 'How to use them',
    how: [
      'Search on the board, not in your head: the position is playable, click the piece then its destination square.',
      'Give yourself three minutes per position before looking at the solution — roughly the time you will have in a game.',
      'A position you miss is worth more than one you solve: note what you did not see, that is where the work is.',
    ],
    archiveTitle: 'All sets',
    ctaLabel: 'Book a lesson',
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
  const positions = weeks.reduce((n, w) => n + w.puzzles.length, 0)

  return (
    <>
      <Seo title={s.title} description={s.desc} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />
      <PageHero eyebrow={s.eyebrow} title={s.heroTitle} lead={s.heroLead} primaryCta={{ to: pathFor('reserver', locale), label: s.ctaLabel }} />

      {/* La page n'avait qu'un hero et une grille de vignettes : une soixantaine
          de mots, dont rien ne disait ce que ces positions sont ni à quoi elles
          servent. */}
      <Section>
        <Container>
          <div className="grid gap-10 lg:grid-cols-2">
            <div>
              <h2 className="font-display text-2xl font-bold text-ink-900">{s.introTitle}</h2>
              <p className="mt-4 leading-relaxed text-ink-600">{s.intro(weeks.length, positions)}</p>
            </div>
            <div>
              <h2 className="font-display text-2xl font-bold text-ink-900">{s.howTitle}</h2>
              <ul className="mt-4 space-y-3">
                {s.how.map((item) => (
                  <li key={item} className="flex items-start gap-3 leading-relaxed text-ink-600">
                    <span aria-hidden className="mt-2 h-1.5 w-1.5 flex-none rounded-full bg-gold-500" />
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </Container>
      </Section>

      <Section className="border-t border-ink-100 bg-cream-100">
        <Container>
          <h2 className="mb-8 font-display text-2xl font-bold text-ink-900">{s.archiveTitle}</h2>
          {weeks.length === 0 ? (
            <p className="text-center text-ink-500">{s.empty}</p>
          ) : (
            <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {weeks.map((w) => (
                <li key={w.slug}>
                  <Link
                    to={`${path}/${w.slug}`}
                    className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-paper p-6 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
                  >
                    <span className="text-xs font-semibold uppercase tracking-[0.1em] text-gold-700">{s.eyebrow}</span>
                    <h3 className="mt-2 font-display text-xl font-bold text-ink-900">{formatWeek(w.slug, locale)}</h3>
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
