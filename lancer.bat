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

where wails >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo Demarrage en mode developpement (Wails)...
    wails dev
    exit /b 0
)

if exist "%USERPROFILE%\go\bin\wails.exe" (
    echo Demarrage en mode developpement (Wails)...
    "%USERPROFILE%\go\bin\wails.exe" dev
    exit /b 0
)

echo [INFO] Wails CLI n'est pas encore detecte dans le PATH.
echo Installation automatique en cours...
go install github.com/wailsapp/wails/v2/cmd/wails@latest
if exist "%USERPROFILE%\go\bin\wails.exe" (
    "%USERPROFILE%\go\bin\wails.exe" dev
    exit /b 0
)

echo [ERREUR] Impossible de lancer l'application.
pause
