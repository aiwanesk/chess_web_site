// Régénère les assets de marque (favicon + logo header) depuis la source
// haute résolution assets/logo-source.png vers public/.
//
// sharp n'est PAS une dépendance permanente (inutile au build). Pour régénérer :
//   (cd frontend && npm i -D sharp)
//   node scripts/gen-brand-assets.mjs
//   (cd frontend && npm uninstall sharp)
//
// Les PNG générés sont commités dans public/.
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'
import { mkdirSync } from 'node:fs'
import { createRequire } from 'node:module'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
// sharp is installed in frontend/node_modules → resolve it from there.
const require = createRequire(join(root, 'frontend', 'package.json'))
const sharp = require('sharp')
const src = join(root, 'assets/logo-source.png')
const out = (name) => join(root, 'public', name)

// Match the logo's navy background so the OG canvas is seamless.
const { data } = await sharp(src).extract({ left: 0, top: 0, width: 1, height: 1 }).raw().toBuffer({ resolveWithObject: true })
const bg = { r: data[0], g: data[1], b: data[2] }

mkdirSync(out('og'), { recursive: true })

await Promise.all([
  sharp(src).resize(32, 32).png().toFile(out('favicon-32.png')),
  sharp(src).resize(48, 48).png().toFile(out('favicon-48.png')),
  sharp(src).resize(180, 180).png().toFile(out('apple-touch-icon.png')),
  sharp(src).resize(128, 128).png().toFile(out('logo-128.png')),
  // Open Graph / Twitter card: 1200×630, king centered on the brand navy.
  sharp(src).resize(1200, 630, { fit: 'contain', background: bg }).png().toFile(out('og/default.png')),
])
console.log('Brand assets regenerated in public/')
