/**
 * Blog metadata WITHOUT the Markdown renderer. Safe to import from money pages
 * for internal linking (related articles) — it never pulls `marked` into those
 * chunks. Front matter only; the HTML body lives in content.ts.
 */
import { DEFAULT_CATEGORY } from './categories'

const rawFiles = import.meta.glob('../../../content/blog/*.md', {
  query: '?raw',
  import: 'default',
  eager: true,
}) as Record<string, string>

export interface PostMeta {
  slug: string
  title: string
  description: string
  date: string
  category: string // browsable section (see categories.ts)
  cluster?: string // theme → money page (SEO silo)
  clusterPath?: string
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

export const postsMeta: PostMeta[] = Object.entries(rawFiles)
  .map(([path, raw]) => {
    const fm = frontMatter(raw)
    return {
      slug: fm.slug || path.split('/').pop()!.replace(/\.md$/, ''),
      title: fm.title ?? '',
      description: fm.description ?? '',
      date: fm.date ?? '1970-01-01',
      category: fm.category || DEFAULT_CATEGORY,
      cluster: fm.cluster,
      clusterPath: fm.clusterPath,
    }
  })
  .sort((a, b) => (a.date < b.date ? 1 : -1))

/** Articles belonging to a given content silo (by `cluster` front-matter). */
export function postsByCluster(cluster: string, limit = 3): PostMeta[] {
  return postsMeta.filter((p) => p.cluster === cluster).slice(0, limit)
}

/** Articles in a browsable category (by `category` front-matter). */
export function postsByCategory(category: string, limit?: number): PostMeta[] {
  const list = postsMeta.filter((p) => p.category === category)
  return limit ? list.slice(0, limit) : list
}

/** Count of published articles per category slug. */
export function categoryCounts(): Record<string, number> {
  const counts: Record<string, number> = {}
  for (const p of postsMeta) counts[p.category] = (counts[p.category] ?? 0) + 1
  return counts
}
