#!/usr/bin/env node
/**
 * Vérifie que chaque carnet de tournoi publié est bien câblé — et fait échouer
 * la CI sinon. Sans ce garde-fou, un carnet orphelin passe le build sans un
 * mot : la page existe, mais sa barre dans le calendrier n'est pas cliquable et
 * la version traduite n'est reliée à rien.
 *
 * Ce qui est vérifié (uniquement pour `category: carnet-de-tournoi`) :
 *   1. tout carnet FR est pointé par un `slugFr` de frontend/src/lib/tournaments.ts
 *   2. tout carnet EN est pointé par un `slugEn` du même fichier
 *   3. tout `slugFr`/`slugEn` du calendrier pointe vers un article qui existe
 *   4. un `altSlug` renvoie vers un fichier existant, qui pointe en retour
 *   5. une image `image:` déclarée existe bien dans public/
 *
 * Avertissements (n'échouent pas) : image OG au format SVG, que les réseaux
 * sociaux ne savent pas afficher en aperçu de partage.
 *
 * Usage :  node scripts/check-content-links.mjs
 */
import { readFileSync, readdirSync, existsSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const DIARY = 'carnet-de-tournoi'
const CALENDAR = 'frontend/src/lib/tournaments.ts'

/** Minimal front-matter reader: `key: value`, quotes stripped. */
function frontmatter(file) {
  const block = /^---\r?\n([\s\S]*?)\r?\n---/.exec(readFileSync(file, 'utf8'))
  if (!block) return {}
  const out = {}
  for (const line of block[1].split(/\r?\n/)) {
    const kv = /^([A-Za-z][\w-]*):\s*(.*)$/.exec(line)
    if (kv) out[kv[1]] = kv[2].trim().replace(/^["']|["']$/g, '')
  }
  return out
}

function posts(rel) {
  const dir = join(root, rel)
  if (!existsSync(dir)) return []
  return readdirSync(dir)
    .filter((f) => f.endsWith('.md'))
    .map((f) => ({ slug: f.replace(/\.md$/, ''), path: `${rel}/${f}`, fm: frontmatter(join(dir, f)) }))
}

const fr = posts('content/blog')
const en = posts('content/blog/en')
const bySlug = { fr: new Map(fr.map((p) => [p.slug, p])), en: new Map(en.map((p) => [p.slug, p])) }

// The calendar is hand-edited TypeScript; read the slugs it references straight
// out of the source rather than dragging a TS toolchain into this check.
const calendar = readFileSync(join(root, CALENDAR), 'utf8')
const referenced = (key) =>
  new Set([...calendar.matchAll(new RegExp(`${key}:\\s*'([^']+)'`, 'g'))].map((m) => m[1]))
const linked = { fr: referenced('slugFr'), en: referenced('slugEn') }

const errors = []
const warnings = []

for (const [locale, list] of [
  ['fr', fr],
  ['en', en],
]) {
  const key = locale === 'fr' ? 'slugFr' : 'slugEn'

  for (const p of list) {
    if (p.fm.category !== DIARY) continue

    // 1 & 2 — the diary must be reachable from the calendar.
    if (!linked[locale].has(p.slug)) {
      errors.push(
        `${p.path} — carnet de tournoi absent du calendrier.\n` +
          `    → ajoute \`${key}: '${p.slug}'\` sur l'entrée du tournoi dans ${CALENDAR}.`,
      )
    }

    // 4 — altSlug must point at a file that points back.
    if (p.fm.altSlug) {
      const other = locale === 'fr' ? 'en' : 'fr'
      const counterpart = bySlug[other].get(p.fm.altSlug)
      if (!counterpart) {
        errors.push(
          `${p.path} — altSlug "${p.fm.altSlug}" ne correspond à aucun article ${other.toUpperCase()}.`,
        )
      } else if (counterpart.fm.altSlug !== p.slug) {
        errors.push(
          `${p.path} — appairage FR/EN à sens unique : ${counterpart.path} devrait porter ` +
            `\`altSlug: "${p.slug}"\` (il a ${counterpart.fm.altSlug ? `"${counterpart.fm.altSlug}"` : 'rien'}).\n` +
            `    → sans réciprocité, pas de hreflang.`,
        )
      }
    }

    // 5 — a declared OG image must exist, and should be a PNG.
    if (p.fm.image) {
      if (!existsSync(join(root, 'public', p.fm.image.replace(/^\//, '')))) {
        errors.push(`${p.path} — image "${p.fm.image}" introuvable dans public/.`)
      } else if (p.fm.image.endsWith('.svg')) {
        warnings.push(
          `${p.path} — image OG en SVG : les réseaux sociaux ne l'affichent pas en aperçu.\n` +
            `    → génère un PNG 1200×630 (scripts/gen-og.mjs) et pointe dessus.`,
        )
      }
    }
  }

  // 3 — no calendar entry may point at an article that doesn't exist.
  for (const slug of linked[locale]) {
    if (!bySlug[locale].has(slug)) {
      errors.push(
        `${CALENDAR} — \`${key}: '${slug}'\` ne correspond à aucun article ${locale.toUpperCase()}.\n` +
          `    → le lien du calendrier mène à une 404.`,
      )
    }
  }
}


// --- Routes du frontend vs content.StaticPages (Go) -------------------------
//
// Les ARTICLES n'ont rien à déclarer : le backend lit content/blog au runtime,
// donc sitemap.xml et llms.txt les reprennent tout seuls. Une PAGE STATIQUE,
// elle, doit être déclarée dans backend/internal/content/site.go — c'est la
// seule liste tenue à la main, et elle alimente les deux fichiers. L'oublier
// rend la page invisible des moteurs sans que rien ne casse.
const ROUTES = 'frontend/src/routes.tsx'
const SITEGO = 'backend/internal/content/site.go'

const routesSrc = readFileSync(join(root, ROUTES), 'utf8')
const siteSrc = readFileSync(join(root, SITEGO), 'utf8')

const declared = new Set(
  [...siteSrc.matchAll(/\{"(\/[^"]*)"/g)].map((m) => m[1]),
)
const allRoutePaths = [...routesSrc.matchAll(/path:\s*'([^']*)'/g)].map((m) =>
  m[1].startsWith('/') ? m[1] : `/${m[1]}`,
)
// Routes dynamiques (`:slug`), joker et page 404 : rien à faire dans un sitemap.
const routed = new Set(
  allRoutePaths.filter((p) => !p.includes(':') && !p.includes('*') && p !== '/404'),
)
// Les préfixes des routes dynamiques : `/blog/categorie/:slug` autorise
// `/blog/categorie/progresser` à figurer dans StaticPages sans route littérale.
const dynamicPrefixes = allRoutePaths
  .filter((p) => p.includes(':'))
  .map((p) => p.slice(0, p.indexOf(':')))

for (const p of routed) {
  if (!declared.has(p)) {
    errors.push(
      `${ROUTES} — la route \`${p}\` n'est pas déclarée dans ${SITEGO}.\n` +
        `    → ajoute-la à StaticPages, sinon elle est absente du sitemap ET du llms.txt.`,
    )
  }
}
for (const p of declared) {
  if (routed.has(p)) continue
  if (dynamicPrefixes.some((prefix) => p.startsWith(prefix))) continue
  warnings.push(
    `${SITEGO} — \`${p}\` est annoncé aux moteurs mais ne correspond à aucune route.\n` +
      `    → page fantôme dans le sitemap et le llms.txt.`,
  )
}

for (const w of warnings) console.warn(`⚠  ${w}`)

if (errors.length) {
  console.error(`\n✗ ${errors.length} problème(s) de câblage du contenu :\n`)
  for (const e of errors) console.error(`  ${e}\n`)
  process.exit(1)
}

const diaries = [...fr, ...en].filter((p) => p.fm.category === DIARY).length
console.log(`✓ ${diaries} carnets de tournoi, tous reliés au calendrier et correctement appairés.`)
