/**
 * Build-time list of blog slugs, used by the router's getStaticPaths so every
 * article is pre-rendered. Kept separate from content.ts (and marked) so the
 * main router chunk stays lean — this only reads file names, not file contents.
 *
 * FR articles live in /content/blog/*.md, EN articles in /content/blog/en/*.md.
 */
// query:'?raw' + non-eager so Vite only exposes the keys (file names), never the
// contents (no Markdown parsed into this chunk).
const frFiles = import.meta.glob('../../../content/blog/*.md', { query: '?raw', import: 'default' })
const enFiles = import.meta.glob('../../../content/blog/en/*.md', { query: '?raw', import: 'default' })

const slugOf = (p: string): string => p.split('/').pop()!.replace(/\.md$/, '')

export const blogSlugs: string[] = Object.keys(frFiles).map(slugOf)
export const blogSlugsEN: string[] = Object.keys(enFiles).map(slugOf)

export const blogStaticPaths = (): string[] => blogSlugs.map((s) => `/blog/${s}`)
export const blogStaticPathsEN = (): string[] => blogSlugsEN.map((s) => `/en/blog/${s}`)
