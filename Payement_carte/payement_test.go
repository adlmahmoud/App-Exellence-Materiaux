package main

import (
	"context"
	"testing"
	"time"

	"exellence-materiaux/Payement_carte/controllers"
	"exellence-materiaux/Payement_carte/drivers"
	"exellence-materiaux/Payement_carte/models"
	"exellence-materiaux/Payement_carte/services"
)

func TestSimulateurSucces(t *testing.T) {
	cfg := models.ConfigTPE{
		Mode:               "simulateur",
		SimulateurScenario: "succes",
		SimulateurDelaiSec: 1,
	}
	driver := drivers.NewSimulateurDriver(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := models.DemandePaiement{
		MontantCentimes: 4550,
		MontantEuros:    45.50,
		Reference:       "FAC-001",
		NumeroCaisse:    1,
	}

	resp, err := driver.InitierPaiement(ctx, req)
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if resp.Statut != models.StatutValide {
		t.Fatalf("statut attendu VALIDE, reçu %s", resp.Statut)
	}
	if resp.CodeReponse != "00" {
		t.Fatalf("code réponse attendu 00, reçu %s", resp.CodeReponse)
	}
}

func TestSimulateurRefus(t *testing.T) {
	cfg := models.ConfigTPE{
		Mode:               "simulateur",
		SimulateurScenario: "refus",
		SimulateurDelaiSec: 1,
	}
	driver := drivers.NewSimulateurDriver(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := models.DemandePaiement{
		MontantCentimes: 12000,
		MontantEuros:    120.00,
	}

	resp, err := driver.InitierPaiement(ctx, req)
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if resp.Statut != models.StatutRefuse {
		t.Fatalf("statut attendu REFUSE, reçu %s", resp.Statut)
	}
	if resp.CodeReponse != "05" {
		t.Fatalf("code réponse attendu 05, reçu %s", resp.CodeReponse)
	}
}

func TestSimulateurAnnulation(t *testing.T) {
	cfg := models.ConfigTPE{
		Mode:               "simulateur",
		SimulateurScenario: "annulation",
		SimulateurDelaiSec: 1,
	}
	driver := drivers.NewSimulateurDriver(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req := models.DemandePaiement{
		MontantCentimes: 500,
		MontantEuros:    5.00,
	}

	resp, err := driver.InitierPaiement(ctx, req)
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if resp.Statut != models.StatutAnnule {
		t.Fatalf("statut attendu ANNULE, reçu %s", resp.Statut)
	}
}

func TestTPEServiceEtController(t *testing.T) {
	tmpDir := t.TempDir()
	srv := services.NewTPEService(tmpDir)

	cfg := srv.GetConfig()
	if cfg.Mode != "simulateur" {
		t.Fatalf("mode attendu simulateur, reçu %s", cfg.Mode)
	}

	cfg.SimulateurDelaiSec = 1
	cfg.SimulateurScenario = "succes"
	if err := srv.SauvegarderConfig(cfg); err != nil {
		t.Fatalf("sauvegarde config: %v", err)
	}

	ctrl := controllers.NewTPEController(srv)
	testConn := ctrl.TesterConnexion()
	if !testConn.Succes {
		t.Fatalf("test connexion simulateur échoué: %s", testConn.Message)
	}

	resp := ctrl.LancerPaiement(25.50, "REF-TEST")
	if resp.Statut != models.StatutValide {
		t.Fatalf("statut attendu VALIDE, reçu %s", resp.Statut)
	}

	ctrl.ChangerScenarioSimulateur("refus")
	respRefus := ctrl.LancerPaiement(10.0, "REF-REFUS")
	if respRefus.Statut != models.StatutRefuse {
		t.Fatalf("statut attendu REFUSE, reçu %s", respRefus.Statut)
	}
}

func TestCalculLRC(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03}
	var lrc byte
	for _, b := range data {
		lrc ^= b
	}
	if lrc != 0x00 {
		t.Fatalf("calcul LRC attendu 0, reçu %x", lrc)
	}
}
