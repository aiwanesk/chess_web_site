# Builds the frontend and embeds it into the Go binary (Windows / PowerShell).
# Usage:  ./scripts/build-local.ps1
$ErrorActionPreference = 'Stop'
$root = Split-Path -Parent $PSScriptRoot

Write-Host '==> Building frontend (SSG)…' -ForegroundColor Cyan
Push-Location "$root/frontend"
npm ci
npm run build
Pop-Location

Write-Host '==> Embedding build into backend/dist…' -ForegroundColor Cyan
$dist = "$root/backend/dist"
if (Test-Path $dist) { Remove-Item -Recurse -Force $dist }
Copy-Item -Recurse "$root/frontend/dist" $dist

Write-Host '==> Building Go binary…' -ForegroundColor Cyan
Push-Location "$root/backend"
$env:CGO_ENABLED = '0'
go build -trimpath -ldflags='-s -w' -o "$root/server.exe" .
Pop-Location

Write-Host "Done. Run: ./server.exe  (set BASE_URL, ADDR as needed)" -ForegroundColor Green
