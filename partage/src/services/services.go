package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"exellence-materiaux/internal/auth"
	"exellence-materiaux/internal/factures"
	"exellence-materiaux/internal/parametres"
	"exellence-materiaux/internal/produits"
	"exellence-materiaux/internal/stock"
)

type Session struct {
	Role      string
	ExpiresAt time.Time
}

type AuthService struct {
	store    *auth.Store
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewAuthService(database *sql.DB) *AuthService {
	return &AuthService{
		store:    auth.NewStore(database),
		sessions: make(map[string]Session),
	}
}

func (s *AuthService) IsConfigured() (bool, error) {
	return s.store.IsConfigured()
}

func (s *AuthService) Setup(adminPassword, userPassword string) error {
	return s.store.SetupPasswords(adminPassword, userPassword)
}

func (s *AuthService) Login(role, password string) (string, time.Time, error) {
	if err := s.store.Authenticate(role, password); err != nil {
		return "", time.Time{}, err
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", time.Time{}, fmt.Errorf("generation token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour)

	s.mu.Lock()
	s.sessions[token] = Session{Role: role, ExpiresAt: expiresAt}
	s.mu.Unlock()

	return token, expiresAt, nil
}

func (s *AuthService) ValidateToken(token string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[token]
	if !exists {
		return "", false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return "", false
	}
	return session.Role, true
}

func (s *AuthService) Logout(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func (s *AuthService) ChangePasswords(currentAdminPassword, adminPassword, userPassword string) error {
	return s.store.ChangePasswords(currentAdminPassword, adminPassword, userPassword)
}

type ProduitService struct {
	store *produits.Store
}

func NewProduitService(database *sql.DB) *ProduitService {
	return &ProduitService{store: produits.NewStore(database)}
}

func (s *ProduitService) List(filter produits.FiltreProduits) (produits.PageProduits, error) {
	return s.store.List(filter)
}

func (s *ProduitService) GetByID(id string) (produits.Produit, error) {
	return s.store.GetByID(id)
}

func (s *ProduitService) GetByCodeBarres(code string) (produits.Produit, error) {
	return s.store.GetByCodeBarres(code)
}

func (s *ProduitService) Create(np produits.NouveauProduit) (produits.Produit, error) {
	return s.store.Create(np)
}

func (s *ProduitService) Update(mp produits.ModificationProduit) (produits.Produit, error) {
	return s.store.Update(mp)
}

func (s *ProduitService) Delete(id string) error {
	return s.store.Delete(id)
}

type FactureService struct {
	store      *factures.Store
	projectDir string
}

func NewFactureService(database *sql.DB, projectDir string) *FactureService {
	return &FactureService{
		store:      factures.NewStore(database),
		projectDir: projectDir,
	}
}

func (s *FactureService) Create(nf factures.NouvelleFacture) (factures.DetailFacture, error) {
	return s.store.CreerFacture(nf)
}

func (s *FactureService) GetByID(id string) (factures.DetailFacture, error) {
	return s.store.GetByID(id)
}

func (s *FactureService) List(filter factures.FiltreFactures) (factures.PageFactures, error) {
	return s.store.List(filter)
}

func (s *FactureService) GeneratePDF(id string) (string, string, error) {
	detail, err := s.store.GetByID(id)
	if err != nil {
		return "", "", err
	}
	logoPath := filepath.Join(s.projectDir, "logo.png")
	if _, err := os.Stat(logoPath); err != nil {
		logoPath = ""
	}

	tempDir := filepath.Join(s.projectDir, "temp_pdf")
	_ = os.MkdirAll(tempDir, 0755)
	filename := fmt.Sprintf("Facture_%s.pdf", detail.Facture.Numero)
	outPath := filepath.Join(tempDir, filename)

	if err := factures.GenererFacturePDF(detail, logoPath, outPath); err != nil {
		return "", "", fmt.Errorf("generer pdf facture: %w", err)
	}
	return outPath, filename, nil
}

type BonService struct {
	store      *factures.Store
	projectDir string
}

func NewBonService(database *sql.DB, projectDir string) *BonService {
	return &BonService{
		store:      factures.NewStore(database),
		projectDir: projectDir,
	}
}

func (s *BonService) Create(nb factures.NouveauBon) (factures.DetailBon, error) {
	return s.store.CreerBon(nb)
}

func (s *BonService) GetByID(id string) (factures.DetailBon, error) {
	return s.store.GetBonByID(id)
}

func (s *BonService) List(filter factures.FiltreBons) (factures.PageBons, error) {
	return s.store.ListBons(filter)
}

func (s *BonService) GeneratePDF(id string) (string, string, error) {
	detail, err := s.store.GetBonByID(id)
	if err != nil {
		return "", "", err
	}
	logoPath := filepath.Join(s.projectDir, "logo.png")
	if _, err := os.Stat(logoPath); err != nil {
		logoPath = ""
	}

	tempDir := filepath.Join(s.projectDir, "temp_pdf")
	_ = os.MkdirAll(tempDir, 0755)
	filename := fmt.Sprintf("Bon_%s.pdf", detail.Bon.Numero)
	outPath := filepath.Join(tempDir, filename)

	if err := factures.GenererBonPDF(detail, logoPath, outPath); err != nil {
		return "", "", fmt.Errorf("generer pdf bon: %w", err)
	}
	return outPath, filename, nil
}

type StockService struct {
	store *stock.Store
}

func NewStockService(database *sql.DB) *StockService {
	return &StockService{store: stock.NewStore(database)}
}

func (s *StockService) CreateMouvement(sm stock.SaisieMouvement) (stock.MouvementStock, error) {
	return s.store.CreerMouvement(sm)
}

func (s *StockService) Ajuster(sa stock.SaisieAjustement) (stock.MouvementStock, error) {
	return s.store.AjusterStock(sa)
}

func (s *StockService) ListMouvements(filter stock.FiltreMouvements) (stock.PageMouvements, error) {
	return s.store.List(filter)
}

func (s *StockService) DeleteMouvement(id string) error {
	return s.store.SupprimerMouvement(id)
}

func (s *StockService) GetAlertes() ([]stock.AlerteStock, error) {
	return s.store.ObtenirAlertes()
}

type ParametresService struct {
	db         *sql.DB
	projectDir string
}

func NewParametresService(database *sql.DB, projectDir string) *ParametresService {
	return &ParametresService{
		db:         database,
		projectDir: projectDir,
	}
}

func (s *ParametresService) Get() (parametres.ParametresEntreprise, error) {
	var p parametres.ParametresEntreprise
	query := `SELECT id, denomination, forme_juridique, capital, adresse, code_postal, ville,
		siren, siret, rcs, numero_tva, telephone, email, site_web, logo_path, devise,
		taux_tva_defaut, taux_tva_disponibles, mentions_legales, conditions_reglement,
		format_numerotation_facture, format_numerotation_bon, updated_at
		FROM parametres_entreprise LIMIT 1`

	var updatedAtStr string
	err := s.db.QueryRow(query).Scan(
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

func (s *ParametresService) Update(mp parametres.ModificationParametres) (parametres.ParametresEntreprise, error) {
	query := `UPDATE parametres_entreprise SET
		denomination = ?, forme_juridique = ?, capital = ?, adresse = ?, code_postal = ?, ville = ?,
		siren = ?, siret = ?, rcs = ?, numero_tva = ?, telephone = ?, email = ?, site_web = ?,
		devise = ?, taux_tva_defaut = ?, mentions_legales = ?, conditions_reglement = ?,
		updated_at = datetime('now')
		WHERE id = '01DEFAULTCONFIG00000000000'`

	_, err := s.db.Exec(query,
		mp.Denomination, mp.FormeJuridique, mp.Capital, mp.Adresse, mp.CodePostal, mp.Ville,
		mp.Siren, mp.Siret, mp.RCS, mp.NumeroTVA, mp.Telephone, mp.Email, mp.SiteWeb,
		mp.Devise, mp.TauxTVADefaut, mp.MentionsLegales, mp.ConditionsReglement,
	)
	if err != nil {
		return parametres.ParametresEntreprise{}, fmt.Errorf("mise a jour parametres: %w", err)
	}
	return s.Get()
}

func (s *ParametresService) SaveLogo(data []byte) (string, error) {
	logoPath := filepath.Join(s.projectDir, "logo.png")
	if err := os.WriteFile(logoPath, data, 0644); err != nil {
		return "", fmt.Errorf("ecriture logo: %w", err)
	}

	_, err := s.db.Exec(`UPDATE parametres_entreprise
		SET logo_path = ?, updated_at = datetime('now')
		WHERE id = '01DEFAULTCONFIG00000000000'`, logoPath)
	if err != nil {
		return "", fmt.Errorf("enregistrement du chemin du logo: %w", err)
	}

	return logoPath, nil
}

type ServiceManager struct {
	DB         *sql.DB
	ProjectDir string
	Auth       *AuthService
	Produits   *ProduitService
	Factures   *FactureService
	Bons       *BonService
	Stock      *StockService
	Parametres *ParametresService
}

func NewServiceManager(database *sql.DB, projectDir string) *ServiceManager {
	return &ServiceManager{
		DB:         database,
		ProjectDir: projectDir,
		Auth:       NewAuthService(database),
		Produits:   NewProduitService(database),
		Factures:   NewFactureService(database, projectDir),
		Bons:       NewBonService(database, projectDir),
		Stock:      NewStockService(database),
		Parametres: NewParametresService(database, projectDir),
	}
}
