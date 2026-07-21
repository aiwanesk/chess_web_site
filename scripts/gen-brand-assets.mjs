// Régénère les assets de marque (favicon + logo header) depuis la source
// haute résolution assets/logo-source.png vers public/.
//
// sharp n'est PAS une dépendance permanente (inutile au build). Pour régénérer :
//   cd frontend
//   npm i -D sharp
//   node ../scripts/gen-brand-assets.mjs
//   npm uninstall sharp
//
// Les PNG générés sont commités dans public/.
import sharp from 'sharp'
import { fileURLToPath } from 'node:url'
import { dirname, join } from 'node:path'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const src = join(root, 'assets/logo-source.png')
const out = (name) => join(root, 'public', name)

await Promise.all([
  sharp(src).resize(32, 32).png().toFile(out('favicon-32.png')),
  sharp(src).resize(48, 48).png().toFile(out('favicon-48.png')),
  sharp(src).resize(180, 180).png().toFile(out('apple-touch-icon.png')),
  sharp(src).resize(128, 128).png().toFile(out('logo-128.png')),
])
console.log('Brand assets regenerated in public/')
