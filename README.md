# Exellence Materiaux

Logiciel de bureau de gestion de stock et facturation pour Exellence Materiaux.

## Stack technique

- Backend : Go, modernc.org/sqlite (sans CGO), go-pdf/fpdf
- Interface : Wails v2, Svelte
- Base de données : SQLite en mode WAL, FTS5 pour la recherche de produits

## Prérequis

- Go >= 1.25
- Node.js & npm
- Wails CLI

## Développement

```bash
wails dev
```

## Compilation

```bash
wails build
```
