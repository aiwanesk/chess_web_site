import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { PageHero } from '../components/PageHero'
import { allPosts } from '../lib/content'
import { breadcrumbSchema, type Crumb } from '../lib/schema'

export function Component() {
  const posts = allPosts()
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Blog', path: '/blog' },
  ]

  return (
    <>
      <Seo
        title="Blog échecs — stratégie, ouvertures, tournois"
        description="Articles d’échecs pour joueurs intermédiaires et avancés : stratégie, ouvertures, préparation de tournoi et méthodes de progression."
        path="/blog"
        jsonLd={[breadcrumbSchema(crumbs)]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <PageHero
        eyebrow="Blog"
        title="Progresser aux échecs, un article à la fois"
        lead="Stratégie, ouvertures, finales et préparation de tournoi — des analyses concrètes pour les joueurs qui veulent avancer."
        primaryCta={{ to: '/cours-echecs-adultes-geneve', label: 'Découvrir les cours' }}
      />

      <Section>
        <Container>
          {posts.length === 0 ? (
            <p className="text-ink-600">Les premiers articles arrivent bientôt.</p>
          ) : (
            <ul className="grid gap-8 sm:grid-cols-2">
              {posts.map((p) => (
                <li key={p.slug}>
                  <article className="hover-lift flex h-full flex-col rounded-2xl border border-ink-200/80 bg-white p-7 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card">
                    <p className="text-xs font-semibold uppercase tracking-[0.1em] text-gold-700">
                      <time dateTime={p.date}>
                        {new Date(p.date).toLocaleDateString('fr-CH', { day: 'numeric', month: 'long', year: 'numeric' })}
                      </time>
                      {' · '}
                      {p.readingMinutes} min de lecture
                    </p>
                    <h2 className="mt-3 font-display text-xl font-bold text-ink-900">
                      <Link to={`/blog/${p.slug}`} className="transition-colors hover:text-gold-700">
                        {p.title}
                      </Link>
                    </h2>
                    <p className="mt-2 flex-1 leading-relaxed text-ink-600">{p.description}</p>
                    <p className="mt-5 text-sm font-semibold text-gold-700">Lire l’article →</p>
                  </article>
                </li>
              ))}
            </ul>
          )}
        </Container>
      </Section>
    </>
  )
}
