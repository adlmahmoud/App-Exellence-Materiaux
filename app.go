package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"

	"exellence-materiaux/internal/db"
	"exellence-materiaux/internal/factures"
	"exellence-materiaux/internal/parametres"
	"exellence-materiaux/internal/produits"
	"exellence-materiaux/internal/stock"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"image/png"
)

type App struct {
	ctx           context.Context
	database      *db.DB
	dataDir       string
	produitsStore *produits.Store
	stockStore    *stock.Store
	facturesStore *factures.Store
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

func resolveDataDir(projectDir, userConfigDir string) (string, error) {
	legacyDBPath := filepath.Join(projectDir, "exellence.db")
	if info, err := os.Stat(legacyDBPath); err == nil && !info.IsDir() {
		return projectDir, nil
	} else if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("verification de l'ancienne base de donnees: %w", err)
	}

	if userConfigDir == "" {
		return "", fmt.Errorf("repertoire de donnees utilisateur indisponible")
	}

	dataDir := filepath.Join(userConfigDir, "Gestion de Stock")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return "", fmt.Errorf("creation du repertoire de donnees: %w", err)
	}

	return dataDir, nil
}

func getDataDir() (string, error) {
	projectDir := getProjectDir()
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		if _, statErr := os.Stat(filepath.Join(projectDir, "exellence.db")); statErr == nil {
			return projectDir, nil
		}
		return "", fmt.Errorf("repertoire de configuration utilisateur indisponible: %w", err)
	}
	return resolveDataDir(projectDir, userConfigDir)
}

func NewApp() *App {
	a := &App{}
	dataDir, err := getDataDir()
	if err != nil {
		println(fmt.Sprintf("application data directory error: %v", err))
		return a
	}
	a.dataDir = dataDir
	dbPath := filepath.Join(a.dataDir, "exellence.db")
	database, err := db.Open(dbPath)
	if err == nil {
		a.database = database
		a.produitsStore = produits.NewStore(database.DB)
		a.stockStore = stock.NewStore(database.DB)
		a.facturesStore = factures.NewStore(database.DB)
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
		"nom":     "Gestion de Stock",
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

func (a *App) ImporterLogo() (string, error) {
	if a.database == nil {
		return "", fmt.Errorf("base de données non initialisée")
	}
	if a.ctx == nil {
		return "", fmt.Errorf("interface de l'application non initialisée")
	}

	selectedPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Choisir le logo des factures et bons de livraison",
		DefaultDirectory: a.dataDir,
		Filters: []runtime.FileFilter{{
			DisplayName: "Images PNG ou JPEG",
			Pattern:     "*.png;*.jpg;*.jpeg",
		}},
	})
	if err != nil {
		return "", fmt.Errorf("sélection du logo: %w", err)
	}
	if selectedPath == "" {
		return "", nil
	}

	source, err := os.ReadFile(selectedPath)
	if err != nil {
		return "", fmt.Errorf("lecture du logo: %w", err)
	}
	if len(source) > 10<<20 {
		return "", fmt.Errorf("le logo ne doit pas dépasser 10 Mo")
	}

	config, _, err := image.DecodeConfig(bytes.NewReader(source))
	if err != nil {
		return "", fmt.Errorf("image de logo invalide: %w", err)
	}
	if config.Width <= 0 || config.Height <= 0 || config.Width > 4096 || config.Height > 4096 {
		return "", fmt.Errorf("les dimensions du logo ne doivent pas dépasser 4096 × 4096 pixels")
	}

	img, _, err := image.Decode(bytes.NewReader(source))
	if err != nil {
		return "", fmt.Errorf("décodage du logo: %w", err)
	}
	var encoded bytes.Buffer
	if err := png.Encode(&encoded, img); err != nil {
		return "", fmt.Errorf("conversion du logo en PNG: %w", err)
	}

	logoPath := filepath.Join(a.dataDir, "logo_facture.png")
	if err := os.WriteFile(logoPath, encoded.Bytes(), 0600); err != nil {
		return "", fmt.Errorf("enregistrement du logo: %w", err)
	}

	_, err = a.database.Exec(`UPDATE parametres_entreprise
		SET logo_path = ?, updated_at = datetime('now')
		WHERE id = '01DEFAULTCONFIG00000000000'`, logoPath)
	if err != nil {
		return "", fmt.Errorf("enregistrement du chemin du logo: %w", err)
	}

	return logoPath, nil
}

func (a *App) UpdateParametres(input parametres.ModificationParametres) error {
	if a.database == nil {
		return fmt.Errorf("base de données non initialisée")
	}

	query := `UPDATE parametres_entreprise SET
		denomination = ?, forme_juridique = ?, capital = ?, adresse = ?, code_postal = ?, ville = ?,
		siren = ?, siret = ?, rcs = ?, numero_tva = ?, telephone = ?, email = ?, site_web = ?,
		devise = ?, taux_tva_defaut = ?, mentions_legales = ?, conditions_reglement = ?,
		updated_at = datetime('now')
		WHERE id = '01DEFAULTCONFIG00000000000'`

	_, err := a.database.Exec(query,
		input.Denomination, input.FormeJuridique, input.Capital, input.Adresse, input.CodePostal, input.Ville,
		input.Siren, input.Siret, input.RCS, input.NumeroTVA, input.Telephone, input.Email, input.SiteWeb,
		input.Devise, input.TauxTVADefaut, input.MentionsLegales, input.ConditionsReglement,
	)
	if err != nil {
		return fmt.Errorf("mise à jour parametres: %w", err)
	}
	return nil
}

func (a *App) SauvegarderBase(dossierDestination string) (string, error) {
	if a.database == nil {
		return "", fmt.Errorf("base de données non initialisée")
	}

	targetDir := dossierDestination
	if targetDir == "" || targetDir == "." {
		targetDir = a.dataDir
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

func (a *App) CreerMouvementStock(input stock.SaisieMouvement) (stock.MouvementStock, error) {
	if a.stockStore == nil {
		return stock.MouvementStock{}, fmt.Errorf("base de données non initialisée")
	}
	return a.stockStore.CreerMouvement(input)
}

func (a *App) AjusterStock(input stock.SaisieAjustement) (stock.MouvementStock, error) {
	if a.stockStore == nil {
		return stock.MouvementStock{}, fmt.Errorf("base de données non initialisée")
	}
	return a.stockStore.AjusterStock(input)
}

func (a *App) ListerMouvementsStock(filtre stock.FiltreMouvements) (stock.PageMouvements, error) {
	if a.stockStore == nil {
		return stock.PageMouvements{}, fmt.Errorf("base de données non initialisée")
	}
	return a.stockStore.List(filtre)
}

func (a *App) ObtenirAlertesStock() ([]stock.AlerteStock, error) {
	if a.stockStore == nil {
		return nil, fmt.Errorf("base de données non initialisée")
	}
	return a.stockStore.ObtenirAlertes()
}

func (a *App) CreerFacture(input factures.NouvelleFacture) (factures.DetailFacture, error) {
	if a.facturesStore == nil {
		return factures.DetailFacture{}, fmt.Errorf("base de données non initialisée")
	}
	detail, err := a.facturesStore.CreerFacture(input)
	if err != nil {
		return detail, err
	}

	settings, err := a.GetParametres()
	if err != nil {
		return detail, fmt.Errorf("recuperation du logo: %w", err)
	}
	pdfPath := filepath.Join(a.dataDir, fmt.Sprintf("Facture_%s.pdf", detail.Facture.Numero))
	if err := factures.GenererFacturePDF(detail, settings.LogoPath, pdfPath); err != nil {
		return detail, fmt.Errorf("generation pdf: %w", err)
	}

	return detail, nil
}

func (a *App) GetFacture(id string) (factures.DetailFacture, error) {
	if a.facturesStore == nil {
		return factures.DetailFacture{}, fmt.Errorf("base de données non initialisée")
	}
	return a.facturesStore.GetByID(id)
}

func (a *App) ListerFactures(filtre factures.FiltreFactures) (factures.PageFactures, error) {
	if a.facturesStore == nil {
		return factures.PageFactures{}, fmt.Errorf("base de données non initialisée")
	}
	return a.facturesStore.List(filtre)
}

func (a *App) ImprimerFacturePDF(id string) (string, error) {
	if a.facturesStore == nil {
		return "", fmt.Errorf("base de données non initialisée")
	}
	detail, err := a.facturesStore.GetByID(id)
	if err != nil {
		return "", err
	}

	settings, err := a.GetParametres()
	if err != nil {
		return "", fmt.Errorf("recuperation du logo: %w", err)
	}
	pdfPath := filepath.Join(a.dataDir, fmt.Sprintf("Facture_%s.pdf", detail.Facture.Numero))
	if err := factures.GenererFacturePDF(detail, settings.LogoPath, pdfPath); err != nil {
		return "", fmt.Errorf("generation pdf: %w", err)
	}

	_ = exec.Command("cmd", "/c", "start", "", pdfPath).Start()

	return pdfPath, nil
}

func (a *App) CreerBon(input factures.NouveauBon) (factures.DetailBon, error) {
	if a.facturesStore == nil {
		return factures.DetailBon{}, fmt.Errorf("base de données non initialisée")
	}
	detail, err := a.facturesStore.CreerBon(input)
	if err != nil {
		return detail, err
	}

	settings, err := a.GetParametres()
	if err != nil {
		return detail, fmt.Errorf("recuperation du logo: %w", err)
	}
	pdfPath := filepath.Join(a.dataDir, fmt.Sprintf("Bon_%s.pdf", detail.Bon.Numero))
	if err := factures.GenererBonPDF(detail, settings.LogoPath, pdfPath); err != nil {
		return detail, fmt.Errorf("generation pdf bon: %w", err)
	}

	return detail, nil
}

func (a *App) GetBon(id string) (factures.DetailBon, error) {
	if a.facturesStore == nil {
		return factures.DetailBon{}, fmt.Errorf("base de données non initialisée")
	}
	return a.facturesStore.GetBonByID(id)
}

func (a *App) ListerBons(filtre factures.FiltreBons) (factures.PageBons, error) {
	if a.facturesStore == nil {
		return factures.PageBons{}, fmt.Errorf("base de données non initialisée")
	}
	return a.facturesStore.ListBons(filtre)
}

func (a *App) ImprimerBonPDF(id string) (string, error) {
	if a.facturesStore == nil {
		return "", fmt.Errorf("base de données non initialisée")
	}
	detail, err := a.facturesStore.GetBonByID(id)
	if err != nil {
		return "", err
	}

	settings, err := a.GetParametres()
	if err != nil {
		return "", fmt.Errorf("recuperation du logo: %w", err)
	}
	pdfPath := filepath.Join(a.dataDir, fmt.Sprintf("Bon_%s.pdf", detail.Bon.Numero))
	if err := factures.GenererBonPDF(detail, settings.LogoPath, pdfPath); err != nil {
		return "", fmt.Errorf("generation pdf bon: %w", err)
	}

	_ = exec.Command("cmd", "/c", "start", "", pdfPath).Start()

	return pdfPath, nil
}
