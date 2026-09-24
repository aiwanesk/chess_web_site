#!/usr/bin/env node
/**
 * Crawl SEO d'un serveur qui tourne, à partir de son propre sitemap.
 *
 * Le contrôle de build (check-seo.mjs) lit des fichiers ; celui-ci interroge le
 * serveur Go. C'est le seul endroit où l'on voit les codes de statut, donc la
 * seule façon de vérifier ce que le build ignore : qu'aucune URL du sitemap ne
 * répond en 301 ou en 404, que les anciennes adresses redirigent bien, et qu'une
 * URL inexistante rend un vrai 404 en noindex. Les règles sur le contenu des
 * pages sont celles du build, importées telles quelles.
 *
 * Usage : node scripts/crawl-seo.mjs [base]        (défaut : http://localhost:8099)
 */
import { checkPages, parseHtml, report } from './lib/seo-html.mjs'

const BASE = (process.argv[2] ?? 'http://localhost:8099').replace(/\/$/, '')

/** Redirections permanentes attendues : d'où, vers quoi. */
const REDIRECTS = [
  // Article placeholder supprimé (6e7270d) : vers l'archive de sa catégorie.
  ['/blog/sortir-du-plateau-1500-elo', '/blog/categorie/progresser'],
  ['/en/blog/breaking-the-1500-elo-plateau', '/en/blog/category/improve'],
  // Slug anglais sous le préfixe français, fabriqué par un lien interne fautif.
  ['/blog/returning-to-chess-after-a-break', '/en/blog/returning-to-chess-after-a-break'],
  // Une page, une adresse.
  ['/cours-echecs-adultes-geneve/', '/cours-echecs-adultes-geneve'],
  ['/index.html', '/'],
]

const problems = []
const fail = (msg) => problems.push(msg)

// ------------------------------------------------------------------ sitemap

const res = await fetch(`${BASE}/sitemap.xml`)
if (!res.ok) {
  console.error(`✗ ${BASE}/sitemap.xml répond ${res.status} — le serveur tourne-t-il ?`)
  process.exit(1)
}
const xml = await res.text()
const locs = [...xml.matchAll(/<loc>([^<]+)<\/loc>/g)].map((m) => m[1])
if (locs.length === 0) {
  console.error('✗ sitemap vide')
  process.exit(1)
}
// Le sitemap annonce les URL de production ; on va les lire sur le serveur local.
const origin = new URL(locs[0]).origin
const local = (url) => BASE + new URL(url).pathname

// ------------------------------------------------------------ pages du sitemap

/** Liens internes d'une page, hors ancres, protocoles et espace privé. */
function internalLinks(html) {
  return [...html.matchAll(/<a\b[^>]*href="(\/[^"#]*)"/g)]
    .map((m) => m[1])
    .filter((href) => !href.startsWith('/api/') && !href.startsWith('/admin'))
}

const pages = []
const linkSources = new Map() // chemin lié → pages qui le citent
for (const loc of locs) {
  const path = new URL(loc).pathname
  // redirect: 'manual' — une URL du sitemap doit répondre 200 elle-même, pas
  // après un saut. Suivre les redirections masquerait exactement ce défaut.
  const r = await fetch(local(loc), { redirect: 'manual' })
  if (r.status !== 200) {
    fail(`${path} — le sitemap déclare une URL qui répond ${r.status}` +
      `${r.headers.get('location') ? ` (→ ${r.headers.get('location')})` : ''}`)
    continue
  }
  const html = await r.text()
  const page = { path, url: loc, ...parseHtml(html) }
  if (/noindex/.test(page.robots)) {
    fail(`${path} — le sitemap déclare une URL en noindex (robots : ${page.robots})`)
  }
  pages.push(page)
  for (const href of internalLinks(html)) {
    if (!linkSources.has(href)) linkSources.set(href, new Set())
    linkSources.get(href).add(path)
  }
}
problems.push(...checkPages(pages))

// --------------------------------------------------------- liens internes

// Tout lien interne doit tomber sur un 200 DIRECT. Un 301 veut dire que le lien
// pointe à côté de la forme canonique ; un 404, qu'il pointe dans le vide.
//
// C'est le seul contrôle qui attrape le défaut d'origine de cet audit : le bloc
// « articles liés » des pages argent construisait /blog/<slug-anglais> depuis les
// pages EN, et rien — ni TypeScript, ni le lint, ni le build — ne le voyait
// passer, les deux branches étant des chaînes valides. La Search Console l'a
// découvert avant nous.
for (const [href, sources] of [...linkSources].sort()) {
  const r = await fetch(BASE + href, { redirect: 'manual' })
  if (r.status === 200) continue
  const from = [...sources].sort().slice(0, 3).join(', ')
  const extra = sources.size > 3 ? ` (+${sources.size - 3} autres)` : ''
  fail(
    `lien interne vers ${href} : ${r.status}` +
      `${r.headers.get('location') ? ` → ${r.headers.get('location')}` : ''}` +
      ` — cité par ${from}${extra}`,
  )
}

// ------------------------------------------------------------------ 301

for (const [from, to] of REDIRECTS) {
  const r = await fetch(BASE + from, { redirect: 'manual' })
  if (r.status !== 301) {
    fail(`${from} — statut ${r.status}, attendu 301`)
    continue
  }
  const got = r.headers.get('location')
  if (got !== to) fail(`${from} — Location : ${got}, attendu ${to}`)
}

// ------------------------------------------------------------------ 404

const missing = '/cette-page-n-a-jamais-existe'
const r404 = await fetch(BASE + missing, { redirect: 'manual' })
if (r404.status !== 404) {
  fail(`${missing} — statut ${r404.status}, attendu 404`)
} else {
  const { robots } = parseHtml(await r404.text())
  if (!/noindex/.test(robots)) fail(`${missing} — page 404 sans noindex (robots : « ${robots} »)`)
}

// ------------------------------------------------------------------ rapport

process.exit(
  report(`${locs.length} URL du sitemap crawlées sur ${BASE} (origine annoncée : ${origin})`, problems, [
    'toutes en 200, indexables, canonical = leur propre URL',
    `${linkSources.size} liens internes distincts, tous en 200 direct`,
    `${REDIRECTS.length} anciennes adresses en 301 vers la bonne cible`,
    'une URL inexistante rend un 404 en noindex',
  ]),
)
