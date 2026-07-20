#!/usr/bin/env bash
# Builds the frontend and embeds it into the Go binary (Linux/macOS).
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"

echo "==> Building frontend (SSG)…"
(cd "$root/frontend" && npm ci && npm run build)

echo "==> Embedding build into backend/dist…"
rm -rf "$root/backend/dist"
cp -r "$root/frontend/dist" "$root/backend/dist"

echo "==> Building Go binary…"
(cd "$root/backend" && CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o "$root/server" .)

echo "Done. Run: ./server  (set BASE_URL, ADDR as needed)"
