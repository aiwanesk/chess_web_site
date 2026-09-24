#!/usr/bin/env node
/**
 * Contrôles SEO sur le HTML PRÉ-RENDU (frontend/dist). Tourne en postbuild.
 *
 * Pourquoi sur la sortie de build et pas sur les composants : ce qui compte est
 * ce qu'un robot reçoit. Les balises passent par react-helmet, plusieurs pages
 * partagent le même composant avec deux configurations, et le hreflang se
 * calcule à partir d'un registre ; une erreur de câblage ne se voit qu'à la fin
 * de la chaîne. C'est d'ailleurs comme ça que /blog/<slug-anglais> a échappé à
 * tout le monde jusqu'à ce que la Search Console le remonte en 404.
 *
 * Les règles vivent dans scripts/lib/seo-html.mjs, partagées avec le crawl
 * (scripts/crawl-seo.mjs) pour que le build et la production jugent pareil.
 *
 * Usage : node scripts/check-seo.mjs [dossier-dist]
 */
import { readFileSync, readdirSync, statSync } from 'node:fs'
import { join, relative, sep } from 'node:path'
import { checkPages, parseHtml, report, TITLE_MAX, DESC_MAX } from './lib/seo-html.mjs'

const DIST = process.argv[2] ?? 'frontend/dist'
const ORIGIN = 'https://iwanesko.ch'

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

const pages = htmlFiles(DIST)
  .map((file) => {
    const path = urlPath(file)
    return { path, url: ORIGIN + (path === '/' ? '/' : path), ...parseHtml(readFileSync(file, 'utf8')) }
  })
  // La page 404 est servie avec un statut 404 et un noindex : elle n'a ni
  // canonical utile ni contrepartie à déclarer.
  .filter((p) => p.path !== '/404')

const problems = checkPages(pages)
const longest = pages.reduce((a, b) => (a.title.length >= b.title.length ? a : b))
const wordiest = pages.reduce((a, b) => (a.description.length >= b.description.length ? a : b))
process.exit(
  report(`${pages.length} pages analysées dans ${DIST}`, problems, [
    `titles ≤ ${TITLE_MAX} (le plus long : ${longest.title.length}, ${longest.path})`,
    `descriptions ≤ ${DESC_MAX} (la plus longue : ${wordiest.description.length}, ${wordiest.path})`,
    'un seul <h1> par page, canonical = URL de la page',
    'hreflang réciproques, fr/en/x-default uniquement',
    `JSON-LD valide, sans référence @id pendante (${pages.reduce((n, p) => n + p.jsonLd.length, 0)} blocs)`,
  ]),
)
