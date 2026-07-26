#!/bin/zsh
# ───────────────────────────────────────────────────────────────
# start-dev-mac.sh — Lance l'environnement de dev complet (macOS)
#
# Usage:  ./scripts/start-dev-mac.sh
# ───────────────────────────────────────────────────────────────
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "══════════════════════════════════════════════════════"
echo "  Chess Web Site — Démarrage dev (macOS)"
echo "══════════════════════════════════════════════════════"

# 1. Install frontend deps if needed
if [ ! -d "frontend/node_modules" ]; then
  echo ""
  echo "→ Installation des dépendances frontend…"
  (cd frontend && npm install)
fi

# 2. Generate chess diagrams for blog
echo ""
echo "→ Génération des diagrammes d'échecs…"
node scripts/gen-diagrams.mjs

# 3. Update image references in article (svg instead of png)
echo ""
echo "→ Mise à jour des références d'images…"
if grep -q '\.png"' content/blog/open-pontevedra-2026.md 2>/dev/null; then
  sed -i '' 's/\.png"/\.svg"/g' content/blog/open-pontevedra-2026.md
  echo "  ✓ Extensions mises à jour (.png → .svg)"
else
  echo "  ✓ Références déjà à jour"
fi

# 4. Start frontend dev server
echo ""
echo "→ Lancement du serveur frontend (Vite)…"
echo "  URL: http://localhost:5173/blog/open-pontevedra-2026"
echo ""
echo "══════════════════════════════════════════════════════"
echo "  Ctrl+C pour arrêter"
echo "══════════════════════════════════════════════════════"
echo ""

cd frontend && npm run dev
