/**
 * Blog metadata WITHOUT the Markdown renderer. Safe to import from money pages
 * for internal linking (related articles) — it never pulls `marked` into those
 * chunks. Front matter only; the HTML body lives in content.ts.
 */
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
  cluster?: string
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
      cluster: fm.cluster,
      clusterPath: fm.clusterPath,
    }
  })
  .sort((a, b) => (a.date < b.date ? 1 : -1))

/** Articles belonging to a given content silo (by `cluster` front-matter). */
export function postsByCluster(cluster: string, limit = 3): PostMeta[] {
  return postsMeta.filter((p) => p.cluster === cluster).slice(0, limit)
}
