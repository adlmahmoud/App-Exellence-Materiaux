package factures

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
)

func GenererFacturePDF(detail DetailFacture, logoPath string, cheminDestination string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 25)
	pdf.AddPage()

	if logoPath != "" {
		if _, err := os.Stat(logoPath); err == nil {
			pdf.ImageOptions(logoPath, 15, 12, 55, 0, false, fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
		}
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(30, 41, 59)
	pdf.SetXY(110, 20)
	pdf.CellFormat(85, 6, fmt.Sprintf("FACTURE N°: %s", detail.Facture.Numero), "", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetXY(110, 27)
	pdf.CellFormat(85, 6, fmt.Sprintf("Date: %s", detail.Facture.DateFacture), "", 1, "R", false, 0, "")

	pdf.SetDrawColor(217, 119, 6)
	pdf.SetLineWidth(0.8)
	pdf.Line(15, 42, 195, 42)

	pdf.SetXY(15, 48)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(30, 41, 59)
	pdf.CellFormat(90, 6, "FACTURÉ À", "", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(51, 65, 85)
	pdf.SetX(15)
	pdf.CellFormat(90, 5, fmt.Sprintf("Client : %s", detail.Facture.ClientNom), "", 1, "L", false, 0, "")

	adresseClient := detail.Facture.ClientAdresse
	if detail.Facture.ClientCodePostal != "" || detail.Facture.ClientVille != "" {
		if adresseClient != "" {
			adresseClient += " "
		}
		adresseClient += fmt.Sprintf("%s %s", detail.Facture.ClientCodePostal, detail.Facture.ClientVille)
	}
	if adresseClient == "" {
		adresseClient = "-"
	}
	pdf.SetX(15)
	pdf.CellFormat(90, 5, fmt.Sprintf("Adresse : %s", adresseClient), "", 1, "L", false, 0, "")

	siretClient := detail.Facture.ClientSirenSiret
	if siretClient == "" {
		siretClient = "[facultatif]"
	}
	pdf.SetX(15)
	pdf.CellFormat(90, 5, fmt.Sprintf("SIRET (si professionnel) : %s", siretClient), "", 1, "L", false, 0, "")

	startY := 75.0
	pdf.SetXY(15, startY)

	pdf.SetFillColor(74, 123, 176)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetDrawColor(180, 205, 230)
	pdf.SetLineWidth(0.2)
	pdf.SetFont("Arial", "B", 10)

	pdf.CellFormat(80, 8, "Désignation", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 8, "Qté", "1", 0, "C", true, 0, "")
	pdf.CellFormat(26, 8, "PU HT", "1", 0, "R", true, 0, "")
	pdf.CellFormat(22, 8, "TVA", "1", 0, "C", true, 0, "")
	pdf.CellFormat(32, 8, "Total HT", "1", 1, "R", true, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(30, 41, 59)

	for i, ligne := range detail.Lignes {
		isEven := (i % 2) == 0
		if isEven {
			pdf.SetFillColor(235, 242, 250)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		puHTStr := fmt.Sprintf("%.2f €", float64(ligne.PrixUnitaireHT)/100.0)
		tvaStr := fmt.Sprintf("%.1f %%", float64(ligne.TauxTVA)/100.0)
		if (ligne.TauxTVA % 100) == 0 {
			tvaStr = fmt.Sprintf("%d %%", ligne.TauxTVA/100)
		}
		totalHTStr := fmt.Sprintf("%.2f €", float64(ligne.MontantHT)/100.0)

		pdf.SetX(15)
		pdf.CellFormat(80, 7, ligne.Designation, "1", 0, "L", true, 0, "")
		pdf.CellFormat(20, 7, fmt.Sprintf("%d", ligne.Quantite), "1", 0, "C", true, 0, "")
		pdf.CellFormat(26, 7, puHTStr, "1", 0, "R", true, 0, "")
		pdf.CellFormat(22, 7, tvaStr, "1", 0, "C", true, 0, "")
		pdf.CellFormat(32, 7, totalHTStr, "1", 1, "R", true, 0, "")
	}

	pdf.Ln(4)
	totalsY := pdf.GetY()
	if totalsY > 215 {
		pdf.AddPage()
		totalsY = 20
	}

	pdf.SetXY(110, totalsY)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(50, 6, "Total HT :", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, fmt.Sprintf("%.2f €", float64(detail.Facture.TotalHT)/100.0), "", 1, "R", false, 0, "")

	pdf.SetX(110)
	pdf.CellFormat(50, 6, "TVA :", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 6, fmt.Sprintf("%.2f €", float64(detail.Facture.TotalTVA)/100.0), "", 1, "R", false, 0, "")

	pdf.SetX(110)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(50, 8, "TOTAL TTC :", "", 0, "R", false, 0, "")
	pdf.CellFormat(35, 8, fmt.Sprintf("%.2f €", float64(detail.Facture.TotalTTC)/100.0), "", 1, "R", false, 0, "")

	reglementY := totalsY + 24
	if reglementY > 235 {
		pdf.AddPage()
		reglementY = 20
	}

	pdf.SetXY(15, reglementY)
	pdf.SetFont("Arial", "", 10)
	modeReglement := string(detail.Facture.ModeReglement)
	switch detail.Facture.ModeReglement {
	case ReglementEspeces:
		modeReglement = "Espèce"
	case ReglementCB:
		modeReglement = "Carte Bancaire"
	case ReglementVirement:
		modeReglement = "Virement Bancaire"
	case ReglementCheque:
		modeReglement = "Chèque"
	}

	pdf.CellFormat(180, 6, fmt.Sprintf("Mode de règlement: %s", modeReglement), "", 1, "L", false, 0, "")

	echeance := detail.Facture.DateEcheance
	if echeance == "" {
		echeance = detail.Facture.DateFacture
	}
	pdf.SetX(15)
	pdf.CellFormat(180, 6, fmt.Sprintf("Échéance: %s", echeance), "", 1, "L", false, 0, "")

	conditions := detail.Facture.ConditionsReglement
	if strings.TrimSpace(conditions) == "" {
		conditions = detail.Facture.MentionsLegales
	}
	pdf.SetX(15)
	pdf.CellFormat(180, 6, fmt.Sprintf("Conditions de règlement / mentions complémentaires: %s", conditions), "", 1, "L", false, 0, "")

	pdf.SetY(272)
	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(100, 116, 139)

	denom := detail.Facture.EntrepriseDenomination
	if denom == "" {
		denom = "EXCELLENCE MATÉRIAUX"
	}
	forme := detail.Facture.EntrepriseFormeJuridique
	if forme == "" {
		forme = "SARL"
	}
	capitalEuros := detail.Facture.EntrepriseCapital / 100
	if capitalEuros == 0 {
		capitalEuros = 30000
	}
	siren := detail.Facture.EntrepriseSiren
	if siren == "" {
		siren = "109 741 124"
	}
	rcs := detail.Facture.EntrepriseRCS
	if rcs == "" {
		rcs = "Pontoise"
	}
	tvaIntra := detail.Facture.EntrepriseNumeroTVA
	if tvaIntra == "" {
		tvaIntra = "FR49 109741124"
	}

	ligneFooter1 := fmt.Sprintf("%s • %s au capital de %d € • SIREN %s • RCS %s", denom, forme, capitalEuros, siren, rcs)
	ligneFooter2 := fmt.Sprintf("• TVA intracommunautaire %s", tvaIntra)

	pdf.CellFormat(180, 4, ligneFooter1, "", 1, "C", false, 0, "")
	pdf.CellFormat(180, 4, ligneFooter2, "", 1, "C", false, 0, "")

	return pdf.OutputFileAndClose(cheminDestination)
}
