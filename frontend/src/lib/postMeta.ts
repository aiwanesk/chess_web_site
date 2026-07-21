/**
 * Blog metadata WITHOUT the Markdown renderer. Safe to import from money pages
 * for internal linking (related articles) — it never pulls `marked` into those
 * chunks. Front matter only; the HTML body lives in content.ts.
 *
 * Bilingual: FR articles live in /content/blog/*.md, EN articles in
 * /content/blog/en/*.md. Language is derived from the path.
 */
import { DEFAULT_CATEGORY } from './categories'
import type { Locale } from './i18n'

const frFiles = import.meta.glob('../../../content/blog/*.md', {
  query: '?raw',
  import: 'default',
  eager: true,
}) as Record<string, string>

const enFiles = import.meta.glob('../../../content/blog/en/*.md', {
  query: '?raw',
  import: 'default',
  eager: true,
}) as Record<string, string>

export interface PostMeta {
  slug: string
  title: string
  description: string
  date: string
  lang: Locale
  category: string // stable category key (see categories.ts)
  cluster?: string // theme → money page (SEO silo)
  clusterPath?: string
  altSlug?: string // counterpart slug in the other locale (for hreflang)
}

function frontMatter(raw: string): Record<string, string> {
  const normalized = raw.replace(/\r\n/g, '\n')
  const m = /^---\n([\s\S]*?)\n---/.exec(normalized)
  const data: Record<string, string> = {}
  if (!m) return data
  for (const line of m[1]!.split('\n')) {
    const i = line.indexOf(':')
    if (i === -1) continue
    const k = line.slice(0, i).trim()
    const v = line.slice(i + 1).trim().replace(/^["']|["']$/g, '')
    if (k) data[k] = v
  }
  return data
}

function toMeta(path: string, raw: string, lang: Locale): PostMeta {
  const fm = frontMatter(raw)
  return {
    slug: fm.slug || path.split('/').pop()!.replace(/\.md$/, ''),
    title: fm.title ?? '',
    description: fm.description ?? '',
    date: fm.date ?? '1970-01-01',
    lang,
    category: fm.category || DEFAULT_CATEGORY,
    cluster: fm.cluster,
    clusterPath: fm.clusterPath,
    altSlug: fm.altSlug,
  }
}

export const postsMeta: PostMeta[] = [
  ...Object.entries(frFiles).map(([p, raw]) => toMeta(p, raw, 'fr')),
  ...Object.entries(enFiles).map(([p, raw]) => toMeta(p, raw, 'en')),
].sort((a, b) => (a.date < b.date ? 1 : -1))

/** Every article in a given locale, newest first. */
export function postsForLocale(locale: Locale): PostMeta[] {
  return postsMeta.filter((p) => p.lang === locale)
}

/** Articles belonging to a given content silo (by `cluster`), for one locale. */
export function postsByCluster(cluster: string, limit = 3, locale: Locale = 'fr'): PostMeta[] {
  return postsMeta.filter((p) => p.cluster === cluster && p.lang === locale).slice(0, limit)
}

/** Articles in a browsable category (by stable key), for one locale. */
export function postsByCategory(key: string, locale: Locale = 'fr', limit?: number): PostMeta[] {
  const list = postsMeta.filter((p) => p.category === key && p.lang === locale)
  return limit ? list.slice(0, limit) : list
}

/** Count of published articles per category key, for one locale. */
export function categoryCounts(locale: Locale = 'fr'): Record<string, number> {
  const counts: Record<string, number> = {}
  for (const p of postsMeta) {
    if (p.lang === locale) counts[p.category] = (counts[p.category] ?? 0) + 1
  }
  return counts
}
