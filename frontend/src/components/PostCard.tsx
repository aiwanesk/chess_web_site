import { Link } from 'vite-react-ssg'
import { catView, categoryPath, getCategory } from '../lib/categories'
import type { PostMeta } from '../lib/postMeta'
import { useLocale, type Locale } from '../lib/i18n'

const dateFmt = (iso: string, locale: Locale) =>
  new Date(iso).toLocaleDateString(locale === 'en' ? 'en-GB' : 'fr-CH', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })

/** Article teaser card used on the blog index and category archives. */
export function PostCard({ post, showCategory = true }: { post: PostMeta; showCategory?: boolean }) {
  const locale = useLocale()
  const cat = getCategory(post.category)
  const postPath = locale === 'en' ? `/en/blog/${post.slug}` : `/blog/${post.slug}`
  return (
    <article className="hover-lift flex h-full flex-col rounded-2xl border border-ink-200/80 bg-white p-6 shadow-soft transition-colors hover:border-gold-300 hover:shadow-card">
      <p className="flex flex-wrap items-center gap-x-2 gap-y-1 text-xs font-semibold uppercase tracking-[0.1em] text-gold-700">
        {showCategory && cat ? (
          <>
            <Link to={categoryPath(cat.key, locale)} className="hover:text-gold-600">
              {catView(cat, locale).short}
            </Link>
            <span aria-hidden className="text-ink-300">·</span>
          </>
        ) : null}
        <time dateTime={post.date} className="text-ink-500">
          {dateFmt(post.date, locale)}
        </time>
      </p>
      <h3 className="mt-2 font-display text-xl font-bold text-ink-900">
        <Link to={postPath} className="hover:text-gold-700">
          {post.title}
        </Link>
      </h3>
      <p className="mt-2 flex-1 leading-relaxed text-ink-600">{post.description}</p>
      <Link
        to={postPath}
        className="mt-4 inline-flex items-center gap-1.5 text-sm font-semibold text-gold-700 hover:text-gold-600"
      >
        {locale === 'en' ? 'Read the article' : 'Lire l’article'}
        <span aria-hidden>→</span>
      </Link>
    </article>
  )
}
