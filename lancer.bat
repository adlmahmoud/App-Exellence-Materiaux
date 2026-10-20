@echo off
title Exellence Materiaux
chcp 65001 >nul
cd /d "%~dp0"

echo ===================================================
echo           Lancement d'Exellence Materiaux
echo ===================================================
echo.

if exist "build\bin\exellence-materiaux.exe" (
    echo Demarrage de l'application...
    start "" "build\bin\exellence-materiaux.exe"
    exit /b 0
)

echo Compilation initiale de l'application...
go build -o build\bin\exellence-materiaux.exe .
if exist "build\bin\exellence-materiaux.exe" (
    echo Demarrage de l'application...
    start "" "build\bin\exellence-materiaux.exe"
    exit /b 0
)

echo [ERREUR] Echec du lancement de l'application.
pause
