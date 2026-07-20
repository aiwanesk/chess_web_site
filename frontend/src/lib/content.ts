/**
 * Blog content pipeline. Markdown lives in /content/blog (repo root) and is
 * compiled to HTML at build time, so the pre-rendered pages ship complete HTML
 * with no client-side Markdown parser on the critical path (this module and
 * `marked` land only in the lazily-loaded blog chunks).
 */
import { marked } from 'marked'
import { DEFAULT_CATEGORY } from './categories'

// Raw Markdown, eagerly inlined at build. Relative to this file: repo/content.
const rawFiles = import.meta.glob('../../../content/blog/*.md', {
  query: '?raw',
  import: 'default',
  eager: true,
}) as Record<string, string>

export interface PostMeta {
  slug: string
  title: string
  description: string
  author: string
  date: string // ISO yyyy-mm-dd
  updated?: string
  category: string // browsable section (see categories.ts)
  cluster?: string // internal-linking silo, e.g. "adultes"
  clusterPath?: string // money page this article links to
  image?: string
  readingMinutes: number
}

export interface Post extends PostMeta {
  html: string
}

marked.setOptions({ gfm: true, breaks: false })

/** Minimal flat YAML front-matter parser (string keys only). */
function parseFrontMatter(raw: string): { data: Record<string, string>; body: string } {
  const normalized = raw.replace(/\r\n/g, '\n')
  const match = /^---\n([\s\S]*?)\n---\n?/.exec(normalized)
  if (!match) return { data: {}, body: normalized }
  const data: Record<string, string> = {}
  for (const line of match[1]!.split('\n')) {
    const idx = line.indexOf(':')
    if (idx === -1) continue
    const key = line.slice(0, idx).trim()
    const value = line
      .slice(idx + 1)
      .trim()
      .replace(/^["']|["']$/g, '')
    if (key) data[key] = value
  }
  return { data, body: normalized.slice(match[0].length) }
}

function slugFromPath(filePath: string): string {
  return filePath.split('/').pop()!.replace(/\.md$/, '')
}

function build(filePath: string, raw: string): Post {
  const { data, body } = parseFrontMatter(raw)
  const words = body.split(/\s+/).filter(Boolean).length
  return {
    slug: data.slug || slugFromPath(filePath),
    title: data.title ?? 'Sans titre',
    description: data.description ?? '',
    author: data.author ?? 'Alexandre Iwanesko',
    date: data.date ?? '1970-01-01',
    updated: data.updated,
    category: data.category || DEFAULT_CATEGORY,
    cluster: data.cluster,
    clusterPath: data.clusterPath,
    image: data.image,
    readingMinutes: Math.max(1, Math.round(words / 200)),
    html: marked.parse(body) as string,
  }
}

const posts: Post[] = Object.entries(rawFiles)
  .map(([path, raw]) => build(path, raw))
  .sort((a, b) => (a.date < b.date ? 1 : -1))

export function allPosts(): Post[] {
  return posts
}

export function getPost(slug: string): Post | undefined {
  return posts.find((p) => p.slug === slug)
}
