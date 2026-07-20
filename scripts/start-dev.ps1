# Lance le serveur de développement (Vite, rechargement à chaud).
# Appelé par Lancer-dev.cmd (double-cliquable).
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot
Set-Location "$root/frontend"

if (-not (Test-Path 'node_modules')) {
  Write-Host '==> Première installation des dépendances (patiente ~30 s)...' -ForegroundColor Cyan
  npm install
}

Write-Host ''
Write-Host '======================================================' -ForegroundColor Green
Write-Host '  Site en dev :  http://localhost:5173' -ForegroundColor Green
Write-Host '  Le navigateur va s ouvrir automatiquement.' -ForegroundColor Green
Write-Host '  Ferme cette fenetre (ou Ctrl+C) pour arreter.' -ForegroundColor Green
Write-Host '======================================================' -ForegroundColor Green
Write-Host ''

# --open : Vite ouvre le navigateur quand le serveur est prêt.
npm run dev -- --open
