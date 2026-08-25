#!/usr/bin/env node
/**
 * Génère les cartes Open Graph des carnets de tournoi.
 *
 * Une carte = 1200×630 : le navy de la marque, un vrai diagramme d'échiquier
 * imbriqué depuis public/images/blog/<dir>/ (donc la position affichée est
 * toujours celle que l'article montre), et le titre. Les réseaux sociaux ne
 * rendent PAS les aperçus SVG : chaque carte est rasterisée en PNG avec Chrome
 * headless, et c'est le PNG qui va dans le front-matter — comme
 * public/og/default.png. Le SVG reste à côté, c'est la source éditable.
 *
 * Usage : node scripts/gen-og.mjs
 */
import { readFileSync, writeFileSync } from 'node:fs'
import { execFileSync } from 'node:child_process'
import { dirname, join } from 'node:path'
import { fileURLToPath, pathToFileURL } from 'node:url'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')

// Palette de marque, échantillonnée sur public/og/default.png.
const NAVY = '#111626'
const GOLD = '#f5c45c'
const GOLD_DEEP = '#e5a53d'
const CREAM = '#f5f0e8'
const MUTED = '#98a3bd'
const FONT = "'Segoe UI',Inter,Helvetica,Arial,sans-serif"

const CHROME = 'C:/Program Files/Google/Chrome/Application/chrome.exe'

const FOOTER_FR = 'Alexandre Iwanesko · Maître FIDE · iwanesko.ch'
const FOOTER_EN = 'Alexandre Iwanesko · FIDE Master · iwanesko.ch'

/** Contenu interne d'un diagramme généré, prêt à être imbriqué. */
function boardBody(dir, file) {
  const svg = readFileSync(join(root, 'public/images/blog', dir, file), 'utf8')
  return svg.replace(/<\?xml[^>]*\?>\s*/, '').replace(/^\s*<svg[^>]*>/, '').replace(/<\/svg>\s*$/, '')
}

function card({ board, eyebrow, title, subtitle, score, footer }) {
  const B = 440 // côté de l'échiquier ; le diagramme source a un viewBox de 360
  const bx = 76
  const by = (630 - B) / 2
  const tx = 596
  const base = by + 92 + title.length * 60
  return `<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="630" viewBox="0 0 1200 630">
<defs>
  <linearGradient id="bg" x1="0" y1="0" x2="1" y2="1">
    <stop offset="0" stop-color="#161d31"/><stop offset="1" stop-color="${NAVY}"/>
  </linearGradient>
  <clipPath id="boardClip"><rect x="${bx}" y="${by}" width="${B}" height="${B}" rx="10"/></clipPath>
</defs>
<rect width="1200" height="630" fill="url(#bg)"/>
<rect x="22" y="22" width="1156" height="586" rx="16" fill="none" stroke="${GOLD}" stroke-opacity="0.16"/>
<rect x="${bx - 6}" y="${by - 6}" width="${B + 12}" height="${B + 12}" rx="14" fill="${GOLD}" fill-opacity="0.13"/>
<g clip-path="url(#boardClip)">
  <g transform="translate(${bx},${by}) scale(${B / 360})">${board}</g>
</g>
<text x="${tx}" y="${by + 26}" font-family="${FONT}" font-size="19" font-weight="600" letter-spacing="4.5" fill="${GOLD}">${eyebrow}</text>
${title
  .map(
    (line, i) =>
      `<text x="${tx}" y="${by + 92 + i * 60}" font-family="${FONT}" font-size="52" font-weight="700" fill="${CREAM}">${line}</text>`,
  )
  .join('\n')}
<rect x="${tx}" y="${base - 18}" width="86" height="4" rx="2" fill="${GOLD_DEEP}"/>
<text x="${tx}" y="${base + 36}" font-family="${FONT}" font-size="25" fill="${MUTED}">${subtitle}</text>
<text x="${tx}" y="${base + 84}" font-family="${FONT}" font-size="29" font-weight="600" fill="${GOLD}">${score}</text>
<text x="${tx}" y="${by + B - 4}" font-family="${FONT}" font-size="21" fill="${MUTED}">${footer}</text>
</svg>
`
}

// Un diagramme marquant par tournoi : celui que l'article met en avant.
const CARDS = [
  {
    out: 'public/images/blog/cse-2026/og-cse.png',
    dir: 'cse-2026',
    board: 'ronde6-diagramme6.svg', // la « position de rêve » après 20.Tfe1
    eyebrow: 'CARNET DE TOURNOI',
    title: ['Championnat suisse', 'par équipes'],
    subtitle: 'Nyon 1 · LNA · rondes 6 et 7',
    score: '4½–3½ · 4½–3½ — maintien assuré',
    footer: FOOTER_FR,
  },
  {
    out: 'public/images/blog/cse-2026/og-cse-en.png',
    dir: 'cse-2026',
    board: 'ronde6-diagramme6.svg',
    eyebrow: 'TOURNAMENT DIARY',
    title: ['Swiss Team', 'Championship'],
    subtitle: 'Nyon 1 · top division · rounds 6 and 7',
    score: '4½–3½ · 4½–3½ — survival secured',
    footer: FOOTER_EN,
  },
  {
    out: 'public/images/blog/badalona-2026/og-badalona.png',
    dir: 'badalona-2026',
    board: 'ronde1-diagramme4.svg', // 27…Tf1# — le mat de la ronde 1
    eyebrow: 'CARNET DE TOURNOI',
    title: ['50ᵉ Open de', 'Badalona'],
    subtitle: 'Badalona, Espagne · 2–10 août 2026',
    score: '4½/8 · dernière ronde non jouée',
    footer: FOOTER_FR,
  },
  {
    out: 'public/images/blog/badalona-2026/og-badalona-en.png',
    dir: 'badalona-2026',
    board: 'ronde1-diagramme4.svg',
    eyebrow: 'TOURNAMENT DIARY',
    title: ['50th Badalona', 'Open'],
    subtitle: 'Badalona, Spain · 2–10 August 2026',
    score: '4½/8 · final round not played',
    footer: FOOTER_EN,
  },
  {
    out: 'public/images/blog/pontevedra-2026/og-pontevedra.png',
    dir: 'pontevedra-2026',
    board: 'ronde5-diagramme4.svg', // Dc2, Fe5… Df5! — le piège se referme
    eyebrow: 'CARNET DE TOURNOI',
    title: ['Open de', 'Pontevedra'],
    subtitle: 'Pontevedra, Espagne · 25–30 juillet 2026',
    score: '9 rondes · 4½/9',
    footer: FOOTER_FR,
  },
  {
    out: 'public/images/blog/pontevedra-2026/og-pontevedra-en.png',
    dir: 'pontevedra-2026',
    board: 'ronde5-diagramme4.svg',
    eyebrow: 'TOURNAMENT DIARY',
    title: ['Pontevedra', 'Open'],
    subtitle: 'Pontevedra, Spain · 25–30 July 2026',
    score: '9 rounds · 4½/9',
    footer: FOOTER_EN,
  },
]

for (const c of CARDS) {
  const svgPath = join(root, c.out.replace(/\.png$/, '.svg'))
  const pngPath = join(root, c.out)
  writeFileSync(svgPath, card({ ...c, board: boardBody(c.dir, c.board) }), 'utf8')
  execFileSync(CHROME, [
    '--headless=new',
    '--disable-gpu',
    '--hide-scrollbars',
    '--force-device-scale-factor=1',
    '--window-size=1200,630',
    `--screenshot=${pngPath}`,
    pathToFileURL(svgPath).href,
  ])
  console.log(`✓ ${c.out}`)
}
