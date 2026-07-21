import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section, Eyebrow } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PostCard } from '../components/PostCard'
import { IconGraduation, IconTrophy, IconArrowRight } from '../components/icons'
import { CATEGORIES } from '../lib/categories'
import { postsMeta, categoryCounts } from '../lib/postMeta'
import { breadcrumbSchema, type Crumb } from '../lib/schema'

const CATEGORY_ICON: Record<string, typeof IconGraduation> = {
  progresser: IconGraduation,
  'carnet-de-tournoi': IconTrophy,
}

export function Component() {
  const counts = categoryCounts()
  const recent = postsMeta.slice(0, 6)
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Blog', path: '/blog' },
  ]

  return (
    <>
      <Seo
        title="Blog échecs — progresser & carnet"
        description="Guides d’un Maître FIDE pour progresser aux échecs (ouvertures, finales, préparation tournoi) et le carnet de tournoi d’Alexandre Iwanesko."
        path="/blog"
        jsonLd={[breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Blog"
        title="Progresser aux échecs, un article à la fois"
        lead="D’un côté des guides concrets pour gagner des points Elo ; de l’autre, mon carnet de tournoi — la compétition vue de l’intérieur."
        primaryCta={{ to: '/cours-echecs-adultes-geneve', label: 'Découvrir les cours' }}
      />

      {/* Category entry points */}
      <Section>
        <Container>
          <Eyebrow>Explorer</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">Deux façons de me suivre</h2>
          <div className="mt-8 grid gap-6 md:grid-cols-2">
            {CATEGORIES.map((c) => {
              const Icon = CATEGORY_ICON[c.slug] ?? IconGraduation
              const n = counts[c.slug] ?? 0
              return (
                <Link
                  key={c.slug}
                  to={`/blog/categorie/${c.slug}`}
                  className="hover-lift group flex flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card"
                >
                  <span aria-hidden className="flex h-12 w-12 items-center justify-center rounded-xl border border-gold-200 bg-gold-50 text-gold-700">
                    <Icon size={24} />
                  </span>
                  <h3 className="mt-5 flex items-center justify-between font-display text-xl font-bold text-ink-900">
                    {c.label}
                    <span aria-hidden className="text-gold-600 transition-transform group-hover:translate-x-1">
                      <IconArrowRight size={18} />
                    </span>
                  </h3>
                  <p className="mt-2 flex-1 leading-relaxed text-ink-600">{c.intro}</p>
                  <p className="mt-4 text-sm font-semibold text-gold-700">
                    {n} article{n > 1 ? 's' : ''}
                  </p>
                </Link>
              )
            })}
          </div>
        </Container>
      </Section>

      {/* Recent posts */}
      <Section className="border-t border-ink-100 bg-cream-100">
        <Container>
          <Eyebrow>Derniers articles</Eyebrow>
          <h2 className="font-display text-3xl font-bold text-ink-900 sm:text-4xl">À lire en ce moment</h2>
          {recent.length === 0 ? (
            <p className="mt-6 text-ink-600">Les premiers articles arrivent bientôt.</p>
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
