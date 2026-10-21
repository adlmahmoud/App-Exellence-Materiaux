package factures

import (
	"path/filepath"
	"testing"

	internaldb "exellence-materiaux/internal/db"
)

func TestCreerFactureAndBonWithLines(t *testing.T) {
	database, err := internaldb.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer database.Close()

	store := NewStore(database.DB)
	ligne := SaisieLigneFacture{
		Designation:    "Produit test",
		Quantite:       2,
		PrixUnitaireHT: 1200,
		TauxTVA:        2000,
	}

	facture, err := store.CreerFacture(NouvelleFacture{
		ClientNom: "Client test",
		Lignes:    []SaisieLigneFacture{ligne},
	})
	if err != nil {
		t.Fatalf("create facture: %v", err)
	}
	if len(facture.Lignes) != 1 || facture.Lignes[0].FactureID != facture.Facture.ID {
		t.Fatalf("facture line was not persisted with its parent: %#v", facture)
	}
	if facture.Facture.TotalHT != 2400 || facture.Facture.TotalTVA != 480 || facture.Facture.TotalTTC != 2880 {
		t.Fatalf("unexpected facture totals: HT=%d TVA=%d TTC=%d",
			facture.Facture.TotalHT, facture.Facture.TotalTVA, facture.Facture.TotalTTC)
	}

	bon, err := store.CreerBon(NouveauBon{
		TypeBon:   BonLivraison,
		ClientNom: "Client test",
		Lignes: []SaisieLigneBon{{
			Designation:    "Produit test",
			Quantite:       2,
			PrixUnitaireHT: 1200,
			TauxTVA:        2000,
		}},
	})
	if err != nil {
		t.Fatalf("create bon: %v", err)
	}
	persistedBon, err := store.GetBonByID(bon.Bon.ID)
	if err != nil {
		t.Fatalf("get created bon: %v", err)
	}
	if len(persistedBon.Lignes) != 1 || persistedBon.Lignes[0].BonID != persistedBon.Bon.ID {
		t.Fatalf("bon line was not persisted with its parent: %#v", persistedBon)
	}
	if bon.Bon.TotalHT != 2400 || bon.Bon.TotalTVA != 480 || bon.Bon.TotalTTC != 2880 {
		t.Fatalf("unexpected bon totals: HT=%d TVA=%d TTC=%d",
			bon.Bon.TotalHT, bon.Bon.TotalTVA, bon.Bon.TotalTTC)
	}
}
