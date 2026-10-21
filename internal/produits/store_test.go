package produits

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestStoreListSearch(t *testing.T) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	schema := `
	CREATE TABLE produits (
		id TEXT PRIMARY KEY,
		code_barres TEXT NOT NULL UNIQUE,
		reference TEXT NOT NULL,
		designation TEXT NOT NULL,
		description TEXT NOT NULL DEFAULT '',
		prix_achat_ht INTEGER NOT NULL DEFAULT 0,
		prix_vente_ht INTEGER NOT NULL DEFAULT 0,
		prix_vente_ttc INTEGER NOT NULL DEFAULT 0,
		taux_tva INTEGER NOT NULL DEFAULT 2000,
		stock_actuel INTEGER NOT NULL DEFAULT 0,
		stock_minimum INTEGER NOT NULL DEFAULT 0,
		unite TEXT NOT NULL DEFAULT 'piece',
		actif INTEGER NOT NULL DEFAULT 1,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);

	CREATE VIRTUAL TABLE produits_fts USING fts5(
		id UNINDEXED,
		code_barres,
		reference,
		designation,
		description,
		content='produits',
		content_rowid='rowid'
	);

	CREATE TRIGGER trg_produits_ai AFTER INSERT ON produits BEGIN
		INSERT INTO produits_fts(rowid, id, code_barres, reference, designation, description)
		VALUES (new.rowid, new.id, new.code_barres, new.reference, new.designation, new.description);
	END;

	CREATE TABLE mouvements_stock (
		id TEXT PRIMARY KEY,
		produit_id TEXT NOT NULL,
		type_mouvement TEXT NOT NULL,
		quantite INTEGER NOT NULL,
		stock_avant INTEGER NOT NULL,
		stock_apres INTEGER NOT NULL,
		prix_unitaire_ht INTEGER NOT NULL,
		reference_document TEXT NOT NULL,
		motif TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatal(err)
	}

	store := NewStore(db)

	_, err = store.Create(NouveauProduit{
		CodeBarres:   "3760123456789",
		Reference:    "CIM-01",
		Designation:  "Sac de ciment gris 35kg NF",
		PrixVenteHT:  650,
		StockInitial: 10,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Create(NouveauProduit{
		CodeBarres:   "3760999999999",
		Reference:    "PL-02",
		Designation:  "Planche de coffrage 27x200mm L.4m",
		PrixVenteHT:  850,
		StockInitial: 5,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.Create(NouveauProduit{
		CodeBarres:   "400500600",
		Reference:    "BET-PR",
		Designation:  "Sac de Béton prêt à l'emploi 25kg",
		PrixVenteHT:  550,
		StockInitial: 20,
	})
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		q   string
		exp int
	}{
		{"ciment", 1},
		{"ciment gris", 1},
		{"gris ciment", 1},
		{"beton", 1},
		{"Béton", 1},
		{"sac", 2},
		{"sac ciment", 1},
		{"sac beton", 1},
		{"beton sac", 1},
		{"3760", 2},
		{"12345", 1},
		{"CIM-01", 1},
		{"PL-02", 1},
		{"coffrage", 1},
		{"planche 4m", 1},
		{"4m planche", 1},
		{"25kg", 1},
	}

	for _, c := range cases {
		res, err := store.List(FiltreProduits{
			Recherche: c.q,
		})
		if err != nil {
			t.Errorf("List(%q) err: %v", c.q, err)
			continue
		}
		if res.Total != int64(c.exp) {
			t.Errorf("List(%q) = %d; want %d", c.q, res.Total, c.exp)
		}
	}
}
