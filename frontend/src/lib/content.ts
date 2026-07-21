/**
 * Blog content pipeline. Markdown lives in /content/blog (repo root) and is
 * compiled to HTML at build time, so the pre-rendered pages ship complete HTML
 * with no client-side Markdown parser on the critical path (this module and
 * `marked` land only in the lazily-loaded blog chunks).
 */
import { marked } from 'marked'
import { DEFAULT_CATEGORY } from './categories'
import type { Locale } from './i18n'

// Raw Markdown, eagerly inlined at build. FR at repo/content/blog, EN under /en.
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
  author: string
  date: string // ISO yyyy-mm-dd
  updated?: string
  lang: Locale
  category: string // browsable section key (see categories.ts)
  cluster?: string // internal-linking silo, e.g. "adultes"
  clusterPath?: string // money page this article links to
  altSlug?: string // counterpart slug in the other locale (for hreflang)
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

function build(filePath: string, raw: string, lang: Locale): Post {
  const { data, body } = parseFrontMatter(raw)
  const words = body.split(/\s+/).filter(Boolean).length
  return {
    slug: data.slug || slugFromPath(filePath),
    title: data.title ?? 'Sans titre',
    description: data.description ?? '',
    author: data.author ?? 'Alexandre Iwanesko',
    date: data.date ?? '1970-01-01',
    updated: data.updated,
    lang,
    category: data.category || DEFAULT_CATEGORY,
    cluster: data.cluster,
    clusterPath: data.clusterPath,
    altSlug: data.altSlug,
    image: data.image,
    readingMinutes: Math.max(1, Math.round(words / 200)),
    html: marked.parse(body) as string,
  }
}

const posts: Post[] = [
  ...Object.entries(frFiles).map(([path, raw]) => build(path, raw, 'fr')),
  ...Object.entries(enFiles).map(([path, raw]) => build(path, raw, 'en')),
].sort((a, b) => (a.date < b.date ? 1 : -1))

export function allPosts(): Post[] {
  return posts
}

export function getPost(slug: string): Post | undefined {
  return posts.find((p) => p.slug === slug)
}
