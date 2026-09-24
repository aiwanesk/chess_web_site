#!/usr/bin/env node
/**
 * Contrôles SEO sur le HTML PRÉ-RENDU (frontend/dist).
 *
 * Pourquoi sur la sortie de build et pas sur les composants : ce qui compte est
 * ce qu'un robot reçoit. Les balises passent par react-helmet, plusieurs pages
 * partagent le même composant avec deux configurations, et le hreflang se
 * calcule à partir d'un registre ; une erreur de câblage ne se voit qu'à la fin
 * de la chaîne. C'est d'ailleurs comme ça que /blog/<slug-anglais> a échappé à
 * tout le monde jusqu'à ce que la Search Console le remonte en 404.
 *
 * Contrôles : longueur des titles et des descriptions, réciprocité du hreflang.
 *
 * Usage : node scripts/check-seo.mjs [dossier-dist]
 * Sort en 1 au premier problème listé, pour être utilisable en CI.
 */
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, relative, sep } from 'node:path'

const DIST = process.argv[2] ?? 'frontend/dist'
const ORIGIN = 'https://iwanesko.ch'

// ---------------------------------------------------------------- collecte

function htmlFiles(dir) {
  const out = []
  for (const entry of readdirSync(dir)) {
    const full = join(dir, entry)
    if (statSync(full).isDirectory()) out.push(...htmlFiles(full))
    else if (entry.endsWith('.html')) out.push(full)
  }
  return out
}

/** Chemin d'URL servi par un fichier : dist/en/about.html → /en/about. */
function urlPath(file) {
  const rel = relative(DIST, file).split(sep).join('/')
  const noExt = rel.replace(/\.html$/, '')
  return noExt === 'index' ? '/' : `/${noExt}`
}

const tag = (html, re) => [...html.matchAll(re)]
const attr = (s, name) => {
  const m = new RegExp(`${name}="([^"]*)"`).exec(s)
  return m ? m[1] : undefined
}
const decode = (s = '') =>
  s
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&amp;/g, '&')

function parse(file) {
  const html = readFileSync(file, 'utf8')
  const alternates = {}
  for (const m of tag(html, /<link\b[^>]*rel="alternate"[^>]*>/g)) {
    const lang = attr(m[0], 'hreflang')
    if (lang) alternates[lang] = attr(m[0], 'href')
  }
  return {
    file,
    path: urlPath(file),
    title: decode(/<title[^>]*>([\s\S]*?)<\/title>/.exec(html)?.[1] ?? ''),
    description: decode(attr(tag(html, /<meta\b[^>]*name="description"[^>]*>/g)[0]?.[0] ?? '', 'content') ?? ''),
    robots: attr(tag(html, /<meta\b[^>]*name="robots"[^>]*>/g)[0]?.[0] ?? '', 'content') ?? '',
    alternates,
  }
}

const pages = htmlFiles(DIST)
  .map(parse)
  .filter((p) => p.path !== '/404')

const byUrl = new Map(pages.map((p) => [ORIGIN + (p.path === '/' ? '/' : p.path), p]))
const problems = []
const fail = (page, msg) => problems.push(`${page.path} — ${msg}`)

// --------------------------------------------------- titles et descriptions

// 60 signes : au-delà, Google tronque dans la SERP (~600 px). Le suffixe de
// marque se retire tout seul quand il ferait déborder (lib/seo.tsx), donc un
// dépassement ici ne peut venir que du titre lui-même — à raccourcir à la main,
// ou via `seoTitle:` en front-matter pour un article.
const TITLE_MAX = 60
// 160 : la limite d'affichage. La cible rédactionnelle est 150–155, mais on ne
// fait échouer un build que sur ce qui est réellement coupé.
const DESC_MAX = 160

for (const page of pages) {
  if (page.title.length === 0) fail(page, 'aucun <title>')
  else if (page.title.length > TITLE_MAX) {
    fail(page, `title de ${page.title.length} signes (max ${TITLE_MAX}) : « ${page.title} »`)
  }
  if (page.description.length === 0) fail(page, 'aucune meta description')
  else if (page.description.length > DESC_MAX) {
    fail(page, `description de ${page.description.length} signes (max ${DESC_MAX})`)
  }
}

// ------------------------------------------------------- hreflang réciproque

// fr, en, x-default et rien d'autre. fr-CH pointait partout vers la même URL
// que fr : un régionalisme en doublon ne dit rien de plus, et c'est une ligne
// de plus à garder cohérente.
const ALLOWED = new Set(['fr', 'en', 'x-default'])

for (const page of pages) {
  const langs = Object.keys(page.alternates)
  for (const lang of langs) {
    if (!ALLOWED.has(lang)) fail(page, `hreflang="${lang}" non autorisé (attendu : fr, en, x-default)`)
  }

  // Toute page du site a une contrepartie dans l'autre langue — routes, blog,
  // catégories, séries hebdomadaires : le routeur les génère par paires. Une
  // page qui ne se déclare qu'elle-même n'est donc pas « sans traduction »,
  // c'est une traduction qu'on a oublié d'annoncer. C'est exactement ce qui
  // était arrivé aux archives de catégorie et aux tactiques de la semaine : les
  // deux côtés existaient, aucun ne citait l'autre, et tolérer ce cas dans le
  // contrôle l'aurait laissé passer une deuxième fois.
  if (langs.length === 0) {
    fail(page, 'aucun hreflang')
    continue
  }
  if (langs.length === 1 && page.alternates[langs[0]] === ORIGIN + page.path) {
    fail(page, `ne déclare qu'elle-même (hreflang="${langs[0]}") : sa contrepartie n'est pas annoncée`)
    continue
  }
  const { fr, en, 'x-default': xd } = page.alternates
  if (!fr || !en) {
    fail(page, `paire incomplète (fr=${fr ?? '—'}, en=${en ?? '—'})`)
    continue
  }
  if (xd !== fr) fail(page, `x-default=${xd ?? '—'} devrait valoir fr=${fr}`)

  // Les deux URL annoncées doivent exister…
  for (const href of [fr, en]) {
    if (!byUrl.has(href)) {
      fail(page, `hreflang pointe vers une URL qui n'existe pas dans le build : ${href}`)
    }
  }
  // …et déclarer exactement la même paire, sinon Google ignore les deux.
  for (const href of [fr, en]) {
    const other = byUrl.get(href)
    if (!other || other.path === page.path) continue
    if (other.alternates.fr !== fr || other.alternates.en !== en) {
      fail(
        page,
        `paire non réciproque avec ${other.path} : ici (fr=${fr}, en=${en}), là-bas (fr=${other.alternates.fr ?? '—'}, en=${other.alternates.en ?? '—'})`,
      )
    }
  }
}

// ------------------------------------------------------------------ rapport

const label = `${pages.length} pages analysées dans ${DIST}`
if (problems.length > 0) {
  console.error(`✗ ${label} — ${problems.length} problème(s) :\n`)
  for (const p of problems) console.error(`  • ${p}`)
  process.exit(1)
}
const longest = pages.reduce((a, b) => (a.title.length >= b.title.length ? a : b))
const wordiest = pages.reduce((a, b) => (a.description.length >= b.description.length ? a : b))
console.log(`✓ ${label}`)
console.log(`  titles ≤ ${TITLE_MAX} (le plus long : ${longest.title.length}, ${longest.path})`)
console.log(`  descriptions ≤ ${DESC_MAX} (la plus longue : ${wordiest.description.length}, ${wordiest.path})`)
console.log('  hreflang réciproques, fr/en/x-default uniquement')
