import { useParams } from 'react-router-dom'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { getPost } from '../lib/content'
import { postsByCluster } from '../lib/postMeta'
import { articleSchema, breadcrumbSchema, type Crumb } from '../lib/schema'

export function Component() {
  const { slug = '' } = useParams()
  const post = getPost(slug)

  if (!post) {
    return (
      <Section>
        <Container>
          <Seo title="Article introuvable" description="Cet article n’existe pas ou a été déplacé." path={`/blog/${slug}`} noindex />
          <h1 className="text-2xl font-bold text-ink-900">Article introuvable</h1>
          <p className="mt-3 text-ink-600">
            Cet article n’existe pas.{' '}
            <Link to="/blog" className="text-gold-600 underline">
              Retour au blog
            </Link>
            .
          </p>
        </Container>
      </Section>
    )
  }

  const path = `/blog/${post.slug}`
  const related = post.cluster
    ? postsByCluster(post.cluster, 4).filter((p) => p.slug !== post.slug).slice(0, 3)
    : []
  const crumbs: Crumb[] = [
    { name: 'Accueil', path: '/' },
    { name: 'Blog', path: '/blog' },
    { name: post.title, path },
  ]

  return (
    <>
      <Seo
        title={post.title}
        description={post.description}
        path={path}
        ogType="article"
        image={post.image}
        jsonLd={[
          breadcrumbSchema(crumbs),
          articleSchema({
            title: post.title,
            description: post.description,
            url: path,
            datePublished: post.date,
            dateModified: post.updated,
            image: post.image,
          }),
        ]}
      />
      <Breadcrumbs crumbs={crumbs} />

      <Section>
        <Container>
          <article>
            <header className="mx-auto max-w-[68ch]">
              <p className="text-xs font-semibold uppercase tracking-[0.1em] text-gold-700">
                <time dateTime={post.date}>
                  {new Date(post.date).toLocaleDateString('fr-CH', { day: 'numeric', month: 'long', year: 'numeric' })}
                </time>
                {' · '}
                {post.readingMinutes} min · Par {post.author}
              </p>
              <h1 className="mt-3 font-display text-[2rem] font-bold leading-[1.12] tracking-tight text-ink-900 sm:text-[2.6rem]">
                {post.title}
              </h1>
              <p className="mt-4 text-lg leading-relaxed text-ink-600">{post.description}</p>
              <div aria-hidden className="mt-6 h-px w-16 bg-gold-400" />
            </header>

            <div
              className="prose mx-auto mt-10"
              // Compiled from Markdown at build time; content is authored/trusted.
              dangerouslySetInnerHTML={{ __html: post.html }}
            />

            {/* Cluster internal link back to the related money page. */}
            {post.clusterPath ? (
              <aside className="mx-auto mt-12 max-w-[68ch] rounded-2xl border border-gold-500/40 bg-gold-500/5 p-6">
                <p className="font-semibold text-ink-900">Envie de passer à la pratique ?</p>
                <p className="mt-1 text-ink-600">
                  Ce sujet se travaille en cours.{' '}
                  <Link to={post.clusterPath} className="font-semibold text-gold-700 underline">
                    Découvrir le coaching associé
                  </Link>
                  .
                </p>
              </aside>
            ) : null}

            {related.length > 0 ? (
              <div className="mx-auto mt-12 max-w-[68ch]">
                <h2 className="text-xl font-bold text-ink-900">À lire aussi</h2>
                <ul className="mt-4 space-y-3">
                  {related.map((r) => (
                    <li key={r.slug}>
                      <Link
                        to={`/blog/${r.slug}`}
                        className="flex items-center justify-between gap-4 rounded-lg border border-ink-200 px-5 py-4 font-medium text-ink-800 transition-colors hover:border-gold-500 hover:text-ink-950"
                      >
                        {r.title}
                        <span aria-hidden className="text-gold-500">→</span>
                      </Link>
                    </li>
                  ))}
                </ul>
              </div>
            ) : null}
          </article>
        </Container>
      </Section>
    </>
  )
}
