import { useParams } from 'react-router-dom'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { PostCard } from '../components/PostCard'
import { getCategory } from '../lib/categories'
import { postsByCategory } from '../lib/postMeta'
import { breadcrumbSchema, type Crumb } from '../lib/schema'

export function Component() {
  const { cat = '' } = useParams()
  const category = getCategory(cat)

  if (!category) {
    return (
      <Section>
        <Container>
          <Seo title="Catégorie introuvable" description="Cette catégorie n’existe pas." path={`/blog/categorie/${cat}`} noindex />
          <h1 className="font-display text-2xl font-bold text-ink-900">Catégorie introuvable</h1>
          <p className="mt-3 text-ink-600">
            <Link to="/blog" className="text-gold-700 underline">Retour au blog</Link>.
          </p>
        </Container>
      </Section>
    )
  }

  const posts = postsByCategory(category.slug)
  const path = `/blog/categorie/${category.slug}`
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Blog', path: '/blog' },
    { name: category.short, path },
  ]

  return (
    <>
      <Seo title={category.label} description={category.description} path={path} jsonLd={[breadcrumbSchema(crumbs)]} />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Blog"
        title={category.label}
        lead={category.intro}
        primaryCta={{ to: '/contact', label: 'Réserver un cours' }}
        secondaryCta={{ to: '/blog', label: 'Tout le blog' }}
      />

      <Section>
        <Container>
          {posts.length === 0 ? (
            <p className="text-ink-600">Les premiers articles de cette catégorie arrivent bientôt.</p>
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
