#!/usr/bin/env node
/**
 * Generates the Open Graph cards for the tournament diaries.
 *
 * A card is a 1200x630 SVG: the brand navy, a real board diagram nested from
 * public/images/blog/<dir>/ (so the position always matches what the article
 * shows), and the headline. Social networks do NOT render SVG previews, so each
 * card is rasterised to PNG with headless Chrome and it's the PNG that goes in
 * the article front-matter — same as public/og/default.png.
 *
 * Usage: node scripts/gen-og.mjs
 */
import { readFileSync, writeFileSync } from 'node:fs'
import { execFileSync } from 'node:child_process'
import { dirname, join } from 'node:path'
import { fileURLToPath, pathToFileURL } from 'node:url'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')

// Brand palette, sampled from public/og/default.png.
const NAVY = '#111626'
const GOLD = '#f5c45c'
const GOLD_DEEP = '#e5a53d'
const CREAM = '#f5f0e8'
const MUTED = '#98a3bd'
const FONT = "'Segoe UI',Inter,Helvetica,Arial,sans-serif"

const CHROME = 'C:/Program Files/Google/Chrome/Application/chrome.exe'

/** Inner markup of a generated diagram, ready to nest in another SVG. */
function boardBody(dir, file) {
  const svg = readFileSync(join(root, 'public/images/blog', dir, file), 'utf8')
  return svg.replace(/<\?xml[^>]*\?>\s*/, '').replace(/^\s*<svg[^>]*>/, '').replace(/<\/svg>\s*$/, '')
}

function card({ dir, board, eyebrow, title, subtitle, score, footer }) {
  const B = 440 // board side; the source diagram is a 360-unit viewBox
  const bx = 76
  const by = (630 - B) / 2
  const tx = 596
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
<rect x="${tx}" y="${by + 92 + title.length * 60 - 18}" width="86" height="4" rx="2" fill="${GOLD_DEEP}"/>
<text x="${tx}" y="${by + 92 + title.length * 60 + 36}" font-family="${FONT}" font-size="25" fill="${MUTED}">${subtitle}</text>
<text x="${tx}" y="${by + 92 + title.length * 60 + 84}" font-family="${FONT}" font-size="29" font-weight="600" fill="${GOLD}">${score}</text>
<text x="${tx}" y="${by + B - 4}" font-family="${FONT}" font-size="21" fill="${MUTED}">${footer}</text>
</svg>
`
}

const boardSvg = boardBody('cse-2026', 'ronde6-diagramme6.svg')

const CARDS = [
  {
    out: 'public/images/blog/cse-2026/og-cse.png',
    eyebrow: 'CARNET DE TOURNOI',
    title: ['Championnat suisse', 'par équipes'],
    subtitle: 'Nyon 1 · LNA · rondes 6 et 7',
    score: '4½–3½ · 4½–3½ — maintien assuré',
    footer: 'Alexandre Iwanesko · Maître FIDE · iwanesko.ch',
  },
  {
    out: 'public/images/blog/cse-2026/og-cse-en.png',
    eyebrow: 'TOURNAMENT DIARY',
    title: ['Swiss Team', 'Championship'],
    subtitle: 'Nyon 1 · top division · rounds 6 and 7',
    score: '4½–3½ · 4½–3½ — survival secured',
    footer: 'Alexandre Iwanesko · FIDE Master · iwanesko.ch',
  },
]

for (const c of CARDS) {
  const svgPath = join(root, c.out.replace(/\.png$/, '.svg'))
  const pngPath = join(root, c.out)
  writeFileSync(svgPath, card({ ...c, board: boardSvg }), 'utf8')
  execFileSync(CHROME, [
    '--headless=new',
    '--disable-gpu',
    '--hide-scrollbars',
    '--force-device-scale-factor=1',
    '--default-background-color=00000000',
    '--window-size=1200,630',
    `--screenshot=${pngPath}`,
    pathToFileURL(svgPath).href,
  ])
  console.log(`✓ ${c.out}`)
}
