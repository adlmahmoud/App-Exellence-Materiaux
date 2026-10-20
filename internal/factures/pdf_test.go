package factures

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenererFacturePDF(t *testing.T) {
	tempDir := t.TempDir()
	outPDF := filepath.Join(tempDir, "facture_test.pdf")

	detail := DetailFacture{
		Facture: Facture{
			Numero:                   "FAC-2026-00001",
			DateFacture:              "20/10/2026",
			DateEcheance:             "20/11/2026",
			ModeReglement:            ReglementCB,
			ClientNom:                "Société Générale de Bâtiment",
			ClientAdresse:            "14 Rue des Artisans",
			ClientCodePostal:         "95500",
			ClientVille:              "Gonesse",
			ClientSirenSiret:         "845 123 678 00019",
			TotalHT:                  15800,
			TotalTVA:                 3160,
			TotalTTC:                 18960,
			ConditionsReglement:      "Paiement par virement ou carte bancaire à 30 jours.",
			EntrepriseDenomination:   "EXCELLENCE MATÉRIAUX",
			EntrepriseFormeJuridique: "SARL",
			EntrepriseCapital:        3000000,
			EntrepriseSiren:          "109 741 124",
			EntrepriseRCS:            "Pontoise",
			EntrepriseNumeroTVA:      "FR49 109741124",
		},
		Lignes: []LigneFacture{
			{
				Designation:    "Sac de ciment gris 35kg NF",
				Quantite:       10,
				PrixUnitaireHT: 650,
				TauxTVA:        2000,
				MontantHT:      6500,
				MontantTVA:     1300,
				MontantTTC:     7800,
			},
			{
				Designation:    "Planche de coffrage 27x200mm L.4m",
				Quantite:       8,
				PrixUnitaireHT: 850,
				TauxTVA:        2000,
				MontantHT:      6800,
				MontantTVA:     1360,
				MontantTTC:     8160,
			},
			{
				Designation:    "Brouette de chantier 100L roue increvable",
				Quantite:       1,
				PrixUnitaireHT: 2500,
				TauxTVA:        2000,
				MontantHT:      2500,
				MontantTVA:     500,
				MontantTTC:     3000,
			},
		},
	}

	logoPath := filepath.Join("..", "..", "logo.png")
	if _, err := os.Stat(logoPath); err != nil {
		logoPath = ""
	}

	err := GenererFacturePDF(detail, logoPath, outPDF)
	if err != nil {
		t.Fatalf("GenererFacturePDF failed: %v", err)
	}

	samplePath := filepath.Join("..", "..", "Exemple_Facture_Generee.pdf")
	_ = GenererFacturePDF(detail, logoPath, samplePath)
}
