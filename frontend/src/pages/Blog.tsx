import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PostCard } from '../components/PostCard'
import { IconGraduation, IconTrophy, IconArrowRight } from '../components/icons'
import { CATEGORIES, catView, categoryPath } from '../lib/categories'
import { postsForLocale, categoryCounts } from '../lib/postMeta'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const CATEGORY_ICON: Record<string, typeof IconGraduation> = {
  progresser: IconGraduation,
  'carnet-de-tournoi': IconTrophy,
}

const STR: Record<Locale, {
  title: string; description: string; heroTitle: string; heroLead: string; heroCta: string
  explore: string; exploreTitle: string; recent: string; recentTitle: string
  empty: string; article: (n: number) => string
}> = {
  fr: {
    title: 'Blog échecs — progresser & carnet',
    description:
      'Guides d’un Maître FIDE pour progresser aux échecs (ouvertures, finales, préparation tournoi) et le carnet de tournoi d’Alexandre Iwanesko.',
    heroTitle: 'Progresser aux échecs, un article à la fois',
    heroLead:
      'D’un côté des guides concrets pour gagner des points Elo ; de l’autre, mon carnet de tournoi — la compétition vue de l’intérieur.',
    heroCta: 'Découvrir les cours',
    explore: 'Explorer', exploreTitle: 'Deux façons de me suivre',
    recent: 'Derniers articles', recentTitle: 'À lire en ce moment',
    empty: 'Les premiers articles arrivent bientôt.',
    article: (n) => `${n} article${n > 1 ? 's' : ''}`,
  },
  en: {
    title: 'Chess blog — improve & diary',
    description:
      'Guides from a FIDE Master to improve at chess (openings, endgames, tournament prep) and Alexandre Iwanesko’s tournament diary.',
    heroTitle: 'Getting better at chess, one article at a time',
    heroLead:
      'On one side, concrete guides to gain Elo points; on the other, my tournament diary — competition seen from the inside.',
    heroCta: 'Explore the lessons',
    explore: 'Explore', exploreTitle: 'Two ways to follow me',
    recent: 'Latest articles', recentTitle: 'Reading right now',
    empty: 'The first articles are coming soon.',
    article: (n) => `${n} article${n > 1 ? 's' : ''}`,
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const counts = categoryCounts(locale)
  const recent = postsForLocale(locale).slice(0, 6)
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: 'Blog', path: pathFor('blog', locale) },
  ]

  return (
    <>
      <Seo title={s.title} description={s.description} path={pathFor('blog', locale)} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Blog"
        title={s.heroTitle}
        lead={s.heroLead}
        primaryCta={{ to: pathFor('coursAdultes', locale), label: s.heroCta }}
      />

      {/* Category entry points */}
      <Section>
        <Container>
          <Eyebrow>{s.explore}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{s.exploreTitle}</h2>
          <div className="mt-8 grid gap-6 md:grid-cols-2">
            {CATEGORIES.map((c) => {
              const Icon = CATEGORY_ICON[c.key] ?? IconGraduation
              const view = catView(c, locale)
              const n = counts[c.key] ?? 0
              return (
                <Link
                  key={c.key}
                  to={categoryPath(c.key, locale)}
                  className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
                >
                  <span aria-hidden className="flex h-12 w-12 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700">
                    <Icon size={24} />
                  </span>
                  <h3 className="mt-5 flex items-center justify-between font-display text-xl font-bold text-ink-900">
                    {view.label}
                    <span aria-hidden className="text-gold-600 transition-transform group-hover:translate-x-1">
                      <IconArrowRight size={18} />
                    </span>
                  </h3>
                  <p className="mt-2 flex-1 leading-relaxed text-ink-600">{view.intro}</p>
                  <p className="mt-4 text-sm font-semibold text-gold-700">{s.article(n)}</p>
                </Link>
              )
            })}
          </div>
        </Container>
      </Section>

      {/* Recent posts */}
      <Section className="border-t border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>{s.recent}</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">{s.recentTitle}</h2>
          {recent.length === 0 ? (
            <p className="mt-6 text-ink-600">{s.empty}</p>
          ) : (
            <ul className="mt-8 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
              {recent.map((p) => (
                <li key={p.slug}>
                  <PostCard post={p} />
                </li>
              ))}
            </ul>
          )}
        </Container>
      </Section>
    </>
  )
}
