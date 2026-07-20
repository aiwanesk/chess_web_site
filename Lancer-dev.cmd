@echo off
chcp 65001 >nul
title iwanesko.ch - mode dev (rechargement a chaud)
echo Demarrage du serveur de dev (Vite)...
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\start-dev.ps1"
echo.
echo Le serveur s'est arrete.
pause
