@echo off
title Gestion de Stock
chcp 65001 >nul
cd /d "%~dp0"

echo ===================================================
echo           Lancement de Gestion de Stock
echo ===================================================
echo.

if exist "Gestion-de-Stock.exe" (
    echo Demarrage de l'application...
    start "" "Gestion-de-Stock.exe"
    exit /b 0
)

if exist "build\bin\Gestion-de-Stock.exe" (
    echo Demarrage de l'application...
    start "" "build\bin\Gestion-de-Stock.exe"
    exit /b 0
)

pause