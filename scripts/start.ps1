# Build le frontend puis lance le serveur Go qui sert TOUT (HTML pré-rendu +
# API + sitemap/robots/llms), comme en production, sur un seul port.
# Appelé par Lancer-le-site.cmd (double-cliquable).
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot

# --- 1) Frontend : install (si besoin) + build SSG ---
Set-Location "$root/frontend"
if (-not (Test-Path 'node_modules')) {
  Write-Host '==> Installation des dependances frontend (premiere fois)...' -ForegroundColor Cyan
  npm install
}
Write-Host '==> Build du frontend (SSG)...' -ForegroundColor Cyan
npm run build

# --- 2) Backend : le serveur Go sert frontend/dist ---
Set-Location "$root/backend"
$env:DEV_STATIC_DIR = "$root/frontend/dist"
$env:CONTENT_DIR    = "$root/content/blog"
$env:BASE_URL       = 'http://localhost:8080'
$env:ADDR           = ':8080'

Write-Host ''
Write-Host '======================================================' -ForegroundColor Green
Write-Host '  Site (backend Go sert le front) :' -ForegroundColor Green
Write-Host '     http://localhost:8080' -ForegroundColor Green
Write-Host '  SEO : /sitemap.xml  /robots.txt  /llms.txt' -ForegroundColor Green
Write-Host '  Ferme cette fenetre (ou Ctrl+C) pour arreter.' -ForegroundColor Green
Write-Host '======================================================' -ForegroundColor Green
Write-Host ''

# Ouvre le navigateur quand le serveur est prêt (2 s), puis lance le serveur.
Start-Job { Start-Sleep -Seconds 2; Start-Process 'http://localhost:8080' } | Out-Null
go run .
