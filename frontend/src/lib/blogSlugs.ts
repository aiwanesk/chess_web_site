/**
 * Build-time list of blog slugs, used by the router's getStaticPaths so every
 * article is pre-rendered. Kept separate from content.ts (and marked) so the
 * main router chunk stays lean — this only reads file names, not file contents.
 */
// query:'?raw' so Vite treats these as raw assets (we only read the keys/names,
// never the contents here) instead of trying to parse Markdown as JS.
const files = import.meta.glob('../../../content/blog/*.md', { query: '?raw', import: 'default' })

export const blogSlugs: string[] = Object.keys(files).map((p) =>
  p.split('/').pop()!.replace(/\.md$/, ''),
)

export const blogStaticPaths = (): string[] => blogSlugs.map((s) => `/blog/${s}`)
