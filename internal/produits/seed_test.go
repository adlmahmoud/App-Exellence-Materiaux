package produits

import (
	"path/filepath"
	"testing"

	"exellence-materiaux/internal/db"
)

func TestSeedSampleProducts(t *testing.T) {
	dbPath := filepath.Join("..", "..", "exellence.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open database error: %v", err)
	}
	defer database.Close()

	store := NewStore(database.DB)

	sampleProducts := []NouveauProduit{
		{
			CodeBarres:   "3760123400018",
			Reference:    "CIM-GRIS-35",
			Designation:  "Sac de ciment gris 35kg NF Calcia",
			Description:  "Ciment de haute qualité pour béton et mortier courant",
			PrixAchatHT:  480,
			PrixVenteHT:  650,
			TauxTVA:      2000,
			StockInitial: 120,
			StockMinimum: 20,
			Unite:        "sac",
		},
		{
			CodeBarres:   "3760123400025",
			Reference:    "BOIS-COF-4M",
			Designation:  "Planche de coffrage sapin 27x200mm L.4m",
			Description:  "Planche brute pour travaux de coffrage et charpente",
			PrixAchatHT:  620,
			PrixVenteHT:  850,
			TauxTVA:      2000,
			StockInitial: 45,
			StockMinimum: 15,
			Unite:        "piece",
		},
		{
			CodeBarres:   "3760123400032",
			Reference:    "OUT-BROU-100",
			Designation:  "Brouette de chantier 100L roue increvable Haemmerlin",
			Description:  "Châssis tubulaire renforcé, cuve peinte 100L",
			PrixAchatHT:  4500,
			PrixVenteHT:  6990,
			TauxTVA:      2000,
			StockInitial: 8,
			StockMinimum: 3,
			Unite:        "piece",
		},
		{
			CodeBarres:   "3760123400049",
			Reference:    "BET-UNIV-25",
			Designation:  "Sac de Béton universel prêt à l'emploi 25kg",
			Description:  "Béton sec prêt à gâcher pour dalles, poteaux et scellements",
			PrixAchatHT:  390,
			PrixVenteHT:  550,
			TauxTVA:      2000,
			StockInitial: 80,
			StockMinimum: 25,
			Unite:        "sac",
		},
		{
			CodeBarres:   "3760123400056",
			Reference:    "PLA-BA13-STD",
			Designation:  "Plaque de plâtre Placo BA13 standard 250x120cm",
			Description:  "Plaque de plâtre à bords amincis pour cloisons et plafonds",
			PrixAchatHT:  750,
			PrixVenteHT:  1090,
			TauxTVA:      2000,
			StockInitial: 60,
			StockMinimum: 20,
			Unite:        "plaque",
		},
		{
			CodeBarres:   "3760123400063",
			Reference:    "ISO-LV-100",
			Designation:  "Rouleau laine de verre IBR kraft ép. 100mm (5.4m²)",
			Description:  "Isolation thermique et acoustique des combles et toitures",
			PrixAchatHT:  1850,
			PrixVenteHT:  2590,
			TauxTVA:      2000,
			StockInitial: 3,
			StockMinimum: 10,
			Unite:        "rouleau",
		},
		{
			CodeBarres:   "3760123400070",
			Reference:    "BLOC-BET-20",
			Designation:  "Parpaing creux béton NF 20x20x50cm",
			Description:  "Bloc béton creux pour murs porteurs et clôtures",
			PrixAchatHT:  110,
			PrixVenteHT:  165,
			TauxTVA:      2000,
			StockInitial: 250,
			StockMinimum: 50,
			Unite:        "piece",
		},
		{
			CodeBarres:   "3760123400087",
			Reference:    "ARM-PAF-C",
			Designation:  "Treillis soudé PAF C anti-fissuration 2.4x1.2m",
			Description:  "Panneau d'armature pour dallage béton intérieur et extérieur",
			PrixAchatHT:  820,
			PrixVenteHT:  1200,
			TauxTVA:      2000,
			StockInitial: 2,
			StockMinimum: 8,
			Unite:        "panneau",
		},
		{
			CodeBarres:   "3760123400094",
			Reference:    "COL-SIKA-11FC",
			Designation:  "Mastic colle polyuréthane Sikaflex Pro-11 FC 300ml gris",
			Description:  "Colle et mastic d'étanchéité haute élasticité tout support",
			PrixAchatHT:  580,
			PrixVenteHT:  890,
			TauxTVA:      2000,
			StockInitial: 36,
			StockMinimum: 12,
			Unite:        "cartouche",
		},
		{
			CodeBarres:   "3760123400100",
			Reference:    "COL-CAR-C2S1",
			Designation:  "Mortier colle carrelage C2S1 extérieur/intérieur 25kg",
			Description:  "Mortier colle déformable haute performance tous carreaux",
			PrixAchatHT:  1420,
			PrixVenteHT:  2150,
			TauxTVA:      2000,
			StockInitial: 28,
			StockMinimum: 10,
			Unite:        "sac",
		},
	}

	for _, np := range sampleProducts {
		existing, err := store.GetByCodeBarres(np.CodeBarres)
		if err == nil && existing.ID != "" {
			t.Logf("Produit %s existe deja (%s)", np.CodeBarres, existing.Designation)
			continue
		}

		cree, err := store.Create(np)
		if err != nil {
			t.Errorf("Echec insertion produit %s: %v", np.Designation, err)
			continue
		}
		t.Logf("Insere avec succes : [%s] %s - Stock: %d (Min: %d)", cree.CodeBarres, cree.Designation, cree.StockActuel, cree.StockMinimum)
	}
}

func TestListAllProducts(t *testing.T) {
	dbPath := filepath.Join("..", "..", "exellence.db")
	database, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open database error: %v", err)
	}
	defer database.Close()

	store := NewStore(database.DB)
	pageResult, err := store.List(FiltreProduits{Limit: 100})
	if err != nil {
		t.Fatalf("List error: %v", err)
	}

	t.Logf("Nombre total de produits dans exellence.db: %d", len(pageResult.Produits))
	for i, p := range pageResult.Produits {
		t.Logf("#%d [%s] %s - Ref: %s - Prix TTC: %.2f EUR - Stock: %d (Min: %d)", i+1, p.CodeBarres, p.Designation, p.Reference, float64(p.PrixVenteTTC)/100.0, p.StockActuel, p.StockMinimum)
	}

	searchRes, err := store.List(FiltreProduits{Recherche: "ciment gris"})
	if err != nil {
		t.Fatalf("Search error: %v", err)
	}
	t.Logf("Recherche 'ciment gris' -> %d resultats trouves", len(searchRes.Produits))
	for _, p := range searchRes.Produits {
		t.Logf("Trouve: [%s] %s", p.Reference, p.Designation)
	}

	searchBarcode, err := store.List(FiltreProduits{Recherche: "3760123400063"})
	if err != nil {
		t.Fatalf("Search barcode error: %v", err)
	}
	t.Logf("Recherche code barre '3760123400063' -> %d resultats trouves", len(searchBarcode.Produits))
	for _, p := range searchBarcode.Produits {
		t.Logf("Trouve: [%s] %s", p.Reference, p.Designation)
	}
}
