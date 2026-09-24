/**
 * Lecture et contrôle SEO d'une page HTML, partagés par deux appelants :
 *
 *  - scripts/check-seo.mjs  lit le build (frontend/dist), tourne en postbuild ;
 *  - scripts/crawl-seo.mjs  interroge un serveur qui tourne, via le sitemap.
 *
 * Les deux doivent juger à l'identique. Deux implémentations divergeraient le
 * jour où l'une gagne une règle — et c'est justement le genre d'écart qui laisse
 * passer une erreur de balise : le build dirait oui, la production non.
 */

// 60 signes : au-delà, Google tronque le titre dans la SERP (~600 px).
export const TITLE_MAX = 60
// 160 : la limite d'affichage de la description. La cible rédactionnelle est
// 150–155, mais on ne fait échouer un build que sur ce qui est réellement coupé.
export const DESC_MAX = 160
// fr, en, x-default et rien d'autre. fr-CH pointait partout vers la même URL que
// fr : un régionalisme en doublon ne dit rien de plus, et c'est une valeur de
// plus à garder cohérente.
export const ALLOWED_HREFLANG = new Set(['fr', 'en', 'x-default'])

const tags = (html, re) => [...html.matchAll(re)]
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

/** Tout ce qu'on contrôle sur une page, extrait de son HTML. */
export function parseHtml(html) {
  const alternates = {}
  for (const m of tags(html, /<link\b[^>]*rel="alternate"[^>]*>/g)) {
    const lang = attr(m[0], 'hreflang')
    if (lang) alternates[lang] = attr(m[0], 'href')
  }
  return {
    title: decode(/<title[^>]*>([\s\S]*?)<\/title>/.exec(html)?.[1] ?? ''),
    description: decode(attr(tags(html, /<meta\b[^>]*name="description"[^>]*>/g)[0]?.[0] ?? '', 'content') ?? ''),
    robots: attr(tags(html, /<meta\b[^>]*name="robots"[^>]*>/g)[0]?.[0] ?? '', 'content') ?? '',
    canonical: attr(tags(html, /<link\b[^>]*rel="canonical"[^>]*>/g)[0]?.[0] ?? '', 'href') ?? '',
    h1: tags(html, /<h1\b/g).length,
    alternates,
    jsonLd: tags(html, /<script\b[^>]*application\/ld\+json[^>]*>([\s\S]*?)<\/script>/g).map((m) => decode(m[1])),
  }
}

/** Tous les nœuds objets d'un graphe JSON-LD, à plat. */
function walk(value, out = []) {
  if (Array.isArray(value)) {
    for (const v of value) walk(v, out)
  } else if (value && typeof value === 'object') {
    out.push(value)
    for (const v of Object.values(value)) walk(v, out)
  }
  return out
}

/**
 * Contrôle un lot de pages. `pages` est une liste de { url, path, ...parseHtml }
 * où `url` est l'URL PUBLIQUE de la page (celle que la canonical doit annoncer),
 * quel que soit l'endroit d'où le HTML a été lu.
 *
 * Retourne la liste des problèmes, sous forme de chaînes prêtes à afficher.
 */
export function checkPages(pages) {
  const problems = []
  const fail = (page, msg) => problems.push(`${page.path} — ${msg}`)
  const byUrl = new Map(pages.map((p) => [p.url, p]))

  for (const page of pages) {
    // --- title et description
    if (page.title.length === 0) fail(page, 'aucun <title>')
    else if (page.title.length > TITLE_MAX) {
      fail(page, `title de ${page.title.length} signes (max ${TITLE_MAX}) : « ${page.title} »`)
    }
    if (page.description.length === 0) fail(page, 'aucune meta description')
    else if (page.description.length > DESC_MAX) {
      fail(page, `description de ${page.description.length} signes (max ${DESC_MAX})`)
    }

    // --- un seul H1 : c'est le titre du document, pas un niveau de titre
    if (page.h1 !== 1) fail(page, `${page.h1} balise(s) <h1> (il en faut exactement une)`)

    // --- canonical
    if (!page.canonical) fail(page, 'aucune canonical')
    else if (page.canonical !== page.url) {
      fail(page, `canonical = ${page.canonical}, attendu ${page.url}`)
    }

    // --- JSON-LD
    if (page.jsonLd.length === 0) fail(page, 'aucun bloc JSON-LD')
    const nodes = []
    for (const raw of page.jsonLd) {
      let data
      try {
        data = JSON.parse(raw)
      } catch (err) {
        fail(page, `JSON-LD illisible : ${err.message}`)
        continue
      }
      if (!data['@context']) fail(page, `bloc JSON-LD sans @context (@type=${data['@type'] ?? '?'})`)
      if (!data['@type']) fail(page, 'bloc JSON-LD sans @type')
      nodes.push(...walk(data))
    }
    // Un nœud qui porte un @type (ou un nom) DÉFINIT son @id ; un objet réduit
    // au seul @id n'est qu'une RÉFÉRENCE. Une référence dont la cible n'est
    // définie nulle part sur la page est pendante : un analyseur qui lit cette
    // page isolément voit un fournisseur, un auteur ou un intervenant sans nom.
    const defined = new Set(nodes.filter((n) => n['@id'] && (n['@type'] || n.name)).map((n) => n['@id']))
    for (const node of nodes) {
      const keys = Object.keys(node)
      if (keys.length === 1 && keys[0] === '@id' && !defined.has(node['@id'])) {
        fail(page, `référence @id pendante : ${node['@id']} n'est défini par aucun nœud de la page`)
      }
    }

    // --- hreflang
    const langs = Object.keys(page.alternates)
    for (const lang of langs) {
      if (!ALLOWED_HREFLANG.has(lang)) {
        fail(page, `hreflang="${lang}" non autorisé (attendu : fr, en, x-default)`)
      }
    }
    if (langs.length === 0) {
      fail(page, 'aucun hreflang')
      continue
    }
    // Toute page du site a une contrepartie dans l'autre langue — le routeur les
    // génère par paires. Une page qui ne se déclare qu'elle-même n'est donc pas
    // « sans traduction », c'est une traduction qu'on a oublié d'annoncer : le
    // cas des archives de catégorie et des tactiques de la semaine, que tolérer
    // ici aurait laissé passer une deuxième fois.
    if (langs.length === 1 && page.alternates[langs[0]] === page.url) {
      fail(page, `ne déclare qu'elle-même (hreflang="${langs[0]}") : sa contrepartie n'est pas annoncée`)
      continue
    }
    const { fr, en, 'x-default': xd } = page.alternates
    if (!fr || !en) {
      fail(page, `paire hreflang incomplète (fr=${fr ?? '—'}, en=${en ?? '—'})`)
      continue
    }
    if (xd !== fr) fail(page, `x-default=${xd ?? '—'} devrait valoir fr=${fr}`)
    for (const href of [fr, en]) {
      const other = byUrl.get(href)
      if (!other) {
        fail(page, `hreflang pointe vers une URL absente du lot analysé : ${href}`)
        continue
      }
      if (other.url === page.url) continue
      if (other.alternates.fr !== fr || other.alternates.en !== en) {
        fail(
          page,
          `paire non réciproque avec ${other.path} : ici (fr=${fr}, en=${en}), ` +
            `là-bas (fr=${other.alternates.fr ?? '—'}, en=${other.alternates.en ?? '—'})`,
        )
      }
    }
  }
  return problems
}

/** Affiche le verdict et rend le code de sortie du processus. */
export function report(label, problems, extras = []) {
  if (problems.length > 0) {
    console.error(`✗ ${label} — ${problems.length} problème(s) :\n`)
    for (const p of problems) console.error(`  • ${p}`)
    return 1
  }
  console.log(`✓ ${label}`)
  for (const line of extras) console.log(`  ${line}`)
  return 0
}
