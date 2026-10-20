package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"exellence-materiaux/internal/db"
	"exellence-materiaux/internal/parametres"
	"exellence-materiaux/internal/produits"
)

type App struct {
	ctx           context.Context
	database      *db.DB
	projectDir    string
	produitsStore *produits.Store
}

func getProjectDir() string {
	execPath, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(execPath)
		if filepath.Base(dir) == "bin" {
			parent := filepath.Dir(dir)
			if filepath.Base(parent) == "build" {
				return filepath.Dir(parent)
			}
		}
		return dir
	}
	cwd, err := os.Getwd()
	if err == nil {
		return cwd
	}
	return "."
}

func NewApp() *App {
	a := &App{}
	a.projectDir = getProjectDir()
	dbPath := filepath.Join(a.projectDir, "exellence.db")
	database, err := db.Open(dbPath)
	if err == nil {
		a.database = database
		a.produitsStore = produits.NewStore(database.DB)
	} else {
		println(fmt.Sprintf("database open error: %v", err))
	}
	return a
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	if a.database != nil {
		a.database.Close()
	}
}

func (a *App) GetAPropos() map[string]string {
	return map[string]string{
		"nom":     "Exellence Materiaux",
		"version": "1.0.0",
		"annee":   "2026",
		"auteurs": "Adel Mahmoud et Sidou",
	}
}

func (a *App) GetParametres() (parametres.ParametresEntreprise, error) {
	var p parametres.ParametresEntreprise
	if a.database == nil {
		return p, fmt.Errorf("base de données non initialisée")
	}

	query := `SELECT id, denomination, forme_juridique, capital, adresse, code_postal, ville,
		siren, siret, rcs, numero_tva, telephone, email, site_web, logo_path, devise,
		taux_tva_defaut, taux_tva_disponibles, mentions_legales, conditions_reglement,
		format_numerotation_facture, format_numerotation_bon, updated_at
		FROM parametres_entreprise LIMIT 1`

	var updatedAtStr string
	err := a.database.QueryRow(query).Scan(
		&p.ID, &p.Denomination, &p.FormeJuridique, &p.Capital, &p.Adresse, &p.CodePostal, &p.Ville,
		&p.Siren, &p.Siret, &p.RCS, &p.NumeroTVA, &p.Telephone, &p.Email, &p.SiteWeb, &p.LogoPath, &p.Devise,
		&p.TauxTVADefaut, &p.TauxTVADisponibles, &p.MentionsLegales, &p.ConditionsReglement,
		&p.FormatNumerotationFacture, &p.FormatNumerotationBon, &updatedAtStr,
	)
	if err != nil {
		return p, fmt.Errorf("recuperation parametres: %w", err)
	}

	return p, nil
}

func (a *App) SauvegarderBase(dossierDestination string) (string, error) {
	if a.database == nil {
		return "", fmt.Errorf("base de données non initialisée")
	}

	targetDir := dossierDestination
	if targetDir == "" || targetDir == "." {
		targetDir = a.projectDir
	}

	destPath := filepath.Join(targetDir, "exellence_backup.db")
	if err := a.database.Backup(destPath); err != nil {
		return "", fmt.Errorf("sauvegarde base: %w", err)
	}

	return destPath, nil
}

func (a *App) CreerProduit(input produits.NouveauProduit) (produits.Produit, error) {
	if a.produitsStore == nil {
		return produits.Produit{}, fmt.Errorf("base de données non initialisée")
	}
	return a.produitsStore.Create(input)
}

func (a *App) ModifierProduit(input produits.ModificationProduit) (produits.Produit, error) {
	if a.produitsStore == nil {
		return produits.Produit{}, fmt.Errorf("base de données non initialisée")
	}
	return a.produitsStore.Update(input)
}

func (a *App) SupprimerProduit(id string) error {
	if a.produitsStore == nil {
		return fmt.Errorf("base de données non initialisée")
	}
	return a.produitsStore.Delete(id)
}

func (a *App) GetProduitParCodeBarres(codeBarres string) (produits.Produit, error) {
	if a.produitsStore == nil {
		return produits.Produit{}, fmt.Errorf("base de données non initialisée")
	}
	return a.produitsStore.GetByCodeBarres(codeBarres)
}

func (a *App) ListerProduits(filtre produits.FiltreProduits) (produits.PageProduits, error) {
	if a.produitsStore == nil {
		return produits.PageProduits{}, fmt.Errorf("base de données non initialisée")
	}
	return a.produitsStore.List(filtre)
}
