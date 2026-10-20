@echo off
title Exellence Materiaux
chcp 65001 >nul
cd /d "%~dp0"

echo ===================================================
echo           Lancement d'Exellence Materiaux
echo ===================================================
echo.

if exist "Exellence-Materiaux.exe" (
    echo Demarrage de l'application...
    start "" "Exellence-Materiaux.exe"
    exit /b 0
)

if exist "build\bin\exellence-materiaux.exe" (
    echo Demarrage de l'application...
    start "" "build\bin\exellence-materiaux.exe"
    exit /b 0
)

echo Compilation initiale de l'application avec Wails...
wails build
if exist "build\bin\exellence-materiaux.exe" (
    copy "build\bin\exellence-materiaux.exe" "Exellence-Materiaux.exe" >nul
    start "" "Exellence-Materiaux.exe"
    exit /b 0
)

echo [ERREUR] Echec de la compilation ou du lancement.
pause
