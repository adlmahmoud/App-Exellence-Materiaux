package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"exellence-materiaux/internal/db"
	"exellence-materiaux/partage/src/controllers"
	"exellence-materiaux/partage/src/models"
	"exellence-materiaux/partage/src/routers"
	"exellence-materiaux/partage/src/services"
)

func TestPartageServerAPI(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_partage.db")

	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open database: %v", err)
	}
	defer database.Close()

	sm := services.NewServiceManager(database.DB, tempDir)
	cm := controllers.NewControllerManager(sm)
	router := routers.NewRouter(cm)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Health check code = %d, expected 200", w.Code)
	}

	setupPayload := `{"adminPassword":"admin-password-123","userPassword":"user-password-123"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/auth/setup", bytes.NewBufferString(setupPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Auth setup code = %d (body: %s), expected 201", w.Code, w.Body.String())
	}

	loginPayload := `{"role":"admin","password":"admin-password-123"}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(loginPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Auth login code = %d, expected 200", w.Code)
	}

	var loginResp struct {
		Success bool                 `json:"success"`
		Data    models.LoginResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("Unmarshal login response: %v", err)
	}
	adminToken := loginResp.Data.Token
	if adminToken == "" {
		t.Fatalf("Login token is empty")
	}

	nouveauProduit := `{
		"codeBarres": "3760000000001",
		"reference": "TEST-CIM",
		"designation": "Sac ciment test 35kg",
		"prixAchatHt": 500,
		"prixVenteHt": 700,
		"tauxTva": 2000,
		"stockInitial": 50,
		"stockMinimum": 10
	}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/produits", bytes.NewBufferString(nouveauProduit))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create product code = %d (body: %s), expected 201", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/produits?recherche=ciment", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("List products code = %d, expected 200", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/produits/get?code=3760000000001", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Get product by code = %d, expected 200", w.Code)
	}

	nouvelleFacture := `{
		"clientNom": "Client Chantier Test",
		"modeReglement": "especes",
		"dateFacture": "25/10/2026",
		"lignes": [
			{
				"designation": "Sac ciment test 35kg",
				"quantite": 2,
				"prixUnitaireHt": 700,
				"tauxTva": 2000
			}
		]
	}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/factures", bytes.NewBufferString(nouvelleFacture))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create facture code = %d (body: %s), expected 201", w.Code, w.Body.String())
	}

	nouveauBon := `{
		"clientNom": "Destinataire Chantier Test",
		"typeBon": "livraison",
		"dateBon": "25/10/2026",
		"lignes": [
			{
				"designation": "Sac ciment test 35kg",
				"quantite": 5,
				"prixUnitaireHt": 700,
				"tauxTva": 2000
			}
		]
	}`
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/api/bons", bytes.NewBufferString(nouveauBon))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Create bon code = %d (body: %s), expected 201", w.Code, w.Body.String())
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/stock/alertes", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Stock alertes code = %d, expected 200", w.Code)
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/parametres", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Parametres get code = %d, expected 200", w.Code)
	}
}

func TestPartageServerMainFindsLocalIPs(t *testing.T) {
	ips := getLocalIPs()
	for _, ip := range ips {
		if ip == "" {
			t.Fatal("empty ip found")
		}
	}
}
