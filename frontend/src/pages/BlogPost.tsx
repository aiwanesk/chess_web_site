import { useParams } from 'react-router-dom'
import { Link } from 'vite-react-ssg'
import { Seo } from '../lib/seo'
import { Container } from '../components/Container'
import { Section } from '../components/ui'
import { Breadcrumbs } from '../components/Breadcrumbs'
import { getPost } from '../lib/content'
import { postsByCluster } from '../lib/postMeta'
import { catView, categoryPath, getCategory } from '../lib/categories'
import { articleSchema, breadcrumbSchema, type Crumb } from '../lib/schema'
import { useLocale, homePath, pathFor, t, type Locale } from '../lib/i18n'

const STR: Record<Locale, {
  notFoundTitle: string; notFoundBody: string; back: string; by: string; min: string
  practiceTitle: string; practiceBody: string; practiceLink: string; alsoRead: string
}> = {
  fr: {
    notFoundTitle: 'Article introuvable',
    notFoundBody: 'Cet article n’existe pas.',
    back: 'Retour au blog',
    by: 'Par', min: 'min',
    practiceTitle: 'Envie de passer à la pratique ?',
    practiceBody: 'Ce sujet se travaille en cours.',
    practiceLink: 'Découvrir le coaching associé',
    alsoRead: 'À lire aussi',
  },
  en: {
    notFoundTitle: 'Article not found',
    notFoundBody: 'This article does not exist.',
    back: 'Back to the blog',
    by: 'By', min: 'min',
    practiceTitle: 'Ready to put it into practice?',
    practiceBody: 'This topic is worked on in lessons.',
    practiceLink: 'Discover the related coaching',
    alsoRead: 'Also worth reading',
  },
}

export function Component() {
  const locale = useLocale()
  const s = STR[locale]
  const { slug = '' } = useParams()
  const post = getPost(slug)
  const blogPath = pathFor('blog', locale)

  if (!post) {
    return (
      <Section>
        <Container>
          <Seo title={s.notFoundTitle} description={s.notFoundBody} path={`${blogPath}/${slug}`} noindex />
          <h1 className="text-2xl font-bold text-ink-900">{s.notFoundTitle}</h1>
          <p className="mt-3 text-ink-600">
            {s.notFoundBody}{' '}
            <Link to={blogPath} className="text-gold-600 underline">
              {s.back}
            </Link>
            .
          </p>
        </Container>
      </Section>
    )
  }

  const path = `${blogPath}/${post.slug}`
  const category = getCategory(post.category)
  const related = post.cluster
    ? postsByCluster(post.cluster, 4, locale).filter((p) => p.slug !== post.slug).slice(0, 3)
    : []
  const crumbs: Crumb[] = [
    { name: t(locale).breadcrumbHome, path: homePath(locale) },
    { name: 'Blog', path: blogPath },
    ...(category ? [{ name: catView(category, locale).short, path: categoryPath(category.key, locale) }] : []),
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
                {category ? (
                  <>
                    <Link to={categoryPath(category.key, locale)} className="hover:text-gold-600">
                      {catView(category, locale).short}
                    </Link>
                    {' · '}
                  </>
                ) : null}
                <time dateTime={post.date}>
                  {new Date(post.date).toLocaleDateString(locale === 'en' ? 'en-GB' : 'fr-CH', {
                    day: 'numeric',
                    month: 'long',
                    year: 'numeric',
                  })}
                </time>
                {' · '}
                {post.readingMinutes} {s.min} · {s.by} {post.author}
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
                <p className="font-semibold text-ink-900">{s.practiceTitle}</p>
                <p className="mt-1 text-ink-600">
                  {s.practiceBody}{' '}
                  <Link to={post.clusterPath} className="font-semibold text-gold-700 underline">
                    {s.practiceLink}
                  </Link>
                  .
                </p>
              </aside>
            ) : null}

            {related.length > 0 ? (
              <div className="mx-auto mt-12 max-w-[68ch]">
                <h2 className="text-xl font-bold text-ink-900">{s.alsoRead}</h2>
                <ul className="mt-4 space-y-3">
                  {related.map((r) => (
                    <li key={r.slug}>
                      <Link
                        to={`${blogPath}/${r.slug}`}
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
