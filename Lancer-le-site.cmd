@echo off
chcp 65001 >nul
title iwanesko.ch - site complet (backend Go sert le front)
echo Demarrage... (build du front puis serveur Go)
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0scripts\start.ps1"
echo.
echo Le serveur s'est arrete.
pause
