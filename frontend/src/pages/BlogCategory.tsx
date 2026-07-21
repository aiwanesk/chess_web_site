import { useParams } from 'react-router-dom'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PostCard } from '../components/PostCard'
import { catView, categoryPath, getCategoryBySlug } from '../lib/categories'
import { postsByCategory } from '../lib/postMeta'
import { breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const STR: Record<Locale, { notFound: string; back: string; empty: string; reserve: string; allBlog: string }> = {
  fr: {
    notFound: 'Catégorie introuvable',
    back: 'Retour au blog',
    empty: 'Les premiers articles de cette catégorie arrivent bientôt.',
    reserve: 'Réserver un cours',
    allBlog: 'Tout le blog',
  },
  en: {
    notFound: 'Category not found',
    back: 'Back to the blog',
    empty: 'The first articles in this category are coming soon.',
    reserve: 'Book a lesson',
    allBlog: 'All articles',
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const { cat = '' } = useParams()
  const category = getCategoryBySlug(cat, locale)

  if (!category) {
    return (
      <Section>
        <Container>
          <Seo title={s.notFound} description={s.notFound} path={pathFor('blog', locale)} noindex />
          <h1 className="font-display text-2xl font-bold text-ink-900">{s.notFound}</h1>
          <p className="mt-3 text-ink-600">
            <Link to={pathFor('blog', locale)} className="text-gold-700 underline">{s.back}</Link>.
          </p>
        </Container>
      </Section>
    )
  }

  const view = catView(category, locale)
  const posts = postsByCategory(category.key, locale)
  const path = categoryPath(category.key, locale)
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: 'Blog', path: pathFor('blog', locale) },
    { name: view.short, path },
  ]

  return (
    <>
      <Seo title={view.label} description={view.description} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Blog"
        title={view.label}
        lead={view.intro}
        primaryCta={{ to: pathFor('contact', locale), label: s.reserve }}
        secondaryCta={{ to: pathFor('blog', locale), label: s.allBlog }}
      />

      <Section>
        <Container>
          {posts.length === 0 ? (
            <p className="text-ink-600">{s.empty}</p>
          ) : (
            <ul className="grid gap-8 sm:grid-cols-2">
              {posts.map((p) => (
                <li key={p.slug}>
                  <PostCard post={p} showCategory={false} />
                </li>
              ))}
            </ul>
          )}
        </Container>
      </Section>
    </>
  )
}
